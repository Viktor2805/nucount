// Package routes sets up the application's API routes and integrates the Swagger documentation.
package routes

import (
	"golang/internal/controllers"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupRoutes configures API routes and integrates Swagger documentation for the provided router and controllers.
// It sets up the route groups and binds the handlers for each route.
func SetupRoutes(routerGroup *gin.RouterGroup, controllers *controllers.Controllers) {
	SetupNucleotideRoutes(routerGroup, controllers)

	routerGroup.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
