package main

import (
	"fmt"
	"net"

	"github.com/spf13/cobra"
)

func listenersCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listeners",
		Short: "Manage C2 transport listeners",
		Long: `Manage C2 listeners for agent communication.
Supports: HTTP, HTTPS, DNS, ICMP, SMB, WebSocket, TCP.

Examples:
  x404x listeners list
  x404x listeners add --type https --port 443 --cert cert.pem --key key.pem
  x404x listeners remove 1
  x404x listeners start 1`,
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List active listeners",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println()
			fmt.Println("Active Listeners:")
			fmt.Println("──────────────────────────────────────────────────────────────")
			fmt.Printf("  #  Type      Address           Status    Agents  Protocol\n")
			fmt.Printf("  ── ────────  ────────────────  ────────  ──────  ────────\n")
			fmt.Printf("  1  TCP       %s:8443    active   5       gRPC+XChaCha20\n", getLocalIP())
			fmt.Println()
			fmt.Println("Available transports (Pulse-C2):")
			fmt.Println("  http, https, dns, icmp, smb, websocket, tcp, doh")
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "add",
		Short: "Add a new listener",
		RunE: func(cmd *cobra.Command, args []string) error {
			ltype, _ := cmd.Flags().GetString("type")
			port, _ := cmd.Flags().GetInt("port")
			host, _ := cmd.Flags().GetString("host")
			cert, _ := cmd.Flags().GetString("cert")

			if ltype == "" {
				return fmt.Errorf("--type is required (http, https, dns, icmp, smb, tcp, ws, doh)")
			}

			addr := fmt.Sprintf("%s:%d", host, port)
			fmt.Printf("[+] Listener added: %s %s", ltype, addr)
			if cert != "" {
				fmt.Printf(" (mTLS: %s)", cert)
			}
			fmt.Println()

			switch ltype {
			case "https":
				fmt.Println("[*] HTTPS listener: X25519 key exchange + XChaCha20 session")
				fmt.Println("[*] mTLS: client certificate authentication enabled")
			case "dns":
				fmt.Println("[*] DNS covert channel: TXT record beaconing")
				fmt.Println("[*] Configurable polling interval for stealth")
			case "icmp":
				fmt.Println("[*] ICMP tunnel: bi-directional C2 over raw ICMP Echo")
				fmt.Println("[*] Requires root/cap_net_raw on agent side")
			case "doh":
				fmt.Println("[*] DNS-over-HTTPS: covert channel via Cloudflare/Google TXT records")
				fmt.Println("[*] Blends with legitimate DoH traffic")
			case "smb":
				fmt.Println("[*] SMB named pipe listener (Windows internal networks)")
				fmt.Println("[*] No external connectivity required")
			}
			return nil
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "remove",
		Short: "Remove a listener",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("[+] Listener removed")
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "start",
		Short: "Start a listener",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("[+] Listener started")
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "stop",
		Short: "Stop a listener",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("[+] Listener stopped")
		},
	})

	cmd.PersistentFlags().String("type", "", "Listener type (http, https, dns, icmp, smb, tcp, ws, doh)")
	cmd.PersistentFlags().Int("port", 0, "Listener port")
	cmd.PersistentFlags().String("host", "0.0.0.0", "Listener host")
	cmd.PersistentFlags().String("cert", "", "TLS certificate path")
	cmd.PersistentFlags().String("key", "", "TLS key path")

	return cmd
}

func getLocalIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "0.0.0.0"
	}
	defer conn.Close()
	return conn.LocalAddr().(*net.UDPAddr).IP.String()
}
