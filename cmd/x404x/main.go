// X404X — Semi-Autonomous Red Team Platform
//
// Three operating modes:
//
//	1. TUI (default):    x404x                    → Bubble Tea terminal UI
//	2. Console:          x404x console            → msfconsole-style shell
//	3. CLI (traditional): x404x campaign start...  → Cobra commands
//
// All modes connect to the same orchestrator backend.

package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	args := os.Args[1:]

	// Mode detection
	switch {
	case len(args) == 0:
		// No args → launch TUI
		runTUI()

	case args[0] == "console":
		// Explicit console mode
		runConsole(args[1:])

	case args[0] == "tui":
		// Explicit TUI mode
		runTUI()

	default:
		// Traditional CLI mode (Cobra)
		runCLI()
	}
}

func runTUI() {
	if err := StartTUI(); err != nil {
		fmt.Fprintf(os.Stderr, "x404x TUI error: %v\n", err)
		os.Exit(1)
	}
}

func runConsole(args []string) {
	if err := StartConsole(args); err != nil {
		fmt.Fprintf(os.Stderr, "x404x console error: %v\n", err)
		os.Exit(1)
	}
}

func runCLI() {
	_ = strings.Join(os.Args, " ")
	Execute()
}
