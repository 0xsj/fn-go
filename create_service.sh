#!/bin/bash

# Script to create a new service following the project's patterns
# Usage: ./create_service.sh service-name

set -e

if [ $# -ne 1 ]; then
    echo "Usage: $0 service-name"
    exit 1
fi

SERVICE_NAME=$1
SERVICE_DIR="services/${SERVICE_NAME}"
PKG_DIR="github.com/0xsj/fn-go/services/${SERVICE_NAME}"

echo "Creating new service: ${SERVICE_NAME}"

# Check if service already exists
if [ -d "$SERVICE_DIR" ]; then
    echo "Error: Service directory already exists: $SERVICE_DIR"
    exit 1
fi

# Create service directory structure
mkdir -p "${SERVICE_DIR}/cmd/server"
mkdir -p "${SERVICE_DIR}/internal/config"
mkdir -p "${SERVICE_DIR}/internal/domain"
mkdir -p "${SERVICE_DIR}/internal/dto"
mkdir -p "${SERVICE_DIR}/internal/handlers"
mkdir -p "${SERVICE_DIR}/internal/repository/mysql"
mkdir -p "${SERVICE_DIR}/internal/service"
mkdir -p "${SERVICE_DIR}/internal/validation"
mkdir -p "${SERVICE_DIR}/migrations"
mkdir -p "${SERVICE_DIR}/pkg"

# Create .air.toml for hot reloading
cat > "${SERVICE_DIR}/.air.toml" << EOF
# services/${SERVICE_NAME}/.air.toml
root = "."
tmp_dir = "tmp"

[build]
  cmd = "go build -gcflags='all=-N -l' -o ./tmp/main ./cmd/server/main.go"
  bin = "tmp/main"
  delay = 1000
  exclude_dir = ["assets", "tmp", "vendor", "testdata"]
  exclude_file = []
  exclude_regex = [".*_test.go"]
  exclude_unchanged = false
  follow_symlink = false
  full_bin = "dlv --listen=:40000 --headless=true --api-version=2 --accept-multiclient exec ./tmp/main"
  include_dir = []
  include_ext = ["go", "tpl", "tmpl", "html"]
  kill_delay = "0s"
  log = "build-errors.log"
  send_interrupt = false
  stop_on_error = true

[color]
  app = ""
  build = "yellow"
  main = "magenta"
  runner = "green"
  watcher = "cyan"

[log]
  time = true
  main_only = false

[misc]
  clean_on_exit = true

[screen]
  clear_on_rebuild = true
  keep_scroll = true
EOF

# Create Dockerfile
cat > "${SERVICE_DIR}/Dockerfile" << EOF
FROM golang:1.24-alpine AS builder
WORKDIR /app

# Define build argument for the service name
ARG SERVICE_NAME=${SERVICE_NAME}

# Copy go.work files first to leverage layer caching
COPY go.work go.work.sum ./

# Copy required packages
COPY pkg/ ./pkg/
COPY gateway/ ./gateway/
COPY services/ ./services/

# Build the service using the SERVICE_NAME arg
RUN go build -o service ./services/\${SERVICE_NAME}/cmd/server

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/service .

EXPOSE 8080
CMD ["./service"]
EOF

# Create Dockerfile.dev
cat > "${SERVICE_DIR}/Dockerfile.dev" << EOF
# services/${SERVICE_NAME}/Dockerfile.dev
FROM golang:1.24-alpine

# Install development tools
RUN go install github.com/cosmtrek/air@latest
RUN go install github.com/go-delve/delve/cmd/dlv@latest
RUN apk add --no-cache git curl

WORKDIR /app

# Copy go.work files
COPY go.work go.work.sum ./

# We'll mount the source code as a volume, but copy initially
# to ensure all dependencies are available
COPY pkg/ ./pkg/
COPY services/ ./services/
COPY gateway/ ./gateway/

# Set working directory to the service
WORKDIR /app/services/${SERVICE_NAME}

# Create directory for air
RUN mkdir -p tmp

# Configure delve debugging
ENV DELVE_LISTEN_PORT=40000
EXPOSE 8080
EXPOSE 40000

# Command to run air for hot reloading
CMD ["air", "-c", ".air.toml"]
EOF

# Create go.mod file
cat > "${SERVICE_DIR}/go.mod" << EOF
module github.com/0xsj/fn-go/services/${SERVICE_NAME}

go 1.24.3
EOF

# Create main.go
cat > "${SERVICE_DIR}/cmd/server/main.go" << EOF
// services/${SERVICE_NAME}/cmd/server/main.go
package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/0xsj/fn-go/pkg/common/log"
	"github.com/0xsj/fn-go/pkg/common/nats"
	"github.com/0xsj/fn-go/services/${SERVICE_NAME}/internal/config"
	"github.com/0xsj/fn-go/services/${SERVICE_NAME}/internal/handlers"
)

