package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/SwanHtetAungPhyo/chat-order/internal/model"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
	"gorm.io/gorm"
	"math/rand"
	"time"
)

var RepoModule = fx.Module("repository", fx.Provide(
	NewOrderRepo,
	NewChatRepository,
))

type OrderRepo struct {
	log          *logrus.Logger
	dynamoClient *dynamodb.Client
	gormClient   *gorm.DB
}

func NewOrderRepo(
	log *logrus.Logger,
	dynamoClient *dynamodb.Client,
	gorm *gorm.DB,
) *OrderRepo {
	return &OrderRepo{
		log:          log,
		dynamoClient: dynamoClient,
		gormClient:   gorm,
	}
}

func (r OrderRepo) FindOrCreate(req *model.OrderPlaceRequest) (*model.Order, error) {
	var order model.Order
	err := r.gormClient.
		Model(model.Order{}).
		Where("buyer_id = ? AND seller_id = ? AND package_id = ?", req.BuyerId, req.SellerId, req.ServiceId).
		First(&order).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		order = model.Order{
			BuyerID:       req.BuyerId,
			SellerID:      req.SellerId,
			PackageID:     req.ServiceId,
			OrderNumber:   r.GenerateOrderNumber(),
			Price:         0.0,
			PaymentMethod: "NOT_SET",
			Status:        "PENDING",
		}
		if err := r.gormClient.
			Model(model.Order{}).
			Create(&order).Error; err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err // Some other DB error
	}

	return &order, nil

}
func (r OrderRepo) UpdateOrderStatus(orderNumber string, packageId uuid.UUID, status string) error {
	err := r.gormClient.
		WithContext(context.TODO()).
		Model(&model.Order{}).
		Where("order_number = ? AND package_id = ?", orderNumber, packageId).
		Updates(map[string]interface{}{
			"status": status,
		}).Error
	if err != nil {
		r.log.Error(err.Error)
		return err
	}
	return nil
}

func (r OrderRepo) CompleteOrder(orderNumber string, packageId uuid.UUID, status string) error {
	err := r.gormClient.
		WithContext(context.TODO()).
		Model(&model.Order{}).
		Where("order_number = ? AND package_id = ?", orderNumber, packageId).
		Updates(map[string]interface{}{
			"status":       status,
			"completed_at": time.Now(),
		}).Error
	if err != nil {
		r.log.Error(err.Error)
		return err
	}
	return nil
}

func (r OrderRepo) UpdateOrderPartial(orderNumber string, packageId uuid.UUID, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return errors.New("no fields to update")
	}

	err := r.gormClient.
		WithContext(context.TODO()).
		Model(&model.Order{}).
		Where("order_number = ? AND package_id = ?", orderNumber, packageId).
		Updates(updates).Error

	if err != nil {
		r.log.Errorf("failed to update order partially: %v", err)
		return err
	}

	return nil
}

func (r OrderRepo) GenerateOrderNumber() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const suffixLength = 6

	rand.Seed(time.Now().UnixNano())

	suffix := make([]byte, suffixLength)
	for i := range suffix {
		suffix[i] = charset[rand.Intn(len(charset))]
	}

	timestamp := time.Now().Format("20060102150405")
	return fmt.Sprintf("SN%s-%s", timestamp, string(suffix))
}

func (r OrderRepo) CancelOrder(orderNumber string, packageId uuid.UUID) error {
	err := r.gormClient.
		WithContext(context.TODO()).
		Model(&model.Order{}).
		Where("order_number = ? AND package_id = ?", orderNumber, packageId).
		Updates(map[string]interface{}{
			"status":       "CANCELLED",
			"cancelled_at": time.Now(),
		}).Error

	if err != nil {
		r.log.Errorf("failed to cancel order: %v", err)
		return err
	}
	return nil
}
