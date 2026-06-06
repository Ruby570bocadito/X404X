package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/ruby570bocadito/x404x/core/api"
	"github.com/ruby570bocadito/x404x/core/appstate"
	"github.com/ruby570bocadito/x404x/shared/config"
)

func startDashboard(cfg *config.Config) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create shared state
	state, err := appstate.New(cfg)
	if err != nil {
		return fmt.Errorf("creating state: %w", err)
	}
	globalState = state

	if err := state.Start(ctx); err != nil {
		fmt.Printf("[!] State start warning: %v\n", err)
	}
	defer state.Stop()

	// Connect Python bridge if available
	bridgeScript := "modules/bridge/bridge.py"
	if _, err := os.Stat(bridgeScript); err == nil {
		if err := state.Bridge.StartBridge(ctx, bridgeScript); err != nil {
			fmt.Printf("[!] Python bridge not available: %v (modules will use offline fallback)\n", err)
		} else {
			fmt.Println("[+] Python bridge connected (9 modules)")
		}
	} else {
		fmt.Println("[!] Python bridge script not found — modules use offline fallback")
	}

	// Create API server with orchestrator
	apiServer, err := api.NewWithState(cfg, state)
	if err != nil {
		return fmt.Errorf("creating API: %w", err)
	}

	// Graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		fmt.Println("\n[!] Shutting down X404X...")
		state.Stop()
		apiServer.Shutdown(context.Background())
		cancel()
		os.Exit(0)
	}()

	fmt.Println()
	fmt.Println("  ╔════════════════════════════════════════════════════╗")
	fmt.Println("  ║              X404X — DASHBOARD MODE               ║")
	fmt.Println("  ╠════════════════════════════════════════════════════╣")
	fmt.Printf("  ║  API:       http://localhost:8445                ║\n")
	fmt.Printf("  ║  WS:        ws://localhost:8445/ws               ║\n")
	fmt.Printf("  ║  Health:    http://localhost:8445/api/health     ║\n")
	fmt.Printf("  ║  Dashboard: http://localhost:3000 (npm run dev)  ║\n")
	fmt.Println("  ║  Ctrl+C to stop                                  ║")
	fmt.Println("  ╚════════════════════════════════════════════════════╝")
	fmt.Println()

	if err := apiServer.Start(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("API server: %w", err)
	}
	return nil
}
