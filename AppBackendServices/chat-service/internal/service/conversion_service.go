package service

import (
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"githubc.com/SwanHtetAungPhyo/chat-service/internal/model"
	"githubc.com/SwanHtetAungPhyo/chat-service/internal/repo"
)

type ConversationService interface {
	GetOrCreateConversation(sender uuid.UUID, req *model.StartConversationRequest) (uuid.UUID, error)
}
type ConversationServiceImpl struct {
	logger  *logrus.Logger
	storage repo.DummyRepo
}

var _ ConversationService = (*ConversationServiceImpl)(nil)

func NewConversationService(logger *logrus.Logger, storage repo.DummyRepo) *ConversationServiceImpl {
	return &ConversationServiceImpl{
		logger:  logger,
		storage: storage,
	}
}
func (c ConversationServiceImpl) GetOrCreateConversation(sender uuid.UUID, req *model.StartConversationRequest) (uuid.UUID, error) {
	conversationID, err := c.storage.GetOrCreateConversation(sender, req)
	if err != nil {
		return uuid.Nil, err
	}
	return conversationID, nil
}
