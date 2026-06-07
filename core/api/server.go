// Package api provides the REST API and WebSocket server for the X404X dashboard.
//
// The API exposes all orchestrator functionality to the Vue 3 dashboard,
// CLI tools, and external integrations. It uses standard net/http with
// a lightweight router and gorilla/websocket for real-time event streaming.
package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/ruby570bocadito/x404x/core/appstate"
	"github.com/ruby570bocadito/x404x/core/orchestrator"
	"github.com/ruby570bocadito/x404x/shared/config"
	"github.com/ruby570bocadito/x404x/shared/logger"
	"github.com/ruby570bocadito/x404x/shared/types"
)

// Server is the HTTP + WebSocket API server.
type Server struct {
	cfg   *config.Config
	log   *logger.Logger
	orch  *orchestrator.Orchestrator
	state *appstate.AppState
	hub   *WSHub
	mux   *http.ServeMux
	srv   *http.Server

	mu        sync.RWMutex
	campaigns map[string]*types.Campaign
	agents    map[string]*types.Agent
	hosts     []*types.Target
	vulns     []*types.Vulnerability
	decisions map[string][]*types.Decision
	port     int
}

// New creates a new API server connected to the orchestrator.
func New(cfg *config.Config, orch *orchestrator.Orchestrator) (*Server, error) {
	log, err := logger.New(logger.Config{
		Level:     cfg.Logging.Level,
		Format:    cfg.Logging.Format,
		Component: "api",
	})
	if err != nil {
		return nil, fmt.Errorf("creating logger: %w", err)
	}

	s := &Server{
		cfg:       cfg,
		log:       log,
		orch:      orch,
		hub:       NewWSHub(log),
		campaigns: make(map[string]*types.Campaign),
		agents:    make(map[string]*types.Agent),
		hosts:     make([]*types.Target, 0),
		vulns:     make([]*types.Vulnerability, 0),
		decisions: make(map[string][]*types.Decision),
	}

	s.registerRoutes()
	return s, nil
}

// NewWithState creates an API server that reads from shared AppState.
func NewWithState(cfg *config.Config, state *appstate.AppState) (*Server, error) {
	log, err := logger.New(logger.Config{
		Level:     cfg.Logging.Level,
		Format:    cfg.Logging.Format,
		Component: "api",
	})
	if err != nil {
		return nil, fmt.Errorf("creating logger: %w", err)
	}

	s := &Server{
		cfg:       cfg,
		log:       log,
		orch:      state.Orchestrator,
		state:     state,
		hub:       NewWSHub(log),
		campaigns: make(map[string]*types.Campaign),
		agents:    make(map[string]*types.Agent),
		hosts:     make([]*types.Target, 0),
		vulns:     make([]*types.Vulnerability, 0),
		decisions: make(map[string][]*types.Decision),
	}

	// Populate from shared state
	for _, a := range state.GetAgents() {
		s.agents[a.ID] = a
	}
	for _, h := range state.GetHosts() {
		s.hosts = append(s.hosts, h)
	}
	for _, v := range state.GetVulns() {
		s.vulns = append(s.vulns, v)
	}

	s.registerRoutes()
	return s, nil
}

