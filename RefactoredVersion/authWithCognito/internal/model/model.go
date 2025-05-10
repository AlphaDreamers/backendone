package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Country struct {
	ID    uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name  string    `gorm:"type:varchar(255);not null"`
	Users []User    `gorm:"foreignKey:CountryID"`
}

type User struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	//FirstName       string         `json:"firstName"`
	//LastName        string         `json:"lastName"`
	//Email           string         `json:"email"`
	//Password        string         `json:"password"`
	//Country         string         `json:"country"`
	//BioMetricHash   string         `json:"bioMetricHash"`
	//IsVerified      bool           `gorm:"not null;default:false"`
	WalletCreated   bool           `gorm:"not null;default:false"`
	WalletCreatedAt *time.Time     `gorm:"default:null"`
	CognitoId       string         `gorm:"type:text;not null" json:"cognitoId"`
	Avatar          *string        `gorm:"type:text;default:null"`
	CreatedAt       time.Time      `gorm:"not null;autoCreateTime"`
	UpdatedAt       time.Time      `gorm:"not null;autoUpdateTime"`
	DeletedAt       gorm.DeletedAt `gorm:"index"`
	ReceivedOrders  []Order        `gorm:"foreignKey:BuyerID"`
	YoursOrders     []Order        `gorm:"foreignKey:SellerID"`
	Chats           []Chat         `gorm:"foreignKey:BuyerID"`
}
type UserData struct {
	UserID        string `json:"sub"`
	Username      string `json:"cognito:username"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	FullName      string `json:"fullName"`
	FirstName     string `json:"given_name"`
	LastName      string `json:"family_name"`
	Country       string `json:"custom:country"`
	BioMetricHash string `json:"custom:bio_metric_hash"`
}
type OrderStatus string

const (
	StatusPending         OrderStatus = "pending"
	StatusCompleted       OrderStatus = "completed"
	StatusInProgress      OrderStatus = "in-progress"
	StatusAwaitingPayment OrderStatus = "awaiting-payment"
	StatusPendingDelivery OrderStatus = "pending-delivery"
)

type Order struct {
	ID        uuid.UUID   `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Title     string      `gorm:"type:varchar(255);not null"`
	Price     float64     `gorm:"type:decimal(10,2);not null"`
	Status    OrderStatus `gorm:"type:varchar(32);not null"`
	CreatedAt time.Time   `gorm:"not null;autoCreateTime"`

	BuyerID uuid.UUID `gorm:"type:uuid;not null"`
	Buyer   User      `gorm:"foreignKey:BuyerID"`

	SellerID uuid.UUID `gorm:"type:uuid;not null"`
	Seller   User      `gorm:"foreignKey:SellerID"`

	Deadline *time.Time `gorm:"default:null"`
}

type Chat struct {
	ID      uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	OrderID uuid.UUID `gorm:"type:uuid;not null"`
	Order   Order     `gorm:"foreignKey:OrderID"`

	BuyerID uuid.UUID `gorm:"type:uuid;not null"`
	Buyer   User      `gorm:"foreignKey:BuyerID"`

	SellerID uuid.UUID `gorm:"type:uuid;not null"`
	Seller   User      `gorm:"foreignKey:SellerID"`

	LastMessageID *uuid.UUID `gorm:"type:uuid;default:null"`
	CreatedAt     time.Time  `gorm:"not null;autoCreateTime"`
}

//Swanhtet12@
