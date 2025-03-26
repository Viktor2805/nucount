package server

import (
	"context"
	"fmt"
	"golang/internal/controllers"
	"golang/internal/routes"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	sentrygin "github.com/getsentry/sentry-go/gin"

	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
)

type Server struct {
	router            *gin.Engine
	port              string
	controllers       *controllers.Controllers
	healthCheckActive bool
	httpServer        *http.Server
}

type Option func(*Server)

func WithPort(port string) Option {
	return func(s *Server) {
		s.port = port
	}
}

func WithControllers(controllers *controllers.Controllers) Option {
	return func(s *Server) {
		s.controllers = controllers
	}
}

func NewServer(options ...Option) *Server {
	router := gin.Default()
	routerGroup := router.Group("/api/v1")

	// Default configuration
	server := &Server{
		router: router,
	}

	for _, option := range options {
		option(server)
	}

	// Setup Sentry.
	sentryDsn := os.Getenv("SENTRY_DSN")
	if err := sentry.Init(sentry.ClientOptions{
		Dsn:              sentryDsn,
		EnableTracing:    true,
		TracesSampleRate: 1.0,
	}); err != nil {
		log.Fatal("Init sentry failed")
	}

	router.Use(sentrygin.New(sentrygin.Options{}))

	routes.SetupRoutes(routerGroup, server.controllers)

	server.httpServer = &http.Server{
		Addr:         fmt.Sprintf(":%s", server.port),
		Handler:      server.router,
		ReadTimeout:  10 * time.Minute, // Adjust as needed
		WriteTimeout: 10 * time.Minute,
	}

	return server
}

func (s *Server) StartServer() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	// Start server in a separate goroutine
	go func() {
		log.Printf("Starting server on %s", s.httpServer.Addr)
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %s\n", err)
		}
	}()

	go func() {
		log.Println("Starting pprof server on :6060")
		if err := http.ListenAndServe(":6060", nil); err != nil {
			log.Fatalf("pprof server failed: %s\n", err)
		}
	}()

	<-ctx.Done()
	log.Println("Received shutdown signal. Shutting down gracefully...")

	// Disable health checks
	s.healthCheckActive = false

	// Create a context with a timeout to allow the server to shut down gracefully
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %s\n", err)
	} else {
		log.Println("Server shutdown successfully")
	}
}

func (s *Server) healthCheckHandler(c *gin.Context) {
	if s.healthCheckActive {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	} else {
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unhealthy"})
	}
}