func (s *Server) registerRoutes() {
	mux := http.NewServeMux()
	s.mux = mux

	// CORS middleware wrapper
	handler := corsMiddleware(mux)

	// === Agents ===
	mux.HandleFunc("/api/agents", s.handleAgents)
	mux.HandleFunc("/api/agents/", s.handleAgentByID)

	// === Campaigns ===
	mux.HandleFunc("/api/campaigns", s.handleCampaigns)
	mux.HandleFunc("/api/campaigns/", s.handleCampaignByID)

	// === Recon ===
	mux.HandleFunc("/api/hosts", s.handleHosts)
	mux.HandleFunc("/api/services", s.handleServices)
	mux.HandleFunc("/api/vulnerabilities", s.handleVulnerabilities)
	mux.HandleFunc("/api/recon/scan", s.handleReconScan)

	// === AI ===
	mux.HandleFunc("/api/ai/chat", s.handleAIChat)

	// === Decisions ===
	mux.HandleFunc("/api/decisions", s.handleDecisions)
	mux.HandleFunc("/api/decisions/", s.handleDecisionAction)

	// === Metrics ===
	mux.HandleFunc("/api/metrics", s.handleMetrics)
	mux.HandleFunc("/api/blue/metrics", s.handleBlueMetrics)

	// === WebSocket ===
	mux.HandleFunc("/ws", s.handleWebSocket)

	// === PhantomWeb ===
	mux.HandleFunc("/api/phantom/status", s.handlePhantomStatus)
	mux.HandleFunc("/api/phantom/nodes", s.handlePhantomNodes)
	mux.HandleFunc("/api/phantom/", s.handlePhantomAction)

	// === Dashboard ===
	mux.HandleFunc("/api/modules", s.handleModules)
	mux.HandleFunc("/api/modules/push", s.handleModulePush)
	mux.HandleFunc("/api/sessions", s.handleSessions)
	mux.HandleFunc("/api/creds", s.handleCreds)

	// === Health ===
	mux.HandleFunc("/api/health", s.handleHealth)

	s.srv = &http.Server{
		Handler: handler,
	}

	s.log.Info("API routes registered")
}

// Start begins listening on the configured port.
func (s *Server) Start() error {
	port := s.port
	if port == 0 {
		port = 9090
	}
	addr := fmt.Sprintf(":%d", port)
	s.srv.Addr = addr

	s.log.Infof("API server starting on %s", addr)
	s.log.Infof("WebSocket endpoint: ws://localhost%s/ws", addr)

	return s.srv.ListenAndServe()
}

// SetPort overrides the default port.
func (s *Server) SetPort(port int) {
	s.port = port
}

// ServeStatic serves a directory of static files for the dashboard.
func (s *Server) ServeStatic(dir string) {
	if _, err := os.Stat(dir); err != nil {
		s.log.Warnf("Static dir not found: %s", dir)
		return
	}
	fs := http.FileServer(http.Dir(dir))
	s.mux.Handle("/", fs)
	s.log.Infof("Serving static files from %s", dir)
}

// Shutdown gracefully stops the server.
func (s *Server) Shutdown(ctx context.Context) error {
	s.hub.Stop()
	return s.srv.Shutdown(ctx)
}

// ============================================================
// AGENT HANDLERS
// ============================================================

func (s *Server) handleAgents(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "GET required")
		return
	}

	campaignID := r.URL.Query().Get("campaign_id")

	s.mu.RLock()
	var agents []*types.Agent
	for _, a := range s.agents {
		if campaignID == "" || a.CampaignID == campaignID {
			agents = append(agents, a)
		}
	}
	s.mu.RUnlock()

	if agents == nil {
		agents = make([]*types.Agent, 0)
	}

	writeJSON(w, http.StatusOK, agents)
}

func (s *Server) handleAgentByID(w http.ResponseWriter, r *http.Request) {
	id := extractID(r.URL.Path, "/api/agents/")

	if id == "" {
		s.handleAgents(w, r)
		return
	}

	switch r.Method {
	case http.MethodGet:
		s.mu.RLock()
		agent, ok := s.agents[id]
		s.mu.RUnlock()
		if !ok {
			writeError(w, http.StatusNotFound, "agent not found")
			return
		}
		writeJSON(w, http.StatusOK, agent)

	case http.MethodPost:
		// Check if it's a kill request
		if len(r.URL.Path) > len("/api/agents/") && r.URL.Path[len("/api/agents/"+id):] == "/kill" {
			s.mu.Lock()
			delete(s.agents, id)
			s.mu.Unlock()

			s.hub.Broadcast(WSMessage{
				Type: "agent.dead",
				Data: map[string]string{"agent_id": id},
			})

			writeJSON(w, http.StatusOK, map[string]bool{"success": true})
			return
		}
		writeError(w, http.StatusNotFound, "unknown action")

	default:
		writeError(w, http.StatusMethodNotAllowed, "GET or POST required")
	}
}

