// services/chat-service/internal/config/config.go
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
	Chat      ChatConfig
}

type ServiceConfig struct {
	Name    string
	Version string
}

type ChatConfig struct {
	MaxMessageSize   int
	MaxRoomMembers   int
	MessageRateLimit int
	FileUploadPath   string
	EnableWebsocket  bool
	WebsocketPort    int
}

func Load(logger log.Logger) (*Config, error) {
	provider := config.NewEnvProvider("CHAT_SERVICE")
	
	// In development mode, log missing variables but continue with defaults
	if err := provider.Validate(); err != nil {
		logger.With("error", err.Error()).
			With("missing_vars", provider.MissingVars()).
			Warn("Some environment variables are missing, using defaults")
	}
	
	cfg := &Config{
		Service: ServiceConfig{
			Name:    provider.GetDefault("NAME", "chat-service"),
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
			Database:        provider.GetDefault("DB_NAME", "chat_service"),
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
		Chat: ChatConfig{
			MaxMessageSize:   provider.GetIntDefault("MAX_MESSAGE_SIZE", 4096),
			MaxRoomMembers:   provider.GetIntDefault("MAX_ROOM_MEMBERS", 100),
			MessageRateLimit: provider.GetIntDefault("MESSAGE_RATE_LIMIT", 30),
			FileUploadPath:   provider.GetDefault("FILE_UPLOAD_PATH", "/tmp/uploads"),
			EnableWebsocket:  provider.GetBoolDefault("ENABLE_WEBSOCKET", true),
			WebsocketPort:    provider.GetIntDefault("WEBSOCKET_PORT", 8081),
		},
	}
	
	return cfg, nil
}