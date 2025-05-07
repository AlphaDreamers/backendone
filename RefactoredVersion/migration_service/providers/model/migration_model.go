package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID                  uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	FullName            string         `gorm:"not null"`
	Email               string         `gorm:"unique;not null"`
	Password            string         `gorm:"not null"`
	BioMetricHash       string         `gorm:"not null"`
	IsVerified          bool           `gorm:"not null;default:false"`
	WalletCreated       bool           `gorm:"not null;default:false"`
	WalletPublicAddress string         `gorm:"not null"`
	CreatedAt           time.Time      `gorm:"not null;autoCreateTime"`
	UpdatedAt           time.Time      `gorm:"not null;autoUpdateTime"`
	DeletedAt           gorm.DeletedAt `gorm:"index"`
}

func (User) TableName() string {
	return "users"
}

type FirebaseToken struct {
	ID     uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID uuid.UUID `gorm:"type:uuid;not null"`
	Token  string    `gorm:"not null"`
}

func (FirebaseToken) TableName() string {
	return "firebase_tokens"
}

type ServicePost struct {
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

func (ServicePost) TableName() string {
	return "service_posts"
}

type Review struct {
	ReviewID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ServiceID       uuid.UUID `gorm:"type:uuid;not null"`
	ReviewerID      uuid.UUID `gorm:"type:uuid;not null"`
	ReviewDesc      string    `gorm:"type:text"`
	ReviewSignature string    `gorm:"type:varchar(512)"`
	CreatedAt       time.Time `gorm:"not null;autoCreateTime"`
	UpdatedAt       time.Time `gorm:"not null;autoUpdateTime"`
}

func (Review) TableName() string {
	return "reviews"
}
