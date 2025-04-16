package model

import (
	"github.com/google/uuid"
	"time"
)

type MessageResponse struct {
	MessageID      uuid.UUID `json:"message_id"`
	ConversationID uuid.UUID `json:"conversation_id"`
	SenderID       uuid.UUID `json:"sender_id"`
	RecipientID    uuid.UUID `json:"recipient_id"`
	Message        string    `json:"message"`
	MessageType    string    `json:"message_type"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
}
