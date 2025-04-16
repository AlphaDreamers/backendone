package model

import (
	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
	"time"
)

type (
	Client struct {
		UserID string          // Authenticated user ID
		Conn   *websocket.Conn // WebSocket connection
	}
)

type (
	ChatConversation struct {
		ConversationID uuid.UUID `json:"conversation_id" bson:"conversation_id"`
		User1ID        uuid.UUID `json:"user1_id" bson:"user1_id"`
		User2ID        uuid.UUID `json:"user2_id" bson:"user2_id"`
		CreatedAt      time.Time `json:"created_at" bson:"created_at"`
	}
)

func NewChatConversation(user1ID uuid.UUID, user2ID uuid.UUID) *ChatConversation {
	return &ChatConversation{
		ConversationID: uuid.New(),
		User1ID:        user1ID,
		User2ID:        user2ID,
		CreatedAt:      time.Now(),
	}
}

type (
	ChatMessage struct {
		MessageID      uuid.UUID  `json:"message_id" bson:"message_id"`
		ConversationID uuid.UUID  `json:"conversation_id" bson:"conversation_id"`
		SenderID       uuid.UUID  `json:"sender_id" bson:"sender_id"`
		RecipientID    uuid.UUID  `json:"recipient_id" bson:"recipient_id"`
		Message        string     `json:"message" bson:"message"`
		MessageType    string     `json:"message_type" bson:"message_type"` // "text", "image", etc.
		Status         string     `json:"status" bson:"status"`             // "sent", "delivered", "read"
		CreatedAt      time.Time  `json:"created_at" bson:"created_at"`
		UpdatedAt      time.Time  `json:"updated_at" bson:"updated_at"`
		DeletedAt      *time.Time `json:"deleted_at,omitempty" bson:"deleted_at,omitempty"`
		IsDeleted      bool       `json:"is_deleted" bson:"is_deleted"`
	}
)

type (
	MessageStatus string
)
type (
	MessageType string
)

const (
	StatusSent      MessageStatus = "sent"
	StatusDelivered MessageStatus = "delivered"
	StatusRead      MessageStatus = "read"
)

const (
	TypeText  MessageType = "text"
	TypeImage MessageType = "image"
	TypeVideo MessageType = "video"
	TypeFile  MessageType = "file"
)
