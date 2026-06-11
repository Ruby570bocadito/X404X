package defense

import (
	"testing"
)

func TestNewBlueForgeEngine(t *testing.T) {
	bf := NewBlueForgeEngine()
	if bf == nil {
		t.Fatal("NewBlueForgeEngine() returned nil")
	}
	if len(bf.Techniques) == 0 {
		t.Error("expected non-zero techniques")
	}
}

func TestMarkUsed(t *testing.T) {
	bf := NewBlueForgeEngine()
	bf.MarkUsed("T1566")

	found := false
	for _, tech := range bf.Techniques {
		if tech.ID == "T1566" {
			found = true
			if !tech.Used {
				t.Error("T1566 should be marked Used")
			}
			break
		}
	}
	if !found {
		t.Error("T1566 not found in techniques")
	}
}

func TestMarkUsedUnknown(t *testing.T) {
	bf := NewBlueForgeEngine()
	bf.MarkUsed("T9999")

	for _, tech := range bf.Techniques {
		if tech.Used {
			t.Fatal("no technique should be marked used for unknown ID")
		}
	}
}

func TestSimulateDetection(t *testing.T) {
	bf := NewBlueForgeEngine()
	bf.MarkUsed("T1059")
	bf.MarkUsed("T1003")

	bf.SimulateDetection()

	if len(bf.Undetected)+len(bf.Detected) != 2 {
		t.Errorf("expected 2 total detections (used+undetected), got %d undetected + %d detected",
			len(bf.Undetected), len(bf.Detected))
	}
}

func TestCoverageAfterDetection(t *testing.T) {
	bf := NewBlueForgeEngine()
	bf.MarkUsed("T1566")

	bf.SimulateDetection()

	if bf.CoveragePct < 0 || bf.CoveragePct > 100 {
		t.Errorf("CoveragePct %.1f out of range [0,100]", bf.CoveragePct)
	}
	if bf.Score != bf.CoveragePct {
		t.Errorf("Score %.1f should equal CoveragePct %.1f", bf.Score, bf.CoveragePct)
	}
}

func TestGenerateReport(t *testing.T) {
	bf := NewBlueForgeEngine()
	bf.MarkUsed("T1486")

	report := bf.GenerateReport()
	if report == "" {
		t.Error("GenerateReport() returned empty string")
	}
}

func TestATTACKTechniquesCount(t *testing.T) {
	bf := NewBlueForgeEngine()
	if len(bf.Techniques) != 20 {
		t.Errorf("expected 20 ATT&CK techniques, got %d", len(bf.Techniques))
	}
}

func TestTechniqueFields(t *testing.T) {
	bf := NewBlueForgeEngine()
	for _, tech := range bf.Techniques {
		if tech.ID == "" {
			t.Error("technique with empty ID")
		}
		if tech.Name == "" {
			t.Error("technique with empty Name")
		}
		if tech.Tactic == "" {
			t.Error("technique with empty Tactic")
		}
	}
}

func TestMultipleMarkUsed(t *testing.T) {
	bf := NewBlueForgeEngine()
	markIDs := []string{"T1566", "T1059", "T1547", "T1068", "T1562", "T1003"}

	for _, id := range markIDs {
		bf.MarkUsed(id)
	}

	usedCount := 0
	for _, tech := range bf.Techniques {
		if tech.Used {
			usedCount++
		}
	}
	if usedCount != len(markIDs) {
		t.Errorf("expected %d used techniques, got %d", len(markIDs), usedCount)
	}
}
