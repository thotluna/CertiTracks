package factories

import (
	"certitrack/internal/config"
	"certitrack/internal/handlers"
	"certitrack/internal/repositories"
	"certitrack/internal/services"

	"gorm.io/gorm"
)

type AuthFactory struct {
	db          *gorm.DB
	cfg         *config.Config
	userRepo    repositories.UserRepository
	authService services.AuthService
	authHandler *handlers.AuthHandler
}

var authFactory *AuthFactory

func AuthFactoryInstance(db *gorm.DB, cfg *config.Config) *AuthFactory {
	if authFactory == nil {
		userRepo := repositories.NewUserRepositoryImpl(db)
		authService := services.NewAuthService(cfg, userRepo)
		authHandler := handlers.NewAuthHandler(authService)
		authFactory = &AuthFactory{
			db:          db,
			cfg:         cfg,
			userRepo:    userRepo,
			authService: authService,
			authHandler: authHandler,
		}
	}
	return authFactory
}

func (f *AuthFactory) GetAuthHandler() *handlers.AuthHandler {
	return f.authHandler
}

func (f *AuthFactory) GetAuthService() services.AuthService {
	return f.authService
}
