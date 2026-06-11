package blockz

import (
	"context"
	"testing"
)

func TestDefaultBlockZConfig(t *testing.T) {
	cfg := DefaultBlockZConfig()
	if cfg == nil {
		t.Fatal("DefaultBlockZConfig returned nil")
	}
	if !cfg.Enabled {
		t.Log("BlockZ disabled by default")
	}
	if cfg.Simulation != true {
		t.Error("Simulation should default to true")
	}
	if cfg.DeadMansHours != 48 {
		t.Errorf("expected DeadMansHours=48, got %d", cfg.DeadMansHours)
	}
}

func TestNewBlockZOrchestrator(t *testing.T) {
	cfg := DefaultBlockZConfig()
	bz := NewBlockZOrchestrator(cfg)
	if bz == nil {
		t.Fatal("NewBlockZOrchestrator() returned nil")
	}
}

func TestBlockZOrchestratorSubEngines(t *testing.T) {
	cfg := DefaultBlockZConfig()
	bz := NewBlockZOrchestrator(cfg)

	if bz.GeneticEvo == nil {
		t.Error("GeneticEvo engine not initialized")
	}
	if bz.Deepfake == nil {
		t.Error("Deepfake engine not initialized")
	}
	if bz.SCADACovert == nil {
		t.Error("SCADACovert engine not initialized")
	}
	if bz.FirmwareWorm == nil {
		t.Error("FirmwareWorm engine not initialized")
	}
	if bz.MedicalAttack == nil {
		t.Error("MedicalAttack engine not initialized")
	}
	if bz.ModelPoison == nil {
		t.Error("ModelPoison engine not initialized")
	}
	if bz.Disinfo == nil {
		t.Error("Disinfo engine not initialized")
	}
	if bz.AirGap == nil {
		t.Error("AirGap engine not initialized")
	}
	if bz.PostQuantum == nil {
		t.Error("PostQuantum engine not initialized")
	}
	if bz.DeadMan == nil {
		t.Error("DeadMan engine not initialized")
	}
	if bz.FalseFlag == nil {
		t.Error("FalseFlag engine not initialized")
	}
	if bz.EDRControl == nil {
		t.Error("EDRControl engine not initialized")
	}
	if bz.Financial == nil {
		t.Error("Financial engine not initialized")
	}
	if bz.IoTChain == nil {
		t.Error("IoTChain engine not initialized")
	}
}

func TestBlockZExecute(t *testing.T) {
	cfg := DefaultBlockZConfig()
	bz := NewBlockZOrchestrator(cfg)

	report, err := bz.ExecuteBlockZ(context.Background(), "TestCorp")
	if err != nil {
		t.Logf("ExecuteBlockZ error: %v (expected in simulation)", err)
	}
	if report != nil {
		if report.Success {
			t.Logf("BlockZ report: %d modules executed", report.ModulesExecuted)
		}
		t.Logf("Report: gen=%d fitness=%.2f deepfakes=%d scada=%d firmware=%d",
			report.GenerationCount, report.BestFitness, report.DeepfakesGenerated,
			report.SCADAGradual, report.FirmwareInfections)
	}
}

func TestBlockZGetFullStatusJSON(t *testing.T) {
	cfg := DefaultBlockZConfig()
	bz := NewBlockZOrchestrator(cfg)

	json := bz.GetFullStatusJSON()
	if json == "" {
		t.Error("GetFullStatusJSON returned empty")
	}
	t.Logf("BlockZ JSON status (truncated): %.100s...", json)
}

func TestBlockZDefaultConfigFields(t *testing.T) {
	cfg := DefaultBlockZConfig()

	if cfg.DeepfakeModelPath == "" {
		t.Log("DeepfakeModelPath empty (default)")
	}
	if cfg.KyberVariant == "" {
		t.Log("KyberVariant empty (default)")
	}
	if cfg.UltrasoundFreq <= 0 {
		t.Logf("UltrasoundFreq: %d (default)", cfg.UltrasoundFreq)
	}
	if cfg.APTImpersonate == "" {
		t.Log("APTImpersonate empty (default)")
	}
}

func TestGeneticEvolutionEngine(t *testing.T) {
	cfg := DefaultBlockZConfig()
	bz := NewBlockZOrchestrator(cfg)

	if bz.GeneticEvo != nil {
		agents := bz.GeneticEvo.Evolve()
		if len(agents) == 0 {
			t.Log("Evolve returned 0 agents (expected in simulation)")
		}
	}
}

func TestPostQuantumEngine(t *testing.T) {
	cfg := DefaultBlockZConfig()
	bz := NewBlockZOrchestrator(cfg)

	if bz.PostQuantum != nil {
		kp := bz.PostQuantum.GenerateKyberKeypair()
		if kp == nil {
			t.Log("GenerateKyberKeypair returned nil (expected if no kyber lib)")
		} else {
			t.Logf("Kyber keypair generated: pub=%d bytes, sec=%d bytes",
				len(kp.Public), len(kp.Secret))
		}
	}
}

func TestEDRControlEngine(t *testing.T) {
	cfg := DefaultBlockZConfig()
	bz := NewBlockZOrchestrator(cfg)

	if bz.EDRControl != nil {
		detected := bz.EDRControl.DetectEDRs()
		if len(detected) == 0 {
			t.Log("No EDRs detected (expected in test env)")
		} else {
			t.Logf("Detected EDRs: %v", detected)
		}
	}
}

func TestDeadManSwitchEngine(t *testing.T) {
	cfg := DefaultBlockZConfig()
	bz := NewBlockZOrchestrator(cfg)

	if bz.DeadMan != nil {
		armed := bz.DeadMan.IsArmed()
		t.Logf("DeadMan switch armed: %v", armed)
	}
}

func TestDisinformationEngine(t *testing.T) {
	cfg := DefaultBlockZConfig()
	bz := NewBlockZOrchestrator(cfg)

	if bz.Disinfo != nil {
		campaigns := bz.Disinfo.GenerateCampaign()
		if len(campaigns) == 0 {
			t.Log("Disinfo campaign empty (expected in simulation)")
		}
	}
}

func TestFalseFlagEngine(t *testing.T) {
	cfg := DefaultBlockZConfig()
	bz := NewBlockZOrchestrator(cfg)

	if bz.FalseFlag != nil {
		forgeries := bz.FalseFlag.PlantEvidence()
		if len(forgeries) == 0 {
			t.Log("False flag forgeries empty (expected in simulation)")
		}
	}
}
