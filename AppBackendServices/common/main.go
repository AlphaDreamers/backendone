package main

import (
	"github.com/SwanHtetAungPhyo/common/http_server"
	"github.com/SwanHtetAungPhyo/common/pkg/logutil"
	"github.com/gofiber/fiber/v2"
)

type RoutesConfig struct {
	routes map[string]map[string]fiber.Handler
}

func NewRoutesConfig() *RoutesConfig {
	return &RoutesConfig{
		routes: make(map[string]map[string]fiber.Handler),
	}
}
func (a *RoutesConfig) AddRoute(method string, route string, handler fiber.Handler) {
	if _, exists := a.routes[method]; !exists {
		a.routes[method] = make(map[string]fiber.Handler)
	}
	a.routes[method][route] = handler
}
func main() {
	logutil.InitLog("common")
	log := logutil.GetLogger()
	routeConfig := NewRoutesConfig()
	routeConfig.AddRoute(fiber.MethodGet, "/", func(c *fiber.Ctx) error {
		return c.SendString("Hello World")
	})
	routeConfig.AddRoute(fiber.MethodGet, "/ping", func(c *fiber.Ctx) error {
		return c.SendString("pong")
	})
	commonService := http_server.NewApp(log, "/Users/swanhtet1aungphyo/IdeaProjects/UniBackend/AppBackendServices/cert/auth.crt", "/Users/swanhtet1aungphyo/IdeaProjects/UniBackend/AppBackendServices/cert/auth.key", "6000", false, routeConfig.routes)
	commonService.Start()
}
