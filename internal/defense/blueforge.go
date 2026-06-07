package defense

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"
)

type BlueForgeEngine struct {
	Techniques   []ATTACKTechnique
	Detected     []string
	Undetected   []string
	CoveragePct  float64
	Score        float64
}

type ATTACKTechnique struct {
	ID     string
	Name   string
	Tactic string
	Used   bool
	Detected bool
}

var allTechniques = []ATTACKTechnique{
	{ID: "T1566", Name: "Spear Phishing", Tactic: "Initial Access"},
	{ID: "T1059", Name: "Command & Scripting", Tactic: "Execution"},
	{ID: "T1547", Name: "Boot/Logon AutoStart", Tactic: "Persistence"},
	{ID: "T1068", Name: "Exploitation for PrivEsc", Tactic: "Privilege Escalation"},
	{ID: "T1562", Name: "Impair Defenses", Tactic: "Defense Evasion"},
	{ID: "T1003", Name: "OS Credential Dumping", Tactic: "Credential Access"},
	{ID: "T1083", Name: "File/Directory Discovery", Tactic: "Discovery"},
	{ID: "T1021", Name: "Remote Services", Tactic: "Lateral Movement"},
	{ID: "T1041", Name: "Exfil over C2 Channel", Tactic: "Exfiltration"},
	{ID: "T1486", Name: "Data Encrypted for Impact", Tactic: "Impact"},
	{ID: "T1565", Name: "Data Manipulation", Tactic: "Impact"},
	{ID: "T1490", Name: "Inhibit System Recovery", Tactic: "Impact"},
	{ID: "T1543", Name: "Create/Modify System Process", Tactic: "Persistence"},
	{ID: "T1136", Name: "Create Account", Tactic: "Persistence"},
	{ID: "T1218", Name: "Signed Binary Proxy", Tactic: "Defense Evasion"},
	{ID: "T1055", Name: "Process Injection", Tactic: "Defense Evasion"},
	{ID: "T1070", Name: "Indicator Removal", Tactic: "Defense Evasion"},
	{ID: "T1574", Name: "Hijack Execution Flow", Tactic: "Persistence"},
	{ID: "T1036", Name: "Masquerading", Tactic: "Defense Evasion"},
	{ID: "T1027", Name: "Obfuscated Files", Tactic: "Defense Evasion"},
}

func NewBlueForgeEngine() *BlueForgeEngine {
	return &BlueForgeEngine{Techniques: allTechniques}
}

func (bf *BlueForgeEngine) MarkUsed(techniqueID string) {
	for i := range bf.Techniques {
		if bf.Techniques[i].ID == techniqueID {
			bf.Techniques[i].Used = true
			return
		}
	}
}

func (bf *BlueForgeEngine) SimulateDetection() {
	rand.Seed(time.Now().UnixNano())
	for i := range bf.Techniques {
		if bf.Techniques[i].Used {
			detected := rand.Float64() < 0.15
			bf.Techniques[i].Detected = detected
			if detected {
				bf.Detected = append(bf.Detected, bf.Techniques[i].ID)
			} else {
				bf.Undetected = append(bf.Undetected, bf.Techniques[i].ID)
			}
		}
	}
	usedCount := 0
	for _, t := range bf.Techniques {
		if t.Used { usedCount++ }
	}
	if usedCount > 0 {
		bf.CoveragePct = float64(len(bf.Undetected)) / float64(usedCount) * 100
	}
	bf.Score = bf.CoveragePct
}

func (bf *BlueForgeEngine) GenerateReport() string {
	bf.SimulateDetection()
	report := fmt.Sprintf("=== X404X ATT&CK COVERAGE REPORT ===\nTechniques Used: %d\nDetected: %d\nUndetected: %d\nCoverage: %.1f%%\n",
		len(bf.Undetected)+len(bf.Detected), len(bf.Detected), len(bf.Undetected), bf.CoveragePct)
	reportPath := "/tmp/x404x_attack_report.json"
	data, _ := json.MarshalIndent(map[string]interface{}{
		"score": bf.Score, "detected": bf.Detected, "undetected": bf.Undetected,
	}, "", "  ")
	os.WriteFile(reportPath, data, 0644)
	_ = reportPath
	return report
}
