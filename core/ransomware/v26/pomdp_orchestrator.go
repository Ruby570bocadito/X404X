package v26

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type POMDPOrchestrator struct {
	Config        *V26Config
	StateSpace    []POMDPState     `json:"state_space"`
	ActionSpace   []POMDPAction    `json:"action_space"`
	BeliefState   map[string]float64 `json:"belief_state"`
	TransitionModel map[string]map[string]float64 `json:"-"`
	GodMode       bool             `json:"god_mode"`
	ChaosInjections int            `json:"chaos_injections"`
	mu            sync.Mutex
}

type POMDPState struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	DetectionProb float64 `json:"detection_prob"`
	StealthLevel  float64 `json:"stealth_level"`
	Resources    float64 `json:"resources"`
	AvailActions []string `json:"available_actions"`
}

type POMDPAction struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Cost     float64 `json:"cost"`
	Risk     float64 `json:"risk"`
	Reward   float64 `json:"reward"`
	Duration float64 `json:"duration_ms"`
	Preconditions []string `json:"preconditions"`
}

type ChaosInjection struct {
	Type      string `json:"type"`
	Target    string `json:"target"`
	Effect    string `json:"effect"`
	Triggered bool   `json:"triggered"`
	Planned   bool   `json:"planned"`
}

type POMDPDecision struct {
	Action      string  `json:"action"`
	Confidence  float64 `json:"confidence"`
	ExpectedReward float64 `json:"expected_reward"`
	RiskLevel   string  `json:"risk_level"`
	Rationale   string  `json:"rationale"`
	Timestamp   time.Time `json:"timestamp"`
}

func NewPOMDPOrchestrator(cfg *V26Config) *POMDPOrchestrator {
	po := &POMDPOrchestrator{
		Config:       cfg,
		BeliefState:  make(map[string]float64),
		TransitionModel: make(map[string]map[string]float64),
	}
	po.initStateSpace()
	po.initActionSpace()
	po.initBeliefState()
	po.initTransitionModel()
	return po
}

func (po *POMDPOrchestrator) initStateSpace() {
	po.StateSpace = []POMDPState{
		{ID: "undetected_high_privs", Name: "Undetected (High Privs)", DetectionProb: 0.05, StealthLevel: 0.95, Resources: 1.0, AvailActions: []string{"all"}},
		{ID: "undetected_low_privs", Name: "Undetected (Low Privs)", DetectionProb: 0.15, StealthLevel: 0.85, Resources: 0.5, AvailActions: []string{"scan", "recon", "exfil"}},
		{ID: "partial_detection", Name: "Partial Detection", DetectionProb: 0.40, StealthLevel: 0.60, Resources: 0.7, AvailActions: []string{"evade", "mutate", "hide"}},
		{ID: "soc_alerted", Name: "SOC Alerted", DetectionProb: 0.70, StealthLevel: 0.30, Resources: 0.3, AvailActions: []string{"evacuate", "self_destruct", "deploy_countermeasure"}},
		{ID: "fully_detected", Name: "Fully Detected", DetectionProb: 0.95, StealthLevel: 0.05, Resources: 0.1, AvailActions: []string{"scorched_earth", "apocalypse"}},
	}
}

func (po *POMDPOrchestrator) initActionSpace() {
	po.ActionSpace = []POMDPAction{
		{ID: "encrypt", Name: "Encrypt Files", Cost: 0.3, Risk: 0.4, Reward: 0.8, Duration: 5000, Preconditions: []string{"file_access"}},
		{ID: "exfil", Name: "Exfiltrate Data", Cost: 0.2, Risk: 0.3, Reward: 0.6, Duration: 3000, Preconditions: []string{"network"}},
		{ID: "persist", Name: "Install Persistence", Cost: 0.1, Risk: 0.2, Reward: 0.9, Duration: 2000, Preconditions: []string{"elevated"}},
		{ID: "propagate", Name: "Propagate Laterally", Cost: 0.5, Risk: 0.7, Reward: 0.95, Duration: 8000, Preconditions: []string{"network", "creds"}},
		{ID: "evade", Name: "Evade Detection", Cost: 0.1, Risk: 0.0, Reward: 0.1, Duration: 1000, Preconditions: []string{}},
		{ID: "mutate", Name: "Mutate Binary", Cost: 0.15, Risk: 0.05, Reward: 0.2, Duration: 3000, Preconditions: []string{}},
		{ID: "hide", Name: "Hide & Wait", Cost: 0.0, Risk: 0.0, Reward: 0.0, Duration: 60000, Preconditions: []string{}},
		{ID: "scorched_earth", Name: "Scorched Earth", Cost: 1.0, Risk: 1.0, Reward: 1.0, Duration: 10000, Preconditions: []string{}},
		{ID: "chaos_inject", Name: "Inject Chaos", Cost: 0.05, Risk: 0.1, Reward: 0.3, Duration: 500, Preconditions: []string{"god_mode"}},
	}
}

func (po *POMDPOrchestrator) initBeliefState() {
	for _, s := range po.StateSpace {
		po.BeliefState[s.ID] = 1.0 / float64(len(po.StateSpace))
	}
	po.BeliefState["undetected_low_privs"] = 0.5
	po.BeliefState["undetected_high_privs"] = 0.3
	po.BeliefState["partial_detection"] = 0.15
	po.BeliefState["soc_alerted"] = 0.04
	po.BeliefState["fully_detected"] = 0.01
}

