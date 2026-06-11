package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/ruby570bocadito/x404x/internal/api"
	"github.com/ruby570bocadito/x404x/internal/appstate"
	"github.com/ruby570bocadito/x404x/pkg/shared/config"
)

func startDashboard(cfg *config.Config) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// ── Redirect all structured logs to file — stdout stays clean ───────────
	cfg.Logging.Output = "file"
	if cfg.Logging.File == "" {
		cfg.Logging.File = "x404x.log"
	}

	port := cfg.Server.Port
	if port == 0 {
		port = 9090
	}

	// ── Banner ───────────────────────────────────────────────────────────────
	printDashboardBanner()

	// ── State (starts bridge internally via state.Start) ─────────────────────
	state, err := appstate.New(cfg)
	if err != nil {
		return fmt.Errorf("creating state: %w", err)
	}
	globalState = state

	if err := state.Start(ctx); err != nil {
		// Non-fatal — continue without bridge
		fmt.Printf("\033[38;5;220m  [~]\033[0m State start warning: %v\n", err)
	}
	defer state.Stop()

	// ── Bridge status (already started by state.Start) ───────────────────────
	if state.Bridge != nil && state.Bridge.Connected() {
		mods := len(state.GetModules())
		fmt.Printf("\033[38;5;46m  [+]\033[0m Python bridge \033[38;5;46m●\033[0m connected — %d modules\n", mods)
	} else {
		fmt.Printf("\033[38;5;240m  [~]\033[0m Python bridge \033[38;5;240m○\033[0m offline (local fallback active)\n")
	}

	// ── API server ───────────────────────────────────────────────────────────
	apiServer, err := api.NewWithState(cfg, state)
	if err != nil {
		return fmt.Errorf("creating API: %w", err)
	}
	apiServer.SetPort(port)
	apiServer.ServeStatic("web/dist")

	// Single shared WebSocket terminal (one Console for all browser tabs)
	apiServer.Mux().HandleFunc("/ws/terminal", handleTerminalWS(state))

	// ── Signal handler ───────────────────────────────────────────────────────
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		fmt.Printf("\n\033[38;5;196m  [!]\033[0m Shutting down X404X...\n")
		apiServer.Shutdown(context.Background())
		state.Stop()
		cancel()
		os.Exit(0)
	}()

	// ── Endpoint summary ─────────────────────────────────────────────────────
	d := "\033[38;5;240m"
	p := "\033[38;5;99m"
	w := "\033[38;5;255m"
	r := "\033[0m"
	fmt.Println()
	fmt.Printf("%s  ────────────────────────────────────────────────%s\n", d, r)
	fmt.Printf("%s  [*]%s Dashboard  %shttp://localhost:%d%s\n", p, r, w, port, r)
	fmt.Printf("%s  [*]%s API        %shttp://localhost:%d/api%s\n", p, r, w, port, r)
	fmt.Printf("%s  [*]%s Terminal   %sws://localhost:%d/ws/terminal%s\n", p, r, w, port, r)
	fmt.Printf("%s  ────────────────────────────────────────────────%s\n", d, r)
	fmt.Printf("  %sCtrl+C to stop%s\n\n", d, r)

	if err := apiServer.Start(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("API server: %w", err)
	}
	return nil
}

func printDashboardBanner() {
	g := []string{
		"\033[38;5;57m", "\033[38;5;63m", "\033[38;5;99m",
		"\033[38;5;135m", "\033[38;5;171m", "\033[38;5;207m",
	}
	rs := "\033[0m"
	b := "\033[1m"

	lines := []struct{ c, t string }{
		{g[0], `  ██╗  ██╗██╗  ██╗ ██████╗ ██╗  ██╗██╗  ██╗`},
		{g[1], `  ╚██╗██╔╝██║  ██║██╔═══██╗██║  ██║╚██╗██╔╝`},
		{g[2], `   ╚███╔╝ ███████║██║   ██║███████║ ╚███╔╝ `},
		{g[3], `   ██╔██╗ ╚════██║██║   ██║╚════██║ ██╔██╗ `},
		{g[4], `  ██╔╝ ██╗      ██║╚██████╔╝      ██║██╔╝ ██╗`},
		{g[5], `  ╚═╝  ╚═╝      ╚═╝ ╚═════╝       ╚═╝╚═╝  ╚═╝`},
	}

	fmt.Println()
	for _, l := range lines {
		fmt.Println(l.c + b + l.t + rs)
	}
	fmt.Printf("\n  %s%sSemi-Autonomous Red Team Platform%s  %s·%s  %sv1.0.0%s  %s·%s  %sDASHBOARD%s\n",
		g[2], b, rs, g[5], rs, "\033[38;5;240m", rs, g[5], rs, g[4], rs)
	fmt.Printf("  \033[38;5;240m%s\033[0m\n\n",
		"────────────────────────────────────────────────")
}
