package main

import (
	_ "golang/docs"
	"golang/pkg/config"
	"golang/pkg/controllers"
	"golang/pkg/repository"
	"golang/pkg/routes"
	"golang/pkg/services"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:3000
// @BasePath  /api/v1

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {

	godotenv.Load("../.env")
	config.Connect()

	repos := repository.NewRepositories(config.DB)
	services := services.NewServices(repos)

	controllers := controllers.NewControllers(services)

	router := gin.Default()

	routes.SetupRoutes(router, controllers)

	router.Run(":3000")
}
