// services/user-service/internal/config/config.go
package config

import (
	"github.com/0xsj/fn-go/pkg/common/config"
	"github.com/0xsj/fn-go/pkg/common/log"
)

type Config struct {
	Service  ServiceConfig
	Server   config.ServerConfig
	Database config.DatabaseConfig
	NATS     config.NATSConfig
	Logging  config.LogConfig
}

type ServiceConfig struct {
	Name    string
	Version string
}

func Load(logger log.Logger) (*Config, error) {
	provider := config.NewEnvProvider("USER_SERVICE").
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
			Name:    provider.GetDefault("NAME", "user-service"),
			Version: provider.GetDefault("VERSION", "1.0.0"),
		},
		Server:   config.LoadServerConfigFromProvider(provider, ""),
		Database: config.LoadDatabaseConfigFromProvider(provider, ""),
		NATS:     config.LoadNATSConfigFromProvider(provider, ""),
		Logging: config.LogConfig{
			Level:  provider.GetDefault("LOG_LEVEL", "info"),
			Format: provider.GetDefault("LOG_FORMAT", "text"),
		},
	}
	
	return cfg, nil
}