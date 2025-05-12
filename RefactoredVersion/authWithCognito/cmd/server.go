package cmd

import (
	"context"
	"github.com/SwanHtetAungPhyo/authCognito/cmd/middleware"
	authHandler "github.com/SwanHtetAungPhyo/authCognito/internal/handler/auth"
	userHandler "github.com/SwanHtetAungPhyo/authCognito/internal/handler/user"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

type AppState struct {
	log         *logrus.Logger
	fiberApp    *fiber.App
	v           *viper.Viper
	handler     *authHandler.ConcreteHandler
	userHandler *userHandler.UserHandler
}

func NewAppState(
	log *logrus.Logger,
	fiberApp *fiber.App,
	v *viper.Viper,
	handler *authHandler.ConcreteHandler,
	userHandler *userHandler.UserHandler,
) *AppState {
	return &AppState{
		log:         log,
		fiberApp:    fiberApp,
		v:           v,
		handler:     handler,
		userHandler: userHandler,
	}
}

func (s *AppState) Start() error {
	//pwd, _ := os.Getwd()

	//certPath := pwd + "/cmd" + s.v.GetString("fiber.cert")
	//keyPath := pwd + "/cmd" + s.v.GetString("fiber.key")
	port := s.v.GetString("fiber.port")
	go func() {
		s.fiberApp.Use(cors.New(cors.Config{
			AllowOrigins: "*",
			AllowHeaders: "*",
			AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		}))
		s.Routes()
		err := s.fiberApp.Listen(":" + port)
		if err != nil {
			s.log.WithError(err).Fatal("fiber.app failed to start")
		}
	}()
	return nil
}
func (s *AppState) middlewareSetup() {
	s.fiberApp.Use(limiter.New(limiter.ConfigDefault))
	s.fiberApp.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "*",
	}))
	s.fiberApp.Use(logger.New(logger.Config{
		Format:     "${pid} ${status} - ${method} ${path}\n",
		TimeFormat: "02-Jan-2006",
		TimeZone:   "America/New_York",
	}))
}
func (s *AppState) Routes() {
	s.fiberApp.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"success": true,
		})
	})
	s.fiberApp.Post("/auth/sign-in", s.handler.SignIn)
	s.fiberApp.Post("/auth/confirm", s.handler.Confirm)
	s.fiberApp.Post("/auth/resend/:email", s.handler.ResendConfirmation)
	s.fiberApp.Post("/auth/sign-up", s.handler.SignUp)
	s.fiberApp.Get("/hello", middleware.JwtMiddleware(), func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"success": true,
			"data": fiber.Map{
				"middleware": "work",
			},
		})
	})
	// Forgot password route
	s.fiberApp.Post("/auth/forgot-password", s.handler.ForgotPassword)

	// Reset password confirmation route
	s.fiberApp.Post("/auth/reset-password-confirm", s.handler.ResetPasswordConfirm)

	// Logout route
	s.fiberApp.Post("/auth/logout", s.handler.Logout)
	s.fiberApp.Post("/kyc-verify/:email", s.handler.KYCVerify)

	s.fiberApp.Put("/auth/:cognito_user_name", s.userHandler.AvatarUploadHandler)
}
func (s *AppState) Stop() error {
	err := s.fiberApp.Shutdown()
	if err != nil {
		s.log.WithError(err).Fatal("fiber.app failed to shutdown")
		return err
	}
	return nil
}

func AppLifeCycle(lc fx.Lifecycle, state *AppState) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return state.Start()
		},
		OnStop: func(ctx context.Context) error {
			return state.Stop()
		},
	})
}
