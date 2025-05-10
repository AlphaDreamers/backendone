package handler

import (
	"github.com/SwanHtetAungPhyo/gis/internal/service"
	"github.com/sirupsen/logrus"
)

type HandlerConcrete struct {
	log *logrus.Logger
	srv *service.ServiceConcrete
}

func NewHandlerConcrete(
	log *logrus.Logger,
	srv *service.ServiceConcrete,

) *HandlerConcrete {
	return &HandlerConcrete{
		log: log,
		srv: srv,
	}
}
