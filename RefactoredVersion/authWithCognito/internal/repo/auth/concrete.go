package auth

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type AuthRepositry struct {
	log  *logrus.Logger
	gorm *gorm.DB
}

func NewAuthRepositry(log *logrus.Logger, gorm *gorm.DB) *AuthRepositry {
	return &AuthRepositry{
		log:  log,
		gorm: gorm,
	}
}
