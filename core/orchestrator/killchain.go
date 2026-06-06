package orchestrator

import (
	"fmt"
	"sync"
	"time"

	"github.com/ruby570bocadito/x404x/shared/logger"
	"github.com/ruby570bocadito/x404x/shared/types"
)

// KillChainOrchestrator manages automatic phase transitions for campaigns.
// It listens to events from the agent and advances phases when conditions are met.
type KillChainOrchestrator struct {
	log       *logger.Logger
	orch      *Orchestrator
	mu        sync.Mutex
	transitions map[string]bool // campaignID → transition in progress
}

// NewKillChainOrchestrator creates a kill chain orchestrator.
func NewKillChainOrchestrator(log *logger.Logger, orch *Orchestrator) *KillChainOrchestrator {
	return &KillChainOrchestrator{
		log:         log,
		orch:        orch,
		transitions: make(map[string]bool),
	}
}

// ReconComplete should be called when recon data has been gathered.
// Automatically advances the campaign from Recon → Weaponization.
func (ko *KillChainOrchestrator) ReconComplete(campaignID string, hostsDiscovered int) error {
	return ko.tryAdvance(campaignID, types.PhaseRecon, types.PhaseWeaponization,
		fmt.Sprintf("recon complete: %d hosts discovered", hostsDiscovered))
}

// WeaponizeComplete should be called when a payload is built.
// Automatically advances Weaponization → Delivery.
func (ko *KillChainOrchestrator) WeaponizeComplete(campaignID string) error {
	return ko.tryAdvance(campaignID, types.PhaseWeaponization, types.PhaseDelivery,
		"payload compiled and ready")
}

// DeliveryComplete should be called when the agent is deployed on target.
// Automatically advances Delivery → Exploitation.
func (ko *KillChainOrchestrator) DeliveryComplete(campaignID, targetHost string) error {
	return ko.tryAdvance(campaignID, types.PhaseDelivery, types.PhaseExploitation,
		fmt.Sprintf("agent deployed on %s", targetHost))
}

// ExploitComplete should be called when root is obtained.
// Automatically advances Exploitation → Installation.
func (ko *KillChainOrchestrator) ExploitComplete(campaignID, privescVector string) error {
	return ko.tryAdvance(campaignID, types.PhaseExploitation, types.PhaseInstallation,
		fmt.Sprintf("privilege escalation: %s", privescVector))
}

// InstallComplete should be called when persistence is set.
// Automatically advances Installation → C2.
func (ko *KillChainOrchestrator) InstallComplete(campaignID string, methods []string) error {
	return ko.tryAdvance(campaignID, types.PhaseInstallation, types.PhaseCommandAndControl,
		fmt.Sprintf("persistence installed: %v", methods))
}

// C2Complete should be called when C2 channel is established.
// Automatically advances C2 → Actions on Objective.
func (ko *KillChainOrchestrator) C2Complete(campaignID string) error {
	return ko.tryAdvance(campaignID, types.PhaseCommandAndControl, types.PhaseActionsOnObjective,
		"encrypted C2 channel established")
}

// ObjectiveComplete marks the campaign as complete.
func (ko *KillChainOrchestrator) ObjectiveComplete(campaignID string) error {
	campaign, err := ko.orch.GetCampaign(campaignID)
	if err != nil {
		return err
	}

	campaign.Status = types.CampaignStatusCompleted
	campaign.Progress = 1.0
	campaign.CompletedAt = now()

	ko.log.Infof("campaign %s completed! goal: %s", campaignID, campaign.Goal)

	ko.orch.eventBus.Publish(Event{
		Type:       EventCampaignCompleted,
		CampaignID: campaignID,
		Data:       campaign,
	})

	return nil
}

func (ko *KillChainOrchestrator) tryAdvance(campaignID string, from, to types.KillChainPhase, reason string) error {
	// Prevent concurrent transitions
	ko.mu.Lock()
	if ko.transitions[campaignID] {
		ko.mu.Unlock()
		return fmt.Errorf("transition already in progress for campaign %s", campaignID)
	}
	ko.transitions[campaignID] = true
	ko.mu.Unlock()

	defer func() {
		ko.mu.Lock()
		delete(ko.transitions, campaignID)
		ko.mu.Unlock()
	}()

	campaign, err := ko.orch.GetCampaign(campaignID)
	if err != nil {
		return err
	}

	if campaign.Phase != from {
		return fmt.Errorf("expected phase %s, got %s", from, campaign.Phase)
	}

	if err := ko.orch.AdvancePhase(campaignID, to); err != nil {
		return fmt.Errorf("advancing phase: %w", err)
	}

	ko.log.Infof("kill chain auto-advance: %s → %s (campaign=%s, %s)", from, to, campaignID, reason)
	return nil
}

func now() *time.Time {
	t := time.Now()
	return &t
}
