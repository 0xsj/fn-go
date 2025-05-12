// services/user-service/internal/config/config.go
package config

import (
	"time"
)

// Config represents the application configuration
type Config struct {
	Service  ServiceConfig
	Server   ServerConfig
	Database DatabaseConfig
	NATS     NATSConfig
	CORS     CORSConfig
	Metrics  MetricsConfig
}

// ServiceConfig represents service-specific configuration
type ServiceConfig struct {
	ID      string
	Name    string
	Version string
}

// ServerConfig represents HTTP server configuration
type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// DatabaseConfig represents database configuration
type DatabaseConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	Name            string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

// NATSConfig represents NATS configuration
type NATSConfig struct {
	URL            string
	RequestTimeout time.Duration
}

// CORSConfig represents CORS configuration
type CORSConfig struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	ExposedHeaders   []string
	AllowCredentials bool
	MaxAge           time.Duration
}

// MetricsConfig represents metrics configuration
type MetricsConfig struct {
	Enabled     bool
	ServiceName string
}

// DefaultConfig returns a configuration with default values
func DefaultConfig() *Config {
	return &Config{
		Service: ServiceConfig{
			ID:      "user-service",
			Name:    "User Service",
			Version: "1.0.0",
		},
		Server: ServerConfig{
			Port:         ":8080",
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
		Database: DatabaseConfig{
			Host:            "localhost",
			Port:            "3306",
			User:            "root",
			Password:        "password",
			Name:            "user_service",
			MaxOpenConns:    25,
			MaxIdleConns:    5,
			ConnMaxLifetime: 5 * time.Minute,
		},
		NATS: NATSConfig{
			URL:            "nats://localhost:4222",
			RequestTimeout: 5 * time.Second,
		},
		CORS: CORSConfig{
			AllowedOrigins:   []string{"*"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: true,
			MaxAge:           300 * time.Second,
		},
		Metrics: MetricsConfig{
			Enabled:     true,
			ServiceName: "user-service",
		},
	}
}