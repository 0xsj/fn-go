package handlers

import (
	"net/http"
	"strings"

	"github.com/0xsj/fn-go/gateway/internal/proxy"
	"github.com/0xsj/fn-go/pkg/common/log"
	"github.com/0xsj/fn-go/pkg/common/nats"
	"github.com/0xsj/fn-go/pkg/common/response"
)

type BaseHandler struct {
	proxy   *proxy.NATSProxy
	conn    *nats.Conn
	resp    *response.HTTPHandler
	logger  log.Logger
	basePath string
}


// NewBaseHandler creates a new base handler
func NewBaseHandler(conn *nats.Conn, respHandler *response.HTTPHandler, logger log.Logger, basePath string) *BaseHandler {
	return &BaseHandler{
		proxy:    proxy.NewNATSProxy(conn, respHandler, logger),
		conn:     conn,
		resp:     respHandler,
		logger:   logger,
		basePath: basePath,
	}
}

// HandleRequest handles the HTTP request by proxying it to the appropriate NATS subject
func (h *BaseHandler) HandleRequest(w http.ResponseWriter, r *http.Request, subject string) {
	h.proxy.ProxyRequest(w, r, subject, nil)
}

// ExtractIDFromPath extracts the resource ID from the URL path
func (h *BaseHandler) ExtractIDFromPath(r *http.Request) string {
	// Path pattern: /basePath/id
	path := strings.TrimPrefix(r.URL.Path, "/"+h.basePath+"/")
	
	// If there's a trailing slash, remove it
	path = strings.TrimSuffix(path, "/")
	
	// If the path contains additional segments, only take the first one
	if idx := strings.Index(path, "/"); idx != -1 {
		path = path[:idx]
	}
	
	return path
}

func (h *BaseHandler) SubPath(r *http.Request, id string) (bool, string) {
	path := strings.TrimPrefix(r.URL.Path, "/"+h.basePath+"/"+id+"/")
	return path != "", path
}

// RespondWithError sends an error response
func (h *BaseHandler) RespondWithError(w http.ResponseWriter, code string, message string, statusCode int) {
	h.resp.Error(w, response.ErrorResponse{
		Code:    code,
		Message: message,
	})
}

// RespondWithMethodNotAllowed sends a method not allowed response
func (h *BaseHandler) RespondWithMethodNotAllowed(w http.ResponseWriter) {
	h.RespondWithError(w, "METHOD_NOT_ALLOWED", "Method not allowed", http.StatusMethodNotAllowed)
}