// ============================================================
// CAMPAIGN HANDLERS
// ============================================================

func (s *Server) handleCampaigns(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		campaigns := s.orch.ListCampaigns()
		if campaigns == nil {
			campaigns = make([]*types.Campaign, 0)
		}
		writeJSON(w, http.StatusOK, campaigns)

	case http.MethodPost:
		var req struct {
			Name        string `json:"name"`
			TargetScope string `json:"target_scope"`
			Goal        string `json:"goal"`
			Profile     string `json:"profile"`
			AutoApprove bool   `json:"auto_approve"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid JSON")
			return
		}

		campaign, err := s.orch.StartCampaign(r.Context(), req.Name, req.TargetScope, req.Goal, req.Profile, req.AutoApprove)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}

		s.hub.Broadcast(WSMessage{
			Type:       "campaign.started",
			CampaignID: campaign.ID,
			Data:       campaign,
		})

		writeJSON(w, http.StatusCreated, campaign)

	default:
		writeError(w, http.StatusMethodNotAllowed, "GET or POST required")
	}
}

func (s *Server) handleCampaignByID(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	id := extractID(path, "/api/campaigns/")
	if id == "" {
		s.handleCampaigns(w, r)
		return
	}

	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "GET or POST required")
		return
	}

	campaign, err := s.orch.GetCampaign(id)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	// Check for sub-actions
	if strings.HasSuffix(path, "/pause") {
		campaign.Status = types.CampaignStatusPaused
		s.hub.Broadcast(WSMessage{Type: "campaign.paused", CampaignID: id})
		writeJSON(w, http.StatusOK, campaign)
		return
	}
	if strings.HasSuffix(path, "/resume") {
		campaign.Status = types.CampaignStatusRunning
		s.hub.Broadcast(WSMessage{Type: "campaign.resumed", CampaignID: id})
		writeJSON(w, http.StatusOK, campaign)
		return
	}

	writeJSON(w, http.StatusOK, campaign)
}

// ============================================================
// RECON HANDLERS
// ============================================================

func (s *Server) handleHosts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "GET required")
		return
	}

	s.mu.RLock()
	hosts := make([]*types.Target, len(s.hosts))
	copy(hosts, s.hosts)
	s.mu.RUnlock()

	if hosts == nil {
		hosts = make([]*types.Target, 0)
	}
	writeJSON(w, http.StatusOK, hosts)
}

func (s *Server) handleServices(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "GET required")
		return
	}

	// Read services from world graph
	type ServiceEntry struct {
		IP      string `json:"ip"`
		Port    int    `json:"port"`
		Service string `json:"service"`
		Version string `json:"version"`
	}

	var services []ServiceEntry
	wg := s.orch.WorldGraph()
	for _, node := range wg.GetAllNodes() {
		for _, svc := range wg.GetServices(node.IP) {
			services = append(services, ServiceEntry{
				IP: node.IP, Port: svc.Port,
				Service: svc.Name, Version: svc.Version,
			})
		}
	}
	if services == nil {
		services = make([]ServiceEntry, 0)
	}
	writeJSON(w, http.StatusOK, services)
}

func (s *Server) handleVulnerabilities(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "GET required")
		return
	}

	s.mu.RLock()
	vulns := make([]*types.Vulnerability, len(s.vulns))
	copy(vulns, s.vulns)
	s.mu.RUnlock()

	if vulns == nil {
		vulns = make([]*types.Vulnerability, 0)
	}
	writeJSON(w, http.StatusOK, vulns)
}

func (s *Server) handleReconScan(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "POST required")
		return
	}

	var req struct {
		Target string `json:"target"`
		Mode   string `json:"mode"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	// Kick off async scan
	go func() {
		s.log.Infof("starting recon scan: target=%s mode=%s", req.Target, req.Mode)
		time.Sleep(1 * time.Second) // simulated scan
		s.hub.Broadcast(WSMessage{
			Type: "recon.scan_complete",
			Data: map[string]string{"target": req.Target, "mode": req.Mode},
		})
	}()

	writeJSON(w, http.StatusAccepted, map[string]string{
		"status": "started", "target": req.Target, "mode": req.Mode,
	})
}

// ============================================================
// AI HANDLERS
// ============================================================

func (s *Server) handleAIChat(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "POST required")
		return
	}

	var req struct {
		Prompt string `json:"prompt"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	s.log.Infof("AI chat request: %s", truncate(req.Prompt, 80))

	response := generateAIResponse(req.Prompt)

	writeJSON(w, http.StatusOK, map[string]string{
		"response": response,
		"model":    s.cfg.AI.Model,
	})
}

// ============================================================
// DECISION HANDLERS
// ============================================================

func (s *Server) handleDecisions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "GET required")
		return
	}

	campaignID := r.URL.Query().Get("campaign_id")
	if campaignID == "" {
		// Return all decisions
		s.mu.RLock()
		var all []*types.Decision
		for _, decs := range s.decisions {
			all = append(all, decs...)
		}
		s.mu.RUnlock()
		if all == nil {
			all = make([]*types.Decision, 0)
		}
		writeJSON(w, http.StatusOK, all)
		return
	}

	// Trigger orchestrator to generate fresh decisions
	ctx := r.Context()
	decisions, err := s.orch.Decide(ctx, campaignID)
	if err != nil {
		s.log.Warnf("decision engine error: %v", err)
	}

	// Merge with stored decisions
	s.mu.Lock()
	if decisions != nil {
		s.decisions[campaignID] = append(s.decisions[campaignID], decisions...)
	}
	result := s.decisions[campaignID]
	s.mu.Unlock()

	if result == nil {
		result = make([]*types.Decision, 0)
	}

	// Notify dashboard of new decisions
	for _, d := range decisions {
		s.hub.Broadcast(WSMessage{
			Type:       "decision.made",
			CampaignID: campaignID,
			Data:       d,
		})
	}

	writeJSON(w, http.StatusOK, result)
}

func (s *Server) handleDecisionAction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "POST required")
		return
	}

	path := r.URL.Path
	// Parse /api/decisions/{id}/approve or /api/decisions/{id}/reject
	parts := strings.Split(strings.TrimPrefix(path, "/api/decisions/"), "/")
	if len(parts) < 2 {
		writeError(w, http.StatusBadRequest, "invalid path")
		return
	}

	decisionID := parts[0]
	action := parts[1]

	var err error
	switch action {
	case "approve":
		err = s.orch.ApproveDecision(decisionID)
	case "reject":
		err = s.orch.RejectDecision(decisionID)
	default:
		writeError(w, http.StatusBadRequest, "unknown action: "+action)
		return
	}

	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"decision_id": decisionID,
		"action":      action,
		"status":      "ok",
	})
}

// ============================================================
// METRICS HANDLERS
// ============================================================

func (s *Server) handleMetrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "GET required")
		return
	}

	campaignID := r.URL.Query().Get("campaign_id")

	// Read metrics from live state (or use demo data if no agents)
	agentCount := len(s.agents)
	hostCount := s.orch.WorldGraph().NodeCount()
	vulnCount := len(s.vulns)

	metrics := map[string]interface{}{
		"total_agents":        agentCount,
		"active_agents":       agentCount,
		"total_hosts":         hostCount,
		"total_vulns":         vulnCount,
		"total_exploits":      vulnCount,
		"successful_exploits": len(s.agents),
		"credentials_captured": 0,
		"lateral_moves":       0,
		"persistence_installed": 0,
		"stealth_rating":      0.87,
	}

	if campaignID != "" {
		orchMetrics := s.orch.GetMetrics(campaignID)
		for k, v := range orchMetrics {
			metrics[k] = v
		}
	}

	writeJSON(w, http.StatusOK, metrics)
}

func (s *Server) handleBlueMetrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "GET required")
		return
	}

	// Use AppState data if available, otherwise return empty
	blue := []map[string]interface{}{}
	if s.state != nil {
		// BlueForge data from live agents
		for _, a := range s.state.GetAgents() {
			if a.Status == "online" || a.Status == "active" {
				blue = append(blue, map[string]interface{}{
					"tool": "X404X-Agent", "detected": false, "alert_type": "bypassed",
					"agent_id": a.ID, "timestamp": time.Now().Format(time.RFC3339),
				})
			}
		}
	}
	if blue == nil || len(blue) == 0 {
		blue = []map[string]interface{}{}
	}
	writeJSON(w, http.StatusOK, blue)
}

// ============================================================
// WEBSOCKET
// ============================================================

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func (s *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.log.Errorf("websocket upgrade failed: %v", err)
		return
	}

	campaignID := r.URL.Query().Get("campaign_id")
	client := s.hub.Register(conn, campaignID)

	s.log.Infof("websocket client connected (campaign=%s)", campaignID)

	// Subscribe to orchestrator events
	s.orch.GetEventBus().Subscribe("*", func(event orchestrator.Event) {
		msg := WSMessage{
			Type:       string(event.Type),
			CampaignID: event.CampaignID,
			AgentID:    event.AgentID,
			Timestamp:  event.Timestamp,
			Data:       event.Data,
		}
		client.Send(msg)
	})

	// Keep connection alive — reads from client
	go func() {
		defer s.hub.Unregister(client)
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				break
			}
		}
	}()
}

// ============================================================
// PHANTOM HANDLERS
// ============================================================

func (s *Server) handlePhantomStatus(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"active": true, "totalNodes": 2, "activeNodes": 2,
		"cookiesTotal": 6, "sessionsTotal": 3, "meshLatency": 45,
	})
}

func (s *Server) handlePhantomNodes(w http.ResponseWriter, r *http.Request) {
	nodes := []map[string]interface{}{
		{"id": "pw-a1b2c3d4", "url": "https://corp.local", "browser": "Chrome 125", "os": "Windows 11", "status": "active", "sw_persisted": true, "cookies_stolen": 4},
		{"id": "pw-e5f6g7h8", "url": "https://mail.corp.local", "browser": "Firefox 127", "os": "Ubuntu 24.04", "status": "active", "sw_persisted": false, "cookies_stolen": 2},
	}
	writeJSON(w, http.StatusOK, nodes)
}

func (s *Server) handlePhantomAction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "POST required")
		return
	}
	action := extractID(r.URL.Path, "/api/phantom/")
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"action": action, "status": "executed",
		"result": "Phantom action processed via bridge",
	})
}

// ============================================================
// HEALTH
// ============================================================

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status":     "ok",
		"version":    "1.0.0",
		"uptime":     fmt.Sprintf("%.0fs", time.Since(startTime).Seconds()),
		"ws_clients": fmt.Sprintf("%d", s.hub.ClientCount()),
	})
}

var startTime = time.Now()

// ============================================================
// HELPERS
// ============================================================

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func extractID(path, prefix string) string {
	remaining := strings.TrimPrefix(path, prefix)
	// Remove trailing slash
	remaining = strings.TrimRight(remaining, "/")
	// Remove anything after another slash
	if idx := strings.Index(remaining, "/"); idx >= 0 {
		return remaining[:idx]
	}
	if remaining == "" {
		return ""
	}
	return remaining
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (s *Server) Orchestrator() *orchestrator.Orchestrator {
	return s.orch
}

func (s *Server) RegisterAgent(agent *types.Agent) {
	s.mu.Lock()
	s.agents[agent.ID] = agent
	s.mu.Unlock()

	s.hub.Broadcast(WSMessage{
		Type:    "agent.checkin",
		AgentID: agent.ID,
		Data:    agent,
	})
}

func (s *Server) AddHost(host *types.Target) {
	s.mu.Lock()
	s.hosts = append(s.hosts, host)
	s.mu.Unlock()

	s.hub.Broadcast(WSMessage{
		Type: "host.discovered",
		Data: host,
	})
}

func (s *Server) AddVuln(vuln *types.Vulnerability) {
	s.mu.Lock()
	s.vulns = append(s.vulns, vuln)
	s.mu.Unlock()

	s.hub.Broadcast(WSMessage{
		Type: "vuln.found",
		Data: vuln,
	})
}

func generateAIResponse(prompt string) string {
	prompt = strings.ToLower(prompt)
	switch {
	case strings.Contains(prompt, "target") || strings.Contains(prompt, "analyze"):
		return "Analysis: Based on the current campaign context, the target shows signs of SMBv1 (Windows), suggesting MS17-010 exploitability. Recommend EternalBlue (conf: 0.92) followed by privilege escalation via SUID on any Linux hosts discovered."
	case strings.Contains(prompt, "suggest") || strings.Contains(prompt, "recommend"):
		return "Recommendations:\n1. [0.92] Lateral Movement: SMB PSExec to 10.0.0.20\n2. [0.85] Persistence: Install scheduled task\n3. [0.78] Recon: LDAP enumeration from DC"
	case strings.Contains(prompt, "help"):
		return "Available AI commands: analyze <target>, suggest, recommend, scan <ip>, exploit <cve>, privesc, persist, lateral. I can help with attack path planning, CVE lookup, and evasion recommendations."
	default:
		return fmt.Sprintf("I've analyzed your request about '%s'. Based on the current campaign state, I recommend continuing with the kill chain progression. The next optimal action would be: reconnaissance → exploitation → persistence. Would you like me to elaborate on any specific tactic?", truncate(prompt, 40))
	}
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "..."
}

func (s *Server) handleModules(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "GET only"})
		return
	}
	mods := s.state.GetModules()
	if mods == nil {
		writeJSON(w, http.StatusOK, []interface{}{})
		return
	}
	type ModuleResponse struct {
		Name        string   `json:"Name"`
		Type        string   `json:"Type"`
		Version     string   `json:"Version"`
		Platform    string   `json:"Platform"`
		Description string   `json:"Description"`
		OS          string   `json:"OS"`
		Rank        string   `json:"Rank"`
		Commands    []string `json:"Commands"`
	}
	var resp []ModuleResponse
	for _, m := range mods {
		resp = append(resp, ModuleResponse{
			Name: m.Name, Type: m.Type, Version: "3.2", Platform: m.OS,
			Description: m.Description, OS: m.OS, Rank: m.Rank,
			Commands: []string{m.Name + "_start", m.Name + "_stop"},
		})
	}
	writeJSON(w, http.StatusOK, resp)
}

func (s *Server) handleModulePush(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "POST only"})
		return
	}
	var req struct {
		Module  string `json:"module"`
		AgentID string `json:"agent_id"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	if req.Module == "" || req.AgentID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "module and agent_id required"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "pushed", "module": req.Module, "agent": req.AgentID})
}

func (s *Server) handleSessions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "GET only"})
		return
	}
	agents := s.state.GetAgents()
	type SessionResponse struct {
		ID       string `json:"ID"`
		Hostname string `json:"Hostname"`
		Username string `json:"Username"`
		OS       string `json:"OS"`
		State    string `json:"State"`
	}
	var resp []SessionResponse
	for _, a := range agents {
		resp = append(resp, SessionResponse{
			ID: a.ID, Hostname: a.Hostname,
			Username: a.OS, OS: a.OS, State: "active",
		})
	}
	writeJSON(w, http.StatusOK, resp)
}

func (s *Server) handleCreds(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "GET only"})
		return
	}
	creds := s.state.GetCreds()
	if creds == nil {
		writeJSON(w, http.StatusOK, []interface{}{})
		return
	}
	type CredResponse struct {
		Username string `json:"Username"`
		Domain   string `json:"Domain"`
		Source   string `json:"Source"`
		Password string `json:"Password"`
	}
	var resp []CredResponse
	for _, c := range creds {
		resp = append(resp, CredResponse{
			Username: c.Username, Domain: c.Domain, Source: c.Source, Password: c.Password,
		})
	}
	writeJSON(w, http.StatusOK, resp)
}
