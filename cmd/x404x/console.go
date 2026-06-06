package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Console struct {
	reader    *bufio.Reader
	running   bool
	workspace string
	sessions  map[string]*Session
	modules   map[string]*ModuleInfo
	context   *ModuleContext
}

type Session struct {
	ID       string
	Target   string
	OS       string
	User     string
	Status   string
}

type ModuleInfo struct {
	Name        string
	Type        string
	Description string
	CVE         string
	Rank        string
	OS          string
}

type ModuleContext struct {
	Name    string
	Options map[string]string
}

func NewConsole() *Console {
	c := &Console{
		reader:    bufio.NewReader(os.Stdin),
		running:   true,
		workspace: "default",
		sessions: map[string]*Session{
			"1": {ID: "1", Target: "10.0.0.10", OS: "Windows 2019", User: "NT AUTHORITY\\SYSTEM", Status: "active"},
		},
		modules: map[string]*ModuleInfo{
			"exploit/eternalblue": {
				Name: "exploit/eternalblue", Type: "exploit",
				Description: "MS17-010 EternalBlue SMB Remote Code Execution",
				CVE: "MS17-010", Rank: "great", OS: "Windows 7/2008/2019",
			},
			"exploit/bluekeep": {
				Name: "exploit/bluekeep", Type: "exploit",
				Description: "CVE-2019-0708 BlueKeep RDP Remote Code Execution",
				CVE: "CVE-2019-0708", Rank: "great", OS: "Windows 7/2008",
			},
			"exploit/kerberoast": {
				Name: "exploit/kerberoast", Type: "exploit",
				Description: "Kerberoasting — TGS ticket extraction",
				CVE: "", Rank: "normal", OS: "Windows AD",
			},
			"exploit/privesc_suid": {
				Name: "exploit/privesc_suid", Type: "exploit",
				Description: "SUID binary privilege escalation via GTFOBins",
				CVE: "", Rank: "excellent", OS: "Linux",
			},
			"auxiliary/recon_tcp": {
				Name: "auxiliary/recon_tcp", Type: "auxiliary",
				Description: "TCP port scanner",
				CVE: "", Rank: "normal", OS: "any",
			},
		},
	}
	return c
}

func (c *Console) Run() error {
	c.printBanner()

	for c.running {
		prompt := "x404x"
		if c.context != nil {
			prompt = fmt.Sprintf("x404x (%s)", c.context.Name)
		}
		fmt.Printf("\n%s > ", prompt)

		input, err := c.reader.ReadString('\n')
		if err != nil {
			break
		}

		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		c.dispatch(input)
	}

	return nil
}

func (c *Console) printBanner() {
	fmt.Print(accentStyle.Render(`
██╗  ██╗ ██╗  ██╗  ██████╗  ██╗  ██╗ ██╗  ██╗
╚██╗██╔╝ ██║  ██║ ██╔═══██╗ ██║  ██║ ╚██╗██╔╝
 ╚███╔╝  ███████║ ██║   ██║ ███████║  ╚███╔╝
 ██╔██╗  ╚════██║ ██║   ██║ ╚════██║  ██╔██╗
██╔╝ ██╗     ██╔╝ ╚██████╔╝     ██╔╝ ██╔╝ ██╗
╚═╝  ╚═╝     ╚═╝   ╚═════╝      ╚═╝  ╚═╝  ╚═╝
`))
	fmt.Println(mutedStyle.Render("     X404X — Autonomous Red Team Platform v1.0"))
	fmt.Println(mutedStyle.Render("     Rafael Gálvez | Cisco NetAcad | TFG Cybersecurity"))
	fmt.Println()
	fmt.Println(mutedStyle.Render(`Type "help" for available commands.`))
}

