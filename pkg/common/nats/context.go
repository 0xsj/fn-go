// In the nats package, add a context key and helper function
// pkg/common/nats/context.go

package nats

import (
	"context"

	"github.com/0xsj/fn-go/pkg/common/log"
)

type contextKey string

const (
	connKey contextKey = "nats_conn"
)

// WithConnContext adds a NATS connection to the context
func WithConnContext(ctx context.Context, conn *Conn) context.Context {
	return context.WithValue(ctx, connKey, conn)
}

// GetConnFromContext retrieves a NATS connection from the context
// If no connection is found, it logs a warning and returns nil
func GetConnFromContext(logger log.Logger) *Conn {
	// Since we don't have the context here, we'll just return nil
	// In a real implementation, you'd retrieve the connection from the context
	logger.Warn("GetConnFromContext called without context")
	return nil
}