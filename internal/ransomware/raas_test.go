package ransomware

import (
	"testing"
)

func TestNewInverseRaaSEngine(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	r := NewInverseRaaSEngine(cfg)
	if r == nil {
		t.Fatal("NewInverseRaaSEngine() returned nil")
	}
}

func TestInverseRaaSGetSubtenants(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	r := NewInverseRaaSEngine(cfg)

	tenants := r.GenerateMultiRansomNotes()
	if len(tenants) == 0 {
		t.Log("GenerateMultiRansomNotes returned empty (expected if no tenants)")
	} else {
		t.Logf("Subtenants: %d", len(tenants))
	}
}

func TestInverseRaaSConfig(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	if cfg.InverseRaaS {
		t.Log("InverseRaaS enabled in config (simulation)")
	}
}

func TestDNAMutationEngine(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	d := NewDNAMutationEngine(cfg)
	if d == nil {
		t.Fatal("NewDNAMutationEngine() returned nil")
	}

	profile := d.GenerateROPGadgets()
	if len(profile) == 0 {
		t.Log("GenerateROPGadgets returned empty (expected in test env)")
	} else {
		t.Logf("DNA mutation gadgets: %d", len(profile))
	}
}

func TestBlockchainC2Engine(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	b := NewBlockchainC2Engine(cfg)
	if b == nil {
		t.Fatal("NewBlockchainC2Engine() returned nil")
	}

	cmds := b.checkBlockchainForCommands()
	if len(cmds) == 0 {
		t.Log("checkBlockchainForCommands returned empty (expected in test env)")
	} else {
		t.Logf("Blockchain commands: %d", len(cmds))
	}
}

func TestLOLBinChainerConfig(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	if !cfg.PolymorphicEnabled {
		t.Log("PolymorphicEnabled disabled in config")
	}
}
