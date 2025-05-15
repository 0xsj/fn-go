package middleware

import (
	"net/http"
	"time"

	"github.com/0xsj/fn-go/pkg/common/log"
)

type Chain struct {
	middlewares []func(http.Handler) http.Handler
	logger log.Logger
}

func NewChain(middlewares ...func(http.Handler) http.Handler) Chain {
	return Chain{
		middlewares: middlewares,
	}
}

// WithLogger adds a logger to the chain
func (c Chain) WithLogger(logger log.Logger) Chain {
	c.logger = logger
	return c
}

// Then applies the middleware chain to a handler
func (c Chain) Then(h http.Handler) http.Handler {
	for i := len(c.middlewares) - 1; i >= 0; i-- {
		h = c.middlewares[i](h)
	}
	return h
}

func (c Chain) ThenFunc(fn http.HandlerFunc) http.Handler {
	return c.Then(http.HandlerFunc(fn))
}

func Logger(logger log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			
			// Create a response wrapper to capture the status code
			rw := newResponseWriter(w)
			
			// Log the request
			reqLogger := logger.With("method", r.Method).
				With("path", r.URL.Path).
				With("remote_addr", r.RemoteAddr).
				With("user_agent", r.UserAgent())
			
			reqLogger.Info("Request started")
			
			// Call the next handler
			next.ServeHTTP(rw, r)
			
			// Log the response
			duration := time.Since(start)
			reqLogger.With("status", rw.status).
				With("duration_ms", duration.Milliseconds()).
				Info("Request completed")
		})
	}
}

func Recovery(logger log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					logger.With("error", err).
						With("path", r.URL.Path).
						With("method", r.Method).
						Error("Panic recovered in HTTP handler")
					
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
			}()
			
			next.ServeHTTP(w, r)
		})
	}
}

type responseWriter struct {
	http.ResponseWriter
	status int
}

// newResponseWriter creates a new response writer
func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{
		ResponseWriter: w,
		status:         http.StatusOK,
	}
}

// WriteHeader captures the status code
func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}