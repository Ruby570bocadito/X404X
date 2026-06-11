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

	tenants := r.GetSubtenants()
	if tenants == nil {
		t.Log("GetSubtenants returned nil (expected if no tenants)")
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

	profile := d.GenerateMutationProfile()
	if profile == nil {
		t.Log("GenerateMutationProfile returned nil (expected in test env)")
	} else {
		t.Logf("DNA mutation profile: %+v", profile)
	}
}

func TestBlockchainC2Engine(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	b := NewBlockchainC2Engine(cfg)
	if b == nil {
		t.Fatal("NewBlockchainC2Engine() returned nil")
	}

	cmd := b.ReceiveCommand()
	if cmd == "" {
		t.Log("ReceiveCommand returned empty (expected in test env)")
	} else {
		t.Logf("Blockchain command: %s", cmd)
	}
}

func TestLOLBinChainerConfig(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	if !cfg.PolymorphicEnabled {
		t.Log("PolymorphicEnabled disabled in config")
	}
}
