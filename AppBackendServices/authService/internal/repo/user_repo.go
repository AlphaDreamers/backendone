package repo

import (
	"errors"
	"fmt"
	"github.com/SwanHtetAungPhyo/common/models"
	"github.com/SwanHtetAungPhyo/common/pkg/logutil"
	db "github.com/SwanHtetAungPhyo/database"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Impl struct {
	db  *gorm.DB
	log *logrus.Logger
}

func NewImpl() *Impl {
	return &Impl{
		db:  db.GetDB(),
		log: logutil.GetLogger(),
	}
}

func (i *Impl) Create(user *models.UserInDB, biometricHash string) error {
	return i.db.Transaction(func(tx *gorm.DB) error {
		// Create user
		user.Verified = false
		user.WalletCreated = false
		if err := tx.Create(user).Error; err != nil {
			i.log.Error("User creation failed:", err)
			return fmt.Errorf("user creation failed: %w", err)
		}

		biometric := models.UserBiometric{
			UserID:        user.ID,
			BioMetricHash: biometricHash,
		}
		if err := tx.Create(&biometric).Error; err != nil {
			i.log.Error("Biometric creation failed:", err)
			return fmt.Errorf("biometric creation failed: %w", err)
		}

		return nil
	})
}

func (i *Impl) Login(req *models.LoginRequest) (*models.UserInDB, error) {
	var userInDB models.UserInDB
	result := i.db.First(&userInDB, "email = ?", req.Email)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			i.log.Error(result.Error.Error())
			return nil, fmt.Errorf("user not found")
		}
		return nil, result.Error
	}
	i.log.Info("User in DB: ", userInDB.Email)
	return &userInDB, nil
}

func (i *Impl) UpdateUserStatus(email string) error {
	result := i.db.Model(&models.UserInDB{}).
		Where("email = ?", email).
		Update("verified", true)

	if result.Error != nil {
		i.log.Errorf("Update failed for email %s: %v", email, result.Error)
		return fmt.Errorf("failed to update user status")
	}

	if result.RowsAffected == 0 {
		i.log.Warnf("No user found with email %s", email)
		return fmt.Errorf("user not found")
	}

	i.log.Infof("Successfully updated verification status for %s", email)
	return nil
}
func (i *Impl) GetByEmail(email string) (*models.UserInDB, error) {
	var user models.UserInDB
	if err := i.db.First(&user, "email = ?", email).Error; err != nil {
		i.log.Error(err.Error())
		return nil, fmt.Errorf("user not found")
	}
	return &user, nil
}

func (i *Impl) UpdateWalletStatus(email string) error {
	result := i.db.Model(&models.UserInDB{}).
		Where("email = ?", email).
		Update("wallet_created", true)

	if result.Error != nil {
		i.log.Errorf("Update failed for email %s: %v", email, result.Error)
		return fmt.Errorf("failed to update user status")
	}

	if result.RowsAffected == 0 {
		i.log.Warnf("No user found with email %s", email)
		return fmt.Errorf("user not found")
	}

	i.log.Infof("Successfully updated verification status for %s", email)
	return nil
}
