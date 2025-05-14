// services/notification-service/internal/config/config.go
package config

import (
	"github.com/0xsj/fn-go/pkg/common/config"
	"github.com/0xsj/fn-go/pkg/common/log"
)

type Config struct {
	Service      ServiceConfig
	Server       config.ServerConfig
	Database     config.DatabaseConfig
	NATS         config.NATSConfig
	Logging      config.LogConfig
	Notification NotificationConfig
}

type ServiceConfig struct {
	Name    string
	Version string
}

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	FromName string
	FromEmail string
	UseTLS   bool
}

type SMSConfig struct {
	Provider string
	APIKey   string
	FromNumber string
}

type PushConfig struct {
	Enabled  bool
	Provider string
	APIKey   string
}

type NotificationConfig struct {
	DefaultChannel   string
	TemplatesPath    string
	BatchSize        int
	RetryAttempts    int
	RetryDelay       string
	SMTP             SMTPConfig
	SMS              SMSConfig
	Push             PushConfig
}

func Load(logger log.Logger) (*Config, error) {
	provider := config.NewEnvProvider("NOTIFICATION_SERVICE").
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
			Name:    provider.GetDefault("NAME", "notification-service"),
			Version: provider.GetDefault("VERSION", "1.0.0"),
		},
		Server:   config.LoadServerConfigFromProvider(provider, ""),
		Database: config.LoadDatabaseConfigFromProvider(provider, ""),
		NATS:     config.LoadNATSConfigFromProvider(provider, ""),
		Logging: config.LogConfig{
			Level:  provider.GetDefault("LOG_LEVEL", "info"),
			Format: provider.GetDefault("LOG_FORMAT", "text"),
		},
		Notification: NotificationConfig{
			DefaultChannel:   provider.GetDefault("DEFAULT_CHANNEL", "email"),
			TemplatesPath:    provider.GetDefault("TEMPLATES_PATH", "./templates"),
			BatchSize:        provider.GetIntDefault("BATCH_SIZE", 50),
			RetryAttempts:    provider.GetIntDefault("RETRY_ATTEMPTS", 3),
			RetryDelay:       provider.GetDefault("RETRY_DELAY", "1m"),
			SMTP: SMTPConfig{
				Host:      provider.GetDefault("SMTP_HOST", "smtp.example.com"),
				Port:      provider.GetIntDefault("SMTP_PORT", 587),
				Username:  provider.Get("SMTP_USERNAME"),
				Password:  provider.Get("SMTP_PASSWORD"),
				FromName:  provider.GetDefault("SMTP_FROM_NAME", "Field Nexus"),
				FromEmail: provider.GetDefault("SMTP_FROM_EMAIL", "notifications@example.com"),
				UseTLS:    provider.GetBoolDefault("SMTP_USE_TLS", true),
			},
			SMS: SMSConfig{
				Provider:   provider.GetDefault("SMS_PROVIDER", "twilio"),
				APIKey:     provider.Get("SMS_API_KEY"),
				FromNumber: provider.GetDefault("SMS_FROM_NUMBER", ""),
			},
			Push: PushConfig{
				Enabled:  provider.GetBoolDefault("PUSH_ENABLED", false),
				Provider: provider.GetDefault("PUSH_PROVIDER", "firebase"),
				APIKey:   provider.Get("PUSH_API_KEY"),
			},
		},
	}

	return cfg, nil
}