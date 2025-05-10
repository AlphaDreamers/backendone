package service

import (
	"github.com/SwanHtetAungPhyo/gis/internal/repository"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/sirupsen/logrus"
)

type ServiceConcrete struct {
	log  *logrus.Logger
	repo *repository.RepositoryConcrete
	s3   *s3.Client
}

func NewServiceConcrete(
	log *logrus.Logger,
	repo *repository.RepositoryConcrete,
	s3Client *s3.Client,
) *ServiceConcrete {
	return &ServiceConcrete{
		log:  log,
		repo: repo,
		s3:   s3Client,
	}
}
