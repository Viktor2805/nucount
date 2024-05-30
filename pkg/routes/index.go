package routes

import (
	"golang/pkg/controllers"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRoutes(router *gin.Engine, controllers *controllers.Controllers) *gin.Engine {

	v1 := router.Group("/api/v1")

	v1.POST("/transactions/upload", controllers.TransactionController.UploadCSVFile)
	v1.GET("/transactions/export-csv", controllers.TransactionController.ExportTransactionsCSV)
	v1.GET("/transactions/export-json", controllers.TransactionController.ExportTransactionsJSON)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}