func main() {
	// Initialize logger
	logger := log.Default()
	logger = logger.WithLayer("${SERVICE_NAME}")
	logger.Info("Initializing ${SERVICE_NAME}")

	// Load configuration
	logger.Info("Loading configuration")
	cfg, err := config.Load(logger)
	if err != nil {
		logger.With("error", err.Error()).Fatal("Failed to load configuration")
	}

	// Log configuration details (omitting sensitive information)
	logger.With("service_name", cfg.Service.Name).
		With("service_version", cfg.Service.Version).
		With("db_host", cfg.Database.Host).
		With("db_name", cfg.Database.Database).
		With("nats_url", cfg.NATS.URL).
		With("log_level", cfg.Logging.Level).
		Info("Configuration loaded")

	// Initialize NATS client
	logger.Info("Connecting to NATS server")
	natsConfig := nats.Config{
		URLs:          []string{cfg.NATS.URL},
		MaxReconnect:  cfg.NATS.MaxReconnects,
		ReconnectWait: cfg.NATS.ReconnectWait,
		Timeout:       cfg.NATS.RequestTimeout,
	}

	client, err := nats.NewClient(logger, natsConfig)
	if err != nil {
		logger.With("error", err.Error()).Fatal("Failed to connect to NATS")
	}
	defer client.Close()
	logger.Info("Successfully connected to NATS server")

	// Initialize handlers
	healthHandler := handlers.NewHealthHandler(logger)
	${SERVICE_NAME}Handler := handlers.New${SERVICE_PASCAL_NAME}HandlerWithMocks(logger)

	// Register handlers
	logger.Info("Setting up request handlers")
	healthHandler.RegisterHandlers(client.Conn())
	${SERVICE_NAME}Handler.RegisterHandlers(client.Conn())
	logger.Info("Handlers registered, service is ready")

	// Wait for termination signal
	logger.Info("Waiting for termination signal")
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)
	<-signalCh

	logger.Info("Shutting down")
}
EOF

# Create config.go
cat > "${SERVICE_DIR}/internal/config/config.go" << EOF
// services/${SERVICE_NAME}/internal/config/config.go
package config

import (
	"time"

	"github.com/0xsj/fn-go/pkg/common/config"
	"github.com/0xsj/fn-go/pkg/common/log"
)

type Config struct {
	Service     ServiceConfig
	Server      config.ServerConfig
	Database    config.DatabaseConfig
	NATS        config.NATSConfig
	Logging     config.LogConfig
	${SERVICE_PASCAL_NAME}  ${SERVICE_PASCAL_NAME}Config
}

type ServiceConfig struct {
	Name    string
	Version string
}

type ${SERVICE_PASCAL_NAME}Config struct {
	// Add service-specific configuration here
}

