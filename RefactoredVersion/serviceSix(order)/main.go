package main

import (
	"github.com/SwanHtetAungPhyo/order_service/cmd/server"
	"github.com/SwanHtetAungPhyo/order_service/providers"
	"go.uber.org/fx"
)

func main() {
	fxApp := fx.New(
		providers.StartUpProvider,
		fx.Provide(server.NewAppState),
		fx.Invoke(
			server.AppLifeCycle,
		))
	fxApp.Run()
}
