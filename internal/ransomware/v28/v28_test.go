package v28

import (
	"testing"
)

func TestDefaultV28Config(t *testing.T) {
	cfg := DefaultV28Config()
	if cfg == nil {
		t.Fatal("DefaultV28Config returned nil")
	}
	if cfg.Simulation != true {
		t.Error("Simulation should default to true")
	}
}

func TestV28ConfigFlags(t *testing.T) {
	cfg := DefaultV28Config()

	flags := []struct {
		name string
		val  bool
	}{
		{"IoTIdentityTheft", cfg.IoTIdentityTheft},
		{"FalseMemoryInjection", cfg.FalseMemoryInjection},
		{"DeathByThousandCuts", cfg.DeathByThousandCuts},
		{"PatchGuardBypass", cfg.PatchGuardBypass},
		{"KeyboardLEDExfil", cfg.KeyboardLEDExfil},
		{"ZombieArmyPolitical", cfg.ZombieArmyPolitical},
		{"LegacyPoison", cfg.LegacyPoison},
		{"SEOSabotage", cfg.SEOSabotage},
		{"FakeVulnInjection", cfg.FakeVulnInjection},
		{"InceptionHypervisor", cfg.InceptionHypervisor},
		{"ISPBGPSubversion", cfg.ISPBGPSubversion},
		{"AntiAttributionClone", cfg.AntiAttributionClone},
		{"PowerGridHarmonics", cfg.PowerGridHarmonics},
		{"TimeLockExtortion", cfg.TimeLockExtortion},
		{"VRSpyware", cfg.VRSpyware},
		{"GlobalAIPoison", cfg.GlobalAIPoison},
		{"CDNMalwareInjection", cfg.CDNMalwareInjection},
		{"BioCyberDNA", cfg.BioCyberDNA},
		{"BrowserParasite", cfg.BrowserParasite},
		{"FakeDocumentsGen", cfg.FakeDocumentsGen},
		{"SoundPanicAttack", cfg.SoundPanicAttack},
		{"EmotionalEncryption", cfg.EmotionalEncryption},
		{"FalseRedemption", cfg.FalseRedemption},
	}

	for _, f := range flags {
		if !f.val {
			t.Logf("Flag %s: %v (simulation)", f.name, f.val)
		}
	}
}

func TestNewV28ConfigAllEnabled(t *testing.T) {
	cfg := DefaultV28Config()
	cfg.Enabled = true
	cfg.IoTIdentityTheft = true
	cfg.FalseMemoryInjection = true
	cfg.DeathByThousandCuts = true
	cfg.PatchGuardBypass = true
	cfg.KeyboardLEDExfil = true
	cfg.ZombieArmyPolitical = true
	cfg.LegacyPoison = true
	cfg.SEOSabotage = true
	cfg.FakeVulnInjection = true
	cfg.InceptionHypervisor = true
	cfg.ISPBGPSubversion = true
	cfg.AntiAttributionClone = true
	cfg.PowerGridHarmonics = true
	cfg.TimeLockExtortion = true
	cfg.VRSpyware = true
	cfg.GlobalAIPoison = true
	cfg.CDNMalwareInjection = true
	cfg.BioCyberDNA = true
	cfg.BrowserParasite = true
	cfg.FakeDocumentsGen = true
	cfg.SoundPanicAttack = true
	cfg.EmotionalEncryption = true
	cfg.FalseRedemption = true

	// Verify all 24 modules are enabled
	count := 0
	if cfg.IoTIdentityTheft { count++ }
	if cfg.FalseMemoryInjection { count++ }
	if cfg.DeathByThousandCuts { count++ }
	if cfg.PatchGuardBypass { count++ }
	if cfg.KeyboardLEDExfil { count++ }
	if cfg.ZombieArmyPolitical { count++ }
	if cfg.LegacyPoison { count++ }
	if cfg.SEOSabotage { count++ }
	if cfg.FakeVulnInjection { count++ }
	if cfg.InceptionHypervisor { count++ }
	if cfg.ISPBGPSubversion { count++ }
	if cfg.AntiAttributionClone { count++ }
	if cfg.PowerGridHarmonics { count++ }
	if cfg.TimeLockExtortion { count++ }
	if cfg.VRSpyware { count++ }
	if cfg.GlobalAIPoison { count++ }
	if cfg.CDNMalwareInjection { count++ }
	if cfg.BioCyberDNA { count++ }
	if cfg.BrowserParasite { count++ }
	if cfg.FakeDocumentsGen { count++ }
	if cfg.SoundPanicAttack { count++ }
	if cfg.EmotionalEncryption { count++ }
	if cfg.FalseRedemption { count++ }

	if count != 23 {
		t.Errorf("expected 23 enabled flags (all except Simulation), got %d", count)
	}
}

