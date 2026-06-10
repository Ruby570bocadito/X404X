package legacy

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type LegacyAdapter struct {
	RepoName    string
	SourcePath  string
	TargetPath  string
	Fused       bool
	Modules     []LegacyModule
}

type LegacyModule struct {
	Name    string
	Repo    string
	Path    string
	Handler string
	Fused   bool
}

var fusionTargets = []LegacyAdapter{
	{RepoName: "Pulse-C2", SourcePath: "core/c2/src/go/", TargetPath: "internal/c2/legacy/"},
	{RepoName: "Rise-Privilege", SourcePath: "core/privesc/", TargetPath: "internal/privesc/"},
	{RepoName: "Vault-Kernel", SourcePath: "core/kernel/", TargetPath: "internal/kernel/"},
	{RepoName: "Breach-Entry", SourcePath: "core/breach/", TargetPath: "internal/breach/"},
	{RepoName: "Horizon-Intel", SourcePath: "plugins/recon/", TargetPath: "internal/recon/"},
	{RepoName: "Specter-Terminal", SourcePath: "plugins/ai/specter/", TargetPath: "internal/ai/specter/"},
	{RepoName: "Apex-Automation", SourcePath: "plugins/ai/apex/", TargetPath: "internal/ai/apex/"},
	{RepoName: "Wormy-ML", SourcePath: "plugins/worm/", TargetPath: "internal/worm/"},
	{RepoName: "Link-Relay", SourcePath: "plugins/relay/", TargetPath: "internal/relay/"},
	{RepoName: "Titan-Operations", SourcePath: "plugins/operations/", TargetPath: "internal/ops/"},
	{RepoName: "BlueForge-Suite", SourcePath: "plugins/blue/", TargetPath: "internal/defense/"},
}

func NewFusionManager() *FusionManager {
	return &FusionManager{targets: fusionTargets}
}

type FusionManager struct {
	targets []LegacyAdapter
	fused   int
}

func (fm *FusionManager) ScanAll() []LegacyAdapter {
	for i := range fm.targets {
		if _, err := os.Stat(fm.targets[i].SourcePath); err == nil {
			fm.targets[i].Fused = fm.scanModules(&fm.targets[i])
		}
	}
	return fm.targets
}

func (fm *FusionManager) scanModules(adapter *LegacyAdapter) bool {
	exts := map[string]bool{".go": true, ".py": true, ".c": true, ".rs": true}
	count := 0
	filepath.Walk(adapter.SourcePath, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() { return nil }
		if exts[filepath.Ext(path)] {
			adapter.Modules = append(adapter.Modules, LegacyModule{
				Name: filepath.Base(path), Repo: adapter.RepoName, Path: path, Fused: true,
			})
			count++
		}
		return nil
	})
	fm.fused += count
	return count > 0
}

func (fm *FusionManager) GenerateFusionReport() string {
	fm.ScanAll()
	var report strings.Builder
	report.WriteString("=== X404X LEGACY FUSION REPORT ===\n")
	totalModules := 0
	for _, t := range fm.targets {
		status := "PENDING"
		if t.Fused { status = "FUSED" }
		fmt.Fprintf(&report, "  [%s] %s (%d modules)\n", status, t.RepoName, len(t.Modules))
		totalModules += len(t.Modules)
	}
	fmt.Fprintf(&report, "\nTotal: %d modules from %d repos\n", totalModules, len(fm.targets))
	fmt.Fprintf(&report, "Fusion status: %d/%d repos fused\n", fm.fused, len(fm.targets))
	reportPath := filepath.Join(os.TempDir(), "x404x_fusion_report.txt")
	os.WriteFile(reportPath, []byte(report.String()), 0644)
	return report.String()
}

func (fm *FusionManager) GetStatusJSON() string {
	fm.ScanAll()
	status := make([]map[string]interface{}, 0)
	for _, t := range fm.targets {
		status = append(status, map[string]interface{}{
			"repo": t.RepoName, "fused": t.Fused, "modules": len(t.Modules),
			"source": t.SourcePath, "target": t.TargetPath,
		})
	}
	data, _ := json.Marshal(map[string]interface{}{"targets": status, "fused_count": fm.fused})
	return string(data)
}

func init() { _ = exec.Command; _ = strings.Builder{}; _ = filepath.Walk }