func Load(logger log.Logger) (*Config, error) {
	provider := config.NewEnvProvider("${SERVICE_NAME_UPPER}")
	
	// In development mode, log missing variables but continue with defaults
	if err := provider.Validate(); err != nil {
		logger.With("error", err.Error()).
			With("missing_vars", provider.MissingVars()).
			Warn("Some environment variables are missing, using defaults")
	}
	
	cfg := &Config{
		Service: ServiceConfig{
			Name:    provider.GetDefault("NAME", "${SERVICE_NAME}"),
			Version: provider.GetDefault("VERSION", "1.0.0"),
		},
		Server: config.ServerConfig{
			Port:         provider.GetIntDefault("PORT", 8080),
			Host:         provider.GetDefault("HOST", ""),
			ReadTimeout:  provider.GetDurationDefault("READ_TIMEOUT", 10*time.Second),
			WriteTimeout: provider.GetDurationDefault("WRITE_TIMEOUT", 10*time.Second),
			IdleTimeout:  provider.GetDurationDefault("IDLE_TIMEOUT", 30*time.Second),
		},
		Database: config.DatabaseConfig{
			Host:            provider.GetDefault("DB_HOST", "localhost"),
			Port:            provider.GetIntDefault("DB_PORT", 3306),
			Username:        provider.GetDefault("DB_USER", "appuser"),
			Password:        provider.GetDefault("DB_PASSWORD", "apppassword"),
			Database:        provider.GetDefault("DB_NAME", "${SERVICE_NAME_SNAKE}"),
			MaxOpenConns:    provider.GetIntDefault("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    provider.GetIntDefault("DB_MAX_IDLE_CONNS", 10),
			ConnMaxLifetime: provider.GetDurationDefault("DB_CONN_MAX_LIFETIME", 5*time.Minute),
			ConnMaxIdleTime: provider.GetDurationDefault("DB_CONN_MAX_IDLE_TIME", 5*time.Minute),
			Timeout:         provider.GetDurationDefault("DB_TIMEOUT", 10*time.Second),
		},
		NATS: config.NATSConfig{
			URL:            provider.GetDefault("NATS_URL", "nats://localhost:4222"),
			MaxReconnects:  provider.GetIntDefault("NATS_MAX_RECONNECTS", 10),
			ReconnectWait:  provider.GetDurationDefault("NATS_RECONNECT_WAIT", 1*time.Second),
			RequestTimeout: provider.GetDurationDefault("NATS_REQUEST_TIMEOUT", 5*time.Second),
		},
		Logging: config.LogConfig{
			Level:      provider.GetDefault("LOG_LEVEL", "info"),
			Format:     provider.GetDefault("LOG_FORMAT", "text"),
			Output:     provider.GetDefault("LOG_OUTPUT", "stdout"),
			TimeFormat: provider.GetDefault("LOG_TIME_FORMAT", "2006-01-02 15:04:05"),
		},
		${SERVICE_PASCAL_NAME}: ${SERVICE_PASCAL_NAME}Config{
			// Initialize service-specific configuration here
		},
	}
	
	return cfg, nil
}
EOF

# Create health_handler.go
cat > "${SERVICE_DIR}/internal/handlers/health_handler.go" << EOF
// services/${SERVICE_NAME}/internal/handlers/health_handler.go
package handlers

import (
	"time"

	"github.com/0xsj/fn-go/pkg/common/log"
	"github.com/0xsj/fn-go/pkg/common/nats"
	"github.com/0xsj/fn-go/pkg/common/nats/patterns"
)

// HealthHandler handles health-related requests
type HealthHandler struct {
	logger log.Logger
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(logger log.Logger) *HealthHandler {
	return &HealthHandler{
		logger: logger.WithLayer("health-handler"),
	}
}

// RegisterHandlers registers health-related handlers with NATS
func (h *HealthHandler) RegisterHandlers(conn *nats.Conn) {
	// Health check handler
	patterns.HandleRequest(conn, "service.${SERVICE_NAME_SNAKE}.health", h.HealthCheck, h.logger)
}

// HealthCheck handles health check requests
func (h *HealthHandler) HealthCheck(data []byte) (interface{}, error) {
	handlerLogger := h.logger.With("subject", "service.${SERVICE_NAME_SNAKE}.health")
	handlerLogger.Info("Received health check request")
	
	response := map[string]interface{}{
		"service": "${SERVICE_NAME}",
		"status":  "ok",
		"time":    time.Now().Format(time.RFC3339),
		"version": "1.0.0",
	}
	
	handlerLogger.Info("Returning health check response")
	return response, nil
}
EOF

# Create service handler (template)
cat > "${SERVICE_DIR}/internal/handlers/${SERVICE_NAME}_handler.go" << EOF
// services/${SERVICE_NAME}/internal/handlers/${SERVICE_NAME}_handler.go
package handlers

import (
	"encoding/json"
	"time"

	"github.com/0xsj/fn-go/pkg/common/errors"
	"github.com/0xsj/fn-go/pkg/common/log"
	"github.com/0xsj/fn-go/pkg/common/nats"
	"github.com/0xsj/fn-go/pkg/common/nats/patterns"
)

// ${SERVICE_PASCAL_NAME}Handler handles ${SERVICE_NAME}-related requests
type ${SERVICE_PASCAL_NAME}Handler struct {
	logger log.Logger
	// ${SERVICE_NAME}Service would normally be here
}

// New${SERVICE_PASCAL_NAME}HandlerWithMocks creates a new ${SERVICE_NAME} handler using mock data
func New${SERVICE_PASCAL_NAME}HandlerWithMocks(logger log.Logger) *${SERVICE_PASCAL_NAME}Handler {
	return &${SERVICE_PASCAL_NAME}Handler{
		logger: logger.WithLayer("${SERVICE_NAME}-handler"),
	}
}

// RegisterHandlers registers ${SERVICE_NAME}-related handlers with NATS
func (h *${SERVICE_PASCAL_NAME}Handler) RegisterHandlers(conn *nats.Conn) {
	// Register your handlers here
	patterns.HandleRequest(conn, "${SERVICE_NAME}.get", h.Get${SERVICE_PASCAL_NAME}, h.logger)
	patterns.HandleRequest(conn, "${SERVICE_NAME}.list", h.List${SERVICE_PASCAL_NAME}s, h.logger)
	patterns.HandleRequest(conn, "${SERVICE_NAME}.create", h.Create${SERVICE_PASCAL_NAME}, h.logger)
}

// Get${SERVICE_PASCAL_NAME} handles requests to get a ${SERVICE_NAME} by ID
func (h *${SERVICE_PASCAL_NAME}Handler) Get${SERVICE_PASCAL_NAME}(data []byte) (interface{}, error) {
	handlerLogger := h.logger.With("subject", "${SERVICE_NAME}.get")
	handlerLogger.Info("Received ${SERVICE_NAME}.get request")
	
	var req struct {
		ID string \`json:"id"\`
	}

	if err := json.Unmarshal(data, &req); err != nil {
		handlerLogger.With("error", err.Error()).Error("Failed to unmarshal request")
		return nil, errors.NewBadRequestError("Invalid request format", err)
	}

	handlerLogger = handlerLogger.With("${SERVICE_NAME}_id", req.ID)
	handlerLogger.Info("Looking up ${SERVICE_NAME} by ID")

	if req.ID == "" {
		handlerLogger.Warn("Empty ${SERVICE_NAME} ID provided")
		return nil, errors.NewBadRequestError("ID is required", nil)
	}

	// Mock ${SERVICE_NAME} data - replace with actual implementation
	result := map[string]interface{}{
		"id":        req.ID,
		"name":      "Example ${SERVICE_PASCAL_NAME}",
		"createdAt": time.Now().Add(-24 * time.Hour),
		"updatedAt": time.Now(),
	}

	handlerLogger.Info("${SERVICE_PASCAL_NAME} found, returning response")
	return result, nil
}

// List${SERVICE_PASCAL_NAME}s handles requests to list all ${SERVICE_NAME}s
func (h *${SERVICE_PASCAL_NAME}Handler) List${SERVICE_PASCAL_NAME}s(data []byte) (interface{}, error) {
	handlerLogger := h.logger.With("subject", "${SERVICE_NAME}.list")
	handlerLogger.Info("Received ${SERVICE_NAME}.list request")
	
	// Mock ${SERVICE_NAME} list - replace with actual implementation
	results := []map[string]interface{}{
		{
			"id":        "1",
			"name":      "${SERVICE_PASCAL_NAME} One",
			"createdAt": time.Now().Add(-48 * time.Hour),
			"updatedAt": time.Now(),
		},
		{
			"id":        "2",
			"name":      "${SERVICE_PASCAL_NAME} Two",
			"createdAt": time.Now().Add(-24 * time.Hour),
			"updatedAt": time.Now(),
		},
	}

	handlerLogger.With("count", len(results)).Info("Returning ${SERVICE_NAME} list")
	return results, nil
}

// Create${SERVICE_PASCAL_NAME} handles requests to create a new ${SERVICE_NAME}
func (h *${SERVICE_PASCAL_NAME}Handler) Create${SERVICE_PASCAL_NAME}(data []byte) (interface{}, error) {
	handlerLogger := h.logger.With("subject", "${SERVICE_NAME}.create")
	handlerLogger.Info("Received ${SERVICE_NAME}.create request")
	
	var ${SERVICE_NAME}Data map[string]interface{}
	if err := json.Unmarshal(data, &${SERVICE_NAME}Data); err != nil {
		handlerLogger.With("error", err.Error()).Error("Failed to unmarshal ${SERVICE_NAME} data")
		return nil, errors.NewBadRequestError("Invalid ${SERVICE_NAME} data", err)
	}

	// Add created/updated times
	${SERVICE_NAME}Data["createdAt"] = time.Now()
	${SERVICE_NAME}Data["updatedAt"] = time.Now()

	handlerLogger.Info("${SERVICE_PASCAL_NAME} created successfully")
	return ${SERVICE_NAME}Data, nil
}
EOF

# Create repository interface
cat > "${SERVICE_DIR}/internal/repository/interfaces.go" << EOF
// services/${SERVICE_NAME}/internal/repository/interfaces.go
package repository

import (
	"context"
)

// ${SERVICE_PASCAL_NAME}Repository defines the contract for ${SERVICE_NAME} data operations
type ${SERVICE_PASCAL_NAME}Repository interface {
	// CRUD operations
	Create(ctx context.Context, ${SERVICE_NAME} map[string]interface{}) error
	GetByID(ctx context.Context, id string) (map[string]interface{}, error)
	List(ctx context.Context, filter map[string]interface{}) ([]map[string]interface{}, int, error)
	Update(ctx context.Context, ${SERVICE_NAME} map[string]interface{}) error
	Delete(ctx context.Context, id string) error
}
EOF

# Create domain service interface
cat > "${SERVICE_DIR}/internal/service/interfaces.go" << EOF
// services/${SERVICE_NAME}/internal/service/interfaces.go
package service

import (
	"context"
)

// ${SERVICE_PASCAL_NAME}Service defines the contract for ${SERVICE_NAME}-related operations
type ${SERVICE_PASCAL_NAME}Service interface {
	// ${SERVICE_PASCAL_NAME} management
	Create${SERVICE_PASCAL_NAME}(ctx context.Context, data map[string]interface{}) (map[string]interface{}, error)
	Get${SERVICE_PASCAL_NAME}(ctx context.Context, id string) (map[string]interface{}, error)
	List${SERVICE_PASCAL_NAME}s(ctx context.Context, filter map[string]interface{}) ([]map[string]interface{}, int, error)
	Update${SERVICE_PASCAL_NAME}(ctx context.Context, id string, data map[string]interface{}) (map[string]interface{}, error)
	Delete${SERVICE_PASCAL_NAME}(ctx context.Context, id string) error
}

// HealthService defines health check operations
type HealthService interface {
	Check(ctx context.Context) (map[string]interface{}, error)
	DeepCheck(ctx context.Context) (map[string]interface{}, error)
}
EOF

# Create empty migration files
touch "${SERVICE_DIR}/migrations/000001_init.up.sql"
touch "${SERVICE_DIR}/migrations/000001_init.down.sql"

# Convert service name to different formats for template substitution
SERVICE_PASCAL_NAME=$(echo $SERVICE_NAME | sed -r 's/(^|-)([a-z])/\U\2/g')
SERVICE_NAME_UPPER=$(echo $SERVICE_NAME | tr '[:lower:]' '[:upper:]' | tr '-' '_')
SERVICE_NAME_SNAKE=$(echo $SERVICE_NAME | tr '[:upper:]' '[:lower:]' | tr '-' '_')

# Perform template substitutions
find "${SERVICE_DIR}" -type f -exec sed -i -e "s/\${SERVICE_NAME}/${SERVICE_NAME}/g" {} \;
find "${SERVICE_DIR}" -type f -exec sed -i -e "s/\${SERVICE_PASCAL_NAME}/${SERVICE_PASCAL_NAME}/g" {} \;
find "${SERVICE_DIR}" -type f -exec sed -i -e "s/\${SERVICE_NAME_UPPER}/${SERVICE_NAME_UPPER}/g" {} \;
find "${SERVICE_DIR}" -type f -exec sed -i -e "s/\${SERVICE_NAME_SNAKE}/${SERVICE_NAME_SNAKE}/g" {} \;

# Update go.work (add new service module)
if [ -f "go.work" ]; then
    # Check if the service is already in go.work
    if ! grep -q "${SERVICE_DIR}" "go.work"; then
        # Add the service to go.work
        awk -v service="${SERVICE_DIR}" '/use \(/ { print; print "\t" service; next } 1' go.work > go.work.tmp
        mv go.work.tmp go.work
        echo "Added ${SERVICE_DIR} to go.work"
    fi
fi

# Update docker-compose.yml
if [ -f "docker-compose.yml" ]; then
    # Create a temporary file for docker-compose.yml
    cp docker-compose.yml docker-compose.yml.tmp

    # Add the new service to docker-compose.yml
    cat >> docker-compose.yml.tmp << EOF

  # ${SERVICE_PASCAL_NAME} Service
  ${SERVICE_NAME}:
    build:
      context: .
      dockerfile: services/${SERVICE_NAME}/Dockerfile
      args:
        SERVICE_NAME: ${SERVICE_NAME}
    container_name: fn-${SERVICE_NAME}
    environment:
      - NATS_URL=\${NATS_URL:-nats://nats:4222}
      - DB_HOST=\${DB_HOST:-mysql}
      - DB_PORT=\${DB_PORT:-3306}
      - DB_USER=\${DB_USER:-appuser}
      - DB_PASSWORD=\${DB_PASSWORD:-apppassword}
      - DB_NAME=\${${SERVICE_NAME_UPPER}_DB_NAME:-${SERVICE_NAME_SNAKE}}
      - SERVICE_PORT=\${${SERVICE_NAME_UPPER}_SERVICE_PORT:-8080}
    depends_on:
      mysql:
        condition: service_healthy
      nats:
        condition: service_healthy
    healthcheck:
      test:
        [
          "CMD",
          "wget",
          "--no-verbose",
          "--tries=1",
          "--spider",
          "http://localhost:8080/health",
        ]
      interval: 10s
      timeout: 5s
      retries: 5
      start_period: 15s
    restart: unless-stopped
    networks:
      - fn-network
EOF

    # Replace the original file
    mv docker-compose.yml.tmp docker-compose.yml
    echo "Added ${SERVICE_NAME} to docker-compose.yml"
fi

# Update docker-compose.dev.yml
if [ -f "docker-compose.dev.yml" ]; then
    # Create a temporary file for docker-compose.dev.yml
    cp docker-compose.dev.yml docker-compose.dev.yml.tmp

    # Add the new service to docker-compose.dev.yml
    cat >> docker-compose.dev.yml.tmp << EOF

  # ${SERVICE_PASCAL_NAME} Service (Dev Mode)
  ${SERVICE_NAME}:
    build:
      context: .
      dockerfile: services/${SERVICE_NAME}/Dockerfile.dev
    volumes:
      - ./services/${SERVICE_NAME}:/app/services/${SERVICE_NAME}
      - ./pkg:/app/pkg
    environment:
      - GO_ENV=development
      - DEBUG=true
    ports:
      - "900${SERVICE_COUNT}:8080" # Expose port for direct access
      - "4000${SERVICE_COUNT}:40000" # Delve debugger port
EOF

    # Replace the original file
    mv docker-compose.dev.yml.tmp docker-compose.dev.yml
    echo "Added ${SERVICE_NAME} to docker-compose.dev.yml"
fi

# Update infra/db/init/init-databases.sql
if [ -f "infra/db/init/init-databases.sql" ]; then
    # Add the new database to init-databases.sql
    echo "CREATE DATABASE IF NOT EXISTS ${SERVICE_NAME_SNAKE};" >> infra/db/init/init-databases.sql
    echo "GRANT ALL PRIVILEGES ON ${SERVICE_NAME_SNAKE}.* TO 'appuser'@'%';" >> infra/db/init/init-databases.sql
    echo "Added ${SERVICE_NAME_SNAKE} database to init-databases.sql"
fi

echo "Service ${SERVICE_NAME} created successfully!"
echo "Next steps:"
echo "1. Update the Makefile to include the new service"
echo "2. Implement the service-specific logic"
echo "3. Add the service to the gateway handlers"
echo "4. Update the .env file with service-specific configuration"