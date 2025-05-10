package repository

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type RepositoryConcrete struct {
	log    *logrus.Logger
	gormDb *gorm.DB
}

func NewRepositoryConcrete(
	log *logrus.Logger,
	db *gorm.DB,
) *RepositoryConcrete {
	return &RepositoryConcrete{
		log:    log,
		gormDb: db,
	}
}
