// services/incident-service/internal/config/config.go
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
	Incident  IncidentConfig
}

type ServiceConfig struct {
	Name    string
	Version string
}

type IncidentConfig struct {
	DefaultPriority          string
	AutoAssign               bool
	AutoAssignRoundRobin     bool
	NotifyOnStatusChange     bool
	RequireResolutionComment bool
	AttachmentStoragePath    string
	MaxAttachmentSize        int
}

func Load(logger log.Logger) (*Config, error) {
	provider := config.NewEnvProvider("INCIDENT_SERVICE").
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
			Name:    provider.GetDefault("NAME", "incident-service"),
			Version: provider.GetDefault("VERSION", "1.0.0"),
		},
		Server:   config.LoadServerConfigFromProvider(provider, ""),
		Database: config.LoadDatabaseConfigFromProvider(provider, ""),
		NATS:     config.LoadNATSConfigFromProvider(provider, ""),
		Logging: config.LogConfig{
			Level:  provider.GetDefault("LOG_LEVEL", "info"),
			Format: provider.GetDefault("LOG_FORMAT", "text"),
		},
		Incident: IncidentConfig{
			DefaultPriority:          provider.GetDefault("DEFAULT_PRIORITY", "medium"),
			AutoAssign:               provider.GetBoolDefault("AUTO_ASSIGN", false),
			AutoAssignRoundRobin:     provider.GetBoolDefault("AUTO_ASSIGN_ROUND_ROBIN", true),
			NotifyOnStatusChange:     provider.GetBoolDefault("NOTIFY_ON_STATUS_CHANGE", true),
			RequireResolutionComment: provider.GetBoolDefault("REQUIRE_RESOLUTION_COMMENT", true),
			AttachmentStoragePath:    provider.GetDefault("ATTACHMENT_STORAGE_PATH", "/attachments"),
			MaxAttachmentSize:        provider.GetIntDefault("MAX_ATTACHMENT_SIZE", 10485760), // 10MB
		},
	}

	return cfg, nil
}