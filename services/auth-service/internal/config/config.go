// services/auth-service/internal/config/config.go
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
	Auth      AuthConfig
}

type ServiceConfig struct {
	Name    string
	Version string
}

type AuthConfig struct {
	JWTSecret           string
	AccessTokenExpiry   time.Duration
	RefreshTokenExpiry  time.Duration
	PasswordHashCost    int
	MaxLoginAttempts    int
	LoginLockoutPeriod  time.Duration
}

func Load(logger log.Logger) (*Config, error) {
	provider := config.NewEnvProvider("AUTH_SERVICE").
		WithRequiredVars(
			"DB_HOST",
			"DB_USER",
			"DB_PASSWORD",
			"DB_NAME",
			"NATS_URL",
			"JWT_SECRET",
		)

	if err := provider.Validate(); err != nil {
		logger.With("error", err.Error()).
			With("missing_vars", provider.MissingVars()).
			Error("Missing required configuration")
		return nil, err
	}

	cfg := &Config{
		Service: ServiceConfig{
			Name:    provider.GetDefault("NAME", "auth-service"),
			Version: provider.GetDefault("VERSION", "1.0.0"),
		},
		Server:   config.LoadServerConfigFromProvider(provider, ""),
		Database: config.LoadDatabaseConfigFromProvider(provider, ""),
		NATS:     config.LoadNATSConfigFromProvider(provider, ""),
		Logging: config.LogConfig{
			Level:  provider.GetDefault("LOG_LEVEL", "info"),
			Format: provider.GetDefault("LOG_FORMAT", "text"),
		},
		Auth: AuthConfig{
			JWTSecret:           provider.Get("JWT_SECRET"),
			AccessTokenExpiry:   provider.GetDurationDefault("ACCESS_TOKEN_EXPIRY", 15*time.Minute),
			RefreshTokenExpiry:  provider.GetDurationDefault("REFRESH_TOKEN_EXPIRY", 7*24*time.Hour),
			PasswordHashCost:    provider.GetIntDefault("PASSWORD_HASH_COST", 10),
			MaxLoginAttempts:    provider.GetIntDefault("MAX_LOGIN_ATTEMPTS", 5),
			LoginLockoutPeriod:  provider.GetDurationDefault("LOGIN_LOCKOUT_PERIOD", 15*time.Minute),
		},
	}

	return cfg, nil
}