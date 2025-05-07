package handler

import (
	"context"
	"github.com/SwanHtetAungPhyo/srvc/internal/service"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	log     *logrus.Logger
	service *service.Service
	ctx     context.Context
}

func NewHandler(log *logrus.Logger, srv *service.Service) *Handler {
	return &Handler{
		log:     log,
		service: srv,
		ctx:     context.Background(),
	}
}
