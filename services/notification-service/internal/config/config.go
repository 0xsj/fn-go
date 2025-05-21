// services/notification-service/internal/config/config.go
package config

import (
	"time"

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
	provider := config.NewEnvProvider("NOTIFICATION_SERVICE")
	
	// In development mode, log missing variables but continue with defaults
	if err := provider.Validate(); err != nil {
		logger.With("error", err.Error()).
			With("missing_vars", provider.MissingVars()).
			Warn("Some environment variables are missing, using defaults")
	}
	
	cfg := &Config{
		Service: ServiceConfig{
			Name:    provider.GetDefault("NAME", "notification-service"),
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
			Database:        provider.GetDefault("DB_NAME", "notification_service"),
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
		Notification: NotificationConfig{
			DefaultChannel:   provider.GetDefault("DEFAULT_CHANNEL", "email"),
			TemplatesPath:    provider.GetDefault("TEMPLATES_PATH", "./templates"),
			BatchSize:        provider.GetIntDefault("BATCH_SIZE", 50),
			RetryAttempts:    provider.GetIntDefault("RETRY_ATTEMPTS", 3),
			RetryDelay:       provider.GetDefault("RETRY_DELAY", "1m"),
			SMTP: SMTPConfig{
				Host:      provider.GetDefault("SMTP_HOST", "localhost"),
				Port:      provider.GetIntDefault("SMTP_PORT", 1025), // Default to mailhog in dev
				Username:  provider.GetDefault("SMTP_USERNAME", ""),
				Password:  provider.GetDefault("SMTP_PASSWORD", ""),
				FromName:  provider.GetDefault("SMTP_FROM_NAME", "FN-GO Dev"),
				FromEmail: provider.GetDefault("SMTP_FROM_EMAIL", "dev@example.com"),
				UseTLS:    provider.GetBoolDefault("SMTP_USE_TLS", false),
			},
			SMS: SMSConfig{
				Provider:   provider.GetDefault("SMS_PROVIDER", "twilio"),
				APIKey:     provider.GetDefault("SMS_API_KEY", "dev-key"),
				FromNumber: provider.GetDefault("SMS_FROM_NUMBER", "+15551234567"),
			},
			Push: PushConfig{
				Enabled:  provider.GetBoolDefault("PUSH_ENABLED", false),
				Provider: provider.GetDefault("PUSH_PROVIDER", "firebase"),
				APIKey:   provider.GetDefault("PUSH_API_KEY", "dev-key"),
			},
		},
	}
	
	return cfg, nil
}