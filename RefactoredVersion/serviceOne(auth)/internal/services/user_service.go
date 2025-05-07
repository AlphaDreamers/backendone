package services

import (
	"context"
	"github.com/SwanHtetAungPhyo/service-one/auth/internal/model"
	"github.com/SwanHtetAungPhyo/service-one/auth/internal/repo"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type UserService struct {
	log  *logrus.Logger
	repo *repo.UserRepo
}

func NewUserService(log *logrus.Logger, repo *repo.UserRepo) *UserService {
	return &UserService{
		log:  log,
		repo: repo,
	}
}

func (us *UserService) GetMeInfo(
	ctx context.Context,
	userId uuid.UUID,
) (*model.User, error) {

	userInfo, err := us.repo.FindByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}
	return userInfo, nil
}

func (us *UserService) GetDashboardData(page int, pageSize int, userID uuid.UUID) (*model.DashboardResponse, error) {
	return us.repo.GetDashboardData(page, pageSize, userID)
}
