package user

import (
	"github.com/SwanHtetAungPhyo/authCognito/internal/model"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

var UserRepoModule = fx.Module("user_repo", fx.Provide(
	NewUserRepository,
))

type UserRepository struct {
	log *logrus.Logger
	db  *gorm.DB
}

func NewUserRepository(log *logrus.Logger, db *gorm.DB) *UserRepository {
	return &UserRepository{
		log: log,
		db:  db,
	}
}
func (u *UserRepository) UpdateAvatar(email string, avatarUrl string) error {
	if err := u.db.Model(&model.User{}).Where("cognito_user_name = ?", email).Update("avatar", avatarUrl).Error; err != nil {
		u.log.WithError(err).Error("failed to update avatar")
		return err
	}
	return nil
}
