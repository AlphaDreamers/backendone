package http_server

import (
	"context"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// RoutesConfig holds routes and middleware mappings
type RoutesConfig struct {
	routes     map[string]map[string]fiber.Handler
	middleware map[string]map[string][]fiber.Handler
}

func NewRoutesConfig() *RoutesConfig {
	return &RoutesConfig{
		routes:     make(map[string]map[string]fiber.Handler),
		middleware: make(map[string]map[string][]fiber.Handler),
	}
}

func (rc *RoutesConfig) AddRoute(method string, route string, handler fiber.Handler, middlewares ...fiber.Handler) {
	if _, exists := rc.routes[method]; !exists {
		rc.routes[method] = make(map[string]fiber.Handler)
		rc.middleware[method] = make(map[string][]fiber.Handler)
	}
	rc.routes[method][route] = handler
	rc.middleware[method][route] = middlewares
}

func (rc *RoutesConfig) GetRoutes() map[string]map[string]fiber.Handler {
	return rc.routes
}

func (rc *RoutesConfig) GetMiddlewares() map[string]map[string][]fiber.Handler {
	return rc.middleware
}

// App struct to encapsulate Fiber, logger, and database
type App struct {
	logger      *logrus.Logger
	app         *fiber.App
	db          *gorm.DB
	port        string
	dropDB      bool
	routes      map[string]map[string]fiber.Handler
	middlewares map[string]map[string][]fiber.Handler
}

// NewApp initializes the application
func NewApp(logger *logrus.Logger, port string, db *gorm.DB, dropDB bool, routes map[string]map[string]fiber.Handler, middlewares map[string]map[string][]fiber.Handler) *App {
	return &App{
		logger:      logger,
		app:         fiber.New(),
		db:          db,
		port:        port,
		dropDB:      dropDB,
		routes:      routes,
		middlewares: middlewares,
	}
}

// LifeCycle runs the server lifecycle
func (a *App) LifeCycle() {
	go a.Start()
	go a.Recover()
}

// Start runs the Fiber server
func (a *App) Start() {
	// CORS configuration
	a.app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:8085",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
		ExposeHeaders:    "Content-Type, Authorization",
	}))

	// Register routes with middleware
	a.Routes()

	a.logger.Infof("Starting service on port %s", a.port)
	if err := a.app.Listen(":" + a.port); err != nil {
		a.logger.Panic("Failed to start server: ", err)
	}
}

// Routes registers routes with middleware
func (a *App) Routes() {
	// Group routes by their base path to avoid duplicate middleware
	routeGroups := make(map[string]fiber.Router)

	for method, paths := range a.routes {
		for path, handler := range paths {
			// Split the path to get the base path (first segment)
			pathParts := strings.Split(strings.Trim(path, "/"), "/")
			if len(pathParts) == 0 {
				a.logger.Warnf("Invalid path: %s", path)
				continue
			}

			// Get or create the group for this base path
			basePath := "/" + pathParts[0]
			group, exists := routeGroups[basePath]
			if !exists {
				group = a.app.Group(basePath)
				routeGroups[basePath] = group
			}

			// Get the remaining path segments
			var remainingPath string
			if len(pathParts) > 1 {
				remainingPath = "/" + strings.Join(pathParts[1:], "/")
			}

			// Apply middleware only if it exists for this specific route
			if mws, exists := a.middlewares[method][path]; exists && len(mws) > 0 {
				// Create a sub-group for this specific route with its middleware
				routeGroup := group.Group(remainingPath)
				for _, mw := range mws {
					routeGroup.Use(mw)
				}

				// Register the route with the sub-group
				switch method {
				case fiber.MethodGet:
					routeGroup.Get("", handler)
				case fiber.MethodPost:
					routeGroup.Post("", handler)
				case fiber.MethodPut:
					routeGroup.Put("", handler)
				case fiber.MethodDelete:
					routeGroup.Delete("", handler)
				case fiber.MethodPatch:
					routeGroup.Patch("", handler)
				case fiber.MethodOptions:
					routeGroup.Options("", handler)
				default:
					a.logger.Warnf("Unsupported method: %s for path: %s", method, path)
				}
			} else {
				// Register the route directly with the base group if no middleware
				switch method {
				case fiber.MethodGet:
					group.Get(remainingPath, handler)
				case fiber.MethodPost:
					group.Post(remainingPath, handler)
				case fiber.MethodPut:
					group.Put(remainingPath, handler)
				case fiber.MethodDelete:
					group.Delete(remainingPath, handler)
				case fiber.MethodPatch:
					group.Patch(remainingPath, handler)
				case fiber.MethodOptions:
					group.Options(remainingPath, handler)
				default:
					a.logger.Warnf("Unsupported method: %s for path: %s", method, path)
				}
			}
		}
	}
}

// Recover handles panics and prevents the app from crashing
func (a *App) Recover() {
	defer func() {
		if err := recover(); err != nil {
			a.logger.Error("Recovered from panic: ", err)
		}
	}()
}

// Stop gracefully shuts down the Fiber app
func (a *App) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := a.app.ShutdownWithContext(ctx); err != nil {
		a.logger.Fatal("Failed to shutdown cleanly:", err)
	}
	a.logger.Info("Server shut down gracefully")
}
