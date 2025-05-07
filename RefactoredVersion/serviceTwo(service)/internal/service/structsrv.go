package service

import (
	"github.com/SwanHtetAungPhyo/srvc/internal/repo"
	"github.com/sirupsen/logrus"
)

type Service struct {
	log  *logrus.Logger
	repo *repo.Repository
}

func NewService(log *logrus.Logger, repo *repo.Repository) *Service {
	return &Service{
		log:  log,
		repo: repo,
	}
}
