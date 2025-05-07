package repo

import (
	"context"
	"github.com/SwanHtetAungPhyo/srvc/internal/model"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type RepoBehaviour interface {
	CreateService(ctx context.Context, srv *model.ServicePost) error
	GetService(ctx context.Context, serviceId uuid.UUID) (*model.ServicePost, error)
	GetServiceWithRelatedUser(ctx context.Context, userId uuid.UUID) []*model.ServicePost
	UpdateService(ctx context.Context, srv *model.ServicePost, userId string) error
	DeleteService(ctx context.Context, serviceId uuid.UUID, userId string) error
}
type Repository struct {
	log *logrus.Logger
	db  *gorm.DB
}

func NewRepository(log *logrus.Logger, db *gorm.DB) *Repository {
	return &Repository{
		log: log,
		db:  db,
	}
}

func (r *Repository) CreateService(ctx context.Context, srv *model.ServicePost) (*model.ServicePost, error) {
	if err := r.db.WithContext(ctx).Create(srv).Error; err != nil {
		r.log.WithError(err).Error("error creating service")
		return nil, err
	}

	var createdService model.ServicePost
	if err := r.db.WithContext(ctx).First(&createdService, srv.ServiceID).Error; err != nil {
		r.log.WithError(err).Error("error fetching created service")
		return nil, err
	}
	return &createdService, nil
}
func (r *Repository) GetService(ctx context.Context, serviceId uuid.UUID) (*model.ServicePost, error) {
	var service model.ServicePost
	if err := r.db.WithContext(ctx).First(&service, "service_id = ?", serviceId).Error; err != nil {
		r.log.WithError(err).Error("error getting service")
		return nil, err
	}
	return &service, nil
}

func (r *Repository) UpdateService(ctx context.Context, srv *model.ServicePost, userId string) error {
	if err := r.db.WithContext(ctx).
		Where("owner_id = ? AND service_id = ?", userId, srv.ServiceID).
		Updates(srv).Error; err != nil {
		r.log.WithError(err).Error("error updating service")
		return err
	}
	return nil
}
func (r *Repository) GetServiceWithRelatedUser(
	ctx context.Context,
	userId uuid.UUID,
	page int,
	pageSize int,
) ([]*model.ServicePost, int64, error) {
	var (
		services []*model.ServicePost
		total    int64
	)

	offset := (page - 1) * pageSize

	if err := r.db.WithContext(ctx).
		Model(&model.ServicePost{}).
		Where("owner_id = ?", userId).
		Count(&total).Error; err != nil {
		r.log.WithError(err).Error("error counting services by user")
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).
		Preload("Reviews").
		Where("owner_id = ?", userId).
		Limit(pageSize).
		Offset(offset).
		Find(&services).Error; err != nil {
		r.log.WithError(err).Error("error getting paginated services by user with reviews")
		return nil, 0, err
	}

	return services, total, nil
}

func (r *Repository) DeleteService(ctx context.Context, serviceId uuid.UUID, userId string) error {
	if err := r.db.WithContext(ctx).
		Where("owner_id = ? AND service_id = ?", userId, serviceId).
		Delete(&model.ServicePost{}).Error; err != nil {
		r.log.WithError(err).Error("error deleting service")
		return err
	}
	return nil
}
