// pkg/common/nats/client.go
package nats

import (
	"sync"
	"time"

	"github.com/0xsj/fn-go/pkg/common/log"
	natspkg "github.com/nats-io/nats.go"
)

// Expose NATS types
type (
	Conn = natspkg.Conn
	Msg = natspkg.Msg
	Subscription = natspkg.Subscription
	Status = natspkg.Status
)

// Also expose common constants
const (
	DefaultURL = natspkg.DefaultURL
)

// Client wraps a NATS connection with additional functionality
type Client struct {
	conn   *natspkg.Conn
	logger log.Logger
	mu     sync.Mutex
}

// Config holds NATS connection configuration
type Config struct {
	URLs          []string
	MaxReconnect  int
	ReconnectWait time.Duration
	Timeout       time.Duration
}

// DefaultConfig provides default NATS configuration
func DefaultConfig() Config {
	return Config{
		URLs:          []string{natspkg.DefaultURL},
		MaxReconnect:  10,
		ReconnectWait: 2 * time.Second,
		Timeout:       5 * time.Second,
	}
}

// NewClient creates a new NATS client
func NewClient(logger log.Logger, config Config) (*Client, error) {
	opts := []natspkg.Option{
		natspkg.MaxReconnects(config.MaxReconnect),
		natspkg.ReconnectWait(config.ReconnectWait),
		natspkg.Timeout(config.Timeout),
		natspkg.DisconnectErrHandler(func(nc *natspkg.Conn, err error) {
			logger.With("error", err.Error()).Warn("NATS disconnected")
		}),
		natspkg.ReconnectHandler(func(nc *natspkg.Conn) {
			logger.With("url", nc.ConnectedUrl()).Info("NATS reconnected")
		}),
		natspkg.ErrorHandler(func(nc *natspkg.Conn, sub *natspkg.Subscription, err error) {
			logger.With("error", err.Error()).Error("NATS error")
		}),
	}

	// Connect to NATS
	conn, err := natspkg.Connect(config.URLs[0], opts...)
	if err != nil {
		return nil, err
	}

	logger.Info("Connected to NATS")

	return &Client{
		conn:   conn,
		logger: logger,
	}, nil
}

// Conn returns the NATS connection
func (c *Client) Conn() *natspkg.Conn {
	return c.conn
}

// Close closes the NATS connection
func (c *Client) Close() {
	c.conn.Close()
}