package config

import "github.com/0xsj/fn-go/pkg/common/config"

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Auth     AuthConfig
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

type AuthConfig struct {
	Secret        string
	TokenDuration int 
}

func LoadConfig(provider config.Provider) *Config {
	return &Config{
		Server: ServerConfig{
			Port: provider.GetInt("SERVER_PORT"),
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
		Auth: AuthConfig{
			Secret: provider.Get("AUTH_SECRET"),
			TokenDuration: provider.GetInt("AUTH_TOKEN_DURATION"),
		},
	}
}