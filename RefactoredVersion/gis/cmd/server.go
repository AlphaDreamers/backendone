package cmd

import (
	"context"
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
	handler  *handler.HandlerConcrete
}

func NewAppState(
	log *logrus.Logger,
	fiberApp *fiber.App,
	v *viper.Viper,
) *AppState {
	return &AppState{
		log:      log,
		fiberApp: fiberApp,
		v:        v,
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
	api := s.fiberApp.Group("/api/gig")
	api.Post("/", s.handler.Create)
	api.Put("/:id", s.handler.Update)
	api.Delete("/:id", s.handler.Delete)
	api.Get("/:id", s.handler.GetById)
	api.Get("/", s.handler.List)
	api.Get("/user/:userId", s.handler.GetByUserId)
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
