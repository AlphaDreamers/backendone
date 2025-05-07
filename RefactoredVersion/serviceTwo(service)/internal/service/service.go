package service

import (
	"context"
	"github.com/SwanHtetAungPhyo/srvc/internal/model"
	"github.com/google/uuid"
)

type ServiceBehaviour interface {
	CreateService(ctx context.Context, srv *model.ServicePost) error
	GetService(ctx context.Context, serviceId uuid.UUID) (*model.ServicePost, error)
	UpdateService(ctx context.Context, srv *model.ServicePost, userId string) error
	DeleteService(ctx context.Context, serviceId uuid.UUID, userId string) error
	GetServicesByUser(ctx context.Context, userId uuid.UUID, page, pageSize int) ([]*model.ServicePost, int64, error)
}
