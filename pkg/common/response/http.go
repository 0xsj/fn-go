// pkg/common/response/http.go
package response

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"

	"github.com/0xsj/fn-go/pkg/common/log"
)

type FormatType string

const (
	JSON FormatType = "application/json"
	XML FormatType = "application/xml"
)

type HTTPOptions struct {
	DefaultFormat FormatType
	DefaultEncoding string
	EnableCompression bool
	AllowCORS bool
	CORSOrigins []string
}

func DefaultHTTPOptions() HTTPOptions {
	return HTTPOptions{
		DefaultFormat:     JSON,
		DefaultEncoding:   "utf-8",
		EnableCompression: false,
		AllowCORS:         false,
		CORSOrigins:       []string{"*"},
	}
}

type HTTPHandler struct {
	options  HTTPOptions
	logger   log.Logger
	errorMgr *ErrorManager
}

func NewHTTP(logger log.Logger, opts ...HTTPOptions) *HTTPHandler {
	options := DefaultHTTPOptions()
	if len(opts) > 0 {
		options = opts[0]
	}
	
	return &HTTPHandler{
		options:  options,
		logger:   logger,
		errorMgr: NewErrorManager(logger),
	}
}

func (h *HTTPHandler) Write(w http.ResponseWriter, statusCode int, data interface{}, format FormatType) error {
	contentType := string(format)
	if h.options.DefaultEncoding != "" {
		contentType = fmt.Sprintf("%s; charset=%s", format, h.options.DefaultEncoding)
	}
	
	w.Header().Set("Content-Type", contentType)
	
	if h.options.AllowCORS {
		origin := "*"
		if len(h.options.CORSOrigins) > 0 {
			origin = h.options.CORSOrigins[0]
		}
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	}
	
	w.WriteHeader(statusCode)
	
	switch format {
	case JSON:
		return json.NewEncoder(w).Encode(data)
	case XML:
		return xml.NewEncoder(w).Encode(data)
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}

func (h *HTTPHandler) JSON(w http.ResponseWriter, statusCode int, data interface{}) error {
	return h.Write(w, statusCode, data, JSON)
}

func (h *HTTPHandler) XML(w http.ResponseWriter, statusCode int, data interface{}) error {
	return h.Write(w, statusCode, data, XML)
}

func (h *HTTPHandler) Success(w http.ResponseWriter, data interface{}, message string, statusCode ...int) error {
	resp := Response{
		Success: true,
		Data:    data,
	}
	
	if message != "" {
		resp.Message = message
	}
	
	code := http.StatusOK
	if len(statusCode) > 0 {
		code = statusCode[0]
	}
	
	format := h.options.DefaultFormat
	return h.Write(w, code, resp, format)
}

func (h *HTTPHandler) Created(w http.ResponseWriter, data interface{}, message string) error {
	return h.Success(w, data, message, http.StatusCreated)
}

func (h *HTTPHandler) Accepted(w http.ResponseWriter, data interface{}, message string) error {
	return h.Success(w, data, message, http.StatusAccepted)
}

func (h *HTTPHandler) WithPagination(w http.ResponseWriter, data interface{}, meta PaginationMeta) error {
	resp := Response{
		Success: true,
		Data:    data,
		Meta:    meta,
	}
	
	format := h.options.DefaultFormat
	return h.Write(w, http.StatusOK, resp, format)
}

func (h *HTTPHandler) NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

func (h *HTTPHandler) Error(w http.ResponseWriter, err ErrorResponse, details ...interface{}) error {
	if len(details) > 0 {
		err.Details = details[0]
	}
	
	statusCode := http.StatusInternalServerError
	switch err.Code {
	case "BAD_REQUEST", "VALIDATION_ERROR", "INVALID_URL":
		statusCode = http.StatusBadRequest
	case "UNAUTHORIZED":
		statusCode = http.StatusUnauthorized
	case "FORBIDDEN":
		statusCode = http.StatusForbidden
	case "NOT_FOUND":
		statusCode = http.StatusNotFound
	case "CONFLICT":
		statusCode = http.StatusConflict
	case "RATE_LIMITED":
		statusCode = http.StatusTooManyRequests
	case "SERVICE_UNAVAILABLE":
		statusCode = http.StatusServiceUnavailable
	}
	
	return h.Write(w, statusCode, err, h.options.DefaultFormat)
}

func (h *HTTPHandler) HandleError(w http.ResponseWriter, err error) error {
	errResp, ok := h.errorMgr.MapError(err)
	if !ok {
		h.logger.With("error", err.Error()).Error("Unhandled error type")
		return h.Error(w, ErrInternalServerResponse)
	}
	
	statusCode := http.StatusInternalServerError
	switch errResp.Code {
	case "BAD_REQUEST", "VALIDATION_ERROR", "INVALID_URL":
		statusCode = http.StatusBadRequest
	case "UNAUTHORIZED":
		statusCode = http.StatusUnauthorized
	case "FORBIDDEN":
		statusCode = http.StatusForbidden
	case "NOT_FOUND":
		statusCode = http.StatusNotFound
	case "CONFLICT":
		statusCode = http.StatusConflict
	case "RATE_LIMITED":
		statusCode = http.StatusTooManyRequests
	case "SERVICE_UNAVAILABLE":
		statusCode = http.StatusServiceUnavailable
	}
	
	h.logger.With("error", err.Error()).
		With("error_code", errResp.Code).
		With("status_code", statusCode).
		Error("HTTP request error")
	
	return h.Write(w, statusCode, errResp, h.options.DefaultFormat)
}

func (h *HTTPHandler) Stream(w http.ResponseWriter, data []byte, contentType string) error {
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(data)))
	w.WriteHeader(http.StatusOK)
	_, err := w.Write(data)
	return err
}

func (h *HTTPHandler) Redirect(w http.ResponseWriter, r *http.Request, url string, statusCode ...int) {
	code := http.StatusFound 
	if len(statusCode) > 0 {
		code = statusCode[0]
	}
	http.Redirect(w, r, url, code)
}

func (h *HTTPHandler) File(w http.ResponseWriter, data []byte, filename string, contentType string) error {
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	return h.Stream(w, data, contentType)
}

func (h *HTTPHandler) FormatMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accept := r.Header.Get("Accept")
		format := h.options.DefaultFormat
		
		if accept != "" {
			switch accept {
			case "application/json":
				format = JSON
			case "application/xml":
				format = XML
			}
		}
		
		if r.URL.Query().Get("format") == "xml" {
			format = XML
		} else if r.URL.Query().Get("format") == "json" {
			format = JSON
		}
		
		ctx := r.Context()
		ctx = contextWithFormat(ctx, format)
		
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

type contextKey string
const formatContextKey contextKey = "response_format"

func contextWithFormat(ctx context.Context, format FormatType) context.Context {
	return context.WithValue(ctx, formatContextKey, format)
}

func GetFormatFromContext(r *http.Request) FormatType {
	if format, ok := r.Context().Value(formatContextKey).(FormatType); ok {
		return format
	}
	return JSON
}