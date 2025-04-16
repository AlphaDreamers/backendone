package http_server

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/sirupsen/logrus"
)

const (
	// DefaultReadTimeout is the default timeout for reading the entire request
	DefaultReadTimeout = 60 * time.Second
	// DefaultWriteTimeout is the default timeout for writing the response
	DefaultWriteTimeout = 60 * time.Second
	// DefaultIdleTimeout is the default timeout for idle connections
	DefaultIdleTimeout = 60 * time.Second
	// DefaultMaxHeaderBytes is the default maximum size of request headers
	DefaultMaxHeaderBytes = 8192
)

// Config holds the server configuration
type Config struct {
	Port               string
	ReadTimeout        time.Duration
	WriteTimeout       time.Duration
	IdleTimeout        time.Duration
	MaxHeaderBytes     int
	AllowedOrigins     []string
	AllowedMethods     []string
	AllowedHeaders     []string
	ExposedHeaders     []string
	AllowCredentials   bool
	MaxRequestBodySize int
}

// DefaultConfig returns a default server configuration
func DefaultConfig() *Config {
	return &Config{
		Port:               "8080",
		ReadTimeout:        DefaultReadTimeout,
		WriteTimeout:       DefaultWriteTimeout,
		IdleTimeout:        DefaultIdleTimeout,
		MaxHeaderBytes:     DefaultMaxHeaderBytes,
		AllowedOrigins:     []string{"*"},
		AllowedMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposedHeaders:     []string{"Content-Length"},
		AllowCredentials:   true,
		MaxRequestBodySize: 10 * 1024 * 1024, // 10MB
	}
}

// Server represents the HTTP server
type Server struct {
	app    *fiber.App
	config *Config
	logger *logrus.Logger
}

// NewServer creates a new HTTP server instance
func NewServer(config *Config, logger *logrus.Logger) *Server {
	if config == nil {
		config = DefaultConfig()
	}

	app := fiber.New(fiber.Config{
		ReadTimeout:             config.ReadTimeout,
		WriteTimeout:            config.WriteTimeout,
		IdleTimeout:             config.IdleTimeout,
		BodyLimit:               config.MaxRequestBodySize,
		ErrorHandler:            customErrorHandler,
		EnableTrustedProxyCheck: true,
		ProxyHeader:             fiber.HeaderXForwardedFor,
		AppName:                 "UniBackend API",
	})

	return &Server{
		app:    app,
		config: config,
		logger: logger,
	}
}

// SetupMiddleware configures the server middleware
func (s *Server) SetupMiddleware() {
	// Recover from panics
	s.app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))

	// CORS middleware
	s.app.Use(cors.New(cors.Config{
		AllowOrigins:     strings.Join(s.config.AllowedOrigins, ","),
		AllowMethods:     strings.Join(s.config.AllowedMethods, ","),
		AllowHeaders:     strings.Join(s.config.AllowedHeaders, ","),
		ExposeHeaders:    strings.Join(s.config.ExposedHeaders, ","),
		AllowCredentials: s.config.AllowCredentials,
		MaxAge:           300,
	}))

	// Request logging
	s.app.Use(logger.New(logger.Config{
		Format:     "${time} | ${status} | ${latency} | ${method} | ${path}\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "Local",
	}))
}

// Start starts the HTTP server
func (s *Server) Start() error {
	addr := fmt.Sprintf(":%s", s.config.Port)
	s.logger.Infof("Starting HTTP server on %s", addr)
	return s.app.Listen(addr)
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("Shutting down HTTP server...")
	return s.app.ShutdownWithContext(ctx)
}

// App returns the Fiber app instance
func (s *Server) App() *fiber.App {
	return s.app
}

// customErrorHandler handles errors in a consistent way
func customErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	return c.Status(code).JSON(fiber.Map{
		"error": fiber.Map{
			"code":    code,
			"message": err.Error(),
		},
	})
}