func (c *Console) dispatch(cmd string) {
	parts := strings.Fields(cmd)
	if len(parts) == 0 {
		return
	}

	action := parts[0]
	args := parts[1:]

	switch action {
	case "?", "help":
		c.cmdHelp(args)
	case "banner":
		c.printBanner()
	case "exit", "quit":
		fmt.Println("[*] Exiting X404X console...")
		c.running = false
	case "version":
		fmt.Println("X404X v1.0.0 — Go 1.22 — Linux/amd64")

	// Campaign
	case "campaign":
		c.cmdCampaign(args)
	case "workspace":
		c.cmdWorkspace(args)

	// Sessions
	case "sessions":
		c.cmdSessions(args)

	// Modules
	case "use":
		c.cmdUse(args)
	case "search":
		c.cmdSearch(args)
	case "show":
		c.cmdShow(args)
	case "set":
		c.cmdSet(args)
	case "unset":
		c.cmdUnset(args)
	case "exploit", "run":
		c.cmdExploit(args)
	case "back":
		c.cmdBack()
	case "info":
		c.cmdInfo(args)

	// AI
	case "ai":
		c.cmdAI(args)
	case "suggest":
		c.cmdSuggest(args)

	// Database
	case "db_status":
		c.cmdDBStatus()
	case "hosts":
		c.cmdHosts()
	case "services":
		c.cmdServices()
	case "creds":
		c.cmdCreds()
	case "vulns":
		c.cmdVulns()

	// Lab
	case "lab":
		c.cmdLab(args)

	default:
		fmt.Printf("%s Unknown command: %s\n", dangerStyle.Render("[-]"), cmd)
		fmt.Println(mutedStyle.Render("    Type 'help' for available commands."))
	}
}

func (c *Console) cmdHelp(args []string) {
	modules := c.context != nil

	fmt.Println()
	fmt.Println(titleStyle.Render("Core Commands"))
	fmt.Println(strings.Repeat("=", 50))
	fmt.Println("  ?              Help menu")
	fmt.Println("  banner         Display banner")
	fmt.Println("  exit           Exit console")
	fmt.Println("  version        Show version")

	fmt.Println()
	fmt.Println(titleStyle.Render("Campaign Commands"))
	fmt.Println(strings.Repeat("=", 50))
	fmt.Println("  campaign       Manage campaigns (list, new, load)")
	fmt.Println("  sessions       List active agent sessions")
	fmt.Println("  sessions -i N  Interact with session N")
	fmt.Println("  workspace      Switch workspaces")

	fmt.Println()
	fmt.Println(titleStyle.Render("Module Commands"))
	fmt.Println(strings.Repeat("=", 50))
	if modules {
		fmt.Println(accentStyle.Render("  show options   Show current module options"))
		fmt.Println(accentStyle.Render("  set <var> <val>  Set module option"))
		fmt.Println(accentStyle.Render("  unset <var>    Unset module option"))
		fmt.Println(accentStyle.Render("  exploit        Execute current module (same as run)"))
		fmt.Println(accentStyle.Render("  run            Execute current module"))
		fmt.Println(accentStyle.Render("  back           Unload current module"))
	}
	fmt.Println("  use <module>   Interact with a module")
	fmt.Println("  search <term>  Search modules by CVE, service, OS")
	fmt.Println("  info <module>  Show module details")

	fmt.Println()
	fmt.Println(titleStyle.Render("AI Commands"))
	fmt.Println(strings.Repeat("=", 50))
	fmt.Println("  ai <prompt>    Ask AI assistant (context-aware)")
	fmt.Println("  suggest        Get AI attack suggestions")

	fmt.Println()
	fmt.Println(titleStyle.Render("Database Commands"))
	fmt.Println(strings.Repeat("=", 50))
	fmt.Println("  db_status      Show database connection")
	fmt.Println("  hosts          List discovered hosts")
	fmt.Println("  services       List discovered services")
	fmt.Println("  creds          List captured credentials")
	fmt.Println("  vulns          List discovered vulnerabilities")

	fmt.Println()
	fmt.Println(titleStyle.Render("Lab Commands"))
	fmt.Println(strings.Repeat("=", 50))
	fmt.Println("  lab up         Start lab environment")
	fmt.Println("  lab down       Stop lab environment")
	fmt.Println("  lab status     Show lab status")
}

