package service

import (
	"context"
	"github.com/SwanHtetAungPhyo/chat-order/internal/model"
	"github.com/SwanHtetAungPhyo/chat-order/internal/repository"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type ChatService struct {
	log  *logrus.Logger
	v    *viper.Viper
	repo *repository.ChatRepository
}

func NewChatService(log *logrus.Logger, v *viper.Viper, repo *repository.ChatRepository) *ChatService {
	return &ChatService{
		log:  log,
		v:    v,
		repo: repo,
	}
}

func (s ChatService) GetChatRoomByOrderId(ctx context.Context, id uuid.UUID) {

}

func (s ChatService) GetAllChatRoomByUserId(ctx context.Context, id uuid.UUID) ([]*model.ChatRoom, error) {
	chatsForUser, err := s.repo.GetAllChatRoomByUserId(ctx, id)
	if err != nil {
		return nil, err
	}
	return chatsForUser, nil
}
