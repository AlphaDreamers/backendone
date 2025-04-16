package models

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type (
	UserInDB struct {
		ID              uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
		FullName        string    `gorm:"not null" json:"fullname" validate:"required,min=3,max=50"`
		Email           string    `gorm:"unique;not null" json:"email" validate:"required,email"`
		Password        string    `gorm:"not null" json:"-" validate:"required,min=8"` // Omit password from JSON
		Verified        bool      `gorm:"not null" json:"verified"`
		CreatedAt       time.Time `gorm:"autoCreateTime" json:"created_at"`
		WalletCreated   bool      `gorm:"not null" json:"wallet_created"`
		WalletCreatedAt time.Time `gorm:"autoCreateTime" json:"wallet_created_at"`
	}

	UserBiometric struct {
		ID            uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
		UserID        uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
		BioMetricHash string    `gorm:"not null" json:"biometric_hash"`
	}

	UserRegisterRequest struct {
		FullName      string ` json:"fullname" validate:"required,min=3,max=50"`
		Email         string ` json:"email" validate:"required,email"`
		Country       string ` validate:"omitempty,min=2,max=50"`
		BioMetricHash string ` json:"biometric_hash" validate:"required"`
		Password      string ` json:"password" validate:"required,min=8"`
	}
)

var userValidator = validator.New()

func (u *UserInDB) Validate() error {
	return userValidator.Struct(u)
}
func (*UserInDB) TableName() string {
	return "users"
}

func (*UserBiometric) TableName() string {
	return "user_biometrics"
}
