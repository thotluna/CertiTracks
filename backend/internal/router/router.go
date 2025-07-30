package router

import (
	"certitrack/internal/config"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB, cfg *config.Config, router *gin.Engine) {
	setupGlobalMiddleware(router)

	v1 := router.Group("/api/v1")
	{
		setupAuthRoutes(v1, db, cfg)

		protected := v1.Group("")
		protected.Use(getAuthMiddleware(db, cfg))
		{
			setupProtectedAuthRoutes(protected, db, cfg) // auth
			setupUserRoutes(protected, db, cfg)          // users
			setupPeopleRoutes(protected, db, cfg)        // people
			setupEquipmentRoutes(protected, db, cfg)     // equipment
			setupCertificationRoutes(protected, db, cfg) // certifications
		}
	}
}
