// services/monitoring-service/internal/config/config.go
package config

import (
	"github.com/0xsj/fn-go/pkg/common/config"
	"github.com/0xsj/fn-go/pkg/common/log"
)

type Config struct {
	Service    ServiceConfig
	Server     config.ServerConfig
	Database   config.DatabaseConfig
	NATS       config.NATSConfig
	Logging    config.LogConfig
	Monitoring MonitoringConfig
}

type ServiceConfig struct {
	Name    string
	Version string
}

type AlertingConfig struct {
	Enabled           bool
	SendEmail         bool
	SendSMS           bool
	DefaultRecipients []string
}

type MonitoringConfig struct {
	PrometheusEnabled bool
	PrometheusPort    int
	CollectionFrequency string
	HealthCheckPath    string
	ServiceTimeout     string
	RetentionPeriod    string
	Alerting           AlertingConfig
}

func Load(logger log.Logger) (*Config, error) {
	provider := config.NewEnvProvider("MONITORING_SERVICE").
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
			Name:    provider.GetDefault("NAME", "monitoring-service"),
			Version: provider.GetDefault("VERSION", "1.0.0"),
		},
		Server:   config.LoadServerConfigFromProvider(provider, ""),
		Database: config.LoadDatabaseConfigFromProvider(provider, ""),
		NATS:     config.LoadNATSConfigFromProvider(provider, ""),
		Logging: config.LogConfig{
			Level:  provider.GetDefault("LOG_LEVEL", "info"),
			Format: provider.GetDefault("LOG_FORMAT", "text"),
		},
		Monitoring: MonitoringConfig{
			PrometheusEnabled:   provider.GetBoolDefault("PROMETHEUS_ENABLED", true),
			PrometheusPort:      provider.GetIntDefault("PROMETHEUS_PORT", 9090),
			CollectionFrequency: provider.GetDefault("COLLECTION_FREQUENCY", "15s"),
			HealthCheckPath:     provider.GetDefault("HEALTH_CHECK_PATH", "/health"),
			ServiceTimeout:      provider.GetDefault("SERVICE_TIMEOUT", "5s"),
			RetentionPeriod:     provider.GetDefault("RETENTION_PERIOD", "15d"),
			Alerting: AlertingConfig{
				Enabled:           provider.GetBoolDefault("ALERTING_ENABLED", true),
				SendEmail:         provider.GetBoolDefault("ALERTING_SEND_EMAIL", true),
				SendSMS:           provider.GetBoolDefault("ALERTING_SEND_SMS", false),
				DefaultRecipients: provider.GetSlice("ALERTING_DEFAULT_RECIPIENTS", ","),
			},
		},
	}

	return cfg, nil
}