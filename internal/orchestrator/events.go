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

type subscriber struct {
	id      int64
	handler EventHandler
}

// EventBus provides a pub/sub event system.
type EventBus struct {
	subscribers map[EventType][]subscriber
	mu          sync.RWMutex
	nextID      int64
}

func NewEventBus() *EventBus {
	return &EventBus{
		subscribers: make(map[EventType][]subscriber),
	}
}

// Subscribe registers a handler for a specific event type.
// Use "*" to subscribe to all events. Returns an unsubscribe function.
func (eb *EventBus) Subscribe(eventType EventType, handler EventHandler) func() {
	eb.mu.Lock()
	id := eb.nextID
	eb.nextID++
	eb.subscribers[eventType] = append(eb.subscribers[eventType], subscriber{id: id, handler: handler})
	eb.mu.Unlock()

	return func() {
		eb.mu.Lock()
		defer eb.mu.Unlock()
		subs := eb.subscribers[eventType]
		for i, s := range subs {
			if s.id == id {
				eb.subscribers[eventType] = append(subs[:i], subs[i+1:]...)
				break
			}
		}
	}
}

// Publish sends an event to all subscribers, including wildcard "*" subscribers.
func (eb *EventBus) Publish(event Event) {
	eb.mu.RLock()
	handlers := make([]EventHandler, 0)
	for _, s := range eb.subscribers[event.Type] {
		handlers = append(handlers, s.handler)
	}
	for _, s := range eb.subscribers[EventType("*")] {
		handlers = append(handlers, s.handler)
	}
	eb.mu.RUnlock()

	for _, handler := range handlers {
		go handler(event)
	}
}
