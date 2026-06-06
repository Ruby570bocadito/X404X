package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/ruby570bocadito/x404x/shared/config"
	"github.com/ruby570bocadito/x404x/shared/logger"
)

var (
	cfg     *config.Config
	log     *logger.Logger
	cfgPath string
)

var rootCmd = &cobra.Command{
	Use:   "x404x",
	Short: "X404X — Semi-Autonomous Red Team Platform",
	Long: banner() + `

X404X is a semi-autonomous Red Team platform covering the complete
cyber kill chain — from reconnaissance to exfiltration.

Modes:
  x404x              Launch interactive TUI
  x404x console      Launch msfconsole-style shell
  x404x <command>    Traditional CLI mode`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var err error
		cfg, err = config.Load(cfgPath)
		if err != nil {
			cfg = config.Default()
		}
		log, _ = logger.New(logger.Config{
			Level:     cfg.Logging.Level,
			Format:    cfg.Logging.Format,
			Component: "cli",
		})
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func banner() string {
	return fmt.Sprintf(`
██╗  ██╗ ██╗  ██╗  ██████╗  ██╗  ██╗ ██╗  ██╗
╚██╗██╔╝ ██║  ██║ ██╔═══██╗ ██║  ██║ ╚██╗██╔╝
 ╚███╔╝  ███████║ ██║   ██║ ███████║  ╚███╔╝
 ██╔██╗  ╚════██║ ██║   ██║ ╚════██║  ██╔██╗
██╔╝ ██╗     ██╔╝ ╚██████╔╝     ██╔╝ ██╔╝ ██╗
╚═╝  ╚═╝     ╚═╝   ╚═════╝      ╚═╝  ╚═╝  ╚═╝

     Autonomous Red Team Platform v1.0`)
}

func Execute() {
	rootCmd.PersistentFlags().StringVarP(&cfgPath, "config", "c", "config.yaml", "Path to config file")
	rootCmd.AddCommand(campaignCmd())
	rootCmd.AddCommand(reconCmd())
	rootCmd.AddCommand(agentCmd())
	rootCmd.AddCommand(exploitCmd())
	rootCmd.AddCommand(aiCmd())
	rootCmd.AddCommand(lateralCmd())
	rootCmd.AddCommand(dashboardCmd())
	rootCmd.AddCommand(dbCmd())
	rootCmd.AddCommand(labCmd())
	rootCmd.AddCommand(&cobra.Command{
		Use: "version", Short: "Show version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("X404X v1.0.0 — Go 1.23+ — linux/amd64")
			fmt.Println("Rafael Gálvez | Cisco NetAcad | TFG Cybersecurity")
		},
	})

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
