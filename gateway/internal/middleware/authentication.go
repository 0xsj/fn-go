// gateway/internal/middleware/authentication.go
package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/0xsj/fn-go/pkg/common/errors"
	"github.com/0xsj/fn-go/pkg/common/log"
	"github.com/0xsj/fn-go/pkg/common/nats"
	"github.com/0xsj/fn-go/pkg/common/nats/patterns"
	"github.com/0xsj/fn-go/pkg/common/response"
)

type userContextKey string

const UserKey userContextKey = "user"

// Define auth-specific error codes
const (
	ErrCodeMissingToken   = "MISSING_TOKEN"
	ErrCodeInvalidToken   = "INVALID_TOKEN"
	ErrCodeExpiredToken   = "EXPIRED_TOKEN"
	ErrCodeValidationFail = "TOKEN_VALIDATION_FAIL"
)

// Register custom error types if needed
func init() {
	// Register custom error factories if not already defined in the errors package
	if _, exists := errors.GetErrorFactory(ErrCodeMissingToken); !exists {
		errors.RegisterErrorCode(ErrCodeMissingToken, func(message string, err error) *errors.AppError {
			return errors.NewUnauthorizedError(message, err).WithField("error_code", ErrCodeMissingToken)
		})
	}
	
	if _, exists := errors.GetErrorFactory(ErrCodeInvalidToken); !exists {
		errors.RegisterErrorCode(ErrCodeInvalidToken, func(message string, err error) *errors.AppError {
			return errors.NewUnauthorizedError(message, err).WithField("error_code", ErrCodeInvalidToken)
		})
	}
	
	if _, exists := errors.GetErrorFactory(ErrCodeExpiredToken); !exists {
		errors.RegisterErrorCode(ErrCodeExpiredToken, func(message string, err error) *errors.AppError {
			return errors.NewUnauthorizedError(message, err).WithField("error_code", ErrCodeExpiredToken)
		})
	}
	
	if _, exists := errors.GetErrorFactory(ErrCodeValidationFail); !exists {
		errors.RegisterErrorCode(ErrCodeValidationFail, func(message string, err error) *errors.AppError {
			return errors.NewInternalError(message, err).WithField("error_code", ErrCodeValidationFail)
		})
	}
}

// Authentication middleware authenticates requests
func Authentication(conn *nats.Conn, respHandler *response.HTTPHandler, logger log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip authentication for specific paths
			if isPublicPath(r.URL.Path) {
				next.ServeHTTP(w, r)
				return
			}
			
			// Extract token from Authorization header
			token := extractToken(r)
			if token == "" {
				err := errors.ErrorFromCode(ErrCodeMissingToken, "Authentication required", nil)
				logger.With("path", r.URL.Path).Warn("Request missing authentication token")
				respHandler.HandleError(w, err)
				return
			}
			
			// Validate token with auth service
			var result struct {
				Success bool        `json:"success"`
				Data    interface{} `json:"data,omitempty"`
				Error   interface{} `json:"error,omitempty"`
			}
			
			authLogger := logger.With("operation", "token_validation")
			authLogger.Debug("Validating token with auth service")
			
			err := patterns.Request(conn, "auth.validate", map[string]string{"token": token}, &result, 5*time.Second, logger)
			
			if err != nil {
				validationErr := errors.ErrorFromCode(ErrCodeValidationFail, "Failed to validate authentication", err)
				authLogger.With("error", err.Error()).Error("Token validation request failed")
				respHandler.HandleError(w, validationErr)
				return
			}
			
			if !result.Success {
				// Try to determine the exact error type from the response
				var errMessage string
				if errMsg, ok := result.Error.(string); ok {
					errMessage = errMsg
				} else {
					// If it's a complex object, convert to string
					errMessage = "Invalid or expired authentication"
				}
				
				var authErr error
				if strings.Contains(strings.ToLower(errMessage), "expired") {
					authErr = errors.ErrorFromCode(ErrCodeExpiredToken, errMessage, nil)
				} else {
					authErr = errors.ErrorFromCode(ErrCodeInvalidToken, errMessage, nil)
				}
				
				authLogger.With("error", result.Error).Warn("Token validation failed")
				respHandler.HandleError(w, authErr)
				return
			}
			
			// Add user info to request context
			ctx := context.WithValue(r.Context(), UserKey, result.Data)
			
			// Call the next handler with the authenticated context
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserFromContext gets the user from the request context
func GetUserFromContext(r *http.Request) (map[string]interface{}, error) {
	user, ok := r.Context().Value(UserKey).(map[string]interface{})
	if !ok {
		return nil, errors.ErrorFromCode(
			"UNAUTHORIZED", 
			"User not found in context", 
			errors.ErrUnauthorized,
		)
	}
	return user, nil
}

// extractToken extracts the token from the Authorization header
func extractToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}
	
	// Check if it's a Bearer token
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}
	
	return parts[1]
}

// isPublicPath checks if the path is public (doesn't require authentication)
func isPublicPath(path string) bool {
	publicPaths := []string{
		"/health",
		"/metrics",
		"/auth/login",
		"/auth/register",
		"/auth/refresh",
		"/auth/forgot-password",
		"/auth/reset-password",
		"/auth/verify-email",
		"/docs",
		"/swagger",
	}
	
	for _, publicPath := range publicPaths {
		if strings.HasPrefix(path, publicPath) {
			return true
		}
	}
	
	return false
}