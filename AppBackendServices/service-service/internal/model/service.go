package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Service represents the Services table in the database
type Service struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;column:id"`
	ServiceName string         `gorm:"type:varchar(100);not null;column:service_name"`
	Description string         `gorm:"type:text;not null;column:description"`
	OfferedBy   uuid.UUID      `gorm:"type:uuid;not null;column:offeredBy"`
	CreatedAt   time.Time      `gorm:"type:timestamp;default:CURRENT_TIMESTAMP;column:createdAt"`
	UpdatedAt   *time.Time     `gorm:"type:timestamp;column:updatedAt"`
	Available   bool           `gorm:"type:bool;not null;column:available"`
	Minimum     float64        `gorm:"type:float8;not null;column:minimum"`
	Maximum     float64        `gorm:"type:float8;not null;column:maximum"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

// TableName specifies the table name for the Service model
func (Service) TableName() string {
	return "Services"
}

// BeforeCreate is a GORM hook that runs before creating a new record
func (s *Service) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}

// BeforeUpdate is a GORM hook that runs before updating a record
func (s *Service) BeforeUpdate(tx *gorm.DB) error {
	now := time.Now()
	s.UpdatedAt = &now
	return nil
}
