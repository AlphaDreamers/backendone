package main

import (
	"fmt"
	"github.com/SwanHtetAungPhyo/service-one/auth/cmd/server"
	"github.com/SwanHtetAungPhyo/service-one/auth/internal/handler"
	"github.com/SwanHtetAungPhyo/service-one/auth/internal/repo"
	"github.com/SwanHtetAungPhyo/service-one/auth/internal/services"
	"github.com/SwanHtetAungPhyo/service-one/auth/utils"
	"github.com/nats-io/nats.go"
	"go.uber.org/fx"
)

func main() {
	app := fx.New(

		server.ProviderModule,
		fx.Provide(
			repo.NewUserRepo,
			services.NewUserService,
		),
		fx.Provide(
			utils.NewJwtTokenGenerator,
		),
		handler.ProviderModule,

		fx.Provide(
			server.NewAppState,
		),
		fx.Invoke(
			Hello,
			//server.Migration,
			server.ServerLifeCycleReg,
		),
	)
	app.Run()
}

func Hello() {
	fmt.Println("Hello World")
}

func NatsProvider() *nats.Conn {
	nc, err := nats.Connect("nats://localhost:4222")
	if err != nil {
		return nil
	}
	return nc
}
