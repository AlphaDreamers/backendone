package req

import "github.com/google/uuid"

type UpdateGigRequest struct {
	GigId       uuid.UUID `json:"gigId" validate:"required"`
	SellerId    uuid.UUID `json:"sellerId" validate:"required"`
	Title       string    `json:"title" validate:"min=10,max=100"`
	Description string    `json:"description" validate:"min=50,max=1000"`
	IsActive    *bool     `json:"isActive"`
}
