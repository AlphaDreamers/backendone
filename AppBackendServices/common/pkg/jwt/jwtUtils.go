package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

// CustomClaims Use RegisteredClaims instead of StandardClaims
type CustomClaims struct {
	UserID     string `json:"id"`
	Email      string `json:"email"`
	DeviceInfo Device `json:"device_info"`
	jwt.RegisteredClaims
}
type Device struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	UserID     uuid.UUID      `json:"user_id" gorm:"index;foreignKey:UserID"`
	Browser    string         `json:"browser" gorm:"index"`
	OS         string         `json:"os" gorm:"index"`
	DeviceType string         `json:"device_type" gorm:"index"`
	FirstLogin time.Time      `json:"first_login" gorm:"type:timestamp"`
	LastLogin  time.Time      `json:"last_login" gorm:"type:timestamp"`
	DetectedAt time.Time      `json:"detected_at" gorm:"autoCreateTime"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

func JwtIssuer(secret []byte, appName string, userID, email, tokenType string, device Device) (string, error) {
	claims := CustomClaims{
		UserID: userID,
		Email:  email,
		DeviceInfo: Device{
			ID:         device.ID,
			UserID:     device.UserID,
			Browser:    device.Browser,
			OS:         device.OS,
			DeviceType: device.DeviceType,
			FirstLogin: device.FirstLogin,
			LastLogin:  device.LastLogin,
			DetectedAt: device.DetectedAt,
			DeletedAt:  device.DeletedAt,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    appName,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(getExpiration(tokenType))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func getExpiration(tokenType string) time.Duration {
	switch tokenType {
	case "access":
		return 2 * time.Minute
	case "refresh":
		return 30 * 24 * time.Hour
	default:
		return 0
	}
}
