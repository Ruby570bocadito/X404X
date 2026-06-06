package orchestrator

import (
	"sync"
	"time"
)

// EventType categorizes events in the system.
type EventType string

const (
	EventCampaignStarted    EventType = "campaign.started"
	EventCampaignPaused     EventType = "campaign.paused"
	EventCampaignCompleted  EventType = "campaign.completed"
	EventPhaseChanged       EventType = "phase.changed"
	EventAgentRegistered    EventType = "agent.registered"
	EventAgentCheckin       EventType = "agent.checkin"
	EventAgentDead          EventType = "agent.dead"
	EventDecisionMade       EventType = "decision.made"
	EventDecisionApproved   EventType = "decision.approved"
	EventKillChainUpdate    EventType = "killchain.update"
	EventVulnerabilityFound EventType = "vuln.found"
	EventExploitSuccess     EventType = "exploit.success"
	EventExploitFailure     EventType = "exploit.failure"
	EventCredentialCaptured EventType = "credential.captured"
	EventPersistenceSet     EventType = "persistence.set"
	EventLateralMove        EventType = "lateral.move"
	EventBlueAlert          EventType = "blue.alert"
	EventWildcard           EventType = "*"
)

// Event represents a system event.
type Event struct {
	Type       EventType   `json:"type"`
	CampaignID string      `json:"campaign_id,omitempty"`
	AgentID    string      `json:"agent_id,omitempty"`
	Timestamp  time.Time   `json:"timestamp"`
	Data       interface{} `json:"data,omitempty"`
}

// EventHandler is a callback for processing events.
type EventHandler func(event Event)

// EventBus provides a pub/sub event system.
type EventBus struct {
	subscribers map[EventType][]EventHandler
	mu          sync.RWMutex
}

func NewEventBus() *EventBus {
	return &EventBus{
		subscribers: make(map[EventType][]EventHandler),
	}
}

// Subscribe registers a handler for a specific event type.
// Use "*" to subscribe to all events.
func (eb *EventBus) Subscribe(eventType EventType, handler EventHandler) {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	eb.subscribers[eventType] = append(eb.subscribers[eventType], handler)
}

// Publish sends an event to all subscribers, including wildcard "*" subscribers.
func (eb *EventBus) Publish(event Event) {
	eb.mu.RLock()
	handlers := make([]EventHandler, 0)
	handlers = append(handlers, eb.subscribers[event.Type]...)
	handlers = append(handlers, eb.subscribers[EventType("*")]...)
	eb.mu.RUnlock()

	for _, handler := range handlers {
		go handler(event)
	}
}
