package router

import (
	"certitrack/internal/config"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func setupCertificationRoutes(rg *gin.RouterGroup, db *gorm.DB, cfg *config.Config) {
	certifications := rg.Group("/certifications")
	{
		certifications.GET("", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{"message": "Certifications"})
		})
		certifications.POST("", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{"message": "Create Certification"})
		})

		certRoutes := certifications.Group("/:id")
		{
			certRoutes.GET("", func(ctx *gin.Context) {
				ctx.JSON(http.StatusOK, gin.H{"message": "Get Certification"})
			})
			certRoutes.PUT("", func(ctx *gin.Context) {
				ctx.JSON(http.StatusOK, gin.H{"message": "Update Certification"})
			})
			certRoutes.DELETE("", func(ctx *gin.Context) {
				ctx.JSON(http.StatusOK, gin.H{"message": "Delete Certification"})
			})
		}
	}
}
