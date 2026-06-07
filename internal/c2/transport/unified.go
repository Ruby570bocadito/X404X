package transport

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"sync"
)

type UnifiedTransport struct {
	mu        sync.RWMutex
	connections map[string]net.Conn
	useWebRTC bool
	useMQTT   bool
	active    bool
}

func NewUnifiedTransport() *UnifiedTransport {
	return &UnifiedTransport{
		connections: make(map[string]net.Conn),
		useWebRTC:   false,
		active:      true,
	}
}

func (ut *UnifiedTransport) EnableWebRTC() {
	ut.mu.Lock()
	defer ut.mu.Unlock()
	ut.useWebRTC = true
}

func (ut *UnifiedTransport) EnableMQTT() {
	ut.mu.Lock()
	defer ut.mu.Unlock()
	ut.useMQTT = true
}

func (ut *UnifiedTransport) Dial(ctx context.Context, address string) (net.Conn, error) {
	ut.mu.RLock()
	defer ut.mu.RUnlock()

	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("transport dial: %w", err)
	}

	if ut.useWebRTC {
		return ut.wrapWebRTC(conn), nil
	}

	return conn, nil
}

func (ut *UnifiedTransport) DialTLS(ctx context.Context, address string, tlsConfig *tls.Config) (net.Conn, error) {
	rawConn, err := ut.Dial(ctx, address)
	if err != nil {
		return nil, err
	}

	tlsConn := tls.Client(rawConn, tlsConfig)
	if err := tlsConn.HandshakeContext(ctx); err != nil {
		rawConn.Close()
		return nil, fmt.Errorf("tls handshake: %w", err)
	}

	return tlsConn, nil
}

func (ut *UnifiedTransport) wrapWebRTC(conn net.Conn) net.Conn {
	return conn
}

func (ut *UnifiedTransport) GetStats() map[string]interface{} {
	ut.mu.RLock()
	defer ut.mu.RUnlock()
	return map[string]interface{}{
		"connections": len(ut.connections),
		"webrtc":      ut.useWebRTC,
		"mqtt":        ut.useMQTT,
		"active":      ut.active,
	}
}
