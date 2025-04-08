package jwt

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"time"
)

type Claims struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.StandardClaims
}

func JwtIssuer(secret []byte, appName string, userID, email, tokenType string) (string, error) {
	var expirationTime time.Time

	switch tokenType {
	case "access":
		expirationTime = time.Now().Add(1 * time.Hour)
	case "refresh":
		expirationTime = time.Now().Add(30 * 24 * time.Hour)
	default:
		return "", errors.New("invalid token type")
	}

	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"type":    tokenType,
		"iss":     appName,
		"exp":     expirationTime.Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
