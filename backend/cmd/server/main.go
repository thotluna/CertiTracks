package main

import (
	"certitrack/internal/database"
	"certitrack/internal/di"
	"certitrack/internal/router"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	deps, err := di.InitializeServer()
	if err != nil {
		log.Fatal("Failed to initialize dependencies:", err)
	}

	switch {
	case deps.Config.IsProduction():
		gin.SetMode(gin.ReleaseMode)
	case os.Getenv("GO_ENV") == "test":
		gin.SetMode(gin.TestMode)
	}

	if err := database.AutoMigrate(deps.DB); err != nil {
		log.Fatal("Failed to run database migrations:", err)
	}
	if err := database.CreateDefaultAdmin(deps.DB); err != nil {
		log.Fatal("Failed to create default admin:", err)
	}

	r := setupRouter(deps)

	srv := &http.Server{
		Addr:    ":" + deps.Config.App.Port,
		Handler: r,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if os.Getenv("GO_ENV") != "test" {
			log.Printf("🚀 CertiTrack API server starting on port %s\n", deps.Config.App.Port)
			log.Printf("📊 Health check: http://localhost:%s/health\n", deps.Config.App.Port)
			log.Printf("🔗 API base URL: http://localhost:%s/api/v1\n", deps.Config.App.Port)
			log.Println("🔐 Auth endpoints:")
			log.Println("   POST /api/v1/auth/register")
			log.Println("   POST /api/v1/auth/login")
			log.Println("   POST /api/v1/auth/refresh")
			log.Println("   GET  /api/v1/profile (protected)")
		}

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server:", err)
		}
	}()

	<-quit
	log.Println("\n🔴 Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("✅ Server exited properly")
}

func setupRouter(deps *di.ServerDependencies) *gin.Engine {
	r := gin.Default()

	authMiddleware := deps.Middleware
	authHandler := deps.AuthHandler

	routerDeps := &router.RouterDeps{
		AuthHandler: authHandler,
		Middleware:  authMiddleware,
	}

	router.SetupRouter(routerDeps, r)

	return r
}
