package config

import "github.com/0xsj/fn-go/pkg/common/config"

type Config struct {
	Server	ServerConfig
	Database DatabaseConfig
}

type ServerConfig struct {
	Port int
}

type DatabaseConfig struct {
	Driver   string
    Host     string
    Port     int
    Username string
    Password string
    DBName   string
    Params   string
}

func LoadConfig(provider config.Provider) *Config {
	return &Config{
		Server: ServerConfig{
			Port: provider.GetInt("PORT"),
		},
		Database: DatabaseConfig{
			Driver: provider.Get("DB_DRIVER"),
			Host: provider.Get("DB_HOST"),
			Port: provider.GetInt("DB_PORT"),
			Username: provider.Get("DB_USERNAME"),
			Password: provider.Get("DB_PASSWORD"),
			DBName: provider.Get("DB_NAME"),
			Params: provider.Get("DB_PARAMS"),
		},
	}
}