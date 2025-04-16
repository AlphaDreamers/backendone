package model

import "github.com/google/uuid"

type StartConversationRequest struct {
	RecipientID uuid.UUID `json:"recipient_id" validate:"required"`
}
