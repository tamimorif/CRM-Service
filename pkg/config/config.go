package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Auth     AuthConfig
	Logging  LoggingConfig
	Redis    RedisConfig
	Metrics  MetricsConfig
}

// ServerConfig holds HTTP server configuration
type ServerConfig struct {
	Port            string
	Host            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
	Environment     string // dev, staging, production
}

// DatabaseConfig holds database connection configuration
type DatabaseConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	Name            string
	SSLMode         string
	Timezone        string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

// AuthConfig holds authentication service configuration
type AuthConfig struct {
	ServiceAddr    string
	Timeout        time.Duration
	SkipAuth       bool
	MaxRetries     int
	RetryBackoff   time.Duration
}

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	Level      string // debug, info, warn, error
	Format     string // json, pretty
	Output     string // stdout, file, both
	FilePath   string
	MaxSize    int  // megabytes
	MaxBackups int
	MaxAge     int  // days
	Compress   bool
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Enabled  bool
	Host     string
	Port     string
	Password string
	DB       int
}

// MetricsConfig holds metrics configuration
type MetricsConfig struct {
	Enabled bool
	Port    string
	Path    string
}

// Load loads configuration from environment variables and config files
func Load() (*Config, error) {
	v := viper.New()

	// Set default values
	setDefaults(v)

	// Read from .env file
	v.SetConfigName(".env")
	v.SetConfigType("env")
	v.AddConfigPath(".")
	v.AddConfigPath("./")
	
	// Allow viper to read from environment variables
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Read config file (optional)
	if err := v.ReadInConfig(); err != nil {
		// It's okay if config file doesn't exist
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	var cfg Config

	// Server configuration
	cfg.Server = ServerConfig{
		Port:            v.GetString("APP_PORT"),
		Host:            v.GetString("APP_HOST"),
		ReadTimeout:     v.GetDuration("SERVER_READ_TIMEOUT"),
		WriteTimeout:    v.GetDuration("SERVER_WRITE_TIMEOUT"),
		ShutdownTimeout: v.GetDuration("SERVER_SHUTDOWN_TIMEOUT"),
		Environment:     v.GetString("ENVIRONMENT"),
	}

	// Database configuration
	cfg.Database = DatabaseConfig{
		Host:            v.GetString("DB_HOST"),
		Port:            v.GetString("DB_PORT"),
		User:            v.GetString("DB_USER"),
		Password:        v.GetString("DB_PASSWORD"),
		Name:            v.GetString("DB_NAME"),
		SSLMode:         v.GetString("DB_SSLMODE"),
		Timezone:        v.GetString("DB_TIMEZONE"),
		MaxOpenConns:    v.GetInt("DB_MAX_OPEN_CONNS"),
		MaxIdleConns:    v.GetInt("DB_MAX_IDLE_CONNS"),
		ConnMaxLifetime: v.GetDuration("DB_CONN_MAX_LIFETIME"),
	}

	// Auth configuration
	cfg.Auth = AuthConfig{
		ServiceAddr:  v.GetString("AUTH_SERVICE_ADDR"),
		Timeout:      v.GetDuration("AUTH_TIMEOUT"),
		SkipAuth:     v.GetBool("SKIP_AUTH"),
		MaxRetries:   v.GetInt("AUTH_MAX_RETRIES"),
		RetryBackoff: v.GetDuration("AUTH_RETRY_BACKOFF"),
	}

	// Logging configuration
	cfg.Logging = LoggingConfig{
		Level:      v.GetString("LOG_LEVEL"),
		Format:     v.GetString("LOG_FORMAT"),
		Output:     v.GetString("LOG_OUTPUT"),
		FilePath:   v.GetString("LOG_FILE_PATH"),
		MaxSize:    v.GetInt("LOG_MAX_SIZE"),
		MaxBackups: v.GetInt("LOG_MAX_BACKUPS"),
		MaxAge:     v.GetInt("LOG_MAX_AGE"),
		Compress:   v.GetBool("LOG_COMPRESS"),
	}

	// Redis configuration
	cfg.Redis = RedisConfig{
		Enabled:  v.GetBool("REDIS_ENABLED"),
		Host:     v.GetString("REDIS_HOST"),
		Port:     v.GetString("REDIS_PORT"),
		Password: v.GetString("REDIS_PASSWORD"),
		DB:       v.GetInt("REDIS_DB"),
	}

	// Metrics configuration
	cfg.Metrics = MetricsConfig{
		Enabled: v.GetBool("METRICS_ENABLED"),
		Port:    v.GetString("METRICS_PORT"),
		Path:    v.GetString("METRICS_PATH"),
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return &cfg, nil
}

// setDefaults sets default configuration values
func setDefaults(v *viper.Viper) {
	// Server defaults
	v.SetDefault("APP_PORT", "8080")
	v.SetDefault("APP_HOST", "0.0.0.0")
	v.SetDefault("SERVER_READ_TIMEOUT", 10*time.Second)
	v.SetDefault("SERVER_WRITE_TIMEOUT", 10*time.Second)
	v.SetDefault("SERVER_SHUTDOWN_TIMEOUT", 30*time.Second)
	v.SetDefault("ENVIRONMENT", "development")

	// Database defaults
	v.SetDefault("DB_HOST", "localhost")
	v.SetDefault("DB_PORT", "5432")
	v.SetDefault("DB_SSLMODE", "disable")
	v.SetDefault("DB_TIMEZONE", "Asia/Dushanbe")
	v.SetDefault("DB_MAX_OPEN_CONNS", 25)
	v.SetDefault("DB_MAX_IDLE_CONNS", 5)
	v.SetDefault("DB_CONN_MAX_LIFETIME", 5*time.Minute)

	// Auth defaults
	v.SetDefault("AUTH_SERVICE_ADDR", "localhost:50051")
	v.SetDefault("AUTH_TIMEOUT", 5*time.Second)
	v.SetDefault("SKIP_AUTH", false)
	v.SetDefault("AUTH_MAX_RETRIES", 3)
	v.SetDefault("AUTH_RETRY_BACKOFF", 100*time.Millisecond)

	// Logging defaults
	v.SetDefault("LOG_LEVEL", "info")
	v.SetDefault("LOG_FORMAT", "json")
	v.SetDefault("LOG_OUTPUT", "stdout")
	v.SetDefault("LOG_FILE_PATH", "logs/app.log")
	v.SetDefault("LOG_MAX_SIZE", 100)
	v.SetDefault("LOG_MAX_BACKUPS", 3)
	v.SetDefault("LOG_MAX_AGE", 28)
	v.SetDefault("LOG_COMPRESS", true)

	// Redis defaults
	v.SetDefault("REDIS_ENABLED", false)
	v.SetDefault("REDIS_HOST", "localhost")
	v.SetDefault("REDIS_PORT", "6379")
	v.SetDefault("REDIS_DB", 0)

	// Metrics defaults
	v.SetDefault("METRICS_ENABLED", true)
	v.SetDefault("METRICS_PORT", "9090")
	v.SetDefault("METRICS_PATH", "/metrics")
}

// Validate validates the configuration
func (c *Config) Validate() error {
	// Validate required fields
	if c.Database.User == "" {
		return fmt.Errorf("DB_USER is required")
	}
	if c.Database.Password == "" {
		return fmt.Errorf("DB_PASSWORD is required")
	}
	if c.Database.Name == "" {
		return fmt.Errorf("DB_NAME is required")
	}

	// Validate port
	if c.Server.Port == "" {
		return fmt.Errorf("APP_PORT is required")
	}

	// Validate log level
	validLogLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}
	if !validLogLevels[c.Logging.Level] {
		return fmt.Errorf("invalid LOG_LEVEL: %s (must be debug, info, warn, or error)", c.Logging.Level)
	}

	// Validate environment
	validEnvs := map[string]bool{
		"development": true,
		"staging":     true,
		"production":  true,
	}
	if !validEnvs[c.Server.Environment] {
		return fmt.Errorf("invalid ENVIRONMENT: %s (must be development, staging, or production)", c.Server.Environment)
	}

	return nil
}

// IsDevelopment returns true if running in development mode
func (c *Config) IsDevelopment() bool {
	return c.Server.Environment == "development"
}

// IsProduction returns true if running in production mode
func (c *Config) IsProduction() bool {
	return c.Server.Environment == "production"
}

// GetDSN returns the database connection string
func (c *Config) GetDSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		c.Database.Host,
		c.Database.User,
		c.Database.Password,
		c.Database.Name,
		c.Database.Port,
		c.Database.SSLMode,
		c.Database.Timezone,
	)
}
