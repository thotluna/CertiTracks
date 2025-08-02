// Package di provides dependency injection setup using Google Wire.
// It defines and wires together all the application's components,
// including services, repositories, and their dependencies.
package di

import (
	"certitrack/internal/cache/redis"
	"certitrack/internal/config"
	"certitrack/internal/database"
	"certitrack/internal/handlers"
	"certitrack/internal/middleware"
	"certitrack/internal/repositories"
	"certitrack/internal/services"

	"github.com/google/wire"
	"gorm.io/gorm"
)

func provideRedisConfig(cfg *config.Config) *config.RedisConfig {
	return &cfg.Redis
}

var (
	redisClientSet = wire.NewSet(
		provideRedisConfig,
		redis.NewClient,
		wire.Bind(new(repositories.RedisClient), new(*redis.Client)),
	)

	tokenRepositorySet = wire.NewSet(
		repositories.NewTokenRepository,
	)

	repositorySet = wire.NewSet(
		repositories.NewUserRepositoryImpl,
	)

	serviceSet = wire.NewSet(
		services.NewAuthService,
		wire.Bind(new(services.AuthService), new(*services.AuthServiceImpl)),
	)

	handlerSet = wire.NewSet(
		handlers.NewAuthHandler,
	)

	middlewareSet = wire.NewSet(
		middleware.NewMiddleware,
	)
)

type ServerDependencies struct {
	Config      *config.Config
	DB          *gorm.DB
	AuthHandler *handlers.AuthHandler
	Middleware  *middleware.Middleware
}

func InitializeServer() (*ServerDependencies, error) {
	wire.Build(
		config.Load,
		redisClientSet,
		database.Connect,
		tokenRepositorySet,
		repositorySet,
		serviceSet,
		handlerSet,
		middlewareSet,

		wire.Struct(new(ServerDependencies), "*"),
	)

	return &ServerDependencies{}, nil
}
