package service

import (
	"context"

	"github.com/SwanHtetAungPhyo/service-service/internal/model"
)

// Service interface defines the service layer methods
type Service interface {
	// GetUserServices User Services
	GetUserServices(ctx context.Context, userId string, pagination model.PaginationRequest, filter model.FilterRequest) (*model.ServiceListResponse, error)
	CreateService(ctx context.Context, userId string, service model.ServiceRequest) (*model.ServiceResponse, error)
	UpdateService(ctx context.Context, userId string, serviceId string, service model.ServiceRequest) (*model.ServiceResponse, error)
	DeleteService(ctx context.Context, userId string, serviceId string) error

	//// GetOrderHistory Orders
	//GetOrderHistory(ctx context.Context, userId string, pagination model.PaginationRequest) (*model.OrderListResponse, error)
	//
	//// GetAllTransactions Transactions
	//GetAllTransactions(ctx context.Context, userId string, pagination model.PaginationRequest, filter model.TransactionFilterRequest) (*model.TransactionListResponse, error)
	//GetFiatTransactions(ctx context.Context, userId string, pagination model.PaginationRequest, filter model.TransactionFilterRequest) (*model.TransactionListResponse, error)
	//GetCryptoTransactions(ctx context.Context, userId string, pagination model.PaginationRequest, filter model.TransactionFilterRequest) (*model.TransactionListResponse, error)
	//
	// GetChatHistory Chat
	//GetChatHistory(ctx context.Context, userId string, pagination model.PaginationRequest) (*model.ChatHistoryResponse, error)
	//
	// ListAllServices Public Services
	ListAllServices(ctx context.Context, pagination model.PaginationRequest, filter model.FilterRequest) (*model.ServiceListResponse, error)
	ListServicesByCategory(ctx context.Context, category string, pagination model.PaginationRequest) (*model.ServiceListResponse, error)
}
