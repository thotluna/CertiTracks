//go:build wireinject

package di

import (
	"certitrack/internal/config"
	"certitrack/internal/database"
	"certitrack/internal/handlers"
	"certitrack/internal/repositories"
	"certitrack/internal/router"
	"certitrack/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"gorm.io/gorm"
)

var (
	// Conjunto de dependencias para autenticación
	authDeps = wire.NewSet(
		handlers.NewAuthHandler,
		services.NewAuthService,
		repositories.NewUserRepositoryImpl,
		wire.Bind(new(services.AuthService), new(*services.AuthServiceImpl)),
	)

	// Conjunto de dependencias de la base de datos
	dbDeps = wire.NewSet(
		config.Load,
		database.Connect,
	)
)

// ServerDependencies contiene todas las dependencias necesarias para el servidor
type ServerDependencies struct {
	Config *config.Config
	DB     *gorm.DB
	Router *gin.Engine
}

// InitializeServer inicializa todas las dependencias del servidor
func InitializeServer() (*ServerDependencies, error) {
	wire.Build(
		dbDeps,
		authDeps,
		provideRouter,
		wire.Struct(new(ServerDependencies), "*"),
	)

	return &ServerDependencies{}, nil
}

// provideRouter crea y configura el router con sus dependencias
func provideRouter(
	cfg *config.Config,
	db *gorm.DB,
	authHandler *handlers.AuthHandler,
) (*gin.Engine, error) {
	r := gin.Default()

	// Configura las rutas
	router.SetupRouter(db, cfg, r)

	return r, nil
}
