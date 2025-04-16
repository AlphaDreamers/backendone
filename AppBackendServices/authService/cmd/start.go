package cmd

import (
	"github.com/SwanHtetAungPhyo/auth/internal/config"
	"github.com/SwanHtetAungPhyo/auth/internal/handler"
	"github.com/SwanHtetAungPhyo/auth/internal/services"
	"github.com/SwanHtetAungPhyo/common/http_server"
	middleware2 "github.com/SwanHtetAungPhyo/common/middleware"
	"github.com/SwanHtetAungPhyo/common/pkg/logutil"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Start(port string, database *gorm.DB) {
	logutil.InitLog("auth")
	log := logutil.GetLogger()

	redisClient := config.GetRedisClient()
	service := services.NewServiceImpl()

	cfg := config.GetConfig()
	handlers := handler.NewAuthHandler(*service, redisClient, cfg)

	routeConfig := http_server.NewRoutesConfig()
	routeConfig.AddRoute(
		fiber.MethodPost,
		"/register",
		handlers.Register,
	)

	routeConfig.AddRoute(
		fiber.MethodPost,
		"/login",
		handlers.Login,
	)

	routeConfig.AddRoute(
		fiber.MethodPost,
		"/logout",
		handlers.Logout)

	routeConfig.AddRoute(
		fiber.MethodPost,
		"/verify",
		handlers.Verify,
	)

	routeConfig.AddRoute(
		fiber.MethodGet,
		"/me",
		handlers.Me,
		middleware2.JwtMiddleware(redisClient),
		middleware2.JWTBlacklistMiddleware(redisClient),
	)

	//routeConfig.AddRoute(
	//	fiber.MethodPost,
	//	"/wallet",
	//	handlers.StoreInVault,
	//	middleware2.JwtMiddleware(),
	//)

	routeConfig.AddRoute(
		fiber.MethodGet,
		"/refresh",
		handlers.RefreshToken,
		middleware2.JWTBlacklistMiddleware(redisClient),
	)

	routeConfig.AddRoute(
		fiber.MethodPost,
		"/forgot-password",
		handlers.ForgotPassword,
	)

	// Health check endpoint
	routeConfig.AddRoute(
		fiber.MethodGet,
		"/health",
		func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{
				"status": "ok",
			})
		},
	)

	commonService := http_server.NewApp(
		log,
		port,
		database,
		true,
		routeConfig.GetRoutes(),
		routeConfig.GetMiddlewares(),
	)

	commonService.Start()
	defer commonService.Stop()
}
