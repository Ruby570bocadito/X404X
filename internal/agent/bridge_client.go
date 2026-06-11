// Package agent provides the BridgeClient for Go↔Python IPC communication.
//
// The BridgeClient connects to a Python BridgeServer over TCP or Unix socket,
// sending JSON-framed requests and receiving responses. This enables the
// Go agent to invoke Python modules (Horizon-Intel, Specter-Terminal,
// Apex-Automation, Wormy-ML, etc.) without embedding Python.
//
// Protocol (over TCP or Unix socket):
//
//	Request:  [4-byte MSB length prefix][JSON body]
//	Response: [4-byte MSB length prefix][JSON body]
//
// JSON format:
//
//	Request:  {"module": "recon", "function": "scan", "params": {...}, "timeout_ms": 5000}
//	Response: {"success": true, "result": {...}, "elapsed_ms": 123}
package agent

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os/exec"
	"sync"
	"time"

	"github.com/ruby570bocadito/x404x/pkg/shared/config"
	"github.com/ruby570bocadito/x404x/pkg/shared/logger"
)

// BridgeRequest is a request to the Python bridge.
type BridgeRequest struct {
	Module    string                 `json:"module"`
	Function  string                 `json:"function"`
	Params    map[string]interface{} `json:"params"`
	TimeoutMS int64                  `json:"timeout_ms"`
}

// BridgeResponse is a response from the Python bridge.
type BridgeResponse struct {
	Success   bool                   `json:"success"`
	Result    map[string]interface{} `json:"result,omitempty"`
	Error     string                 `json:"error,omitempty"`
	ElapsedMS int64                  `json:"elapsed_ms"`
}

// BridgeClient communicates with the Python BridgeServer over TCP.
type BridgeClient struct {
	cfg      *config.Config
	log      *logger.Logger
	address  string
	conn     net.Conn
	mu       sync.Mutex
	cmd      *exec.Cmd
}

// NewBridgeClient creates a new bridge client.
// If bridgePath is empty, connects to an already-running bridge server.
// If bridgePath is set, starts the bridge as a subprocess.
func NewBridgeClient(cfg *config.Config, log *logger.Logger) *BridgeClient {
	port := 9100
	if cfg != nil && cfg.Agent.BridgePort > 0 {
		port = cfg.Agent.BridgePort
	}
	addr := fmt.Sprintf("127.0.0.1:%d", port)

	return &BridgeClient{
		cfg:     cfg,
		log:     log,
		address: addr,
	}
}

// Connect establishes a TCP connection to the Python bridge server.
func (bc *BridgeClient) Connect(ctx context.Context) error {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	if bc.conn != nil {
		return fmt.Errorf("already connected")
	}

	var d net.Dialer
	d.Timeout = 5 * time.Second

	conn, err := d.DialContext(ctx, "tcp", bc.address)
	if err != nil {
		return fmt.Errorf("connecting to bridge at %s: %w", bc.address, err)
	}

	bc.conn = conn
	bc.log.Infof("bridge connected at %s", bc.address)

	return nil
}

// StartBridge starts the Python bridge as a subprocess and connects to it.
func (bc *BridgeClient) StartBridge(ctx context.Context, bridgeScript string) error {
	bc.log.Infof("starting Python bridge: %s", bridgeScript)

	port := 9100
	if bc.cfg != nil && bc.cfg.Agent.BridgePort > 0 {
		port = bc.cfg.Agent.BridgePort
	}

	bc.cmd = exec.CommandContext(ctx, "python3", bridgeScript,
		"--host", "127.0.0.1",
		"--port", fmt.Sprintf("%d", port),
	)

	bc.cmd.Stderr = nil // discard stderr in production

	if err := bc.cmd.Start(); err != nil {
		return fmt.Errorf("starting bridge process: %w", err)
	}

	// Wait for bridge to be ready
	time.Sleep(500 * time.Millisecond)

	// Connect
	if err := bc.Connect(ctx); err != nil {
		bc.cmd.Process.Kill()
		return fmt.Errorf("connecting to bridge after start: %w", err)
	}

	bc.log.Info("Python bridge started and connected")
	return nil
}

