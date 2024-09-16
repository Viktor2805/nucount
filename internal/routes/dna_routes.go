package routes

import (
	"golang/internal/controllers"

	"github.com/gin-gonic/gin"
)

// SetupDNARoutes sets up DNA-related routes.
func SetupDNARoutes(router *gin.RouterGroup, controllers *controllers.Controllers) {
	dna := router.Group("/dna")
	{
		dna.POST("/analyze", controllers.DNAController.AnalyzeDNASeq)
	}
}
