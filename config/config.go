// Package config provides configuration management for WeKnora.
// It loads and validates application settings from environment variables
// and configuration files using the viper library.
package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config holds all application configuration.
type Config struct {
	App      AppConfig
	Database DatabaseConfig
	Redis    RedisConfig
	Auth     AuthConfig
	Server   ServerConfig
}

// AppConfig holds general application settings.
type AppConfig struct {
	Name    string
	Env     string // development, staging, production
	Debug   bool
	Version string
}

// DatabaseConfig holds database connection settings.
type DatabaseConfig struct {
	Driver          string
	Host            string
	Port            int
	Name            string
	User            string
	Password        string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

// RedisConfig holds Redis connection settings.
type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

// AuthConfig holds authentication and authorization settings.
type AuthConfig struct {
	JWTSecret          string
	JWTExpiry          time.Duration
	RefreshTokenExpiry time.Duration
}

// ServerConfig holds HTTP server settings.
type ServerConfig struct {
	Host         string
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// Load reads configuration from environment variables and returns a Config instance.
// Environment variables take precedence over defaults.
func Load() (*Config, error) {
	v := viper.New()

	// Set defaults
	v.SetDefault("app.name", "WeKnora")
	v.SetDefault("app.env", "development")
	v.SetDefault("app.debug", false)
	v.SetDefault("app.version", "0.1.0")

	v.SetDefault("server.host", "0.0.0.0")
	v.SetDefault("server.port", 8080)
	v.SetDefault("server.read_timeout", "15s")
	v.SetDefault("server.write_timeout", "15s")
	v.SetDefault("server.idle_timeout", "60s")

	v.SetDefault("database.driver", "postgres")
	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", 5432)
	v.SetDefault("database.sslmode", "disable")
	v.SetDefault("database.max_open_conns", 25)
	v.SetDefault("database.max_idle_conns", 5)
	v.SetDefault("database.conn_max_lifetime", "5m")

	v.SetDefault("redis.host", "localhost")
	v.SetDefault("redis.port", 6379)
	v.SetDefault("redis.db", 0)

	v.SetDefault("auth.jwt_expiry", "24h")
	v.SetDefault("auth.refresh_token_expiry", "168h") // 7 days

	// Bind environment variables
	v.SetEnvPrefix("WEKNORA")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// Attempt to read a .env file if present (non-fatal if missing)
	v.SetConfigName(".env")
	v.SetConfigType("dotenv")
	v.AddConfigPath(".")
	_ = v.ReadInConfig()

	cfg := &Config{}

	if err := v.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("config: failed to unmarshal configuration: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("config: validation failed: %w", err)
	}

	return cfg, nil
}

// validate checks that required configuration values are present.
func (c *Config) validate() error {
	if c.Auth.JWTSecret == "" {
		return fmt.Errorf("AUTH_JWT_SECRET must be set")
	}
	if c.Database.Name == "" {
		return fmt.Errorf("DATABASE_NAME must be set")
	}
	if c.Database.User == "" {
		return fmt.Errorf("DATABASE_USER must be set")
	}
	return nil
}

// DSN returns the database connection string for the configured driver.
func (c *DatabaseConfig) DSN() string {
	switch c.Driver {
	case "postgres":
		return fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			c.Host, c.Port, c.User, c.Password, c.Name, c.SSLMode,
		)
	default:
		return ""
	}
}
