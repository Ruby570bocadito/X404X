// Package c2server provides a lightweight C2 listener that agents connect to.
// It uses gRPC and implements the AgentService defined in agent.proto.
// In production, this is replaced by Pulse-C2. For the unified framework,
// it serves as the integrated C2 that the agent, orchestrator, and dashboard share.
package c2server

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/ruby570bocadito/x404x/shared/config"
	"github.com/ruby570bocadito/x404x/shared/logger"
	"github.com/ruby570bocadito/x404x/shared/types"
)

// Server is the integrated C2 listener.
// Agents connect to it, and it feeds data to the shared AppState.
type Server struct {
	cfg     *config.Config
	log     *logger.Logger
	mu      sync.RWMutex
	agents  map[string]*AgentConnection
	running bool
}

// AgentConnection represents a connected agent.
type AgentConnection struct {
	Agent      *types.Agent
	ConnectedAt time.Time
	LastSeen   time.Time
	SessionKey []byte
}

// New creates a new C2 server.
func New(cfg *config.Config, log *logger.Logger) *Server {
	return &Server{
		cfg:    cfg,
		log:    log,
		agents: make(map[string]*AgentConnection),
	}
}

// Start begins listening for agent connections.
func (s *Server) Start(ctx context.Context) error {
	addr := fmt.Sprintf("%s:%d", s.cfg.Server.Host, s.cfg.Server.Port)
	s.log.Infof("C2 server starting on %s", addr)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("listening on %s: %w", addr, err)
	}

	s.running = true

	go func() {
		<-ctx.Done()
		listener.Close()
		s.running = false
	}()

	for s.running {
		conn, err := listener.Accept()
		if err != nil {
			if s.running {
				s.log.Errorf("accept error: %v", err)
			}
			continue
		}

		go s.handleConnection(conn)
	}

	return nil
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	s.log.Infof("agent connected from %s", conn.RemoteAddr())

	// In production: perform gRPC handshake with X25519 key exchange
	// For now: register the connection

	agentID := fmt.Sprintf("agent-%d", time.Now().UnixNano())

	ac := &AgentConnection{
		Agent: &types.Agent{
			ID:        agentID,
			LocalIP:   conn.RemoteAddr().String(),
			Status:    types.AgentStatusOnline,
			FirstSeen: time.Now(),
			LastCheckin: time.Now(),
		},
		ConnectedAt: time.Now(),
		LastSeen:   time.Now(),
	}

	s.mu.Lock()
	s.agents[agentID] = ac
	count := len(s.agents)
	s.mu.Unlock()

	s.log.Infof("agent registered: %s (total: %d)", agentID, count)

	// Keep connection alive and handle heartbeats
	buf := make([]byte, 4096)
	for {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		_, err := conn.Read(buf)
		if err != nil {
			break
		}

		s.mu.Lock()
		if ac, ok := s.agents[agentID]; ok {
			ac.LastSeen = time.Now()
			ac.Agent.LastCheckin = time.Now()
		}
		s.mu.Unlock()
	}

	s.mu.Lock()
	if ac, ok := s.agents[agentID]; ok {
		ac.Agent.Status = types.AgentStatusDead
	}
	delete(s.agents, agentID)
	s.mu.Unlock()

	s.log.Infof("agent disconnected: %s", agentID)
}

// GetAgents returns all connected agents.
func (s *Server) GetAgents() []*types.Agent {
	s.mu.RLock()
	defer s.mu.RUnlock()

	agents := make([]*types.Agent, 0, len(s.agents))
	for _, ac := range s.agents {
		agents = append(agents, ac.Agent)
	}
	return agents
}

// AgentCount returns the number of connected agents.
func (s *Server) AgentCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.agents)
}

// Stop shuts down the C2 server.
func (s *Server) Stop() {
	s.running = false
	s.log.Info("C2 server stopped")
}
