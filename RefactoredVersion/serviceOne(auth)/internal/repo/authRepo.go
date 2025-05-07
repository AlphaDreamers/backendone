package repo

import (
	"context"
	"github.com/SwanHtetAungPhyo/service-one/auth/internal/model"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type AuthRepoBehaviour interface {
	Create(ctx context.Context, a *model.User) (bool, error)
	GetByID(ctx context.Context, id string) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, a *model.User) (bool, error)
	PartialUpdateAfterEmailCode(ctx context.Context, email string) (bool, error)
	UpdatePassword(ctx context.Context, email string, newPassword string) (bool, error)
}

type AuthRepo struct {
	log *logrus.Logger
	db  *gorm.DB
}

func NewAuthRepo(log *logrus.Logger, db *gorm.DB) *AuthRepo {
	return &AuthRepo{
		log: log,
		db:  db,
	}
}

func (ar *AuthRepo) Create(ctx context.Context, a *model.User) (bool, error) {
	err := ar.db.WithContext(ctx).Create(a).Error
	if err != nil {
		ar.log.WithError(err).Error("Failed to create user")
		return false, err
	}
	ar.log.WithField("user", a).Info("Created user")
	return true, nil
}

func (ar *AuthRepo) GetByID(ctx context.Context, id string) (*model.User, error) {
	var req *model.User
	err := ar.db.WithContext(ctx).Where("id = ?", id).First(&req).Error
	if err != nil {
		ar.log.WithError(err).Error("Failed to get user")
		return nil, err
	}
	return req, nil
}

func (ar *AuthRepo) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var result *model.User
	err := ar.db.WithContext(ctx).Where("email = ?", email).First(&result).Error
	if err != nil {
		ar.log.WithError(err).Error("Failed to get user by email")
		return nil, err
	}
	return result, nil
}

func (ar *AuthRepo) Delete(ctx context.Context, id string) error {
	err := ar.db.WithContext(ctx).Where("id = ?", id).Delete(&model.User{}).Error
	if err != nil {
		ar.log.WithError(err).Error("Failed to delete user")
		return err
	}
	return nil
}

func (ar *AuthRepo) Update(ctx context.Context, a *model.User) (bool, error) {
	err := ar.db.WithContext(ctx).Where("email = ?", a.Email).Updates(a).Error
	if err != nil {
		ar.log.WithError(err).Error("Failed to update user")
		return false, err
	}
	return true, nil
}

func (ar *AuthRepo) PartialUpdateAfterEmailCode(ctx context.Context, email string) (bool, error) {
	err := ar.db.WithContext(ctx).
		Model(&model.User{}).
		Where("email = ?", email).
		UpdateColumn("is_verified", true).Error
	if err != nil {
		return false, err
	}
	return true, nil
}

func (ar *AuthRepo) UpdatePassword(ctx context.Context, email string, newPassword string) (bool, error) {
	err := ar.db.WithContext(ctx).
		Model(&model.User{}).
		Where("email = ?", email).
		Update("password", newPassword).Error

	if err != nil {
		ar.log.WithError(err).Error("failed to update password")
		return false, err
	}
	return true, nil
}

var _ AuthRepoBehaviour = (*AuthRepo)(nil)
