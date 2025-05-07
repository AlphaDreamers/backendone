package utils

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"time"
)

type JwtTokenGenerator struct {
	v   *viper.Viper
	log *logrus.Logger
}

func NewJwtTokenGenerator(v *viper.Viper, log *logrus.Logger) *JwtTokenGenerator {
	return &JwtTokenGenerator{
		v:   v,
		log: log,
	}
}

func (n *JwtTokenGenerator) GenerateJwtToken(tokenType, userId, deviceId string) string {
	claims := jwt.MapClaims{
		"user_id":   userId,
		"device_id": deviceId,
		"iat":       time.Now().Unix(),
	}
	switch tokenType {
	case "access_token":
		claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	case "refresh_token":
		claims["exp"] = time.Now().Add(time.Hour * 1).Unix()
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(n.v.GetString("secret")))
	if err != nil {
		n.log.Error(err.Error())
		return ""
	}
	return tokenString
}

func (n *JwtTokenGenerator) ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(n.v.GetString("secret")), nil
	})
}
