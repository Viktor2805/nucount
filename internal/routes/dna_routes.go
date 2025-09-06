package routes

import (
	"golang/internal/controllers"
	"golang/internal/middleware"
	basesCounter "golang/internal/services/dnaBaseCounter"

	"github.com/gin-gonic/gin"
)

// SetupDNARoutes sets up DNA-related routes.
func SetupDNARoutes(router *gin.RouterGroup, controllers *controllers.Controllers) {
	dna := router.Group("/nucleotides")
	{
		dna.POST(
			"/count",
			middleware.MaxUploadSizeMiddleware(basesCounter.MaxUploadSizeFASTA),
			controllers.BasesCounterController.CountBases,
		)
	}
}
