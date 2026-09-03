package main

// helpers.go — ANSI color constants, output helpers and a minimal table renderer.
//
// These definitions are shared across the CLI (root.go, commands.go, console.go,
// tui.go, listeners.go, payload*.go). They were removed during a cleanup commit
// and recreated during the security audit so that cmd/x404x compiles again.

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// ─── ANSI escape codes ────────────────────────────────────────────────────────

const (
	ansiR  = "\033[0m"    // reset
	ansiB  = "\033[1m"    // bold
	ansiIt = "\033[3m"    // italic

	cSuccess = "\033[32m" // green
	cDanger  = "\033[31m" // red
	cWarn    = "\033[33m" // yellow
	cOrange  = "\033[38;5;208m"
	cInfo    = "\033[36m" // cyan
	cCyan    = "\033[36m" // cyan
	cPrimary = "\033[35m" // magenta
	cMuted   = "\033[90m" // bright black / gray
	cWhite   = "\033[97m" // bright white

	// Gradient colors for the X404X banner (purple → cyan).
	g1 = "\033[38;5;129m"
	g2 = "\033[38;5;135m"
	g3 = "\033[38;5;141m"
	g4 = "\033[38;5;147m"
	g5 = "\033[38;5;153m"
	g6 = "\033[38;5;159m"
)

// ─── Output helpers ───────────────────────────────────────────────────────────

func printInfo(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Fprintf(ConsoleOut, "\n  %s[*]%s %s\n", cInfo, ansiR, msg)
}

func printOK(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Fprintf(ConsoleOut, "  %s[+]%s %s\n", cSuccess, ansiR, msg)
}

func printWarn(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Fprintf(ConsoleOut, "  %s[!]%s %s\n", cWarn, ansiR, msg)
}

func printErr(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Fprintf(ConsoleOut, "  %s[x]%s %s\n", cDanger, ansiR, msg)
}

// maskPassword masks a credential value for safe display.
func maskPassword(pw string) string {
	if pw == "" {
		return "(none)"
	}
	return strings.Repeat("•", len(pw))
}

func printSection(title string) {
	fmt.Fprintf(ConsoleOut, "\n  %s%s%s%s\n  %s%s%s\n",
		cPrimary, ansiB, title, ansiR,
		cMuted, strings.Repeat("─", 52), ansiR)
}

func printPanel(title, body string) {
	fmt.Fprintf(ConsoleOut, "\n  %s┌─ %s%s%s %s%s\n", cMuted, cPrimary, ansiB, title, strings.Repeat("─", 50), ansiR)
	for _, line := range strings.Split(body, "\n") {
		fmt.Fprintf(ConsoleOut, "  %s│%s %s\n", cMuted, ansiR, line)
	}
	fmt.Fprintf(ConsoleOut, "  %s└%s%s%s\n", cMuted, strings.Repeat("─", 52), ansiR, ansiR)
}

func printBigBanner() {
	rows := []struct{ c, t string }{
		{g1, `  ▄▄▄   ▄▄▄  ▄▄   ▄▄   ▄▄▄▄▄▄   ▄▄   ▄▄  ▄▄▄   ▄▄▄ `},
		{g2, `  ▀██▄ ▄██▀  ██   ██  ██▀  ▀██  ██   ██  ▀██▄ ▄██▀ `},
		{g3, `    ▀███▀    ███████  ██    ██  ███████    ▀███▀   `},
		{g4, `  ▄██▀ ▀██▄       ██  ██▄  ▄██       ██  ▄██▀ ▀██▄ `},
		{g5, `  ██▀   ▀██       ██   ▀████▀        ██  ██▀   ▀██ `},
		{g6, `  ▀▀     ▀▀       ▀▀                 ▀▀  ▀▀     ▀▀ `},
	}
	fmt.Fprintln(ConsoleOut)
	for _, r := range rows {
		fmt.Fprintf(ConsoleOut, "%s%s%s\n", r.c, r.t, ansiR)
	}
	fmt.Fprintf(ConsoleOut, "\n  %s%sSemi-Autonomous Red Team Platform%s  %sv1.0.0%s\n",
		cPrimary, ansiB, ansiR, cMuted, ansiR)
}

