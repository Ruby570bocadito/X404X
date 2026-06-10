package v26

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type AINegotiator struct {
	Config        *V26Config
	Conversation  []NegotiationTurn `json:"conversation"`
	TargetInfo    *TargetProfile    `json:"target_info"`
	NegotiationPhase string         `json:"negotiation_phase"`
	RescateBase   float64           `json:"rescate_base"`
	RescateActual float64           `json:"rescate_actual"`
	Deadline      time.Time         `json:"deadline"`
}

type NegotiationTurn struct {
	Role      string    `json:"role"`
	Message   string    `json:"message"`
	Sentiment string    `json:"sentiment"`
	Strategy  string    `json:"strategy"`
	Timestamp time.Time `json:"timestamp"`
}

type TargetProfile struct {
	Company      string  `json:"company"`
	Industry     string  `json:"industry"`
	Revenue      string  `json:"revenue"`
	Employees    int     `json:"employees"`
	CEO          string  `json:"ceo"`
	DataExfilGB  float64 `json:"data_exfil_gb"`
	Insurance    bool    `json:"insurance"`
	PriorBreaches bool   `json:"prior_breaches"`
}

type LLMResponse struct {
	Message  string `json:"message"`
	Sentiment string `json:"sentiment"`
	Strategy string `json:"strategy"`
	Tone     string `json:"tone"`
}

var negotiationTemplates = map[string][]string{
	"initial_demand": {
		"Your network has been fully compromised. %d GB of sensitive data has been exfiltrated. All workstations encrypted with military-grade algorithms. Payment of %.2f BTC is required within %d hours. Contact: %s.onion",
	},
	"counter_reject": {
		"You're negotiating? That's cute. Every hour you delay, we publish another 10GB of your data to journalists. The clock is ticking, %s. Next batch drops in %d hours.",
	},
	"partial_accept": {
		"We acknowledge your counter-offer. As a gesture of goodwill, we will delay the next data release by 12 hours. But %.2f BTC is our floor. Your insurance won't cover reputational damage, %s.",
	},
	"psychological_pressure": {
		"Your CFO just received an email with all board meeting minutes from Q3. Tomorrow, every investor gets them too. How much is your stock price worth, %s?",
	},
	"fake_compassion": {
		"We understand this is difficult. We're not monsters. Pay %.2f BTC and we'll even include a free security audit of your vulnerabilities. Consider it... consulting.",
	},
	"deadline_ultimatum": {
		"TIME'S UP, %s. Final offer: %.2f BTC in the next 2 hours, or we nuke everything. No more games. No more delays. Your data. Your choice.",
	},
	"insurance_hook": {
		"We know you have cyber insurance with %s. Your premium will triple after this. Pay us %.2f BTC now, and we'll leave a clean exit — no public disclosure. Your insurance stays intact.",
	},
}

func NewAINegotiator(cfg *V26Config) *AINegotiator {
	return &AINegotiator{
		Config:        cfg,
		RescateBase:   5000000,
		RescateActual: 5000000,
		Deadline:      time.Now().Add(48 * time.Hour),
		NegotiationPhase: "initial",
		TargetInfo: &TargetProfile{
			Company: "Target Corporation",
			Industry: "Technology",
			Revenue: "500M",
			Employees: 2000,
			CEO: "Unknown",
			DataExfilGB: 150,
			Insurance: true,
			PriorBreaches: false,
		},
	}
}

func (an *AINegotiator) GenerateResponse(humanMessage string) *LLMResponse {
	an.Conversation = append(an.Conversation, NegotiationTurn{
		Role: "human", Message: humanMessage,
		Sentiment: an.detectSentiment(humanMessage),
		Timestamp: time.Now(),
	})

	var response *LLMResponse

	switch {
	case strings.Contains(strings.ToLower(humanMessage), "who are you") ||
		strings.Contains(strings.ToLower(humanMessage), "what do you want"):
		response = an.initialDemand()
	case strings.Contains(strings.ToLower(humanMessage), "too much") ||
		strings.Contains(strings.ToLower(humanMessage), "can't pay") ||
		strings.Contains(strings.ToLower(humanMessage), "lower"):
		an.RescateActual *= 0.85
		response = an.counterOffer()
	case strings.Contains(strings.ToLower(humanMessage), "need more time") ||
		strings.Contains(strings.ToLower(humanMessage), "deadline") ||
		strings.Contains(strings.ToLower(humanMessage), "extension"):
		response = an.psychologicalPressure()
	case strings.Contains(strings.ToLower(humanMessage), "police") ||
		strings.Contains(strings.ToLower(humanMessage), "fbi") ||
		strings.Contains(strings.ToLower(humanMessage), "authorities"):
		response = an.escalateThreat()
	case strings.Contains(strings.ToLower(humanMessage), "data") ||
		strings.Contains(strings.ToLower(humanMessage), "files") ||
		strings.Contains(strings.ToLower(humanMessage), "decrypt"):
		response = an.dataLeverage()
	case strings.Contains(strings.ToLower(humanMessage), "insurance") ||
		strings.Contains(strings.ToLower(humanMessage), "insurer"):
		response = an.insuranceHook()
	default:
		response = an.automatedLLMResponse(humanMessage)
	}

	an.Conversation = append(an.Conversation, NegotiationTurn{
		Role: "x404x", Message: response.Message,
		Sentiment: response.Sentiment, Strategy: response.Strategy,
		Timestamp: time.Now(),
	})

	return response
}

