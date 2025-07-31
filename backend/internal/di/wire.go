//go:build wireinject

package di

import (
	"certitrack/internal/config"
	"certitrack/internal/database"
	"certitrack/internal/handlers"
	"certitrack/internal/middleware"
	"certitrack/internal/repositories"
	"certitrack/internal/services"

	"github.com/google/wire"
	"gorm.io/gorm"
)

var (
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
		database.Connect,
		repositorySet,
		serviceSet,
		handlerSet,
		middlewareSet,

		wire.Struct(new(ServerDependencies), "*"),
	)

	return &ServerDependencies{}, nil
}
