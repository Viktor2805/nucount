// Package main initializes and runs the DNA analyzer application, setting up routes,
// controllers, and database connections necessary for processing and analyzing DNA sequences.
package main

import (
	_ "golang/docs"
	"golang/internal/config"
	"golang/internal/controllers"
	"golang/internal/db"
	"golang/internal/logger"
	"golang/internal/repository"
	"golang/internal/server"
	"golang/internal/services"
	_ "net/http/pprof"
	"os"

	"go.uber.org/zap"
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
	log := logger.Init()
	defer log.Sync()

	config.LoadConfig()

	dbLogger := log.With(zap.String("component", "db"))

	db, err := db.New(dbLogger)
	if err != nil {
		log.Fatal("db init failed", zap.Error(err))
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Error("db close failed", zap.Error(err))
		} else {
			log.Info("db closed")
		}
	}()

	repositories := repository.NewRepositories(db.DB)
	services := services.NewServices(repositories)
	controllers := controllers.NewControllers(services)

	server := server.NewServer(
		server.WithLogger(log),
		server.WithPort(os.Getenv("PORT")),
		server.WithControllers(controllers),
	)

	server.StartServer()
}