func TestIoTIdentityTheftEngine(t *testing.T) {
	cfg := DefaultV28Config()
	engine := NewIoTIdentityTheftEngine(cfg)
	if engine == nil {
		t.Fatal("NewIoTIdentityTheftEngine() returned nil")
	}
}

func TestFalseMemoryInjectionEngine(t *testing.T) {
	cfg := DefaultV28Config()
	engine := NewFalseMemoryInjectionEngine(cfg)
	if engine == nil {
		t.Fatal("NewFalseMemoryInjectionEngine() returned nil")
	}
}

func TestPatchGuardBypassEngine(t *testing.T) {
	cfg := DefaultV28Config()
	engine := NewPatchGuardBypassEngine(cfg)
	if engine == nil {
		t.Fatal("NewPatchGuardBypassEngine() returned nil")
	}
}

func TestKeyboardLEDExfilEngine(t *testing.T) {
	cfg := DefaultV28Config()
	engine := NewKeyboardLEDExfilEngine(cfg)
	if engine == nil {
		t.Fatal("NewKeyboardLEDExfilEngine() returned nil")
	}
}

func TestISPBGPSubversionEngine(t *testing.T) {
	cfg := DefaultV28Config()
	engine := NewISPBGPSubversionEngine(cfg)
	if engine == nil {
		t.Fatal("NewISPBGPSubversionEngine() returned nil")
	}
}

func TestPowerGridHarmonicsEngine(t *testing.T) {
	cfg := DefaultV28Config()
	engine := NewPowerGridHarmonicsEngine(cfg)
	if engine == nil {
		t.Fatal("NewPowerGridHarmonicsEngine() returned nil")
	}
}

func TestGlobalAIPoisonEngine(t *testing.T) {
	cfg := DefaultV28Config()
	engine := NewGlobalAIPoisonEngine(cfg)
	if engine == nil {
		t.Fatal("NewGlobalAIPoisonEngine() returned nil")
	}
}

func TestBioCyberDNAEngine(t *testing.T) {
	cfg := DefaultV28Config()
	engine := NewBioCyberDNAEngine(cfg)
	if engine == nil {
		t.Fatal("NewBioCyberDNAEngine() returned nil")
	}
}

func TestSoundPanicAttackEngine(t *testing.T) {
	cfg := DefaultV28Config()
	engine := NewSoundPanicAttackEngine(cfg)
	if engine == nil {
		t.Fatal("NewSoundPanicAttackEngine() returned nil")
	}
}

func TestFalseRedemptionEngine(t *testing.T) {
	cfg := DefaultV28Config()
	engine := NewFalseRedemptionEngine(cfg)
	if engine == nil {
		t.Fatal("NewFalseRedemptionEngine() returned nil")
	}
}

func TestV28ConfigC2Endpoint(t *testing.T) {
	cfg := DefaultV28Config()
	if cfg.C2Endpoint == "" {
		t.Log("C2Endpoint empty (default)")
	}
	cfg.C2Endpoint = "https://c2.x404x.local"
	if cfg.C2Endpoint != "https://c2.x404x.local" {
		t.Errorf("expected custom C2 endpoint, got %s", cfg.C2Endpoint)
	}
}
