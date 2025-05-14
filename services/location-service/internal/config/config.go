// services/location-service/internal/config/config.go
package config

import (
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
	provider := config.NewEnvProvider("LOCATION_SERVICE").
		WithRequiredVars(
			"DB_HOST",
			"DB_USER",
			"DB_PASSWORD",
			"DB_NAME",
			"NATS_URL",
		)

	if err := provider.Validate(); err != nil {
		logger.With("error", err.Error()).
			With("missing_vars", provider.MissingVars()).
			Error("Missing required configuration")
		return nil, err
	}

	cfg := &Config{
		Service: ServiceConfig{
			Name:    provider.GetDefault("NAME", "location-service"),
			Version: provider.GetDefault("VERSION", "1.0.0"),
		},
		Server:   config.LoadServerConfigFromProvider(provider, ""),
		Database: config.LoadDatabaseConfigFromProvider(provider, ""),
		NATS:     config.LoadNATSConfigFromProvider(provider, ""),
		Logging: config.LogConfig{
			Level:  provider.GetDefault("LOG_LEVEL", "info"),
			Format: provider.GetDefault("LOG_FORMAT", "text"),
		},
		Location: LocationConfig{
			GeocodeProvider:     provider.GetDefault("GEOCODE_PROVIDER", "google"),
			GeocodeAPIKey:       provider.Get("GEOCODE_API_KEY"),
			EnableGeocoding:     provider.GetBoolDefault("ENABLE_GEOCODING", true),
			DefaultCoordinates:  provider.GetDefault("DEFAULT_COORDINATES", "0.0,0.0"),
			MaxNestingLevel:     provider.GetIntDefault("MAX_NESTING_LEVEL", 5),
			EnableSpatialSearch: provider.GetBoolDefault("ENABLE_SPATIAL_SEARCH", true),
		},
	}

	return cfg, nil
}