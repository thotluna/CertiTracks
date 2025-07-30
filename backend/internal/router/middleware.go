package router

import (
	"certitrack/internal/config"
	"certitrack/internal/factories"
	"certitrack/internal/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func setupGlobalMiddleware(r *gin.Engine) {
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(corsMiddleware())

	r.GET("/health", healthCheck)
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "certitrack-api",
		"version": "1.0.0",
	})
}

func getAuthMiddleware(db *gorm.DB, cfg *config.Config) gin.HandlerFunc {
	authFactory := factories.AuthFactoryInstance(db, cfg)
	authService := authFactory.GetAuthService()

	return middleware.AuthMiddleware(authService)
}
