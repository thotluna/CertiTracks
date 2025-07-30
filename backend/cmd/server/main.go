package main

import (
	"fmt"
	"log"
	"os"

	"certitrack/internal/config"
	"certitrack/internal/database"
	"certitrack/internal/router"

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

	route := gin.Default()
	router.SetupRouter(db, cfg, route)

	fmt.Printf("🚀 CertiTrack API server starting on port %s\n", cfg.App.Port)
	fmt.Printf("📊 Health check: http://localhost:%s/health\n", cfg.App.Port)
	fmt.Printf("🔗 API base URL: http://localhost:%s/api/v1\n", cfg.App.Port)
	fmt.Printf("🔐 Auth endpoints:\n")
	fmt.Printf("   POST /api/v1/auth/register\n")
	fmt.Printf("   POST /api/v1/auth/login\n")
	fmt.Printf("   POST /api/v1/auth/refresh\n")
	fmt.Printf("   GET  /api/v1/profile (protected)\n")

	if err := route.Run(":" + cfg.App.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