func (po *POMDPOrchestrator) initTransitionModel() {
	for _, s := range po.StateSpace {
		po.TransitionModel[s.ID] = make(map[string]float64)
		total := 0.0
		for _, t := range po.StateSpace {
			if s.ID == t.ID {
				po.TransitionModel[s.ID][t.ID] = 0.7
			} else {
				po.TransitionModel[s.ID][t.ID] = 0.3 / float64(len(po.StateSpace)-1)
			}
			total += po.TransitionModel[s.ID][t.ID]
		}
		_ = total
	}
}

func (po *POMDPOrchestrator) Decide(ctx context.Context) *POMDPDecision {
	po.mu.Lock()
	defer po.mu.Unlock()

	if po.GodMode {
		po.injectChaos()
	}

	bestAction := po.ActionSpace[0]
	bestScore := -1.0

	for _, action := range po.ActionSpace {
		if action.ID == "chaos_inject" && !po.GodMode {
			continue
		}
		score := po.expectedValue(action) * po.beliefWeight(action)
		if score > bestScore {
			bestScore = score
			bestAction = action
		}
	}

	riskLevel := "low"
	if bestAction.Risk > 0.5 {
		riskLevel = "high"
	} else if bestAction.Risk > 0.3 {
		riskLevel = "medium"
	}

	detectionProb := po.currentDetectionProb()
	rationale := fmt.Sprintf("detection_prob=%.2f belief_high_privs=%.2f action_score=%.2f",
		detectionProb, po.BeliefState["undetected_high_privs"], bestScore)

	if po.GodMode {
		rationale += " [CHAOS MODE ACTIVE]"
	}

	return &POMDPDecision{
		Action:      bestAction.Name,
		Confidence:  bestScore,
		ExpectedReward: bestAction.Reward,
		RiskLevel:   riskLevel,
		Rationale:   rationale,
		Timestamp:   time.Now(),
	}
}

func (po *POMDPOrchestrator) expectedValue(action POMDPAction) float64 {
	reward := action.Reward * (1.0 - po.currentDetectionProb())
	riskPenalty := action.Risk * po.currentDetectionProb()
	return reward - riskPenalty - action.Cost
}

func (po *POMDPOrchestrator) beliefWeight(action POMDPAction) float64 {
	if len(action.Preconditions) == 0 {
		return 1.0
	}
	hasNeeded := false
	for _, prec := range action.Preconditions {
		if prec == "file_access" || prec == "network" {
			hasNeeded = true
		}
	}
	if !hasNeeded {
		return 0.1
	}
	return 1.0
}

func (po *POMDPOrchestrator) currentDetectionProb() float64 {
	prob := 0.0
	for _, s := range po.StateSpace {
		prob += s.DetectionProb * po.BeliefState[s.ID]
	}
	return prob
}

func (po *POMDPOrchestrator) UpdateBelief(observation string, success bool) {
	po.mu.Lock()
	defer po.mu.Unlock()

	decay := 0.95
	for _, s := range po.StateSpace {
		if success {
			if s.ID == "undetected_high_privs" || s.ID == "undetected_low_privs" {
				po.BeliefState[s.ID] *= 1.05
			} else {
				po.BeliefState[s.ID] *= decay
			}
		} else {
			if s.ID == "fully_detected" || s.ID == "soc_alerted" {
				po.BeliefState[s.ID] *= 1.1
			} else {
				po.BeliefState[s.ID] *= decay
			}
		}
	}

	total := 0.0
	for _, v := range po.BeliefState {
		total += v
	}
	for k := range po.BeliefState {
		po.BeliefState[k] /= total
	}
}

func (po *POMDPOrchestrator) EnableGodMode() {
	po.mu.Lock()
	defer po.mu.Unlock()
	po.GodMode = true
}

func (po *POMDPOrchestrator) injectChaos() {
	if rand.Float64() > 0.3 {
		return
	}

	injections := []ChaosInjection{
		{Type: "fake_av_alert", Target: "SOC", Effect: "Simulate detection of random malware to test SOC response", Triggered: true, Planned: true},
		{Type: "fake_system_crash", Target: "server_01", Effect: "Simulate crash of non-critical server to observe failover", Triggered: true, Planned: true},
		{Type: "decoy_movement", Target: "finance_dept", Effect: "Generate fake lateral movement noise to distract", Triggered: true, Planned: true},
		{Type: "fake_data_leak", Target: "pastebin", Effect: "Post fake data leak to observe PR team response", Triggered: true, Planned: true},
		{Type: "phantom_beacon", Target: "C2", Effect: "Spawn fake C2 beacon to test network monitoring", Triggered: true, Planned: true},
		{Type: "siren_test", Target: "all", Effect: "Fake Windows Defender popup on all screens", Triggered: true, Planned: true},
	}

	idx := rand.Intn(len(injections))
	po.ChaosInjections++
	_ = injections[idx]
}

func (po *POMDPOrchestrator) PlanBSwitch(ctx context.Context) *POMDPDecision {
	po.mu.Lock()
	po.BeliefState["fully_detected"] = 0.5
	po.BeliefState["soc_alerted"] = 0.4
	po.BeliefState["partial_detection"] = 0.1
	po.mu.Unlock()

	return po.Decide(ctx)
}

func (po *POMDPOrchestrator) GetStatusJSON() string {
	po.mu.Lock()
	defer po.mu.Unlock()
	data, _ := json.Marshal(map[string]interface{}{
		"god_mode":   po.GodMode,
		"chaos_injections": po.ChaosInjections,
		"detection_prob": po.currentDetectionProb(),
		"belief_undetected": po.BeliefState["undetected_high_privs"] + po.BeliefState["undetected_low_privs"],
		"states": len(po.StateSpace),
		"actions": len(po.ActionSpace),
	})
	return string(data)
}
