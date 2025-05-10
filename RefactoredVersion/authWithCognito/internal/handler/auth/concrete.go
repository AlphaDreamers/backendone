package auth

import (
	authsrv "github.com/SwanHtetAungPhyo/authCognito/internal/service/auth"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
)

var HandlerLayerModule = fx.Module("handler_module",
	fx.Provide(
		NewConcreteHandler,
	),
)

type ConcreteHandler struct {
	log *logrus.Logger
	srv *authsrv.AuthConcrete
}

func NewConcreteHandler(
	log *logrus.Logger,
	srv *authsrv.AuthConcrete) *ConcreteHandler {
	return &ConcreteHandler{
		log: log,
		srv: srv,
	}
}
