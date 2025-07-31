package router

import (
	"certitrack/internal/handlers"
	"certitrack/internal/middleware"

	"github.com/gin-gonic/gin"
)

type RouterDeps struct {
	AuthHandler *handlers.AuthHandler
	// TODO: Agregar los demás handlers cuando estén definidos
	// UserHandler      *handlers.UserHandler
	// PeopleHandler    *handlers.PeopleHandler
	// EquipmentHandler *handlers.EquipmentHandler
	// CertHandler      *handlers.CertificationHandler
	Middleware *middleware.Middleware
}

func SetupRouter(deps *RouterDeps, r *gin.Engine) {
	setupGlobalMiddleware(r)
	r.GET("/health", healthCheck)

	v1 := r.Group("/api/v1")
	{
		setupAuthRoutes(v1, deps)

		protected := v1.Group("")
		protected.Use(deps.Middleware.AuthMiddleware())
		{
			setupProtectedAuthRoutes(protected, deps)

			adminProtected := protected.Group("")
			adminProtected.Use(deps.Middleware.AdminMiddleware())
			{
				setupUserRoutes(adminProtected, deps)
			}
		}
	}
}

func setupGlobalMiddleware(router *gin.Engine) {
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	router.Use(gin.Logger())

	router.Use(gin.Recovery())
}

func healthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"status":  "ok",
		"service": "certitrack-api",
	})
}
