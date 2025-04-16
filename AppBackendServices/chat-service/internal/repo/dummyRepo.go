package repo

import (
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"githubc.com/SwanHtetAungPhyo/chat-service/internal/model"
)

type DummyRepo interface {
	GetOrCreateConversation(sender uuid.UUID, req *model.StartConversationRequest) (uuid.UUID, error)
}
type DummyRepoImpl struct {
	logger  *logrus.Logger
	storage ChatConversationInterface
}

var _ DummyRepo = (*DummyRepoImpl)(nil)
var err error

func NewDummyRepo(logger *logrus.Logger, storage ChatConversationInterface) *DummyRepoImpl {
	return &DummyRepoImpl{
		logger:  logger,
		storage: storage,
	}
}
func (c *DummyRepoImpl) GetOrCreateConversation(sender uuid.UUID, req *model.StartConversationRequest) (uuid.UUID, error) {
	var conversationID *uuid.UUID
	conversationID, err = c.storage.GetID(sender, *req)
	if err != nil {
		c.logger.Warn("Failed to get conversation id: %v", err)
	}
	c.logger.Info("Trying to create conversation with id")
	c.storage.Add(sender, req.RecipientID)
	conversationID, err = c.storage.GetID(sender, *req)
	if err != nil {
		c.logger.Warn("Failed to get conversation id again: %v", err)
	}
	c.logger.Debug("Got conversation id: %v", conversationID)
	return *conversationID, nil
}
