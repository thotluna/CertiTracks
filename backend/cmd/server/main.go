package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"certitrack/internal/config"
	"certitrack/internal/database"
	"certitrack/internal/handlers"
	"certitrack/internal/middleware"
	"certitrack/internal/repositories"
	"certitrack/internal/services"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	switch {
	case cfg.IsProduction():
		gin.SetMode(gin.ReleaseMode)
	case os.Getenv("GO_ENV") == "test":
		gin.SetMode(gin.TestMode)
	}

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err := database.AutoMigrate(db); err != nil {
		log.Fatal("Failed to run database migrations:", err)
	}

	if err := database.CreateDefaultAdmin(db); err != nil {
		log.Fatal("Failed to create default admin:", err)
	}

	userRepo := repositories.NewUserRepositoryImpl(db)
	authService := services.NewAuthService(cfg, userRepo)
	authHandler := handlers.NewAuthHandler(authService)
	router := gin.Default()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": "certitrack-api",
			"version": "1.0.0",
		})
	})

	v1 := router.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/register", authHandler.Register)
			auth.POST("/refresh", authHandler.RefreshToken)
			auth.POST("/logout", authHandler.Logout)
		}

		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(authService))
		{
			protected.GET("/profile", authHandler.GetProfile)
			protected.POST("/change-password", authHandler.ChangePassword)

			users := protected.Group("/users")
			users.Use(middleware.AdminMiddleware())
			{
				users.GET("", func(c *gin.Context) {
					c.JSON(http.StatusOK, gin.H{
						"message": "Users list endpoint - coming soon",
					})
				})
				users.POST("", func(c *gin.Context) {
					c.JSON(http.StatusOK, gin.H{
						"message": "Create user endpoint - coming soon",
					})
				})
			}

			// People routes
			people := protected.Group("/people")
			{
				people.GET("", func(c *gin.Context) {
					c.JSON(http.StatusOK, gin.H{
						"message": "People list endpoint - coming soon",
					})
				})
				people.POST("", func(c *gin.Context) {
					c.JSON(http.StatusOK, gin.H{
						"message": "Create person endpoint - coming soon",
					})
				})
			}

			// Equipment routes
			equipment := protected.Group("/equipment")
			{
				equipment.GET("", func(c *gin.Context) {
					c.JSON(http.StatusOK, gin.H{
						"message": "Equipment list endpoint - coming soon",
					})
				})
				equipment.POST("", func(c *gin.Context) {
					c.JSON(http.StatusOK, gin.H{
						"message": "Create equipment endpoint - coming soon",
					})
				})
			}

			certifications := protected.Group("/certifications")
			{
				certifications.GET("", func(c *gin.Context) {
					c.JSON(http.StatusOK, gin.H{
						"message": "Certifications list endpoint - coming soon",
					})
				})
				certifications.POST("", func(c *gin.Context) {
					c.JSON(http.StatusOK, gin.H{
						"message": "Create certification endpoint - coming soon",
					})
				})
			}
		}
	}

	fmt.Printf("🚀 CertiTrack API server starting on port %s\n", cfg.App.Port)
	fmt.Printf("📊 Health check: http://localhost:%s/health\n", cfg.App.Port)
	fmt.Printf("🔗 API base URL: http://localhost:%s/api/v1\n", cfg.App.Port)
	fmt.Printf("🔐 Auth endpoints:\n")
	fmt.Printf("   POST /api/v1/auth/register\n")
	fmt.Printf("   POST /api/v1/auth/login\n")
	fmt.Printf("   POST /api/v1/auth/refresh\n")
	fmt.Printf("   GET  /api/v1/profile (protected)\n")

	// Start server
	if err := router.Run(":" + cfg.App.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
