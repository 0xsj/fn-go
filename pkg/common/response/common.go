// pkg/common/response/common.go
package response

import (
	"errors"
	"strings"

	appErrors "github.com/0xsj/fn-go/pkg/common/errors"
	"github.com/0xsj/fn-go/pkg/common/log"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

type ErrorResponse struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

type PaginationMeta struct {
	CurrentPage  int `json:"current_page"`
	TotalPages   int `json:"total_pages"`
	PerPage      int `json:"per_page"`
	TotalRecords int `json:"total_records"`
}

var (
	ErrBadRequestResponse = ErrorResponse{
		Code:    "BAD_REQUEST",
		Message: "The request was invalid",
	}
	ErrUnauthorizedResponse = ErrorResponse{
		Code:    "UNAUTHORIZED",
		Message: "Authentication is required",
	}
	ErrForbiddenResponse = ErrorResponse{
		Code:    "FORBIDDEN",
		Message: "You don't have permission to access this resource",
	}
	ErrNotFoundResponse = ErrorResponse{
		Code:    "NOT_FOUND",
		Message: "The requested resource was not found",
	}
	ErrConflictResponse = ErrorResponse{
		Code:    "CONFLICT",
		Message: "The resource already exists",
	}
	ErrInternalServerResponse = ErrorResponse{
		Code:    "INTERNAL_SERVER_ERROR",
		Message: "An unexpected error occurred",
	}
	ErrServiceUnavailableResponse = ErrorResponse{
		Code:    "SERVICE_UNAVAILABLE",
		Message: "The service is currently unavailable",
	}
	ErrRateLimitedResponse = ErrorResponse{
		Code:    "RATE_LIMITED",
		Message: "You have exceeded the rate limit",
	}
)

type ErrorManager struct {
	logger log.Logger
}

func NewErrorManager(logger log.Logger) *ErrorManager {
	return &ErrorManager{
		logger: logger,
	}
}

func (em *ErrorManager) MapError(err error) (ErrorResponse, bool) {
	var appErr *appErrors.AppError
	if errors.As(err, &appErr) {
		return ErrorResponse{
			Code:    appErr.Code,
			Message: appErr.Message,
			Details: appErr.Fields,
		}, true
	}
	
	em.logger.With("error", err.Error()).Debug("Mapping error")
	
	errMsg := err.Error()
	
	switch {
	case isErrorType(errMsg, "invalid input", "validation", "bad request"):
		em.logger.With("error", errMsg).Warn("Bad request error")
		return ErrBadRequestResponse, true
		
	case isErrorType(errMsg, "invalid url", "url format"):
		em.logger.With("error", errMsg).Info("Invalid URL error")
		return ErrorResponse{
			Code:    "INVALID_URL",
			Message: "The provided URL is invalid or not supported",
		}, true
		
	case isErrorType(errMsg, "unauthorized", "authentication"):
		em.logger.With("error", errMsg).Warn("Unauthorized error")
		return ErrUnauthorizedResponse, true
		
	case isErrorType(errMsg, "forbidden", "permission"):
		em.logger.With("error", errMsg).Warn("Forbidden error")
		return ErrForbiddenResponse, true
		
	case isErrorType(errMsg, "not found", "missing"):
		em.logger.With("error", errMsg).Info("Not found error")
		return ErrNotFoundResponse, true
		
	case isErrorType(errMsg, "conflict", "duplicate", "already exists"):
		em.logger.With("error", errMsg).Warn("Conflict error")
		return ErrConflictResponse, true
		
	case isErrorType(errMsg, "rate limit", "too many requests"):
		em.logger.With("error", errMsg).Warn("Rate limited error")
		return ErrRateLimitedResponse, true
		
	default:
		em.logger.With("error", errMsg).Error("Unhandled error")
		return ErrInternalServerResponse, true
	}
}

func isErrorType(errMsg string, keywords ...string) bool {
	lowerMsg := strings.ToLower(errMsg)
	for _, keyword := range keywords {
		if strings.Contains(lowerMsg, strings.ToLower(keyword)) {
			return true
		}
	}
	return false
}