// services/location-service/internal/config/config.go
package config

import (
	"time"

	"github.com/0xsj/fn-go/pkg/common/config"
	"github.com/0xsj/fn-go/pkg/common/log"
)

type Config struct {
	Service   ServiceConfig
	Server    config.ServerConfig
	Database  config.DatabaseConfig
	NATS      config.NATSConfig
	Logging   config.LogConfig
	Location  LocationConfig
}

type ServiceConfig struct {
	Name    string
	Version string
}

type LocationConfig struct {
	GeocodeProvider     string
	GeocodeAPIKey       string
	EnableGeocoding     bool
	DefaultCoordinates  string
	MaxNestingLevel     int
	EnableSpatialSearch bool
}

func Load(logger log.Logger) (*Config, error) {
	provider := config.NewEnvProvider("LOCATION_SERVICE")
	
	// In development mode, log missing variables but continue with defaults
	if err := provider.Validate(); err != nil {
		logger.With("error", err.Error()).
			With("missing_vars", provider.MissingVars()).
			Warn("Some environment variables are missing, using defaults")
	}
	
	cfg := &Config{
		Service: ServiceConfig{
			Name:    provider.GetDefault("NAME", "location-service"),
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
			Database:        provider.GetDefault("DB_NAME", "location_service"),
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
		Location: LocationConfig{
			GeocodeProvider:     provider.GetDefault("GEOCODE_PROVIDER", "google"),
			GeocodeAPIKey:       provider.GetDefault("GEOCODE_API_KEY", "dev-api-key"),
			EnableGeocoding:     provider.GetBoolDefault("ENABLE_GEOCODING", false), // Disabled in dev for simplicity
			DefaultCoordinates:  provider.GetDefault("DEFAULT_COORDINATES", "0.0,0.0"),
			MaxNestingLevel:     provider.GetIntDefault("MAX_NESTING_LEVEL", 5),
			EnableSpatialSearch: provider.GetBoolDefault("ENABLE_SPATIAL_SEARCH", true),
		},
	}
	
	return cfg, nil
}