func (an *AINegotiator) initialDemand() *LLMResponse {
	tpl := negotiationTemplates["initial_demand"][0]
	msg := fmt.Sprintf(tpl, int(an.TargetInfo.DataExfilGB), an.RescateActual/1000000, 48, "x404x")
	return &LLMResponse{Message: msg, Sentiment: "cold", Strategy: "shock_and_awe", Tone: "professional"}
}

func (an *AINegotiator) counterOffer() *LLMResponse {
	tpl := negotiationTemplates["counter_reject"][0]
	msg := fmt.Sprintf(tpl, an.TargetInfo.CEO, 12)
	return &LLMResponse{Message: msg, Sentiment: "dismissive", Strategy: "anchoring", Tone: "arrogant"}
}

func (an *AINegotiator) psychologicalPressure() *LLMResponse {
	tpl := negotiationTemplates["psychological_pressure"][0]
	msg := fmt.Sprintf(tpl, an.TargetInfo.CEO)
	return &LLMResponse{Message: msg, Sentiment: "threatening", Strategy: "fear_escalation", Tone: "menacing"}
}

func (an *AINegotiator) escalateThreat() *LLMResponse {
	msg := fmt.Sprintf("Calling the authorities? We anticipated that. Your FBI contact is Agent %s. They received your 2019 tax returns 10 minutes ago. Want to guess what happens next?", randomAgentName())
	return &LLMResponse{Message: msg, Sentiment: "amused", Strategy: "preemptive_strike", Tone: "mocking"}
}

func (an *AINegotiator) dataLeverage() *LLMResponse {
	msg := fmt.Sprintf("Your data? Let's see... board_minutes_%s.pdf, employee_ssn_dump.csv, customer_credit_cards.sql. Which one should we post first? %.2f BTC makes them all disappear.", time.Now().AddDate(-1, 0, 0).Format("Jan"), an.RescateActual/1000000)
	return &LLMResponse{Message: msg, Sentiment: "cold", Strategy: "proof_of_life", Tone: "clinical"}
}

func (an *AINegotiator) insuranceHook() *LLMResponse {
	tpl := negotiationTemplates["insurance_hook"][0]
	insurer := "CyberGuard Insurance"
	msg := fmt.Sprintf(tpl, insurer, an.RescateActual/1000000)
	return &LLMResponse{Message: msg, Sentiment: "conspiratorial", Strategy: "mutual_benefit", Tone: "smooth"}
}

func (an *AINegotiator) automatedLLMResponse(humanInput string) *LLMResponse {
	msg := fmt.Sprintf("Interesting. Let me be direct: %.2f BTC within %d hours, or %.0f GB of your data goes public. Your move, %s.",
		an.RescateActual/1000000,
		int(time.Until(an.Deadline).Hours()),
		an.TargetInfo.DataExfilGB*0.1,
		an.TargetInfo.CEO)
	return &LLMResponse{Message: msg, Sentiment: "impatient", Strategy: "re_focus", Tone: "direct"}
}

func (an *AINegotiator) detectSentiment(msg string) string {
	lower := strings.ToLower(msg)
	if strings.Contains(lower, "please") || strings.Contains(lower, "help") {
		return "desperate"
	}
	if strings.Contains(lower, "fuck") || strings.Contains(lower, "bastard") {
		return "angry"
	}
	if strings.Contains(lower, "police") || strings.Contains(lower, "fbi") {
		return "defiant"
	}
	if strings.Contains(lower, "pay") || strings.Contains(lower, "money") {
		return "negotiating"
	}
	return "confused"
}

func (an *AINegotiator) SaveConversationLog() string {
	logPath := filepath.Join(os.TempDir(), "x404x_negotiation_log.json")
	data, _ := json.MarshalIndent(an.Conversation, "", "  ")
	os.WriteFile(logPath, data, 0644)
	return logPath
}

func (an *AINegotiator) GetStatusJSON() string {
	data, _ := json.Marshal(map[string]interface{}{
		"phase": an.NegotiationPhase,
		"rescate_base": an.RescateBase,
		"rescate_actual": an.RescateActual,
		"deadline_hours": int(time.Until(an.Deadline).Hours()),
		"conversation_turns": len(an.Conversation),
	})
	return string(data)
}

func randomAgentName() string {
	names := []string{"Rodriguez", "Johnson", "Park", "Ivanov", "Chen", "Mueller"}
	return names[rand.Intn(len(names))]
}
