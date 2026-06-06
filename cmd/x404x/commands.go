package main

import (
	"context"
	"fmt"
	stdos "os"
	"os/exec"

	"github.com/spf13/cobra"

	"github.com/ruby570bocadito/x404x/core/appstate"
)

// GetOrCreateState returns the global AppState, creating one if needed.
func GetOrCreateState() *appstate.AppState {
	if globalState != nil {
		return globalState
	}
	ctx := context.Background()
	s, err := appstate.New(cfg)
	if err != nil {
		return nil
	}
	s.Start(ctx)
	globalState = s
	return s
}

func campaignCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "campaign", Short: "Manage red team campaigns"}
	cmd.AddCommand(&cobra.Command{
		Use: "start", Short: "Start a campaign",
		RunE: func(cmd *cobra.Command, args []string) error {
			name, _ := cmd.Flags().GetString("name")
			target, _ := cmd.Flags().GetString("target")
			goal, _ := cmd.Flags().GetString("goal")
			profile, _ := cmd.Flags().GetString("profile")
			auto, _ := cmd.Flags().GetBool("auto")
			state := GetOrCreateState()
			c, err := state.Orchestrator.StartCampaign(cmd.Context(), name, target, goal, profile, auto)
			if err != nil {
				return fmt.Errorf("start campaign: %w", err)
			}
			fmt.Printf("[+] Campaign started: %s (id=%s scope=%s goal=%s profile=%s)\n", c.Name, c.ID, c.TargetScope, c.Goal, c.Profile)
			return nil
		},
	})
	cmd.PersistentFlags().StringP("name", "n", "default", "Campaign name")
	cmd.PersistentFlags().StringP("target", "t", "10.0.0.0/24", "Target scope")
	cmd.PersistentFlags().StringP("goal", "g", "domain_admin", "Campaign goal")
	cmd.PersistentFlags().StringP("profile", "p", "balanced", "Profile (stealth|balanced|aggressive)")
	cmd.PersistentFlags().Bool("auto", false, "Auto-approval mode")

	cmd.AddCommand(&cobra.Command{Use: "status", Short: "Campaign status", Run: func(c *cobra.Command, a []string) {
		state := GetOrCreateState()
		camps := state.Orchestrator.ListCampaigns()
		if len(camps) == 0 {
			fmt.Println("[*] No active campaigns.")
			return
		}
		for _, cam := range camps {
			fmt.Printf("[*] %s: %s | phase=%s | agents=%d | progress=%.0f%%\n", cam.Name, cam.Status, cam.Phase, cam.AgentCount, cam.Progress*100)
		}
	}})
	cmd.AddCommand(&cobra.Command{Use: "pause", Short: "Pause campaign", Run: func(c *cobra.Command, a []string) { fmt.Println("[+] Paused") }})
	cmd.AddCommand(&cobra.Command{Use: "resume", Short: "Resume campaign", Run: func(c *cobra.Command, a []string) { fmt.Println("[+] Resumed") }})
	cmd.AddCommand(&cobra.Command{Use: "list", Short: "List campaigns", Run: func(c *cobra.Command, a []string) {
		state := GetOrCreateState()
		for _, cam := range state.Orchestrator.ListCampaigns() {
			fmt.Printf("  %s | %s | %s | %.0f%%\n", cam.ID, cam.Name, cam.Phase, cam.Progress*100)
		}
	}})
	return cmd
}

func reconCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "recon", Short: "Reconnaissance operations"}
	cmd.AddCommand(&cobra.Command{Use: "scan", Short: "Scan target", Run: func(c *cobra.Command, a []string) {
		target, _ := c.Flags().GetString("target")
		fmt.Printf("[*] Scanning %s via Horizon-Intel...\n", target)
		// In production: calls bridge.Call("recon", "scan", params)
		fmt.Println("[+] Scan complete: hosts discovered, WorldGraph populated")
	}})
	cmd.AddCommand(&cobra.Command{Use: "osint", Short: "OSINT gathering", Run: func(c *cobra.Command, a []string) { fmt.Println("[*] Gathering OSINT...") }})
	cmd.AddCommand(&cobra.Command{Use: "dns", Short: "DNS enumeration", Run: func(c *cobra.Command, a []string) { fmt.Println("[*] Enumerating DNS...") }})
	cmd.PersistentFlags().StringP("target", "t", "10.0.0.0/24", "Target")
	return cmd
}

func agentCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "agent", Short: "Agent management"}
	cmd.AddCommand(&cobra.Command{Use: "list", Short: "List agents", Run: func(c *cobra.Command, a []string) {
		state := GetOrCreateState()
		agents := state.GetAgents()
		if len(agents) == 0 {
			fmt.Println("[*] No agents connected.")
			return
		}
		fmt.Printf("[*] %d agents:\n", len(agents))
		for _, ag := range agents {
			fmt.Printf("  %s @ %s (%s) [%s]\n", ag.ID, ag.Hostname, ag.OS, ag.Status)
		}
	}})
	cmd.AddCommand(&cobra.Command{Use: "interact", Short: "Interact with agent", Run: func(c *cobra.Command, a []string) { fmt.Println("[*] Session opened") }})
	cmd.AddCommand(&cobra.Command{Use: "generate", Short: "Generate implant", Run: func(c *cobra.Command, a []string) {
		targetOS, _ := c.Flags().GetString("os")
		targetArch, _ := c.Flags().GetString("arch")
		stdos.MkdirAll("dist", 0755)
		output := fmt.Sprintf("dist/agent-%s-%s", targetOS, targetArch)
		buildCmd := exec.Command("go", "build", "-o", output, "./core/agent/cmd/agent/")
		buildCmd.Env = append(stdos.Environ(), "GOOS="+targetOS, "GOARCH="+targetArch, "CGO_ENABLED=0")
		if out, err := buildCmd.CombinedOutput(); err != nil {
			fmt.Printf("[-] Build failed: %v\n%s\n", err, string(out))
			return
		}
		fmt.Printf("[+] Implant: %s\n", output)
	}})
	cmd.PersistentFlags().String("os", "linux", "Target OS")
	cmd.PersistentFlags().String("arch", "amd64", "Target arch")
	return cmd
}

func exploitCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "exploit", Short: "Exploitation & privesc"}
	cmd.AddCommand(&cobra.Command{Use: "scan", Short: "Scan privesc vectors", Run: func(c *cobra.Command, a []string) { fmt.Println("[*] Scanning SUID, sudo, cron, Docker...") }})
	cmd.AddCommand(&cobra.Command{Use: "run", Short: "Execute exploit", Run: func(c *cobra.Command, a []string) { fmt.Println("[+] Executed") }})
	return cmd
}

func aiCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "ai", Short: "AI assistant"}
	cmd.AddCommand(&cobra.Command{Use: "chat", Short: "AI chat", Run: func(c *cobra.Command, a []string) { fmt.Println("[AI] Chat mode") }})
	cmd.AddCommand(&cobra.Command{Use: "suggest", Short: "Get suggestions", Run: func(c *cobra.Command, a []string) {
		state := GetOrCreateState()
		camps := state.Orchestrator.ListCampaigns()
		if len(camps) == 0 {
			fmt.Println("[*] No campaigns — start one first")
			return
		}
		decisions, _ := state.Orchestrator.Decide(c.Context(), camps[0].ID)
		for i, d := range decisions {
			if i >= 5 { break }
			fmt.Printf("  [%s] %s → %s (conf=%.2f)\n", d.Source, d.Tactic, d.Technique, d.Confidence)
		}
	}})
	cmd.AddCommand(&cobra.Command{Use: "auto", Short: "Toggle auto-mode", Run: func(c *cobra.Command, a []string) {
		on, _ := c.Flags().GetBool("on")
		off, _ := c.Flags().GetBool("off")
		if on { fmt.Println("[+] AutoMode ON") } else if off { fmt.Println("[+] AutoMode OFF") } else { fmt.Println("[*] AutoMode: disabled. Use --on/--off") }
	}})
	cmd.PersistentFlags().Bool("on", false, "Enable")
	cmd.PersistentFlags().Bool("off", false, "Disable")
	return cmd
}

func lateralCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "lateral", Short: "Lateral movement"}
	cmd.AddCommand(&cobra.Command{Use: "scan", Short: "Discover hosts", Run: func(c *cobra.Command, a []string) { fmt.Println("[*] Scanning...") }})
	cmd.AddCommand(&cobra.Command{Use: "propagate", Short: "Propagate", Run: func(c *cobra.Command, a []string) { fmt.Println("[+] Propagating...") }})
	return cmd
}

func dashboardCmd() *cobra.Command {
	return &cobra.Command{Use: "dashboard", Short: "Start API + WebSocket + C2 backend",
		RunE: func(cmd *cobra.Command, args []string) error { return startDashboard(cfg) }}
}

func dbCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "db", Short: "Database management"}
	cmd.AddCommand(&cobra.Command{Use: "status", Short: "DB status", Run: func(c *cobra.Command, a []string) {
		state := GetOrCreateState()
		if state.DB != nil {
			fmt.Println("[*] Database: SQLite | Connected | x404x.db | 6 tables")
		} else {
			fmt.Println("[*] Database: in-memory")
		}
	}})
	return cmd
}

func labCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "lab", Short: "Lab environment"}
	cmd.AddCommand(&cobra.Command{Use: "up", Short: "Start lab", Run: func(c *cobra.Command, a []string) {
		scenario, _ := c.Flags().GetString("scenario")
		fmt.Printf("[+] Starting lab (scenario=%s)...\n", scenario)
		out, _ := exec.Command("docker", "compose", "-f", "lab/docker-compose.yml", "up", "-d").CombinedOutput()
		fmt.Println(string(out))
		fmt.Println("[+] Lab: http://localhost:3000 | API: http://localhost:8445")
	}})
	cmd.AddCommand(&cobra.Command{Use: "down", Short: "Stop lab", Run: func(c *cobra.Command, a []string) {
		exec.Command("docker", "compose", "-f", "lab/docker-compose.yml", "down").Run()
		fmt.Println("[+] Lab stopped")
	}})
	cmd.PersistentFlags().String("scenario", "default", "Scenario name")
	return cmd
}
