package repository

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type ChatRepository struct {
	log *logrus.Logger
	v   *viper.Viper
	db  *gorm.DB
}

func NewChatRepository(log *logrus.Logger, v *viper.Viper, db *gorm.DB) *ChatRepository {
	return &ChatRepository{
		log: log,
		v:   v,
		db:  db,
	}
}
