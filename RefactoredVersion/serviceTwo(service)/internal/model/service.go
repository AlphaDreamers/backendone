package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type (
	ServicePost struct {
		ServiceID   uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
		ServiceName string    `gorm:"type:varchar(255);not null"`
		ServiceType string    `gorm:"type:varchar(100);not null"`
		Description string    `gorm:"type:text"`
		PhotoUrl    string    `gorm:"type:text"`
		Fee         float64   `gorm:"type:decimal(10,2)"`
		OwnerID     uuid.UUID `gorm:"type:uuid;not null"`
		//OwnerName   string         `gorm:"type:varchar(255)"`
		CreatedAt time.Time      `gorm:"not null;autoCreateTime"`
		UpdatedAt time.Time      `gorm:"not null;autoUpdateTime"`
		DeletedAt gorm.DeletedAt `gorm:"index"`
	}

	Review struct {
		ReviewID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
		ServiceID       uuid.UUID `gorm:"type:uuid;not null"`
		ReviewerID      uuid.UUID `gorm:"type:uuid;not null"`
		ReviewDesc      string    `gorm:"type:text"`
		ReviewSignature string    `gorm:"type:varchar(512)"`
		CreatedAt       time.Time `gorm:"not null;autoCreateTime"`
		UpdatedAt       time.Time `gorm:"not null;autoUpdateTime"`
	}
	SrvReq struct {
		SrvName string  `json:"srv_name"`
		SrvType string  `json:"srv_type"`
		Desc    string  `json:"desc"`
		Photo   string  `json:"photo"`
		Fee     float64 `json:"fee"`
	}
	SrvPost struct {
		ServiceId   uuid.UUID `json:"service_id"`
		ServiceName string    `json:"service_name"`
		ServiceType string    `json:"service_type"`
		Description string    `json:"description"`
		PhotoUrl    string    `json:"photo_url"`
		OwnerId     uuid.UUID `json:"owner_id"`
		OwnerName   string    `json:"owner_name"`
		CreateAt    time.Time `json:"create_at"`
		UpdateAt    time.Time `json:"update_at"`
		DeleteAt    time.Time `json:"delete_at"`
	}
)

func (ServicePost) TableName() string {
	return "service_posts"
}

func (Review) TableName() string {
	return "reviews"
}
