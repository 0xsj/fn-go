// services/user-service/internal/config/loader.go
package config

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/0xsj/fn-go/pkg/common/log"
)

// Load loads configuration from environment variables and applies defaults
func Load(logger log.Logger) (*Config, error) {
	logger = logger.WithLayer("config-loader")
	logger.Info("Loading configuration")
	
	cfg := DefaultConfig()
	
	// Server configuration
	if port := os.Getenv("SERVER_PORT"); port != "" {
		cfg.Server.Port = port
		logger.With("port", port).Debug("Loaded server port from environment")
	}
	
	if readTimeout := os.Getenv("SERVER_READ_TIMEOUT"); readTimeout != "" {
		if duration, err := time.ParseDuration(readTimeout); err == nil {
			cfg.Server.ReadTimeout = duration
			logger.With("read_timeout", duration).Debug("Loaded server read timeout from environment")
		} else {
			logger.With("error", err.Error()).Warn("Invalid read timeout value, using default")
		}
	}
	
	if writeTimeout := os.Getenv("SERVER_WRITE_TIMEOUT"); writeTimeout != "" {
		if duration, err := time.ParseDuration(writeTimeout); err == nil {
			cfg.Server.WriteTimeout = duration
			logger.With("write_timeout", duration).Debug("Loaded server write timeout from environment")
		} else {
			logger.With("error", err.Error()).Warn("Invalid write timeout value, using default")
		}
	}
	
	// Database configuration
	if dbHost := os.Getenv("DB_HOST"); dbHost != "" {
		cfg.Database.Host = dbHost
		logger.With("host", dbHost).Debug("Loaded database host from environment")
	}
	
	if dbPort := os.Getenv("DB_PORT"); dbPort != "" {
		cfg.Database.Port = dbPort
		logger.With("port", dbPort).Debug("Loaded database port from environment")
	}
	
	if dbUser := os.Getenv("DB_USER"); dbUser != "" {
		cfg.Database.User = dbUser
		logger.Debug("Loaded database user from environment")
	}
	
	if dbPassword := os.Getenv("DB_PASSWORD"); dbPassword != "" {
		cfg.Database.Password = dbPassword
		logger.Debug("Loaded database password from environment")
	}
	
	if dbName := os.Getenv("DB_NAME"); dbName != "" {
		cfg.Database.Name = dbName
		logger.With("name", dbName).Debug("Loaded database name from environment")
	}
	
	if dbMaxOpenConns := os.Getenv("DB_MAX_OPEN_CONNS"); dbMaxOpenConns != "" {
		if maxOpenConns, err := strconv.Atoi(dbMaxOpenConns); err == nil {
			cfg.Database.MaxOpenConns = maxOpenConns
			logger.With("max_open_conns", maxOpenConns).Debug("Loaded database max open connections from environment")
		} else {
			logger.With("error", err.Error()).Warn("Invalid max open connections value, using default")
		}
	}
	
	if dbMaxIdleConns := os.Getenv("DB_MAX_IDLE_CONNS"); dbMaxIdleConns != "" {
		if maxIdleConns, err := strconv.Atoi(dbMaxIdleConns); err == nil {
			cfg.Database.MaxIdleConns = maxIdleConns
			logger.With("max_idle_conns", maxIdleConns).Debug("Loaded database max idle connections from environment")
		} else {
			logger.With("error", err.Error()).Warn("Invalid max idle connections value, using default")
		}
	}
	
	if dbConnMaxLifetime := os.Getenv("DB_CONN_MAX_LIFETIME"); dbConnMaxLifetime != "" {
		if connMaxLifetime, err := time.ParseDuration(dbConnMaxLifetime); err == nil {
			cfg.Database.ConnMaxLifetime = connMaxLifetime
			logger.With("conn_max_lifetime", connMaxLifetime).Debug("Loaded database connection max lifetime from environment")
		} else {
			logger.With("error", err.Error()).Warn("Invalid connection max lifetime value, using default")
		}
	}
	
	// NATS configuration
	if natsURL := os.Getenv("NATS_URL"); natsURL != "" {
		cfg.NATS.URL = natsURL
		logger.With("url", natsURL).Debug("Loaded NATS URL from environment")
	}
	
	if natsRequestTimeout := os.Getenv("NATS_REQUEST_TIMEOUT"); natsRequestTimeout != "" {
		if duration, err := time.ParseDuration(natsRequestTimeout); err == nil {
			cfg.NATS.RequestTimeout = duration
			logger.With("request_timeout", duration).Debug("Loaded NATS request timeout from environment")
		} else {
			logger.With("error", err.Error()).Warn("Invalid NATS request timeout value, using default")
		}
	}
	
	// CORS configuration
	if allowedOrigins := os.Getenv("CORS_ALLOWED_ORIGINS"); allowedOrigins != "" {
		cfg.CORS.AllowedOrigins = strings.Split(allowedOrigins, ",")
		logger.With("allowed_origins", cfg.CORS.AllowedOrigins).Debug("Loaded CORS allowed origins from environment")
	}
	
	// Service-specific configuration
	if serviceID := os.Getenv("SERVICE_ID"); serviceID != "" {
		cfg.Service.ID = serviceID
		logger.With("id", serviceID).Debug("Loaded service ID from environment")
	}
	
	if serviceName := os.Getenv("SERVICE_NAME"); serviceName != "" {
		cfg.Service.Name = serviceName
		logger.With("name", serviceName).Debug("Loaded service name from environment")
	}
	
	logger.Info("Configuration loaded successfully")
	
	return cfg, nil
}