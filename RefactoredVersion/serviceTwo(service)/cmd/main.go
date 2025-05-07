package main

import (
	ap "github.com/SwanHtetAungPhyo/srvc/cmd/server/app"
	"github.com/SwanHtetAungPhyo/srvc/cmd/server/providers"
	"github.com/SwanHtetAungPhyo/srvc/internal/handler"
	"github.com/SwanHtetAungPhyo/srvc/internal/repo"
	"github.com/SwanHtetAungPhyo/srvc/internal/service"
	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		providers.ProviderModule,

		fx.Provide(
			repo.NewRepository,
			service.NewService,
			handler.NewHandler,
			ap.NewAppState,
		),
		fx.Invoke(
			ap.AppLifeCycle,
		),
	)
	app.Run()
}
