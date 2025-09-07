package routes

import (
	"golang/internal/controllers"
	"golang/internal/middleware"
	nucleotide "golang/internal/services/nucleotide"

	"github.com/gin-gonic/gin"
)

func SetupNucleotideRoutes(router *gin.RouterGroup, controllers *controllers.Controllers) {
	dna := router.Group("/nucleotides")
	{
		dna.POST(
			"/count",
			middleware.MaxUploadSizeMiddleware(nucleotide.MaxUploadSizeFASTA),
			controllers.NucleotideController.Count,
		)
	}
}
