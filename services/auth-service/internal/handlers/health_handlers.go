// // services/auth-service/internal/handlers/health_handlers.go
package handlers

// import (
// 	"time"

// 	"github.com/0xsj/fn-go/pkg/common/log"
// 	"github.com/0xsj/fn-go/pkg/common/nats"
// 	"github.com/0xsj/fn-go/pkg/common/nats/patterns"
// 	"github.com/0xsj/fn-go/services/auth-service/internal/service"
// )

// // HealthHandler handles health-related requests
// type HealthHandler struct {
// 	authService service.AuthService
// 	logger      log.Logger
// }

// // NewHealthHandler creates a new health handler
// func NewHealthHandler(authService service.AuthService, logger log.Logger) *HealthHandler {
// 	return &HealthHandler{
// 		authService: authService,
// 		logger:      logger.WithLayer("health-handler"),
// 	}
// }

// // RegisterHandlers registers health-related handlers with NATS
// func (h *HealthHandler) RegisterHandlers(conn *nats.Conn) {
// 	// Basic health check
// 	patterns.HandleRequest(conn, "service.auth.health", h.HealthCheck, h.logger)

// 	// Deep health check with dependencies
// 	patterns.HandleRequest(conn, "service.auth.health.deep", h.DeepHealthCheck, h.logger)

// 	// Service info
// 	patterns.HandleRequest(conn, "service.auth.info", h.ServiceInfo, h.logger)
// }

// // HealthCheck handles basic health check requests
// func (h *HealthHandler) HealthCheck(data []byte) (any, error) {
// 	handlerLogger := h.logger.With("subject", "service.auth.health")
// 	handlerLogger.Debug("Received health check request")

// 	response := map[string]any{
// 		"service": "auth-service",
// 		"status":  "ok",
// 		"time":    time.Now().Format(time.RFC3339),
// 		"version": "1.0.0",
// 		"features": []string{
// 			"authentication",
// 			"authorization",
// 			"token-management",
// 			"session-management",
// 			"password-reset",
// 			"email-verification",
// 		},
// 	}

// 	handlerLogger.Debug("Returning health check response")
// 	return response, nil
// }

// // DeepHealthCheck handles deep health check requests that test dependencies
// func (h *HealthHandler) DeepHealthCheck(data []byte) (any, error) {
// 	handlerLogger := h.logger.With("subject", "service.auth.health.deep")
// 	handlerLogger.Debug("Received deep health check request")

// 	startTime := time.Now()
// 	status := "ok"
// 	checks := make(map[string]any)

// 	// Check auth service stats (tests repository layer)
// 	ctx := patterns.GetContext()
// 	if h.authService != nil {
// 		if stats, err := h.authService.GetAuthStats(ctx); err != nil {
// 			handlerLogger.With("error", err.Error()).Warn("Auth service stats check failed")
// 			status = "degraded"
// 			checks["auth_service"] = map[string]any{
// 				"status": "error",
// 				"error":  err.Error(),
// 			}
// 		} else {
// 			checks["auth_service"] = map[string]any{
// 				"status": "ok",
// 				"stats":  stats,
// 			}
// 		}
// 	} else {
// 		status = "error"
// 		checks["auth_service"] = map[string]any{
// 			"status": "error",
// 			"error":  "auth service not initialized",
// 		}
// 	}

// 	// Check database connectivity by attempting to cleanup (read-only operation)
// 	if deletedTokens, err := h.authService.CleanupExpiredTokens(ctx); err != nil {
// 		handlerLogger.With("error", err.Error()).Warn("Database connectivity check failed")
// 		status = "degraded"
// 		checks["database"] = map[string]any{
// 			"status": "error",
// 			"error":  err.Error(),
// 		}
// 	} else {
// 		checks["database"] = map[string]any{
// 			"status":        "ok",
// 			"cleanup_test":  "passed",
// 			"expired_tokens": deletedTokens,
// 		}
// 	}

// 	// Check user service connectivity (if available)
// 	checks["user_service"] = map[string]any{
// 		"status": "unknown",
// 		"note":   "user service connectivity not directly testable from auth service",
// 	}

// 	duration := time.Since(startTime)

// 	response := map[string]any{
// 		"service":     "auth-service",
// 		"status":      status,
// 		"time":        time.Now().Format(time.RFC3339),
// 		"version":     "1.0.0",
// 		"duration_ms": duration.Milliseconds(),
// 		"checks":      checks,
// 	}

// 	handlerLogger.With("status", status).With("duration_ms", duration.Milliseconds()).Debug("Deep health check completed")
// 	return response, nil
// }

// // ServiceInfo handles service information requests
// func (h *HealthHandler) ServiceInfo(data []byte) (any, error) {
// 	handlerLogger := h.logger.With("subject", "service.auth.info")
// 	handlerLogger.Debug("Received service info request")

// 	ctx := patterns.GetContext()
// 	var stats *service.AuthStatsResponse
// 	if h.authService != nil {
// 		if s, err := h.authService.GetAuthStats(ctx); err != nil {
// 			handlerLogger.With("error", err.Error()).Warn("Failed to get auth stats for service info")
// 		} else {
// 			stats = s
// 		}
// 	}

// 	response := map[string]any{
// 		"service":     "auth-service",
// 		"version":     "1.0.0",
// 		"description": "Authentication and authorization service",
// 		"time":        time.Now().Format(time.RFC3339),
// 		"endpoints": []string{
// 			"auth.login",
// 			"auth.register",
// 			"auth.refresh",
// 			"auth.logout",
// 			"auth.validate",
// 			"auth.revoke",
// 			"auth.change-password",
// 			"auth.forgot-password",
// 			"auth.reset-password",
// 			"auth.verify-email",
// 			"auth.resend-verification",
// 			"auth.sessions.list",
// 			"auth.sessions.revoke",
// 			"auth.sessions.revoke-all",
// 			"auth.permissions.get",
// 			"auth.permissions.check",
// 			"auth.permissions.assign",
// 			"auth.permissions.revoke",
// 			"auth.stats",
// 			"auth.cleanup.tokens",
// 			"auth.cleanup.sessions",
// 		},
// 		"features": map[string]any{
// 			"authentication": map[string]any{
// 				"login":                true,
// 				"registration":         true,
// 				"password_requirements": "8+ characters, 1 uppercase, 1 digit",
// 				"account_lockout":      true,
// 				"session_management":   true,
// 			},
// 			"authorization": map[string]any{
// 				"role_based":       true,
// 				"permission_based": true,
// 				"token_validation": true,
// 			},
// 			"security": map[string]any{
// 				"password_hashing": "bcrypt",
// 				"token_type":       "JWT",
// 				"token_refresh":    true,
// 				"token_revocation": true,
// 				"email_verification": true,
// 				"password_reset":   true,
// 			},
// 		},
// 	}

// 	if stats != nil {
// 		response["statistics"] = stats
// 	}

// 	handlerLogger.Debug("Returning service info")
// 	return response, nil
// }