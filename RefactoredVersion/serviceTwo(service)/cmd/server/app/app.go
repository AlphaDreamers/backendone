package app

import (
	"context"
	"github.com/SwanHtetAungPhyo/srvc/internal/handler"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"gorm.io/gorm"
	"time"
)

type AppState struct {
	log *logrus.Logger
	db  *gorm.DB
	app *fiber.App
	//mongo    *mongo.Client
	v *viper.Viper
	//s3Client *s3.Bucket
	handler *handler.Handler
}

func NewAppState(
	log *logrus.Logger,
	db *gorm.DB,
	app *fiber.App,
	v *viper.Viper,
	handler *handler.Handler,
) *AppState {
	return &AppState{
		log: log,
		db:  db,
		app: app,
		//mongo:   mongo,
		v:       v,
		handler: handler,
	}
}

func (state *AppState) Start() error {
	state.log.Info("starting app")
	port := state.v.GetString("server.port")
	cert := state.v.GetString("server.cert")
	key := state.v.GetString("server.key")
	state.routeSetup()
	var err error
	go func() {
		err = state.app.ListenTLS(port, cert, key)
		if err != nil {
			return
		}
	}()
	return nil
}

func (state *AppState) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return state.app.ShutdownWithContext(ctx)
}

func (state *AppState) routeSetup() {
	state.log.Info("Setting up routes...")

	// Health check route
	state.app.Get("/health", func(ctx *fiber.Ctx) error {
		state.log.Info("Health check called")
		return ctx.JSON(fiber.Map{"status": "ok"})
	})

	srvGrp := state.app.Group("/service")
	srvGrp.Use(func(c *fiber.Ctx) error {
		state.log.Infof("Service route called: %s %s", c.Method(), c.Path())
		return c.Next()
	})

	srvGrp.Get("/:userId", state.handler.GetSpcServiceByUserId)
	srvGrp.Post("/:userId", state.handler.CreateService)
	srvGrp.Put("/:userId/:serviceId", state.handler.UpdateService)
	srvGrp.Delete("/:userId/:serviceId", state.handler.DeleteService)
	srvGrp.Get("/:userId/:serviceId", state.handler.GetServiceByServiceId)
	srvGrp.Get("/", state.handler.GetAllServices)

}

func (state *AppState) middleWareSetup() {
	// Define any middleware logic here
}

func AppLifeCycle(lc fx.Lifecycle, server *AppState) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return server.Start()
		},
		OnStop: func(ctx context.Context) error {
			return server.Stop()
		},
	})
}
