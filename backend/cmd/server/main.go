package main

import (
	"certitrack/internal/database"
	"certitrack/internal/di"
	"certitrack/internal/handlers"
	"certitrack/internal/middleware"
	"certitrack/internal/repositories"
	"certitrack/internal/router"
	"certitrack/internal/services"
	"certitrack/internal/validators"
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
		log.Fatal("Failed to migrate database:", err)
	}

	r := setupRouter(deps)

	srv := &http.Server{
		Addr:    ":" + deps.Config.App.Port,
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}

func setupRouter(deps *di.ServerDependencies) *gin.Engine {
	r := gin.Default()

	if err := validators.RegisterAll(); err != nil {
		log.Fatal("Failed to register validators:", err)
	}

	// Configurar middlewares globales
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Configurar rutas
	setupRoutes(r, deps)

	return r
}

func setupRoutes(r *gin.Engine, deps *di.ServerDependencies) {
	// Inicializar repositorios
	userRepo := repositories.NewUserRepositoryImpl(deps.DB)

	// Inicializar servicios
	authService := services.NewAuthService(deps.Config, userRepo)

	// Inicializar handlers
	authHandler := handlers.NewAuthHandler(authService)

	// Configurar middlewares
	m := middleware.NewMiddleware(authService)

	// Configurar rutas
	routerDeps := &router.RouterDeps{
		AuthHandler: authHandler,
		Middleware:  m,
	}

	router.SetupRouter(routerDeps, r)
}
