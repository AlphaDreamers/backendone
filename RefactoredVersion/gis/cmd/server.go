package cmd

import (
	"context"
	"github.com/SwanHtetAungPhyo/gis/cmd/middleware"
	"github.com/SwanHtetAungPhyo/gis/internal/handler"
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
	handler  *handler.GigHandler
}

func NewAppState(
	log *logrus.Logger,
	fiberApp *fiber.App,
	v *viper.Viper,
	handler *handler.GigHandler,
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
	s.Routes()
	go func() {

		err := s.fiberApp.ListenTLS(":"+port, certPath, keyPath)
		if err != nil {
			s.log.WithError(err).Fatal("fiber.app failed to start")
		}
	}()
	return nil
}

func (s *AppState) Routes() {
	gig := s.fiberApp.Group("/gig")
	gig.Post("/", middleware.AuthMiddleware(), s.handler.CreateGig)
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
