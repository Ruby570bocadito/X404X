// X404X — Semi-Autonomous Red Team Platform
//
// Three operating modes, all sharing the same AppState:
//
//	1. TUI (default):    x404x                    → Bubble Tea terminal UI
//	2. Console:          x404x console            → msfconsole-style shell
//	3. CLI (traditional): x404x campaign start...  → Cobra commands
//
// All modes connect to the same orchestrator, bridge, and C2 backend.

package main

import (
	"context"
	"fmt"
	"os"

	"github.com/ruby570bocadito/x404x/core/appstate"
	"github.com/ruby570bocadito/x404x/shared/config"
)

var globalState *appstate.AppState

func main() {
	args := os.Args[1:]

	// Load config
	var cfg *config.Config
	var err error
	cfgPath := "config.yaml"
	for i, a := range os.Args {
		if a == "--config" || a == "-c" {
			if i+1 < len(os.Args) {
				cfgPath = os.Args[i+1]
			}
		}
	}
	cfg, err = config.Load(cfgPath)
	if err != nil {
		cfg = config.Default()
	}

	// Mode detection
	switch {
	case len(args) == 0:
		runTUI()

	case args[0] == "console":
		runConsoleMode(cfg, args[1:])

	case args[0] == "tui":
		runTUI()

	case args[0] == "dashboard":
		startDashboard(cfg)

	default:
		runCLI()
	}
}

func runConsoleMode(cfg *config.Config, args []string) {
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

	if err := StartConsoleState(state, args); err != nil {
		fmt.Fprintf(os.Stderr, "console error: %v\n", err)
		os.Exit(1)
	}
}

func runTUI() {
	if err := StartTUI(); err != nil {
		fmt.Fprintf(os.Stderr, "x404x TUI error: %v\n", err)
		os.Exit(1)
	}
}

func runCLI() {
	Execute()
}