func (c *Console) cmdCampaign(args []string) {
	if len(args) == 0 {
		fmt.Println()
		fmt.Println(titleStyle.Render("Active Campaigns"))
		fmt.Println(strings.Repeat("=", 50))
		fmt.Println("  ID        Name       Status    Phase         Agents  Progress")
		fmt.Println("  TFG-001   TFG-Demo   running   exploitation  5       67%")
		return
	}
	switch args[0] {
	case "list":
		fmt.Println("[*] 1 campaign(s): TFG-Demo (running)")
	case "new":
		fmt.Println("[+] Campaign created. Use 'campaign start <id>' to begin.")
	default:
		fmt.Printf("[-] Unknown campaign subcommand: %s\n", args[0])
	}
}

func (c *Console) cmdWorkspace(args []string) {
	if len(args) == 0 {
		fmt.Printf("[*] Current workspace: %s\n", c.workspace)
		return
	}
	c.workspace = args[0]
	fmt.Printf("[+] Workspace: %s\n", c.workspace)
}

func (c *Console) cmdSessions(args []string) {
	if len(args) >= 2 && args[0] == "-i" {
		id := args[1]
		sess, ok := c.sessions[id]
		if !ok {
			fmt.Printf("[-] Session %s not found\n", id)
			return
		}
		fmt.Printf("[*] Interacting with Session %s (%s@%s)\n", sess.ID, sess.User, sess.Target)
		fmt.Printf("[*] Use 'background' to return or 'exit' to close.\n")
		return
	}

	fmt.Println()
	fmt.Println(titleStyle.Render("Active Sessions"))
	fmt.Println(strings.Repeat("=", 70))
	fmt.Println(mutedStyle.Render("  Id  Target        OS              User                    Status"))
	fmt.Println(mutedStyle.Render("  ──  ────────────  ──────────────  ──────────────────────  ──────"))
	for _, s := range c.sessions {
		fmt.Printf("  %-3s %-13s %-15s %-23s %s\n",
			s.ID, s.Target, s.OS, s.User, accentStyle.Render(s.Status))
	}
}

func (c *Console) cmdUse(args []string) {
	if len(args) == 0 {
		fmt.Println("[-] Usage: use <module_name>")
		return
	}

	mod, ok := c.modules[args[0]]
	if !ok {
		fmt.Printf("[-] Module not found: %s\n", args[0])
		fmt.Println(mutedStyle.Render("    Use 'search' to find modules."))
		return
	}

	c.context = &ModuleContext{
		Name: mod.Name,
		Options: map[string]string{
			"RHOSTS": "",
			"RPORT":  "445",
			"LHOST":  "",
			"LPORT":  "4444",
		},
	}

	fmt.Printf("\n%s %s\n", accentStyle.Render("[+]"), mod.Name)
	fmt.Printf("    %s\n", mod.Description)
	fmt.Printf("    Type: %s | CVE: %s | Rank: %s\n", mod.Type, mod.CVE, mod.Rank)
	fmt.Println(mutedStyle.Render("\n    Use 'show options' to view configurable options."))
}

func (c *Console) cmdSearch(args []string) {
	if len(args) == 0 {
		fmt.Println("[-] Usage: search <term>")
		return
	}

	term := strings.ToLower(strings.Join(args, " "))
	fmt.Printf("\nMatching Modules (search: \"%s\")\n", term)
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println(mutedStyle.Render("  Name                     Type        CVE            Rank       OS"))
	fmt.Println(mutedStyle.Render("  ──────────────────────── ─────────── ────────────── ────────── ────────────"))

	for _, mod := range c.modules {
		if strings.Contains(strings.ToLower(mod.Name), term) ||
			strings.Contains(strings.ToLower(mod.CVE), term) ||
			strings.Contains(strings.ToLower(mod.OS), term) ||
			strings.Contains(strings.ToLower(mod.Description), term) {
			fmt.Printf("  %-24s %-11s %-14s %-10s %s\n",
				mod.Name, mod.Type, mod.CVE, mod.Rank, mod.OS)
		}
	}
}

