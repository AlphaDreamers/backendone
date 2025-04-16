package model

import "github.com/google/uuid"

type SendMessageRequest struct {
	RecipientID uuid.UUID `json:"recipient_id" validate:"required"`
	Message     string    `json:"message" validate:"required"`
	MessageType string    `json:"message_type" validate:"oneof=text image video file"`
}