// Call invokes a module function on the Python bridge.
func (bc *BridgeClient) Call(ctx context.Context, module, function string, params map[string]interface{}) (*BridgeResponse, error) {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	if bc.conn == nil {
		return nil, fmt.Errorf("not connected to bridge")
	}

	req := BridgeRequest{
		Module:    module,
		Function:  function,
		Params:    params,
		TimeoutMS: 30000,
	}

	data, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshaling request: %w", err)
	}

	// Send: 4-byte length prefix + JSON body
	lengthBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(lengthBuf, uint32(len(data)))

	if _, err := bc.conn.Write(lengthBuf); err != nil {
		return nil, fmt.Errorf("writing length: %w", err)
	}
	if _, err := bc.conn.Write(data); err != nil {
		return nil, fmt.Errorf("writing body: %w", err)
	}

	bc.log.Debugf("bridge call: %s.%s (params=%v)", module, function, params)

	// Receive: 4-byte length prefix + JSON body
	respLengthBuf := make([]byte, 4)
	if _, err := io.ReadFull(bc.conn, respLengthBuf); err != nil {
		return nil, fmt.Errorf("reading response length: %w", err)
	}

	respLen := binary.BigEndian.Uint32(respLengthBuf)
	if respLen > 10*1024*1024 { // 10MB max
		return nil, fmt.Errorf("response too large: %d bytes", respLen)
	}

	respBody := make([]byte, respLen)
	if _, err := io.ReadFull(bc.conn, respBody); err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}

	var resp BridgeResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("unmarshaling response: %w", err)
	}

	if !resp.Success {
		return &resp, fmt.Errorf("bridge error: %s", resp.Error)
	}

	bc.log.Debugf("bridge response: %s.%s → %v (elapsed=%dms)", module, function, resp.Result, resp.ElapsedMS)
	return &resp, nil
}

// CallModule is a convenience method that calls a module's default handler.
func (bc *BridgeClient) CallModule(ctx context.Context, module string, params map[string]interface{}) (*BridgeResponse, error) {
	return bc.Call(ctx, module, "execute", params)
}

// HealthCheck verifies the bridge is responsive.
func (bc *BridgeClient) HealthCheck(ctx context.Context) (*BridgeResponse, error) {
	return bc.Call(ctx, "health", "check", nil)
}

// ListModules returns registered modules from the bridge.
func (bc *BridgeClient) ListModules(ctx context.Context) ([]string, error) {
	resp, err := bc.Call(ctx, "health", "list_modules", nil)
	if err != nil {
		return nil, err
	}

	if modules, ok := resp.Result["modules"].([]interface{}); ok {
		names := make([]string, 0, len(modules))
		for _, m := range modules {
			if s, ok := m.(string); ok {
				names = append(names, s)
			}
		}
		return names, nil
	}

	return nil, nil
}

// Disconnect closes the bridge connection and stops the subprocess.
func (bc *BridgeClient) Disconnect() error {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	if bc.conn != nil {
		bc.conn.Close()
		bc.conn = nil
	}

	if bc.cmd != nil && bc.cmd.Process != nil {
		bc.cmd.Process.Kill()
		bc.cmd = nil
	}

	bc.log.Info("bridge disconnected")
	return nil
}

// Connected returns whether the bridge is connected.
func (bc *BridgeClient) Connected() bool {
	bc.mu.Lock()
	defer bc.mu.Unlock()
	return bc.conn != nil
}

// IsConnected is the interface-compatible alias for Connected().
func (bc *BridgeClient) IsConnected() bool {
	return bc.Connected()
}

// CallRaw calls the bridge and returns the raw response for the dispath interface.
func (bc *BridgeClient) CallRaw(ctx context.Context, module, function string, params map[string]interface{}) (map[string]interface{}, error) {
	resp, err := bc.Call(ctx, module, function, params)
	if err != nil {
		return nil, err
	}
	return resp.Result, nil
}
