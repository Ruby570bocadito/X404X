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

	bootType := b.DetectBootMethod()
	if bootType == "" {
		t.Error("DetectBootMethod returned empty (expected fallback)")
	}
	t.Logf("Detected bootloader: %s", bootType)
}

func TestBootkitGetStatus(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	b := NewBootkitEngine(cfg)

	status := b.CheckBootkitStatus()
	if status == nil {
		t.Fatal("CheckBootkitStatus returned nil")
	}
	for k, v := range status {
		t.Logf("  %s: %v", k, v)
	}
}
