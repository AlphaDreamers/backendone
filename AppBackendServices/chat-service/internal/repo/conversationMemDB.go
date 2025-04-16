package repo

import (
	"errors"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"githubc.com/SwanHtetAungPhyo/chat-service/internal/model"
	"sync"
)

type ChatConversationInterface interface {
	Add(sender, receiver uuid.UUID)
	GetID(sender uuid.UUID, req model.StartConversationRequest) (*uuid.UUID, error)
	Delete(req model.StartConversationRequest) error
}

var _ ChatConversationInterface = (*ConversationMemDB)(nil)

type ConversationMemDB struct {
	logger  *logrus.Logger
	storage map[string]*model.ChatConversation
	mutex   sync.RWMutex
}

var _ ChatConversationInterface = (*ConversationMemDB)(nil)

func NewConversationMemDB(logger *logrus.Logger) *ConversationMemDB {
	return &ConversationMemDB{
		logger:  logger,
		storage: make(map[string]*model.ChatConversation),
	}
}

func (c *ConversationMemDB) Add(sender, receiver uuid.UUID) {
	if c.storage == nil {
		c.storage = make(map[string]*model.ChatConversation)
	}
	key, _ := c.mapKeysProducer(sender, receiver)
	if _, ok := c.storage[key]; !ok {
		c.logger.Infof("New Conversation is created with MapKey %v", key)
		c.storage[key] = model.NewChatConversation(sender, receiver)
	}
}

func (c *ConversationMemDB) GetID(sender uuid.UUID, req model.StartConversationRequest) (*uuid.UUID, error) {
	if c.storage == nil {
		return nil, errors.New("no conversation storage")
	}
	c.mutex.Lock()
	defer c.mutex.Unlock()

	mapKey, reverseKey := c.mapKeysProducer(sender, req.RecipientID)

	if conv, ok := c.storage[mapKey]; ok {
		return &conv.ConversationID, nil
	}
	if conv, ok := c.storage[reverseKey]; ok {
		return &conv.ConversationID, nil
	}
	return nil, errors.New("no conversation found")
}

func (c *ConversationMemDB) mapKeysProducer(sender, req uuid.UUID) (string, string) {
	mapKey := c.mapHelper(sender, req)
	reversedKey := c.reverseMapKey(mapKey)
	return mapKey, reversedKey
}
func (c *ConversationMemDB) Delete(req model.StartConversationRequest) error {
	//TODO implement me
	panic("implement me")
}

func (c *ConversationMemDB) mapHelper(sender, receiver uuid.UUID) string {
	return sender.String() + receiver.String()
}

func (c *ConversationMemDB) reverseMapKey(mapKey string) string {
	runes := []rune(mapKey)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
