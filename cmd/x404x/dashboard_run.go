package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ruby570bocadito/x404x/core/api"
	"github.com/ruby570bocadito/x404x/core/orchestrator"
	"github.com/ruby570bocadito/x404x/shared/config"
	"github.com/ruby570bocadito/x404x/shared/types"
)

func startDashboard(cfg *config.Config) error {
	orch, err := orchestrator.New(cfg)
	if err != nil {
		return fmt.Errorf("creating orchestrator: %w", err)
	}

	server, err := api.New(cfg, orch)
	if err != nil {
		return fmt.Errorf("creating API server: %w", err)
	}

	// Load demo data into world graph
	orch.WorldGraph().GenerateDemoData()

	// Register demo agents
	now := time.Now()
	demoAgents := []types.Agent{
		{ID: "abc123", Hostname: "DC", OS: "Windows 2019", Username: "NT\\SYSTEM", LocalIP: "10.0.0.10", Status: types.AgentStatusOnline, FirstSeen: now, LastCheckin: now, CampaignID: "demo-001", Uptime: 52000, Privileges: []string{"SYSTEM"}},
		{ID: "def456", Hostname: "DB", OS: "Ubuntu 24.04", Username: "root", LocalIP: "10.0.0.20", Status: types.AgentStatusOnline, FirstSeen: now, LastCheckin: now, CampaignID: "demo-001", Uptime: 22800, Privileges: []string{"root"}},
		{ID: "ghi789", Hostname: "WS1", OS: "Windows 11", Username: "user", LocalIP: "10.0.0.50", Status: types.AgentStatusActive, FirstSeen: now, LastCheckin: now, CampaignID: "demo-001", Uptime: 7500, Privileges: []string{"user"}},
	}
	for i := range demoAgents {
		server.RegisterAgent(&demoAgents[i])
	}

	// Start a demo campaign
	campaign, _ := orch.StartCampaign(context.Background(), "TFG-Demo", "10.0.0.0/24", "domain_admin", "balanced", false)
	orch.AdvancePhase(campaign.ID, types.PhaseExploitation)

	// Add demo hosts and vulns
	hosts := []types.Target{
		{IP: "10.0.0.10", Hostname: "DC", OS: "Windows 2019", OpenPorts: []int{445, 3389, 53}, AssetValue: 100},
		{IP: "10.0.0.20", Hostname: "DB", OS: "Ubuntu 24.04", OpenPorts: []int{22, 3306, 6379}, AssetValue: 70},
		{IP: "10.0.0.50", Hostname: "WS1", OS: "Windows 11", OpenPorts: []int{445, 135}, AssetValue: 10},
	}
	for i := range hosts {
		server.AddHost(&hosts[i])
	}

	vulns := []types.Vulnerability{
		{CVE: "MS17-010", Description: "EternalBlue SMB Remote Code Execution", Severity: "critical", Service: "smb", Port: 445, TargetIP: "10.0.0.10"},
		{CVE: "CVE-2024-XXXX", Description: "apport ExecutablePath spoofing on Ubuntu 24.04", Severity: "high", Service: "apport", Port: 0, TargetIP: "10.0.0.20"},
	}
	for i := range vulns {
		server.AddVuln(&vulns[i])
	}

	// Handle graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		fmt.Println("\n[!] Shutting down API server...")
		server.Shutdown(context.Background())
		os.Exit(0)
	}()

	fmt.Printf("[+] X404X API server starting on http://localhost:%d\n", cfg.Dashboard.Port)
	fmt.Printf("[+] WebSocket: ws://localhost:%d/ws\n", cfg.Dashboard.Port)
	fmt.Printf("[+] Health check: http://localhost:%d/api/health\n", cfg.Dashboard.Port)
	fmt.Println("[+] Press Ctrl+C to stop")

	if err := server.Start(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("API server: %w", err)
	}
	return nil
}