func (c *Console) cmdShow(args []string) {
	if c.context == nil {
		fmt.Println("[-] No module selected. Use 'use <module>' first.")
		return
	}

	if len(args) > 0 && args[0] == "options" {
		fmt.Println()
		fmt.Printf("%s %s\n", titleStyle.Render("Module:"), c.context.Name)
		fmt.Println()
		fmt.Println(mutedStyle.Render("  Name        Value          Required  Description"))
		fmt.Println(mutedStyle.Render("  ────        ─────          ────────  ───────────"))
		for k, v := range c.context.Options {
			req := "yes"
			if v == "" {
				v = mutedStyle.Render("(empty)")
			}
			fmt.Printf("  %-11s %-15s %-9s\n", k, v, req)
		}
	} else {
		fmt.Println("[-] Usage: show options")
	}
}

func (c *Console) cmdSet(args []string) {
	if c.context == nil {
		fmt.Println("[-] No module selected. Use 'use <module>' first.")
		return
	}
	if len(args) < 2 {
		fmt.Println("[-] Usage: set <option> <value>")
		return
	}
	c.context.Options[args[0]] = strings.Join(args[1:], " ")
	fmt.Printf("%s %s → %s\n", accentStyle.Render("[+]"), args[0], args[1])
}

func (c *Console) cmdUnset(args []string) {
	if c.context == nil {
		fmt.Println("[-] No module selected.")
		return
	}
	if len(args) == 0 {
		fmt.Println("[-] Usage: unset <option>")
		return
	}
	c.context.Options[args[0]] = ""
	fmt.Printf("%s %s unset\n", accentStyle.Render("[+]"), args[0])
}

func (c *Console) cmdExploit(args []string) {
	if c.context == nil {
		fmt.Println("[-] No module selected. Use 'use <module>' first.")
		return
	}

	fmt.Printf("\n[*] Starting %s...\n", c.context.Name)

	// Simulate exploit execution
	switch c.context.Name {
	case "exploit/eternalblue":
		fmt.Println("[*] Target: " + c.context.Options["RHOSTS"] + ":445")
		fmt.Println("[*] Detected: Windows 2019 — SMBv1 enabled")
		fmt.Println("[+] EternalBlue exploit sent successfully")
		fmt.Println(accentStyle.Render("[+] Session 2 opened — NT AUTHORITY\\SYSTEM@10.0.0.10"))
		c.sessions["2"] = &Session{ID: "2", Target: "10.0.0.10", OS: "Windows 2019", User: "NT AUTHORITY\\SYSTEM", Status: "active"}
	case "exploit/privesc_suid":
		fmt.Println("[*] Scanning SUID binaries...")
		fmt.Println("[+] Found SUID: /usr/bin/python3-suid")
		fmt.Println("[+] GTFOBins match: python → os.execl('/bin/sh', 'sh')")
		fmt.Println(accentStyle.Render("[+] Root shell obtained!"))
	default:
		fmt.Println(accentStyle.Render("[+] Exploit executed successfully"))
	}
}

func (c *Console) cmdBack() {
	if c.context != nil {
		fmt.Printf("[*] Unloading %s\n", c.context.Name)
		c.context = nil
	} else {
		fmt.Println("[-] No module loaded.")
	}
}

func (c *Console) cmdInfo(args []string) {
	if len(args) == 0 {
		fmt.Println("[-] Usage: info <module_name>")
		return
	}

	mod, ok := c.modules[args[0]]
	if !ok {
		fmt.Printf("[-] Module not found: %s\n", args[0])
		return
	}

	fmt.Println()
	fmt.Printf("%s\n", titleStyle.Render(mod.Name))
	fmt.Println(strings.Repeat("=", 60))
	fmt.Printf("  Type:        %s\n", mod.Type)
	fmt.Printf("  CVE:         %s\n", mod.CVE)
	fmt.Printf("  Rank:        %s\n", mod.Rank)
	fmt.Printf("  Platform:    %s\n", mod.OS)
	fmt.Printf("\n  Description:\n    %s\n", mod.Description)
}

