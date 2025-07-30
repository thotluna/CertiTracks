package router

import (
	"certitrack/internal/config"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func setupEquipmentRoutes(rg *gin.RouterGroup, db *gorm.DB, cfg *config.Config) {
	equipment := rg.Group("/equipment")
	{
		equipment.GET("", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{"message": "Equipment"})
		})
		equipment.POST("", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{"message": "Create Equipment"})
		})

		equipmentRoutes := equipment.Group("/:id")
		{
			equipmentRoutes.GET("", func(ctx *gin.Context) {
				ctx.JSON(http.StatusOK, gin.H{"message": "Get Equipment"})
			})
			equipmentRoutes.PUT("", func(ctx *gin.Context) {
				ctx.JSON(http.StatusOK, gin.H{"message": "Update Equipment"})
			})
			equipmentRoutes.DELETE("", func(ctx *gin.Context) {
				ctx.JSON(http.StatusOK, gin.H{"message": "Delete Equipment"})
			})
		}
	}
}
