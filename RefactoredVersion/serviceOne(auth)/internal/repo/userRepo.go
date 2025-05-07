package repo

import (
	"context"
	"fmt"
	"github.com/SwanHtetAungPhyo/service-one/auth/internal/model"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"math"
)

type UserRepo struct {
	log *logrus.Logger
	db  *gorm.DB
}

func NewUserRepo(
	log *logrus.Logger,
	db *gorm.DB,
) *UserRepo {
	return &UserRepo{
		log: log,
		db:  db,
	}
}

func (ur *UserRepo) FindByUserId(ctx context.Context, id uuid.UUID) (*model.User, error) {
	var user model.User
	if err := ur.db.WithContext(ctx).
		Find(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (ur *UserRepo) UserExistence(ctx context.Context, id uuid.UUID) (bool, error) {
	var user model.User
	if err := ur.db.WithContext(ctx).Find(&user, id).Error; err != nil {
		ur.log.WithError(err).Error(ctx, "UserExistence")
		return false, err
	}
	return true, nil
}
func (ur *UserRepo) GetDashboardData(page, pageSize int, userID uuid.UUID) (*model.DashboardResponse, error) {
	var response model.DashboardResponse
	var total int64

	// 1. Get user info
	var user model.User
	if err := ur.db.First(&user, "id = ?", userID).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}
	response.UserInfo = &user
	response.UserInfo.Password = ""

	query := ur.db.Model(&model.ServicePost{}).Where("deleted_at IS NULL")

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count service posts: %w", err)
	}

	// Calculate pagination
	offset := (page - 1) * pageSize
	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	var services []model.ServicePost
	if err := query.
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&services).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve services: %w", err)
	}

	response.Services = services
	response.Pagination.Page = page
	response.Pagination.PageSize = pageSize
	response.Pagination.Total = total
	response.Pagination.TotalPages = totalPages

	return &response, nil
}
