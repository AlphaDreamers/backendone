package service

import (
	"context"
	"github.com/SwanHtetAungPhyo/srvc/internal/model"
	"github.com/google/uuid"
)

func (s *Service) CreateService(ctx context.Context, srv *model.SrvReq, ownerId uuid.UUID) (*model.ServicePost, error) {
	partialTransform := s.ReqToModel(srv)
	partialTransform.OwnerID = ownerId
	createPost, err := s.repo.CreateService(ctx, partialTransform)
	if err != nil {
		return nil, err
	}
	return createPost, nil
}

func (s *Service) GetService(ctx context.Context, serviceId uuid.UUID) (*model.ServicePost, error) {
	return s.repo.GetService(ctx, serviceId)
}

func (s *Service) UpdateService(ctx context.Context, srv *model.ServicePost, userId string) error {
	return s.repo.UpdateService(ctx, srv, userId)
}

func (s *Service) DeleteService(ctx context.Context, serviceId uuid.UUID, userId string) error {
	return s.repo.DeleteService(ctx, serviceId, userId)
}

func (s *Service) GetServicesByUser(ctx context.Context, userId uuid.UUID, page, pageSize int) ([]*model.ServicePost, int64, error) {
	return s.repo.GetServiceWithRelatedUser(ctx, userId, page, pageSize)
}

func (s *Service) ReqToModel(req *model.SrvReq) *model.ServicePost {
	return &model.ServicePost{
		ServiceName: req.SrvName,
		ServiceType: req.SrvType,
		Fee:         req.Fee,
		Description: req.Desc,
		PhotoUrl:    req.Photo,
	}
}
