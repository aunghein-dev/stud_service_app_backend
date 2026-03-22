package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig
	Auth     AuthConfig
	Server   ServerConfig
	Database DatabaseConfig
}

type AppConfig struct {
	Name string
	Env  string
}

type ServerConfig struct {
	Host string
	Port int
}

type AuthConfig struct {
	Secret              string
	AccessTokenTTLHours int
}

type DatabaseConfig struct {
	URL             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime int
}

func Load() (*Config, error) {
	v := viper.New()
	v.SetConfigName(".env")
	v.SetConfigType("env")
	v.AddConfigPath(".")
	v.AddConfigPath("..")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	setDefaults(v)

	_ = v.ReadInConfig()

	cfg := &Config{
		App: AppConfig{
			Name: v.GetString("APP_NAME"),
			Env:  v.GetString("APP_ENV"),
		},
		Auth: AuthConfig{
			Secret:              v.GetString("AUTH_SECRET"),
			AccessTokenTTLHours: v.GetInt("AUTH_ACCESS_TOKEN_TTL_HOURS"),
		},
		Server: ServerConfig{
			Host: v.GetString("SERVER_HOST"),
			Port: v.GetInt("SERVER_PORT"),
		},
		Database: DatabaseConfig{
			URL:             v.GetString("DB_URL"),
			MaxOpenConns:    v.GetInt("DB_MAX_OPEN_CONNS"),
			MaxIdleConns:    v.GetInt("DB_MAX_IDLE_CONNS"),
			ConnMaxLifetime: v.GetInt("DB_CONN_MAX_LIFETIME_MINUTES"),
		},
	}

	if cfg.Database.URL == "" {
		return nil, fmt.Errorf("DB_URL is required")
	}

	return cfg, nil
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("APP_NAME", "student-service-app")
	v.SetDefault("APP_ENV", "development")
	v.SetDefault("AUTH_SECRET", "student-service-app-local-secret")
	v.SetDefault("AUTH_ACCESS_TOKEN_TTL_HOURS", 72)
	v.SetDefault("SERVER_HOST", "0.0.0.0")
	v.SetDefault("SERVER_PORT", 8080)
	v.SetDefault("DB_MAX_OPEN_CONNS", 25)
	v.SetDefault("DB_MAX_IDLE_CONNS", 25)
	v.SetDefault("DB_CONN_MAX_LIFETIME_MINUTES", 5)
}
