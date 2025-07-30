package main

import (
	"fmt"
	"log"
	"os"

	"certitrack/internal/database"
	"certitrack/internal/di"

	"github.com/gin-gonic/gin"
)

func main() {
	// Inicializar dependencias con Wire
	deps, err := di.InitializeServer()
	if err != nil {
		log.Fatal("Failed to initialize dependencies:", err)
	}

	// Configurar modo de Gin
	switch {
	case deps.Config.IsProduction():
		gin.SetMode(gin.ReleaseMode)
	case os.Getenv("GO_ENV") == "test":
		gin.SetMode(gin.TestMode)
	}

	// Ejecutar migraciones
	if err := database.AutoMigrate(deps.DB); err != nil {
		log.Fatal("Failed to run database migrations:", err)
	}

	// Crear administrador por defecto
	if err := database.CreateDefaultAdmin(deps.DB); err != nil {
		log.Fatal("Failed to create default admin:", err)
	}

	// Mostrar información del servidor
	if os.Getenv("GO_ENV") != "test" {
		fmt.Printf("🚀 CertiTrack API server starting on port %s\n", deps.Config.App.Port)
		fmt.Printf("📊 Health check: http://localhost:%s/health\n", deps.Config.App.Port)
		fmt.Printf("🔗 API base URL: http://localhost:%s/api/v1\n", deps.Config.App.Port)
		fmt.Printf("🔐 Auth endpoints:\n")
		fmt.Printf("   POST /api/v1/auth/register\n")
		fmt.Printf("   POST /api/v1/auth/login\n")
		fmt.Printf("   POST /api/v1/auth/refresh\n")
		fmt.Printf("   GET  /api/v1/profile (protected)\n")
	}

	// Iniciar el servidor
	if err := deps.Router.Run(":" + deps.Config.App.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
