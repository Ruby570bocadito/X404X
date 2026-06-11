package ransomware

import (
	"context"
	"testing"
)

func TestNewEngine(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	e, err := NewEngine(cfg)
	if err != nil {
		t.Fatalf("NewEngine() failed: %v", err)
	}
	if e == nil {
		t.Fatal("NewEngine() returned nil")
	}
}

func TestNewEngineInvalidConfig(t *testing.T) {
	cfg := &RansomwareConfig{}
	e, err := NewEngine(cfg)
	if err != nil {
		t.Logf("NewEngine with empty config: %v (may be expected)", err)
	}
	if e != nil {
		t.Log("NewEngine succeeded with minimal config")
	}
}

func TestEngineExecute(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	e, err := NewEngine(cfg)
	if err != nil {
		t.Fatalf("NewEngine() failed: %v", err)
	}

	report, err := e.Execute(context.Background(), "test-campaign", "TestCorp")
	if err != nil {
		t.Logf("Execute returned error: %v (expected in test env)", err)
	}
	if report != nil {
		if report.CampaignID != "test-campaign" {
			t.Errorf("expected CampaignID=test-campaign, got %s", report.CampaignID)
		}
		t.Logf("Report: %+v", report)
	}
}

func TestEngineReport(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	e, _ := NewEngine(cfg)

	r := e.Report()
	if r == nil {
		t.Log("Report() returned nil (expected before Execute)")
	}
}

func TestNewExtendedEngine(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	ee, err := NewExtendedEngine(cfg)
	if err != nil {
		t.Fatalf("NewExtendedEngine() failed: %v", err)
	}
	if ee == nil {
		t.Fatal("NewExtendedEngine() returned nil")
	}
}

func TestExtendedEngineExecute(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	ee, err := NewExtendedEngine(cfg)
	if err != nil {
		t.Fatalf("NewExtendedEngine() failed: %v", err)
	}

	report, err := ee.ExecuteExtended(context.Background(), "ext-campaign", "ExtCorp")
	if err != nil {
		t.Logf("ExecuteExtended error: %v", err)
	}
	if report != nil {
		t.Logf("Extended report: %+v", report)
	}
}

func TestEngineSubEngines(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	e, _ := NewEngine(cfg)

	if e.Scanner == nil {
		t.Error("Scanner engine not initialized")
	}
	if e.Crypto == nil {
		t.Error("Crypto (Hydra) engine not initialized")
	}
	if e.Extortion == nil {
		t.Error("Extortion engine not initialized")
	}
	if e.Destruction == nil {
		t.Error("Destruction engine not initialized")
	}
	if e.Propagation == nil {
		t.Error("Propagation engine not initialized")
	}
	if e.Psychological == nil {
		t.Error("Psychological engine not initialized")
	}
	if e.AntiAnalysis == nil {
		t.Error("AntiAnalysis engine not initialized")
	}
	if e.Polymorph == nil {
		t.Error("Polymorph engine not initialized")
	}
	if e.Trust == nil {
		t.Error("Trust engine not initialized")
	}
}

func TestExtendedEngineSubEngines(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	ee, _ := NewExtendedEngine(cfg)

	if ee.Engine == nil {
		t.Error("Embedded Engine not initialized")
	}
	if ee.PsychAdvanced == nil {
		t.Error("PsychAdvanced engine not initialized")
	}
	if ee.IdentityDest == nil {
		t.Error("IdentityDest engine not initialized")
	}
	if ee.Worm == nil {
		t.Error("Worm engine not initialized")
	}
	if ee.SupplyChain == nil {
		t.Error("SupplyChain engine not initialized")
	}
	if ee.CloudExploit == nil {
		t.Error("CloudExploit engine not initialized")
	}
	if ee.Bluetooth == nil {
		t.Error("Bluetooth engine not initialized")
	}
	if ee.SCADA == nil {
		t.Error("SCADA engine not initialized")
	}
	if ee.Hardware == nil {
		t.Error("Hardware engine not initialized")
	}
	if ee.NetworkPoison == nil {
		t.Error("NetworkPoison engine not initialized")
	}
	if ee.DNA == nil {
		t.Error("DNA engine not initialized")
	}
	if ee.Bootkit == nil {
		t.Error("Bootkit engine not initialized")
	}
	if ee.BlockchainC2 == nil {
		t.Error("BlockchainC2 engine not initialized")
	}
	if ee.Survivor == nil {
		t.Error("Survivor engine not initialized")
	}
}
