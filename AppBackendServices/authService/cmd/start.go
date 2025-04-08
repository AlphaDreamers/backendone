package cmd

import (
	"context"
	"github.com/SwanHtetAungPhyo/auth/cmd/middleware"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"os"
	"os/signal"
	"syscall"

	"github.com/SwanHtetAungPhyo/auth/internal/handler"
	"github.com/SwanHtetAungPhyo/common/pkg/logutil"
	"github.com/gofiber/fiber/v2"
)

func Start(port string) {
	app := fiber.New()
	startlog := logutil.GetLogger()

	routeSetUp(app)
	go func() {
		startlog.Info("AuthService  is started  at the port")
		if err := app.Listen(":" + port); err != nil {
			startlog.Panic(err.Error())
		}

	}()
	osChannel := make(chan os.Signal, 1)
	signal.Notify(osChannel, syscall.SIGABRT, os.Interrupt)
	<-osChannel

	if err := app.ShutdownWithContext(context.Background()); err != nil {
		startlog.Fatal("failed to shutdown")
	}

}

func routeSetUp(app *fiber.App) {
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:8085",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
		ExposeHeaders:    "Content-Type, Authorization",
	}))
	app.Use(middleware.LoggingMiddleware())
	handlers := handler.NewHandler()
	app.Post("/login", handlers.Login)
	app.Post("/register", handlers.Register)
	app.Get("/refresh", handlers.Refresh)
	app.Post("/me", handlers.Me, middleware.JwtMiddleware())
	app.Post("/verify", handlers.Verify)
	app.Post("/wallet", handlers.StoreInVault, middleware.JwtMiddleware())
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})
}
