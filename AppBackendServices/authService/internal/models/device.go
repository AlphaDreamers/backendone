package models

import (
	"time"
)

// Device represents a user's device information
type Device struct {
	ID                string    `json:"id" bson:"_id,omitempty"`
	UserID            string    `json:"user_id" bson:"user_id"`
	DeviceID          string    `json:"device_id" bson:"device_id"`
	DeviceType        string    `json:"device_type" bson:"device_type"` // mobile, tablet, desktop
	DeviceModel       string    `json:"device_model" bson:"device_model"`
	OS                string    `json:"os" bson:"os"`
	OSVersion         string    `json:"os_version" bson:"os_version"`
	Browser           string    `json:"browser" bson:"browser"`
	BrowserVersion    string    `json:"browser_version" bson:"browser_version"`
	IPAddress         string    `json:"ip_address" bson:"ip_address"`
	LastActiveAt      time.Time `json:"last_active_at" bson:"last_active_at"`
	IsActive          bool      `json:"is_active" bson:"is_active"`
	CreatedAt         time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt         time.Time `json:"updated_at" bson:"updated_at"`
	LastLoginAt       time.Time `json:"last_login_at" bson:"last_login_at"`
	LoginCount        int       `json:"login_count" bson:"login_count"`
	Location          string    `json:"location" bson:"location"`
	Timezone          string    `json:"timezone" bson:"timezone"`
	PushToken         string    `json:"push_token,omitempty" bson:"push_token,omitempty"`
	IsPushEnabled     bool      `json:"is_push_enabled" bson:"is_push_enabled"`
	IsTrusted         bool      `json:"is_trusted" bson:"is_trusted"`
	SecurityLevel     string    `json:"security_level" bson:"security_level"` // low, medium, high
	LastSecurityCheck time.Time `json:"last_security_check" bson:"last_security_check"`
}
