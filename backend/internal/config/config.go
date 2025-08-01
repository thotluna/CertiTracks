// Package config provides configuration management for the application.
// It handles loading and validating environment variables from .env files
// and provides typed access to configuration values.
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	App      AppConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	SMTP     SMTPConfig
	Storage  StorageConfig
	Logger   LoggerConfig
}

type AppConfig struct {
	Env  string
	Port string
	URL  string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
	SSLMode  string
}

type RedisConfig struct {
	URL      string
	Password string
}

type JWTConfig struct {
	Secret             string
	AccessTokenExpiry  time.Duration
	RefreshTokenExpiry time.Duration
	Issuer             string
	Audience           string
}

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

type StorageConfig struct {
	Root        string
	MaxSizeMB   int
	AllowedExts []string
}

type LoggerConfig struct {
	Level         string
	EnableMetrics bool
}

func Load() (*Config, error) {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	currentDir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("error getting current directory: %v", err)
	}

	rootDir := currentDir
	backendDir := rootDir
	if filepath.Base(rootDir) == "backend" {
		rootDir = filepath.Dir(rootDir)
	} else {
		backendDir = filepath.Join(rootDir, "backend")
	}

	envPath := filepath.Join(rootDir, ".env")
	fmt.Printf("Looking for .env at: %s\n", envPath)

	if _, err := os.Stat(envPath); err == nil {
		fmt.Printf("Loading .env from: %s\n", envPath)
		if err := godotenv.Load(envPath); err != nil {
			return nil, fmt.Errorf("error loading %s: %v", envPath, err)
		}
	} else {
		fmt.Printf("No .env file found at %s, using environment variables\n", envPath)
	}

	envFilePath := filepath.Join(rootDir, fmt.Sprintf(".env.%s", env))
	if _, err := os.Stat(envFilePath); err == nil {
		if err := godotenv.Overload(envFilePath); err != nil {
			return nil, fmt.Errorf("error loading %s: %v", envFilePath, err)
		}
	}

	localEnvPath := filepath.Join(backendDir, ".env.local")
	if _, err := os.Stat(localEnvPath); err == nil {
		fmt.Printf("Loading local environment file: %s\n", localEnvPath)
		if err := godotenv.Overload(localEnvPath); err != nil {
			return nil, fmt.Errorf("error loading %s: %v", localEnvPath, err)
		}
	}

	config := &Config{
		App: AppConfig{
			Env:  GetEnv("APP_ENV", "development"),
			Port: GetEnv("PORT", "8080"),
			URL:  GetEnv("APP_URL", "http://localhost:3000"),
		},
		Database: DatabaseConfig{
			Host:     GetEnv("DB_HOST", "localhost"),
			Port:     GetEnv("DB_PORT", "5432"),
			Name:     GetEnv("DB_NAME", "certitrack_dev"),
			User:     GetEnv("DB_USER", "certitrack_user"),
			Password: GetEnv("DB_PASSWORD", "dev_password"),
			SSLMode:  GetEnv("DB_SSL_MODE", "disable"),
		},
		Redis: RedisConfig{
			URL:      GetEnv("REDIS_URL", "redis://localhost:6379"),
			Password: GetEnv("REDIS_PASSWORD", ""),
		},
		JWT: JWTConfig{
			Secret:             GetEnv("JWT_SECRET", "development-jwt-secret-key-minimum-32-characters"),
			AccessTokenExpiry:  parseDuration(GetEnv("JWT_ACCESS_TOKEN_EXPIRY", "15m")),
			RefreshTokenExpiry: parseDuration(GetEnv("JWT_REFRESH_TOKEN_EXPIRY", "168h")),
			Issuer:             GetEnv("JWT_ISSUER", "certitrack-api"),
			Audience:           GetEnv("JWT_AUDIENCE", "certitrack-client"),
		},
		SMTP: SMTPConfig{
			Host:     GetEnv("SMTP_HOST", "localhost"),
			Port:     parseInt(GetEnv("SMTP_PORT", "1025")),
			Username: GetEnv("SMTP_USERNAME", ""),
			Password: GetEnv("SMTP_PASSWORD", ""),
			From:     GetEnv("SMTP_FROM", "CertiTrack <noreply@certitrack.local>"),
		},
		Storage: StorageConfig{
			Root:        GetEnv("STORAGE_ROOT", "./storage"),
			MaxSizeMB:   parseInt(GetEnv("MAX_FILE_SIZE_MB", "10")),
			AllowedExts: []string{".pdf", ".doc", ".docx", ".jpg", ".jpeg", ".png", ".gif"},
		},
		Logger: LoggerConfig{
			Level:         GetEnv("LOG_LEVEL", "info"),
			EnableMetrics: parseBool(GetEnv("ENABLE_METRICS", "false")),
		},
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return config, nil
}

func (c *Config) Validate() error {
	if c.JWT.Secret == "" || len(c.JWT.Secret) < 32 {
		return fmt.Errorf("JWT_SECRET must be at least 32 characters long")
	}

	if c.Database.Host == "" {
		return fmt.Errorf("DB_HOST is required")
	}

	if c.Database.Name == "" {
		return fmt.Errorf("DB_NAME is required")
	}

	if c.Database.User == "" {
		return fmt.Errorf("DB_USER is required")
	}

	return nil
}

func (c *Config) GetDatabaseDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Password,
		c.Database.Name,
		c.Database.SSLMode,
	)
}

func (c *Config) IsDevelopment() bool {
	return c.App.Env == "development"
}

func (c *Config) IsProduction() bool {
	return c.App.Env == "production"
}

func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func parseInt(s string) int {
	if i, err := strconv.Atoi(s); err == nil {
		return i
	}
	return 0
}

func parseBool(s string) bool {
	if b, err := strconv.ParseBool(s); err == nil {
		return b
	}
	return false
}

func parseDuration(s string) time.Duration {
	if d, err := time.ParseDuration(s); err == nil {
		return d
	}
	return 0
}
