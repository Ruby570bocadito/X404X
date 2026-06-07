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

	port := cfg.Server.Port
	if port == 0 {
		port = 9090
	}

	state, err := appstate.New(cfg)
	if err != nil {
		return fmt.Errorf("creating state: %w", err)
	}
	globalState = state

	if err := state.Start(ctx); err != nil {
		fmt.Printf("[!] State start warning: %v\n", err)
	}
	defer state.Stop()

	bridgeScript := "modules/bridge/bridge.py"
	moduleCount := 0
	if _, err := os.Stat(bridgeScript); err == nil {
		if err := state.Bridge.StartBridge(ctx, bridgeScript); err != nil {
			fmt.Printf("[!] Python bridge: %v\n", err)
		} else {
			moduleCount = len(state.GetModules())
			fmt.Printf("[+] Python bridge connected (%d modules)\n", moduleCount)
		}
	}

	apiServer, err := api.NewWithState(cfg, state)
	if err != nil {
		return fmt.Errorf("creating API: %w", err)
	}

	apiServer.SetPort(port)
	apiServer.ServeStatic("core/c2/web/dist")

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		fmt.Println("\n[!] Shutting down X404X...")
		apiServer.Shutdown(context.Background())
		state.Stop()
		cancel()
		os.Exit(0)
	}()

	fmt.Println()
	fmt.Println("  ╔══════════════════════════════════════════════════════╗")
	fmt.Println("  ║              X404X — DASHBOARD v3.2                 ║")
	fmt.Println("  ╠══════════════════════════════════════════════════════╣")
	fmt.Printf("  ║  Dashboard:  http://localhost:%d                    ║\n", port)
	fmt.Printf("  ║  API:        http://localhost:%d/api                ║\n", port)
	fmt.Printf("  ║  Health:     http://localhost:%d/api/health         ║\n", port)
	fmt.Printf("  ║  WebSocket:  ws://localhost:%d/ws                   ║\n", port)
	fmt.Printf("  ║  Módulos:    %d (Go) + 85 (Python Bridge)          ║\n", moduleCount)
	fmt.Println("  ║  Ctrl+C to stop                                     ║")
	fmt.Println("  ╚══════════════════════════════════════════════════════╝")
	fmt.Println()

	if err := apiServer.Start(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("API server: %w", err)
	}
	return nil
}
