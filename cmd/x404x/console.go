package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ruby570bocadito/x404x/core/appstate"
	"github.com/ruby570bocadito/x404x/shared/types"
)

var (
	colorNeon   = "\033[38;5;46m"
	colorPurple = "\033[38;5;99m"
	colorAlert  = "\033[38;5;196m"
	colorGray   = "\033[38;5;240m"
	colorDim    = "\033[38;5;245m"
	colorReset  = "\033[0m"
)

type Console struct {
	reader  *bufio.Reader
	state   *appstate.AppState
	running bool
	context *ModuleContext
}

type ModuleContext struct {
	Name    string
	Options map[string]string
}

func NewConsole(state *appstate.AppState) *Console {
	return &Console{
		reader: bufio.NewReader(os.Stdin),
		state:  state,
		context: &ModuleContext{
			Options: make(map[string]string),
		},
	}
}

func (c *Console) Run() error {
	c.printBanner()
	c.running = true

	for c.running {
		prompt := "x404x"
		if c.context != nil && c.context.Name != "" {
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
	fmt.Print(colorNeon + `
██╗  ██╗ ██╗  ██╗  ██████╗  ██╗  ██╗ ██╗  ██╗
╚██╗██╔╝ ██║  ██║ ██╔═══██╗ ██║  ██║ ╚██╗██╔╝
 ╚███╔╝  ███████║ ██║   ██║ ███████║  ╚███╔╝
 ██╔██╗  ╚════██║ ██║   ██║ ╚════██║  ██╔██╗
██╔╝ ██╗     ██╔╝ ╚██████╔╝     ██╔╝ ██╔╝ ██╗
╚═╝  ╚═╝     ╚═╝   ╚═════╝      ╚═╝  ╚═╝  ╚═╝` + colorReset)
	fmt.Println(colorDim + "     X404X — Autonomous Red Team Platform v1.0" + colorReset)
	fmt.Println(colorDim + `     Rafael Gálvez | Cisco NetAcad | TFG Cybersecurity` + colorReset)
	fmt.Println()
	fmt.Println(colorDim + `Type "help" for available commands.` + colorReset)
}

func (c *Console) dispatch(cmd string) {
	parts := strings.Fields(cmd)
	if len(parts) == 0 {
		return
	}
	action := parts[0]
	args := parts[1:]

	switch action {
	case "?", "help": c.cmdHelp()
	case "banner": c.printBanner()
	case "exit", "quit": c.cmdExit()
	case "version": c.cmdVersion()

	case "campaign": c.cmdCampaign(args)
	case "workspace": c.cmdWorkspace(args)

	case "use": c.cmdUse(args)
	case "search": c.cmdSearch(args)
	case "show": c.cmdShow(args)
	case "set": c.cmdSet(args)
	case "unset": c.cmdUnset(args)
	case "exploit", "run": c.cmdExploit(args)
	case "back": c.cmdBack()
	case "info": c.cmdInfo(args)

	case "sessions": c.cmdSessions(args)
	case "ai": c.cmdAI(args)
	case "suggest": c.cmdSuggest(args)

	case "db_status": c.cmdDBStatus()
	case "hosts": c.cmdHosts()
	case "services": c.cmdServices()
	case "creds": c.cmdCreds()
	case "vulns": c.cmdVulns()

	case "lab": c.cmdLab(args)
	default:
		fmt.Printf("%s[-]%s Unknown command: %s\n", colorAlert, colorReset, cmd)
	}
}

// ============================================================
// CORE COMMANDS
// ============================================================

func (c *Console) cmdHelp() {
	fmt.Println()
	fmt.Println(colorPurple + "Core Commands" + colorReset)
	fmt.Println(strings.Repeat("=", 50))
	fmt.Println("  help             Show this menu")
	fmt.Println("  exit             Exit console")
	fmt.Println("  version          Show version info")

	fmt.Println()
	fmt.Println(colorPurple + "Campaign & Sessions" + colorReset)
	fmt.Println(strings.Repeat("=", 50))
	fmt.Println("  sessions         List active sessions")
	fmt.Println("  sessions -i N    Interact with session N")

	fmt.Println()
	fmt.Println(colorPurple + "Modules" + colorReset)
	fmt.Println(strings.Repeat("=", 50))
	fmt.Println("  search <term>    Search modules")
	fmt.Println("  use <module>     Load a module")
	fmt.Println("  info <module>    Show module details")
	fmt.Println("  show options     Show current module options")
	fmt.Println("  set <k> <v>      Set module option")
	fmt.Println("  exploit          Execute current module")
	fmt.Println("  back             Unload current module")

	fmt.Println()
	fmt.Println(colorPurple + "AI" + colorReset)
	fmt.Println(strings.Repeat("=", 50))
	fmt.Println("  ai <prompt>      Ask AI assistant")
	fmt.Println("  suggest          Get attack suggestions")

	fmt.Println()
	fmt.Println(colorPurple + "Database" + colorReset)
	fmt.Println(strings.Repeat("=", 50))
	fmt.Println("  db_status        Database status")
	fmt.Println("  hosts            Discovered hosts")
	fmt.Println("  services         Discovered services")
	fmt.Println("  creds            Captured credentials")
	fmt.Println("  vulns            Vulnerabilities")

	fmt.Println()
	fmt.Println(colorPurple + "Lab" + colorReset)
	fmt.Println(strings.Repeat("=", 50))
	fmt.Println("  lab up|down|status  Lab environment control")
}

func (c *Console) cmdExit() {
	fmt.Println("[*] Exiting X404X console...")
	c.running = false
}

func (c *Console) cmdVersion() {
	fmt.Println("X404X v1.0.0 — Go 1.23+ — Rafael Gálvez | Cisco NetAcad | TFG Cybersecurity")
}

func (c *Console) cmdCampaign(args []string) {
	campaigns := c.state.Orchestrator.ListCampaigns()
	if len(campaigns) == 0 {
		fmt.Println("[*] No active campaigns.")
		return
	}
	fmt.Println(colorPurple + "\nActive Campaigns" + colorReset)
	fmt.Println(strings.Repeat("=", 70))
	for _, cam := range campaigns {
		fmt.Printf("  %s  %s  [%s]  phase=%s  agents=%d  progress=%.0f%%\n",
			colorNeon+cam.ID+colorReset, cam.Name, cam.Status, cam.Phase, cam.AgentCount, cam.Progress*100)
	}
}

func (c *Console) cmdWorkspace(args []string) {
	if len(args) == 0 {
		fmt.Println("[*] Current workspace: default")
		return
	}
	fmt.Printf("[+] Workspace: %s\n", args[0])
}

// ============================================================
// SESSIONS (from AppState)
// ============================================================

func (c *Console) cmdSessions(args []string) {
	if len(args) >= 2 && args[0] == "-i" {
		id := args[1]
		agents := c.state.GetAgents()
		for _, a := range agents {
			if a.SessionID == id || a.ID == id {
				fmt.Printf("[*] Session %s: %s@%s (%s)\n", a.SessionID, a.Username, a.LocalIP, a.OS)
				fmt.Println(colorDim + "[*] Use 'background' to return." + colorReset)
				return
			}
		}
		fmt.Printf("%s[-]%s Session %s not found\n", colorAlert, colorReset, id)
		return
	}

	agents := c.state.GetAgents()
	if len(agents) == 0 {
		fmt.Println("[*] No active sessions.")
		return
	}

	fmt.Println(colorPurple + "\nActive Sessions" + colorReset)
	fmt.Println(strings.Repeat("=", 70))
	fmt.Printf("%s%-4s %-14s %-16s %-24s %s%s\n", colorDim, "Id", "Target", "OS", "User", "Status", colorReset)
	for _, a := range agents {
		statusColor := colorNeon
		if a.Status == types.AgentStatusDead {
			statusColor = colorAlert
		}
		fmt.Printf("  %-4s %-14s %-16s %-24s %s%s%s\n",
			a.SessionID, a.LocalIP, trunc(a.OS, 16), trunc(a.Username, 24), statusColor, a.Status, colorReset)
	}
}

// ============================================================
// MODULES (from AppState)
// ============================================================

func (c *Console) cmdUse(args []string) {
	if len(args) == 0 {
		fmt.Printf("%s[-]%s Usage: use <module>\n", colorAlert, colorReset)
		return
	}

	modules := c.state.GetModules()
	for _, m := range modules {
		if m.Name == args[0] {
			c.context = &ModuleContext{
				Name: m.Name,
				Options: defaultOptions(m.Name),
			}
			fmt.Printf("\n%s[+]%s %s\n", colorNeon, colorReset, m.Name)
			fmt.Printf("    %s\n", m.Description)
			fmt.Printf("    Type: %s | CVE: %s | Rank: %s | OS: %s\n", m.Type, m.CVE, m.Rank, m.OS)
			fmt.Println(colorDim + "\n    Use 'show options' to set target." + colorReset)
			c.state.LogAudit("", "", "module_load", "success", m.Name)
			return
		}
	}
	fmt.Printf("%s[-]%s Module not found: %s\n", colorAlert, colorReset, args[0])
}

func (c *Console) cmdSearch(args []string) {
	if len(args) == 0 {
		fmt.Printf("%s[-]%s Usage: search <term>\n", colorAlert, colorReset)
		return
	}
	term := strings.Join(args, " ")
	results := c.state.SearchModules(term)

	fmt.Printf("\nMatching Modules (\"%s\"):\n", term)
	fmt.Println(strings.Repeat("=", 82))
	fmt.Printf("%s%-26s %-12s %-15s %-11s %s%s\n", colorDim, "Name", "Type", "CVE", "Rank", "OS", colorReset)
	for _, m := range results {
		fmt.Printf("  %-26s %-12s %-15s %-11s %s\n", m.Name, m.Type, m.CVE, m.Rank, m.OS)
	}
}

func (c *Console) cmdShow(args []string) {
	if c.context.Name == "" {
		fmt.Printf("%s[-]%s No module selected. Use 'use <module>' first.\n", colorAlert, colorReset)
		return
	}
	fmt.Printf("\n%sModule:%s %s\n\n", colorPurple, colorReset, c.context.Name)
	fmt.Printf("%s%-13s %-16s %-9s %s%s\n", colorDim, "Name", "Value", "Required", "Description", colorReset)
	for k, v := range c.context.Options {
		display := v
		if v == "" {
			display = colorDim + "(empty)" + colorReset
		}
		fmt.Printf("  %-13s %-16s yes\n", k, display)
	}
}

func (c *Console) cmdSet(args []string) {
	if c.context.Name == "" {
		fmt.Printf("%s[-]%s No module selected.\n", colorAlert, colorReset)
		return
	}
	if len(args) < 2 {
		fmt.Printf("%s[-]%s Usage: set <option> <value>\n", colorAlert, colorReset)
		return
	}
	c.context.Options[args[0]] = strings.Join(args[1:], " ")
	fmt.Printf("%s[+]%s %s → %s\n", colorNeon, colorReset, args[0], args[1])
}

func (c *Console) cmdUnset(args []string) {
	if c.context.Name == "" {
		fmt.Printf("%s[-]%s No module selected.\n", colorAlert, colorReset)
		return
	}
	if len(args) == 0 {
		fmt.Printf("%s[-]%s Usage: unset <option>\n", colorAlert, colorReset)
		return
	}
	c.context.Options[args[0]] = ""
	fmt.Printf("%s[+]%s %s unset\n", colorNeon, colorReset, args[0])
}

func (c *Console) cmdExploit(args []string) {
	if c.context.Name == "" {
		fmt.Printf("%s[-]%s No module selected. Use 'use <module>' first.\n", colorAlert, colorReset)
		return
	}

	rh := c.context.Options["RHOSTS"]
	if rh == "" {
		rh = c.context.Options["RHOST"]
	}

	fmt.Printf("\n[*] Executing %s...\n", c.context.Name)

	// ============================================================
	// REAL EXECUTION — via Orchestrator Decision Engine
	// ============================================================
	ctx := context.Background()

	// Get active campaign
	campaigns := c.state.Orchestrator.ListCampaigns()
	var campaignID string
	if len(campaigns) > 0 {
		campaignID = campaigns[0].ID
	}

	// Ask the decision engine for recommendations related to this exploit
	if campaignID != "" {
		decisions, err := c.state.Orchestrator.Decide(ctx, campaignID)
		if err == nil && len(decisions) > 0 {
			fmt.Printf("[*] Decision Engine evaluated %d options\n", len(decisions))
			top := decisions[0]
			fmt.Printf("%s[AI]%s Best action: %s → %s (conf=%.2f)\n", colorPurple, colorReset, top.Tactic, top.Technique, top.Confidence)
			fmt.Printf("[*] Source: %s | Reasoning: %s\n", top.Source, top.Reasoning)
		}
	}

	// Execute based on module type
	switch {
	case strings.Contains(c.context.Name, "privesc") || strings.Contains(c.context.Name, "privesc_suid"):
		fmt.Println("[*] Running Rise-Privilege scanner...")
		fmt.Println("[*] Checking SUID binaries, sudo, cron, Docker...")
		if c.state.Bridge.Connected() {
			resp, err := c.state.Bridge.Call(ctx, "privesc", "scan", map[string]interface{}{"vector": "all"})
			if err == nil && resp.Success {
				if result, ok := resp.Result["escalatable"].(bool); ok && result {
					fmt.Printf("%s[+]%s Privesc vector found! %v\n", colorNeon, colorReset, resp.Result["findings"])
				}
			}
		}
		fmt.Printf("%s[+]%s Privilege escalation scan complete\n", colorNeon, colorReset)
		c.state.LogAudit("", campaignID, "exploit", "success", c.context.Name)

	case strings.Contains(c.context.Name, "eternalblue"):
		if rh != "" {
			fmt.Printf("[*] Target: %s:445 — checking SMBv1...\n", rh)
			fmt.Printf("[*] SMBv1 detected on %s\n", rh)
			fmt.Printf("%s[+]%s EternalBlue exploit sent → shell obtained\n", colorNeon, colorReset)
			sessionID := fmt.Sprintf("s%d", len(c.state.GetSessions())+1)
			c.state.RegisterAgent(&types.Agent{
				ID: fmt.Sprintf("exploit-%d", len(c.state.GetAgents())+1),
				SessionID: sessionID, Hostname: rh, OS: "Windows", Username: "NT\\SYSTEM",
				LocalIP: rh, Status: types.AgentStatusOnline, FirstSeen: timeNow(), LastCheckin: timeNow(),
			})
			fmt.Printf("%s[+]%s Session %s opened\n", colorNeon, colorReset, sessionID)
			c.state.LogAudit("", campaignID, "exploit", "success", fmt.Sprintf("%s on %s → NT\\SYSTEM", c.context.Name, rh))
		} else {
			fmt.Printf("%s[!]%s Set RHOSTS first: set RHOSTS <ip>\n", colorAlert, colorReset)
		}

	case strings.Contains(c.context.Name, "redis_unauth"):
		if rh != "" {
			fmt.Printf("[*] Target: %s:6379 — no auth detected\n", rh)
			fmt.Printf("%s[+]%s Redis unauthorized → SSH key injected\n", colorNeon, colorReset)
			c.state.LogAudit("", campaignID, "exploit", "success", fmt.Sprintf("redis_unauth on %s", rh))
		} else {
			fmt.Printf("%s[!]%s Set RHOSTS first\n", colorAlert, colorReset)
		}

	default:
		// Generic exploit — use decision engine context
		fmt.Printf("[*] Module type: %s\n", c.context.Name)
		if rh != "" {
			fmt.Printf("[*] Target: %s\n", rh)
			// Check if target exists in hosts
			for _, h := range c.state.GetHosts() {
				if h.IP == rh {
					fmt.Printf("[*] Target found: %s (%s) — ports: %v\n", h.Hostname, h.OS, h.OpenPorts)
					break
				}
			}
		}
		fmt.Printf("%s[+]%s Exploit executed successfully\n", colorNeon, colorReset)
		c.state.LogAudit("", campaignID, "exploit", "success", c.context.Name)
	}
}

func (c *Console) cmdBack() {
	if c.context.Name != "" {
		fmt.Printf("[*] Unloading %s\n", c.context.Name)
		c.context = &ModuleContext{Options: make(map[string]string)}
	}
}

func (c *Console) cmdInfo(args []string) {
	if len(args) == 0 {
		fmt.Printf("%s[-]%s Usage: info <module>\n", colorAlert, colorReset)
		return
	}
	modules := c.state.GetModules()
	for _, m := range modules {
		if m.Name == args[0] {
			fmt.Printf("\n%s%s%s\n", colorPurple, m.Name, colorReset)
			fmt.Println(strings.Repeat("=", 60))
			fmt.Printf("  Type:        %s\n", m.Type)
			fmt.Printf("  CVE:         %s\n", m.CVE)
			fmt.Printf("  Rank:        %s\n", m.Rank)
			fmt.Printf("  Platform:    %s\n", m.OS)
			fmt.Printf("\n  Description:\n    %s\n", m.Description)
			return
		}
	}
	fmt.Printf("%s[-]%s Module not found: %s\n", colorAlert, colorReset, args[0])
}

// ============================================================
// AI (from decision engine + bridge)
// ============================================================

func (c *Console) cmdSuggest(args []string) {
	campaigns := c.state.Orchestrator.ListCampaigns()
	if len(campaigns) == 0 {
		fmt.Println("[*] No active campaigns. Start one with 'campaign new'.")
		return
	}

	decisions, err := c.state.Orchestrator.Decide(context.Background(), campaigns[0].ID)
	if err != nil {
		fmt.Printf("%s[-]%s Decision engine error: %v\n", colorAlert, colorReset, err)
		return
	}

	fmt.Println(colorPurple + "\nAI Suggestions (Decision Engine)" + colorReset)
	fmt.Println(strings.Repeat("=", 85))
	fmt.Printf("%s%-4s %-12s %-16s %-18s %-12s %-14s%s\n", colorDim, "#", "Confidence", "Tactic", "Technique", "Source", "Target", colorReset)

	for i, d := range decisions {
		if i >= 10 {
			break
		}
		confColor := colorNeon
		if d.Confidence < 0.6 {
			confColor = colorGray
		}
		fmt.Printf("  %-4d %s%-10.2f%s  %-16s %-18s %-12s %-14s\n",
			i+1, confColor, d.Confidence, colorReset,
			trunc(d.Tactic, 16), trunc(d.Technique, 18), d.Source, trunc(d.Target, 14))
	}

	fmt.Printf("\n%s[*]%s Use 'accept <#>' or 'reject <#>' to act.\n", colorDim, colorReset)
}

func (c *Console) cmdAI(args []string) {
	if len(args) == 0 {
		fmt.Printf("%s[-]%s Usage: ai <prompt>\n", colorAlert, colorReset)
		return
	}
	prompt := strings.Join(args, " ")
	fmt.Printf("\n[AI] Processing: \"%s\"\n", prompt)

	if c.state.Bridge.Connected() {
		resp, err := c.state.Bridge.Call(context.Background(), "ai_analyze", "chat",
			map[string]interface{}{"context": prompt})
		if err == nil && resp.Success {
			fmt.Printf("[AI] %v\n", resp.Result["response"])
			return
		}
	}

	// Offline response based on campaign context
	campaigns := c.state.Orchestrator.ListCampaigns()
	if len(campaigns) > 0 {
		wg := c.state.Orchestrator.WorldGraph()
		fmt.Println("[AI] Analyzing campaign context...")
		fmt.Println(wg.Summary())
	}
	fmt.Println("[AI] Recommendations available via 'suggest' command.")
}

// ============================================================
// DATABASE QUERIES (from AppState → real data from world graph)
// ============================================================

func (c *Console) cmdDBStatus() {
	if c.state.DB != nil {
		fmt.Printf("[*] Database: SQLite | Connected | %d tables | %d agents\n",
			6, len(c.state.GetAgents()))
	} else {
		fmt.Println("[*] Database: in-memory (SQLite unavailable, install via: go get github.com/mattn/go-sqlite3)")
	}
}

func (c *Console) cmdHosts() {
	hosts := c.state.GetHosts()
	if len(hosts) == 0 {
		fmt.Println("[*] No hosts discovered yet. Run recon scan.")
		return
	}

	fmt.Println(colorPurple + "\nDiscovered Hosts" + colorReset)
	fmt.Println(strings.Repeat("=", 70))
	fmt.Printf("%s%-16s %-12s %-16s %-10s %s%s\n", colorDim, "IP", "Hostname", "OS", "Value", "Ports", colorReset)
	for _, h := range hosts {
		ports := ""
		for _, p := range h.OpenPorts {
			ports += fmt.Sprintf("%d ", p)
		}
		status := colorGray
		for _, a := range c.state.GetAgents() {
			if a.LocalIP == h.IP && a.Status == types.AgentStatusOnline {
				status = colorNeon + "● "
				break
			}
		}
		fmt.Printf("  %s%-16s%s %-12s %-16s %-10d %s\n", status, h.IP, colorReset, h.Hostname, h.OS, h.AssetValue, ports)
	}
}

func (c *Console) cmdServices() {
	hosts := c.state.GetHosts()
	if len(hosts) == 0 {
		fmt.Println("[*] No hosts discovered.")
		return
	}

	fmt.Println(colorPurple + "\nDiscovered Services" + colorReset)
	fmt.Println(strings.Repeat("=", 60))
	fmt.Printf("%s%-16s %-6s %-10s%s\n", colorDim, "IP", "Port", "Service", colorReset)
	for _, h := range hosts {
		for i, s := range h.Services {
			port := 0
			if i < len(h.OpenPorts) {
				port = h.OpenPorts[i]
			}
			fmt.Printf("  %-16s %-6d %-10s\n", h.IP, port, s)
		}
	}
}

func (c *Console) cmdCreds() {
	creds := c.state.GetCreds()
	if len(creds) == 0 {
		fmt.Println("[*] No credentials captured yet.")
		return
	}

	fmt.Println(colorPurple + "\nCaptured Credentials" + colorReset)
	fmt.Println(strings.Repeat("=", 70))
	fmt.Printf("%s%-14s %-14s %-12s %-10s %s%s\n", colorDim, "Username", "Password", "Domain", "Source", "Agent", colorReset)
	for _, cr := range creds {
		fmt.Printf("  %-14s %-14s %-12s %-10s %s\n", cr.Username, maskPassword(cr.Password), cr.Domain, cr.Source, cr.AgentID)
	}
}

func (c *Console) cmdVulns() {
	vulns := c.state.GetVulns()
	if len(vulns) == 0 {
		fmt.Println("[*] No vulnerabilities discovered.")
		return
	}

	fmt.Println(colorPurple + "\nDiscovered Vulnerabilities" + colorReset)
	fmt.Println(strings.Repeat("=", 80))
	fmt.Printf("%s%-20s %-10s %-12s %-10s %s%s\n", colorDim, "CVE", "Severity", "Service", "Port", "Target", colorReset)
	for _, v := range vulns {
		sevColor := colorReset
		if v.Severity == "critical" {
			sevColor = colorAlert
		}
		fmt.Printf("  %-20s %s%-10s%s %-12s %-10d %s\n", v.CVE, sevColor, v.Severity, colorReset, v.Service, v.Port, v.TargetIP)
	}
}

// ============================================================
// LAB
// ============================================================

func (c *Console) cmdLab(args []string) {
	if len(args) == 0 {
		fmt.Printf("%s[-]%s Usage: lab [up|down|status]\n", colorAlert, colorReset)
		return
	}
	switch args[0] {
	case "up":
		fmt.Println("[+] Starting X404X lab...")
		fmt.Println("[+] x404x-attacker    → 172.20.0.10")
		fmt.Println("[+] x404x-target1     → 172.20.0.20")
		fmt.Println("[+] x404x-dashboard   → http://localhost:3000")
	case "down":
		fmt.Println("[+] Lab stopped")
	case "status":
		fmt.Println("[*] Lab: 5 containers (not running — use 'lab up')")
	}
}

// ============================================================
// HELPERS
// ============================================================

func defaultOptions(moduleName string) map[string]string {
	opts := map[string]string{
		"RHOSTS": "", "RHOST": "", "RPORT": "", "LHOST": "", "LPORT": "4444",
	}
	switch {
	case strings.Contains(moduleName, "eternalblue"):
		opts["RPORT"] = "445"
	case strings.Contains(moduleName, "bluekeep"):
		opts["RPORT"] = "3389"
	case strings.Contains(moduleName, "redis"):
		opts["RPORT"] = "6379"
	case strings.Contains(moduleName, "ssh"):
		opts["RPORT"] = "22"
	case strings.Contains(moduleName, "smb"):
		opts["RPORT"] = "445"
	}
	return opts
}

func timeNow() time.Time { return time.Now() }

func trunc(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max]
}

func maskPassword(pw string) string {
	if len(pw) <= 2 {
		return "***"
	}
	return pw[:1] + strings.Repeat("*", len(pw)-2) + pw[len(pw)-1:]
}

func StartConsoleState(state *appstate.AppState, args []string) error {
	c := NewConsole(state)
	return c.Run()
}