func (c *Console) cmdAI(args []string) {
	if len(args) == 0 {
		fmt.Println("[-] Usage: ai <prompt>")
		return
	}
	prompt := strings.Join(args, " ")
	fmt.Printf("\n[AI] Processing: \"%s\"\n", prompt)
	fmt.Println("[AI] Response: Based on current target context")
	fmt.Println("     → Recommended vector: privilege escalation via SUID (conf=0.89)")
}

func (c *Console) cmdSuggest(args []string) {
	fmt.Println("[AI] Analyzing current campaign context...")
	fmt.Println()
	fmt.Println(titleStyle.Render("AI Suggestions:"))
	fmt.Println(mutedStyle.Render("  #  Confidence  Tactic          Technique       Target"))
	fmt.Println(mutedStyle.Render("  ── ─────────── ─────────────── ──────────────── ───────────"))
	fmt.Println(accentStyle.Render("  1  0.92        Lateral Move    SMB PSExec       10.0.0.20"))
	fmt.Println("  2  0.85        Persistence     Scheduled Task   Session 1")
	fmt.Println("  3  0.78        Reconnaissance  LDAP Enumeration 10.0.0.10")
	fmt.Println()
	fmt.Println(mutedStyle.Render("  Use 'accept <#> to execute, or 'reject <#>' to skip."))
}

func (c *Console) cmdDBStatus() {
	fmt.Println("[*] Database: SQLite | Connected | x404x.db | 12 tables | 157 records")
}

func (c *Console) cmdHosts() {
	fmt.Println()
	fmt.Println(titleStyle.Render("Discovered Hosts"))
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println(mutedStyle.Render("  IP            Hostname    OS              Status"))
	fmt.Println(mutedStyle.Render("  ────────────  ──────────  ──────────────  ───────────"))
	fmt.Println("  10.0.0.10     DC          Windows 2019    compromised")
	fmt.Println("  10.0.0.20     DB          Ubuntu 24.04    scanned")
	fmt.Println("  10.0.0.50     WS1         Windows 11      compromised")
}

func (c *Console) cmdServices() {
	fmt.Println()
	fmt.Println(titleStyle.Render("Discovered Services"))
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println(mutedStyle.Render("  IP            Port   Service   Version"))
	fmt.Println(mutedStyle.Render("  ────────────  ─────  ────────  ───────"))
	fmt.Println("  10.0.0.10     445    SMB       Windows 2019")
	fmt.Println("  10.0.0.10     3389   RDP       Windows 2019")
	fmt.Println("  10.0.0.20     22     SSH       OpenSSH 9.6")
}

func (c *Console) cmdCreds() {
	fmt.Println()
	fmt.Println(titleStyle.Render("Captured Credentials"))
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println(mutedStyle.Render("  Username    Password      Hash      Domain   Source"))
	fmt.Println(mutedStyle.Render("  ──────────  ────────────  ────────  ───────  ──────"))
	fmt.Println("  admin       password123   -         CORP     SSH")
	fmt.Println("  svc_mssql   P@ssw0rd!     -         CORP     SMB")
}

func (c *Console) cmdVulns() {
	fmt.Println()
	fmt.Println(titleStyle.Render("Discovered Vulnerabilities"))
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println(mutedStyle.Render("  CVE            Severity  Service   Target"))
	fmt.Println(mutedStyle.Render("  ─────────────  ────────  ────────  ───────────"))
	fmt.Println("  MS17-010       critical  SMB       10.0.0.10")
	fmt.Println("  CVE-2024-XXXX  high      apport    10.0.0.20")
}

func (c *Console) cmdLab(args []string) {
	if len(args) == 0 {
		fmt.Println("[-] Usage: lab [up|down|status]")
		return
	}
	switch args[0] {
	case "up":
		fmt.Println("[+] Starting X404X lab environment...")
		fmt.Println("[+] attacker     → 172.20.0.10 (running)")
		fmt.Println("[+] target1      → 172.20.0.20 (running)")
		fmt.Println("[+] dashboard    → http://localhost:3000")
	case "down":
		fmt.Println("[+] Lab environment stopped")
	case "status":
		fmt.Println("[*] Lab status: 5 containers running")
	}
}

func StartConsole(args []string) error {
	c := NewConsole()
	return c.Run()
}
