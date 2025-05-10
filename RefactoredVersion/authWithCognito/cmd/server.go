package cmd

import (
	"context"
	authHandler "github.com/SwanHtetAungPhyo/authCognito/internal/handler/auth"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"os"
)

type AppState struct {
	log      *logrus.Logger
	fiberApp *fiber.App
	v        *viper.Viper
	handler  *authHandler.ConcreteHandler
}

func NewAppState(
	log *logrus.Logger,
	fiberApp *fiber.App,
	v *viper.Viper,
	handler *authHandler.ConcreteHandler,
) *AppState {
	return &AppState{
		log:      log,
		fiberApp: fiberApp,
		v:        v,
		handler:  handler,
	}
}

func (s *AppState) Start() error {
	pwd, _ := os.Getwd()

	certPath := pwd + "/cmd" + s.v.GetString("fiber.cert")
	keyPath := pwd + "/cmd" + s.v.GetString("fiber.key")
	port := s.v.GetString("fiber.port")
	go func() {
		s.Routes()
		err := s.fiberApp.ListenTLS(":"+port, certPath, keyPath)
		if err != nil {
			s.log.WithError(err).Fatal("fiber.app failed to start")
		}
	}()
	return nil
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

	// Forgot password route
	s.fiberApp.Post("/auth/forgot-password", s.handler.ForgotPassword)

	// Reset password confirmation route
	s.fiberApp.Post("/auth/reset-password-confirm", s.handler.ResetPasswordConfirm)

	// Logout route
	s.fiberApp.Post("/auth/logout", s.handler.Logout)
	s.fiberApp.Post("/kyc-verify", s.handler.KYCVerify)
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
