// X404X — Semi-Autonomous Red Team Platform
//
// Entry point for the X404X CLI. Supports three operating modes:
//
//   1. Console (default):   x404x             → interactive msfconsole-style shell
//   2. Dashboard:           x404x dashboard   → starts the Vue 3 web UI
//   3. CLI:                 x404x <command>   → subcommand execution (campaign, etc.)
//
// Note: Console and TUI modes are under active development.
// See cmd/x404x/commands.go, console.go, tui.go for planned implementations.

package main

import (
	"context"
	"fmt"
	"os"

	"github.com/ruby570bocadito/x404x/internal/appstate"
	"github.com/ruby570bocadito/x404x/pkg/shared/config"
)

var globalState *appstate.AppState

func main() {
	args := os.Args[1:]

	cfg, err := loadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "config warning: %v (using defaults)\n", err)
		cfg = config.Default()
	}

	if len(args) == 0 {
		runConsole(cfg)
		return
	}

	switch args[0] {
	case "console", "--console":
		runConsole(cfg)
	case "dashboard", "--dashboard":
		fmt.Println("Dashboard mode: run 'make lab-up' or see docs/DEPLOYMENT.md")
		fmt.Println("  Web UI:  http://localhost:3000")
		fmt.Println("  API:     http://localhost:8445/api/health")
	case "help", "--help", "-h":
		printHelp()
	case "version", "--version", "-v":
		fmt.Println("X404X v0.1.0")
	default:
		fmt.Fprintf(os.Stderr, "unknown mode: %s\n", args[0])
		printHelp()
		os.Exit(1)
	}
}

func runConsole(cfg *config.Config) {
	cfg.Logging.Output = "file"
	cfg.Logging.File = "x404x.log"

	ctx := context.Background()
	state, err := appstate.New(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create state: %v\n", err)
		os.Exit(1)
	}
	globalState = state

	if err := state.Start(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "state start failed: %v\n", err)
	}
	defer state.Stop()

	fmt.Println("X404X console — interactive mode")
	fmt.Println("  Type 'help' for available commands.")
	fmt.Println("  (Full console implementation in progress — see docs/CONSOLE.md)")
	fmt.Println()
	printHelp()
}

func printHelp() {
	fmt.Println("Usage: x404x [mode]")
	fmt.Println()
	fmt.Println("Modes:")
	fmt.Println("  (default)    Launch interactive console (msfconsole-style)")
	fmt.Println("  dashboard    Show dashboard connection info")
	fmt.Println("  help         Print this help message")
	fmt.Println("  version      Print version")
	fmt.Println()
	fmt.Println("Documentation: docs/")
}

func loadConfig() (*config.Config, error) {
	cfgPath := "config.yaml"
	for i, a := range os.Args {
		if a == "--config" || a == "-c" {
			if i+1 < len(os.Args) {
				cfgPath = os.Args[i+1]
			}
		}
	}
	return config.Load(cfgPath)
}