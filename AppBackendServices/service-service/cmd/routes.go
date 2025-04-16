package cmd

import (
	"github.com/SwanHtetAungPhyo/common/http_server"
	"github.com/SwanHtetAungPhyo/service-service/internal/handler"
	"github.com/gofiber/fiber/v2"
)

// RegisterRoutes registers all routes for the service microservice
func RegisterRoutes(routeConfig *http_server.RoutesConfig, handlers *handler.Impl) {

	// User Services
	routeConfig.AddRoute(
		fiber.MethodGet,
		"/dashboard/:userId/service",
		handlers.GetUserServices,
	)

	routeConfig.AddRoute(
		fiber.MethodPost,
		"/dashboard/:userId/service",
		handlers.CreateService,
	)

	routeConfig.AddRoute(
		fiber.MethodPatch,
		"/dashboard/:userId/service/:serviceId",
		handlers.UpdateService,
	)

	routeConfig.AddRoute(
		fiber.MethodDelete,
		"/dashboard/:userId/service/:serviceId",
		handlers.DeleteService,
	)

	// Orders
	//routeConfig.AddRoute(
	//	fiber.MethodGet,
	//	"/api/dashboard/:userId/orders",
	//	handlers.GetOrderHistory,
	//)
	//
	//// Transactions
	//routeConfig.AddRoute(
	//	fiber.MethodGet,
	//	"/api/dashboard/:userId/tx",
	//	handlers.GetAllTransactions,
	//)
	//
	//routeConfig.AddRoute(
	//	fiber.MethodGet,
	//	"/api/dashboard/:userId/tx/fiat",
	//	handlers.GetFiatTransactions,
	//)
	//
	//routeConfig.AddRoute(
	//	fiber.MethodGet,
	//	"/api/dashboard/:userId/tx/crypto",
	//	handlers.GetCryptoTransactions,
	//)
	//
	//// Chat
	//routeConfig.AddRoute(
	//	fiber.MethodGet,
	//	"/api/dashboard/:userId/chat",
	//	handlers.GetChatHistory,
	//)
	//
	// Public Services API (/api/services/**)
	// These routes are public and don't require authentication
	routeConfig.AddRoute(
		fiber.MethodGet,
		"/services",
		handlers.ListAllServices,
	)

	routeConfig.AddRoute(
		fiber.MethodGet,
		"/services/:category",
		handlers.ListServicesByCategory,
	)
}
