package v27

import (
	"context"
	"testing"
)

func TestDefaultV27Config(t *testing.T) {
	cfg := DefaultV27Config()
	if cfg == nil {
		t.Fatal("DefaultV27Config returned nil")
	}
	if cfg.Simulation != true {
		t.Error("Simulation should default to true")
	}
}

func TestNewV27Orchestrator(t *testing.T) {
	cfg := DefaultV27Config()
	vo := NewV27Orchestrator(cfg)
	if vo == nil {
		t.Fatal("NewV27Orchestrator() returned nil")
	}
}

func TestV27OrchestratorExecute(t *testing.T) {
	cfg := DefaultV27Config()
	vo := NewV27Orchestrator(cfg)

	result := vo.ExecuteAll(context.Background())
	if result == nil {
		t.Fatal("ExecuteAll returned nil")
	}
	t.Logf("V27 execution results: %v", result)
}

func TestV27GetFullStatusJSON(t *testing.T) {
	cfg := DefaultV27Config()
	vo := NewV27Orchestrator(cfg)

	json := vo.GetFullStatusJSON()
	if json == "" {
		t.Error("GetFullStatusJSON returned empty")
	}
	t.Logf("V27 JSON (truncated): %.200s", json)
}

func TestNewUEFIBootkitEngine(t *testing.T) {
	cfg := DefaultV27Config()
	uefi := NewUEFIBootkitEngine(cfg)
	if uefi == nil {
		t.Fatal("NewUEFIBootkitEngine() returned nil")
	}
}

func TestNewHypervisorEngine(t *testing.T) {
	cfg := DefaultV27Config()
	hv := NewHypervisorEngine(cfg)
	if hv == nil {
		t.Fatal("NewHypervisorEngine() returned nil")
	}
}

func TestNewPCIeRootkitEngine(t *testing.T) {
	cfg := DefaultV27Config()
	pcie := NewPCIeRootkitEngine(cfg)
	if pcie == nil {
		t.Fatal("NewPCIeRootkitEngine() returned nil")
	}
}

func TestNewKernelInstrumentEngine(t *testing.T) {
	cfg := DefaultV27Config()
	ki := NewKernelInstrumentEngine(cfg)
	if ki == nil {
		t.Fatal("NewKernelInstrumentEngine() returned nil")
	}
}

func TestNewSecureBootBypassEngine(t *testing.T) {
	cfg := DefaultV27Config()
	sbb := NewSecureBootBypassEngine(cfg)
	if sbb == nil {
		t.Fatal("NewSecureBootBypassEngine() returned nil")
	}
}

func TestNewPhishingInfraEngine(t *testing.T) {
	cfg := DefaultV27Config()
	pi := NewPhishingInfraEngine(cfg)
	if pi == nil {
		t.Fatal("NewPhishingInfraEngine() returned nil")
	}
}

func TestNewSpearPhishAIEngine(t *testing.T) {
	cfg := DefaultV27Config()
	sp := NewSpearPhishAIEngine(cfg)
	if sp == nil {
		t.Fatal("NewSpearPhishAIEngine() returned nil")
	}
}

func TestNewSmishingEngine(t *testing.T) {
	cfg := DefaultV27Config()
	sm := NewSmishingEngine(cfg)
	if sm == nil {
		t.Fatal("NewSmishingEngine() returned nil")
	}
}

func TestNewVishingEngine(t *testing.T) {
	cfg := DefaultV27Config()
	ve := NewVishingEngine(cfg)
	if ve == nil {
		t.Fatal("NewVishingEngine() returned nil")
	}
}

func TestNewAntiPhishEvasionEngine(t *testing.T) {
	cfg := DefaultV27Config()
	ape := NewAntiPhishEvasionEngine(cfg)
	if ape == nil {
		t.Fatal("NewAntiPhishEvasionEngine() returned nil")
	}
}

func TestPhishingInfraGenerateDGA(t *testing.T) {
	cfg := DefaultV27Config()
	pi := NewPhishingInfraEngine(cfg)

	domains := pi.GenerateDGADomains(5)
	if len(domains) == 0 {
		t.Log("DGA domains empty (expected in test env)")
	} else {
		t.Logf("DGA domains: %v", domains)
	}
}

func TestSpearPhishAIGenerateLure(t *testing.T) {
	cfg := DefaultV27Config()
	sp := NewSpearPhishAIEngine(cfg)

	lure := sp.GenerateLureWithLLM()
	if lure == "" {
		t.Log("GenerateLure empty (expected if no LLM)")
	} else {
		t.Logf("Phishing lure: %.100s", lure)
	}
}

func TestVishingEngineCloneVoice(t *testing.T) {
	cfg := DefaultV27Config()
	ve := NewVishingEngine(cfg)

	result := ve.CloneVoice("/tmp/sample.wav")
	t.Logf("Voice clone result: %v", result)
}
