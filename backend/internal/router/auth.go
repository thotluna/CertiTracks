package router

import (
	"certitrack/internal/config"
	"certitrack/internal/factories"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func setupAuthRoutes(rg *gin.RouterGroup, db *gorm.DB, cfg *config.Config) {
	auth := rg.Group("/auth")
	{
		authFactory := factories.AuthFactoryInstance(db, cfg)
		authHandler := authFactory.GetAuthHandler()

		auth.POST("/login", authHandler.Login)
		auth.POST("/register", authHandler.Register)
		auth.POST("/refresh", authHandler.RefreshToken)
		auth.POST("/logout", authHandler.Logout)
	}
}

func setupProtectedAuthRoutes(rg *gin.RouterGroup, db *gorm.DB, cfg *config.Config) {
	authFactory := factories.AuthFactoryInstance(db, cfg)
	authHandler := authFactory.GetAuthHandler()

	rg.GET("/profile", authHandler.GetProfile)
	rg.POST("/change-password", authHandler.ChangePassword)
}
