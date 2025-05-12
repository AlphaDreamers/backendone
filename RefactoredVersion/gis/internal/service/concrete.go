package service

import (
	"context"
	"github.com/SwanHtetAungPhyo/gis/internal/model/model"
	"github.com/SwanHtetAungPhyo/gis/internal/model/req"
	"github.com/SwanHtetAungPhyo/gis/internal/repo"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type GigService struct {
	log  *logrus.Logger
	repo *repo.GigRepository
	ctx  context.Context
}

func NewGigService(
	log *logrus.Logger,
	repo *repo.GigRepository) *GigService {
	return &GigService{
		log:  log,
		repo: repo,
		ctx:  context.Background(),
	}
}

func (gs *GigService) CreateGig(req *req.CreateGigRequest) (*model.Gig, error) {
	createdGig, err := gs.repo.CreateGig(gs.ctx, req)
	if err != nil {
		gs.log.Debug("git create gig err:", err.Error())
		return nil, err
	}
	return createdGig, nil
}

func (gs *GigService) UpdateGig(req *req.UpdateGigRequest) (*model.Gig, error) {

	updateGig, err := gs.repo.UpdateGig(gs.ctx, req)
	if err != nil {
		gs.log.Debug("git update gig err:", err.Error())
		return nil, err
	}
	return updateGig, nil
}

func (gs *GigService) PartialUpdate(gigId uuid.UUID, updates map[string]any) (*model.Gig, error) {
	updatedVersion, err := gs.repo.PartialUpdate(gs.ctx, gigId, updates)
	if err != nil {
		gs.log.Debug("git partial update err:", err.Error())
		return nil, err
	}
	return updatedVersion, nil
}

func (gs *GigService) GetAllGigByPaganition(page int, perPage int) (int64, []*model.Gig, error) {
	total, err := gs.repo.CountTheGig()
	if err != nil {
		gs.log.Debug("git count gig err:", err.Error())
		return 0, nil, err
	}

	offset := (page - 1) * perPage
	byOffset, err := gs.repo.GetByOffset(total, page, perPage, offset)
	if err != nil {
		gs.log.Debug("git get gig by offset err:", err.Error())
		return 0, nil, err
	}

	return *total, byOffset, nil
}

func (gs *GigService) AddPackageToGig(id uuid.UUID, request *req.GigPackageRequest) (*model.GigPackage, error) {
	gs.repo.AddPackageToGig(id, request)
}
