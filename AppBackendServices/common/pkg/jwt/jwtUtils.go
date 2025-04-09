package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

// CustomClaims Use RegisteredClaims instead of StandardClaims
type CustomClaims struct {
	UserID string `json:"id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func JwtIssuer(secret []byte, appName string, userID, email, tokenType string) (string, error) {

	// Create claims with proper structure
	claims := CustomClaims{
		UserID: userID,
		Email:  email,
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
