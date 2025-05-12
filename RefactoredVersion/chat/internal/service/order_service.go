package service

import (
	"github.com/SwanHtetAungPhyo/chat-order/internal/model"
	"github.com/SwanHtetAungPhyo/chat-order/internal/repository"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

var ServiceModule = fx.Module("service", fx.Provide(
	NewOrderService,
))

type OrderService struct {
	log  *logrus.Logger
	v    *viper.Viper
	repo *repository.OrderRepo
}

func NewOrderService(
	log *logrus.Logger,
	v *viper.Viper,
	repo *repository.OrderRepo,
) *OrderService {
	return &OrderService{
		log:  log,
		v:    v,
		repo: repo,
	}
}

func (os *OrderService) FindOrCreate(req *model.OrderPlaceRequest) (*model.Order, error) {
	orderInDB, err := os.repo.FindOrCreate(req)
	if err != nil {
		os.log.Error(err)
		return nil, err
	}
	return orderInDB, nil
}
