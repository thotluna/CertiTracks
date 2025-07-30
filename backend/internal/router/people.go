package router

import (
	"certitrack/internal/config"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func setupPeopleRoutes(rg *gin.RouterGroup, db *gorm.DB, cfg *config.Config) {
	people := rg.Group("/people")
	{
		people.GET("", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{"message": "People"})
		})
		people.POST("", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{"message": "Create Person"})
		})

		personRoutes := people.Group("/:id")
		{
			personRoutes.GET("", func(ctx *gin.Context) {
				ctx.JSON(http.StatusOK, gin.H{"message": "Get Person"})
			})
			personRoutes.PUT("", func(ctx *gin.Context) {
				ctx.JSON(http.StatusOK, gin.H{"message": "Update Person"})
			})
			personRoutes.DELETE("", func(ctx *gin.Context) {
				ctx.JSON(http.StatusOK, gin.H{"message": "Delete Person"})
			})
		}
	}
}
