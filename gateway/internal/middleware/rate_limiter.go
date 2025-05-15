package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/0xsj/fn-go/pkg/common/response"
)

// RateLimiter implements a simple rate limiting middleware
type RateLimiter struct {
	mu        sync.Mutex
	tokens    map[string][]time.Time
	rate      int           // Requests per window
	window    time.Duration // Time window
	respHandler *response.HTTPHandler
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(rate int, window time.Duration, respHandler *response.HTTPHandler) *RateLimiter {
	return &RateLimiter{
		tokens:    make(map[string][]time.Time),
		rate:      rate,
		window:    window,
		respHandler: respHandler,
	}
}

// RateLimit returns a middleware that limits the number of requests per client
func (rl *RateLimiter) RateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get client IP or API key as the token
		token := r.RemoteAddr
		if apiKey := r.Header.Get("X-API-Key"); apiKey != "" {
			token = apiKey
		}
		
		// Check if the client has exceeded the rate limit
		if rl.isLimited(token) {
			rl.respHandler.Error(w, response.ErrorResponse{
				Code:    "RATE_LIMITED",
				Message: "Too many requests, please try again later",
			})
			return
		}
		
		next.ServeHTTP(w, r)
	})
}

// isLimited checks if the client has exceeded the rate limit
func (rl *RateLimiter) isLimited(token string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	
	now := time.Now()
	
	// Clean up old timestamps
	if timestamps, ok := rl.tokens[token]; ok {
		var validTimestamps []time.Time
		for _, ts := range timestamps {
			if now.Sub(ts) < rl.window {
				validTimestamps = append(validTimestamps, ts)
			}
		}
		rl.tokens[token] = validTimestamps
	}
	
	// Check if the client has exceeded the rate limit
	if len(rl.tokens[token]) >= rl.rate {
		return true
	}
	
	// Add the current timestamp
	rl.tokens[token] = append(rl.tokens[token], now)
	
	return false
}