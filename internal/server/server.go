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
	"go.uber.org/zap"

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
	logger            *zap.Logger
}

type Option func(*Server)

func WithLogger(l *zap.Logger) Option {
	return func(s *Server) {
		s.logger = l
	}
}

func WithPort(port string) Option {
	return func(s *Server) {
		s.port = port
	}
}

func WithControllers(c *controllers.Controllers) Option {
	return func(s *Server) {
		s.controllers = c
	}
}

func NewServer(options ...Option) *Server {
	router := gin.Default()
	routerGroup := router.Group("/api/v1")

	server := &Server{
		router: router,
	}

	for _, option := range options {
		option(server)
	}

	sentryDsn := os.Getenv("SENTRY_DSN")

	err := sentry.Init(sentry.ClientOptions{
		Dsn:              sentryDsn,
		EnableTracing:    true,
		TracesSampleRate: 1.0,
	})

	if err != nil {
		server.logger.Fatal("sentry init failed", zap.Error(err))
	}

	if os.Getenv("SENTRY_DSN") != "" {
		router.Use(sentrygin.New(sentrygin.Options{}))
	}

	router.Use(sentrygin.New(sentrygin.Options{}))

	routes.SetupRoutes(routerGroup, server.controllers)

	server.httpServer = &http.Server{
		Addr:         fmt.Sprintf(":%s", server.port),
		Handler:      server.router,
		ReadTimeout:  10 * time.Minute,
		WriteTimeout: 10 * time.Minute,
	}

	return server
}

func (s *Server) StartServer() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	go func() {
		s.logger.Info("starting server", zap.String("addr", s.httpServer.Addr))
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Fatal("server failed", zap.Error(err))
		}
	}()

	go func() {
		log.Println("Starting pprof server on :6060")
		if err := http.ListenAndServe(":6060", nil); err != nil {
			s.logger.Fatal("pprof server failed", zap.Error(err))
		}
	}()

	<-ctx.Done()
	s.logger.Info("Received shutdown signal. Shutting down gracefully...")

	s.healthCheckActive = false

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		s.logger.Error("graceful shutdown failed", zap.Error(err))
	} else {
		s.logger.Info("server stopped cleanly")
	}
}

func (s *Server) healthCheckHandler(c *gin.Context) {
	if s.healthCheckActive {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	} else {
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unhealthy"})
	}
}
