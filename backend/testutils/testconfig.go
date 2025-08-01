package testutils

import (
	"certitrack/internal/config"
)

func GetTestConfig() *config.Config {
	return &config.Config{
		App: config.AppConfig{
			Env:  "test",
			Port: "8080",
			URL:  "http://localhost:8080",
		},
		Database: config.DatabaseConfig{
			Host:     "localhost",
			Port:     "5432",
			Name:     "certitrack_test",
			User:     "testuser",
			Password: "testpassword",
			SSLMode:  "disable",
		},
		Redis: config.RedisConfig{
			URL:      "redis://localhost:6379/0",
			Password: "dev_redis_password",
		},
		JWT: config.JWTConfig{
			Secret:             "test-secret-key",
			AccessTokenExpiry:  3600,  // 1 hora
			RefreshTokenExpiry: 86400, // 24 horas
			Issuer:             "certitrack-test",
			Audience:           "certitrack-clients",
		},
	}
}