// ─── Minimal table renderer ───────────────────────────────────────────────────

type table struct {
	headers []string
	rows    [][]string
	widths  []int
}

func newTable(headers ...string) *table {
	return &table{
		headers: headers,
		widths:  make([]int, len(headers)),
	}
}

func (t *table) addRow(cells ...string) {
	t.rows = append(t.rows, cells)
	for i, c := range cells {
		if i >= len(t.widths) {
			continue
		}
		if w := ansiWidth(c); w > t.widths[i] {
			t.widths[i] = w
		}
	}
}

func (t *table) render() {
	if len(t.headers) == 0 {
		return
	}
	for i, h := range t.headers {
		if w := ansiWidth(h); w > t.widths[i] {
			t.widths[i] = w
		}
	}

	// Header
	for i, h := range t.headers {
		pad := t.widths[i] - ansiWidth(h)
		if i == len(t.headers)-1 {
			fmt.Fprintf(ConsoleOut, "  %s%s%s%s\n", cPrimary, ansiB, h, ansiR)
		} else {
			fmt.Fprintf(ConsoleOut, "  %s%s%s%s%s", cPrimary, ansiB, h, ansiR, strings.Repeat(" ", pad+2))
		}
	}

	// Separator
	fmt.Fprintf(ConsoleOut, "  %s%s%s\n", cMuted, strings.Repeat("─", 52), ansiR)

	// Rows
	for _, row := range t.rows {
		for i, cell := range row {
			pad := 0
			if i < len(t.widths) {
				pad = t.widths[i] - ansiWidth(cell)
			}
			if i == len(row)-1 {
				fmt.Fprintf(ConsoleOut, "  %s%s%s\n", cWhite, cell, ansiR)
			} else {
				fmt.Fprintf(ConsoleOut, "  %s%s%s%s", cWhite, cell, ansiR, strings.Repeat(" ", pad+2))
			}
		}
	}
	fmt.Fprintln(ConsoleOut)
}

// statusTag returns a colored status marker for console/table output.
func statusTag(status string) string {
	switch strings.ToLower(status) {
	case "active", "running", "completed", "success", "healthy", "online":
		return cSuccess + "● " + status + ansiR
	case "paused", "pending", "idle", "waiting":
		return cWarn + "● " + status + ansiR
	case "failed", "error", "stopped", "dead", "terminated", "offline":
		return cDanger + "● " + status + ansiR
	default:
		return cInfo + "● " + status + ansiR
	}
}

// trunc truncates s to at most n runes, appending "…" when it was cut.
func trunc(s string, n int) string {
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return string(r[:n]) + "…"
}

// pbar renders an ASCII progress bar for a 0.0–1.0 progress value.
func pbar(progress float64, width int) string {
	if progress < 0 {
		progress = 0
	}
	if progress > 1 {
		progress = 1
	}
	filled := int(progress * float64(width))
	if filled > width {
		filled = width
	}
	return cSuccess + strings.Repeat("█", filled) + cMuted + strings.Repeat("░", width-filled) +
		ansiR + fmt.Sprintf(" %3.0f%%", progress*100)
}

// killChainOrder maps a kill-chain phase string to its 0-based position in the
// CLI's 7-step chain (Recon … Objectives).
func killChainOrder(phase string) int {
	switch strings.ToLower(phase) {
	case "recon", "reconnaissance":
		return 0
	case "weaponization":
		return 1
	case "delivery":
		return 2
	case "exploitation", "exploit":
		return 3
	case "installation", "persistence":
		return 4
	case "c2", "command and control", "command_and_control":
		return 5
	default: // actions_on_objective, exfiltration, objectives
		return 6
	}
}

// ansiWidth returns the visible (printable) width of a string, ignoring ANSI
// escape sequences so that colored cells align correctly.
func ansiWidth(s string) int {
	var b strings.Builder
	inEsc := false
	for i := 0; i < len(s); i++ {
		if s[i] == '\x1b' {
			inEsc = true
			continue
		}
		if inEsc {
			if s[i] == 'm' {
				inEsc = false
			}
			continue
		}
		b.WriteByte(s[i])
	}
	return utf8.RuneCountInString(b.String())
}
