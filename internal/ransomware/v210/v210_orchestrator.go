package v210

import (
	"context"
	"encoding/json"
	"sync"
)

type V210Orchestrator struct {
	Config      *V210Config
	Apocalipsis *ApocalipsisEngine
	Phantom     *PhantomEvasionEngine
	Status      map[string]bool
	mu          sync.Mutex
}

func NewV210Orchestrator(cfg *V210Config) *V210Orchestrator {
	return &V210Orchestrator{
		Config:      cfg,
		Apocalipsis: NewApocalipsisEngine(cfg),
		Phantom:     NewPhantomEvasionEngine(cfg),
		Status:      make(map[string]bool),
	}
}

func (vo *V210Orchestrator) ExecuteAll(ctx context.Context) map[string]bool {
	vo.mu.Lock()
	defer vo.mu.Unlock()
	vo.Phantom.Initialize()
	vo.Status["phantom_evasion"] = true
	if vo.Config.ApocalipsisEnable {
		apocStatus := vo.Apocalipsis.ExecuteAll()
		for k, v := range apocStatus { vo.Status["apoc_"+k] = v }
	}
	return vo.Status
}

func (vo *V210Orchestrator) GetFullStatusJSON() string {
	data, _ := json.Marshal(map[string]interface{}{
		"status": vo.Status, "apocalipsis": vo.Apocalipsis.GetStatusJSON(),
		"phantom": vo.Phantom.GetStatusJSON(), "modules_executed": len(vo.Status),
	})
	return string(data)
}
