// Package main initializes and runs the DNA analyzer application, setting up routes,
// controllers, and database connections necessary for processing and analyzing DNA sequences.
package main

import (
	_ "golang/docs"
	"golang/internal/controllers"
	"golang/internal/db"
	"golang/internal/repository"
	"golang/internal/server"
	"golang/internal/services"
	"log"
	_ "net/http/pprof"
	"os"

	// "sort"

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
	if err := godotenv.Load(); err != nil {
		log.Fatal(".env file not found, skipping loading", err)
	}

	postgres, err := db.New()

	if err != nil {
		log.Fatal(err)
	}

	defer postgres.Close()

	repositories := repository.NewRepositories(postgres.GetDB())
	services := services.NewServices(repositories)
	controllers := controllers.NewControllers(services)

	server := server.NewServer(
		server.WithPort(os.Getenv("PORT")),
		server.WithControllers(controllers),
	)

	server.StartServer()
}
