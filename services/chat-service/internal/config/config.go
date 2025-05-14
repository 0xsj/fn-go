// services/chat-service/internal/config/config.go
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
	provider := config.NewEnvProvider("CHAT_SERVICE").
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
			Name:    provider.GetDefault("NAME", "chat-service"),
			Version: provider.GetDefault("VERSION", "1.0.0"),
		},
		Server:   config.LoadServerConfigFromProvider(provider, ""),
		Database: config.LoadDatabaseConfigFromProvider(provider, ""),
		NATS:     config.LoadNATSConfigFromProvider(provider, ""),
		Logging: config.LogConfig{
			Level:  provider.GetDefault("LOG_LEVEL", "info"),
			Format: provider.GetDefault("LOG_FORMAT", "text"),
		},
		Chat: ChatConfig{
			MaxMessageSize:   provider.GetIntDefault("MAX_MESSAGE_SIZE", 4096),
			MaxRoomMembers:   provider.GetIntDefault("MAX_ROOM_MEMBERS", 100),
			MessageRateLimit: provider.GetIntDefault("MESSAGE_RATE_LIMIT", 30),
			FileUploadPath:   provider.GetDefault("FILE_UPLOAD_PATH", "/uploads"),
			EnableWebsocket:  provider.GetBoolDefault("ENABLE_WEBSOCKET", true),
			WebsocketPort:    provider.GetIntDefault("WEBSOCKET_PORT", 8081),
		},
	}

	return cfg, nil
}