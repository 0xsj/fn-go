// pkg/common/config/config.go
package config

import (
	"time"
)

// Provider defines the interface for configuration providers
type Provider interface {
	Get(key string) string
	GetDefault(key, defaultValue string) string
	GetInt(key string) int
	GetIntDefault(key string, defaultValue int) int
	GetBool(key string) bool
	GetBoolDefault(key string, defaultValue bool) bool
	GetDuration(key string) (time.Duration, error)
	GetDurationDefault(key string, defaultValue time.Duration) time.Duration
	GetSlice(key string, separator string) []string
	
	// Validation
	WithRequiredVars(vars ...string) Provider
	Validate() error
	MissingVars() []string
}

// Common configuration types that can be embedded in service-specific configs

// DatabaseConfig provides common database configuration settings
type DatabaseConfig struct {
	Host            string
	Port            int
	Username        string
	Password        string
	Database        string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
	Timeout         time.Duration
}

// ServerConfig provides common HTTP server configuration
type ServerConfig struct {
	Port         int
	Host         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// NATSConfig provides common NATS configuration
type NATSConfig struct {
	URL            string
	MaxReconnects  int
	ReconnectWait  time.Duration
	RequestTimeout time.Duration
}

// LogConfig provides common logging configuration
type LogConfig struct {
	Level      string
	Format     string
	Output     string
	TimeFormat string
}

// DefaultDatabaseConfig returns default database configuration
func DefaultDatabaseConfig() DatabaseConfig {
	return DatabaseConfig{
		Host:            "localhost",
		Port:            3306,
		Username:        "root",
		Password:        "",
		Database:        "app",
		MaxOpenConns:    25,
		MaxIdleConns:    10,
		ConnMaxLifetime: 5 * time.Minute,
		ConnMaxIdleTime: 5 * time.Minute,
		Timeout:         10 * time.Second,
	}
}

// DefaultServerConfig returns default server configuration
func DefaultServerConfig() ServerConfig {
	return ServerConfig{
		Port:         8080,
		Host:         "",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}
}

// DefaultNATSConfig returns default NATS configuration
func DefaultNATSConfig() NATSConfig {
	return NATSConfig{
		URL:            "nats://localhost:4222",
		MaxReconnects:  10,
		ReconnectWait:  1 * time.Second,
		RequestTimeout: 5 * time.Second,
	}
}

// DefaultLogConfig returns default logging configuration
func DefaultLogConfig() LogConfig {
	return LogConfig{
		Level:      "info",
		Format:     "text",
		Output:     "stdout",
		TimeFormat: "2006-01-02 15:04:05",
	}
}

// LoadDatabaseConfigFromProvider loads database configuration from a provider
func LoadDatabaseConfigFromProvider(provider Provider, prefix string) DatabaseConfig {
	if prefix != "" && !endsWithSeparator(prefix) {
		prefix = prefix + "_"
	}
	
	return DatabaseConfig{
		Host:            provider.GetDefault(prefix+"DB_HOST", "localhost"),
		Port:            provider.GetIntDefault(prefix+"DB_PORT", 3306),
		Username:        provider.Get(prefix + "DB_USER"),
		Password:        provider.Get(prefix + "DB_PASSWORD"),
		Database:        provider.Get(prefix + "DB_NAME"),
		MaxOpenConns:    provider.GetIntDefault(prefix+"DB_MAX_OPEN_CONNS", 25),
		MaxIdleConns:    provider.GetIntDefault(prefix+"DB_MAX_IDLE_CONNS", 10),
		ConnMaxLifetime: provider.GetDurationDefault(prefix+"DB_CONN_MAX_LIFETIME", 5*time.Minute),
		ConnMaxIdleTime: provider.GetDurationDefault(prefix+"DB_CONN_MAX_IDLE_TIME", 5*time.Minute),
		Timeout:         provider.GetDurationDefault(prefix+"DB_TIMEOUT", 10*time.Second),
	}
}

// LoadServerConfigFromProvider loads server configuration from a provider
func LoadServerConfigFromProvider(provider Provider, prefix string) ServerConfig {
	if prefix != "" && !endsWithSeparator(prefix) {
		prefix = prefix + "_"
	}
	
	return ServerConfig{
		Port:         provider.GetIntDefault(prefix+"PORT", 8080),
		Host:         provider.Get(prefix + "HOST"),
		ReadTimeout:  provider.GetDurationDefault(prefix+"READ_TIMEOUT", 10*time.Second),
		WriteTimeout: provider.GetDurationDefault(prefix+"WRITE_TIMEOUT", 10*time.Second),
		IdleTimeout:  provider.GetDurationDefault(prefix+"IDLE_TIMEOUT", 30*time.Second),
	}
}

// LoadNATSConfigFromProvider loads NATS configuration from a provider
func LoadNATSConfigFromProvider(provider Provider, prefix string) NATSConfig {
	if prefix != "" && !endsWithSeparator(prefix) {
		prefix = prefix + "_"
	}
	
	return NATSConfig{
		URL:            provider.GetDefault(prefix+"NATS_URL", "nats://localhost:4222"),
		MaxReconnects:  provider.GetIntDefault(prefix+"NATS_MAX_RECONNECTS", 10),
		ReconnectWait:  provider.GetDurationDefault(prefix+"NATS_RECONNECT_WAIT", 1*time.Second),
		RequestTimeout: provider.GetDurationDefault(prefix+"NATS_REQUEST_TIMEOUT", 5*time.Second),
	}
}

// Helper function to check if a string ends with a separator
func endsWithSeparator(s string) bool {
	return len(s) > 0 && (s[len(s)-1] == '_' || s[len(s)-1] == '.')
}