package service

import (
	"github.com/SwanHtetAungPhyo/chat-order/internal/repository"
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
