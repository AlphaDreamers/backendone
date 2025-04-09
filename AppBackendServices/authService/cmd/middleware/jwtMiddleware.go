package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"strings"
	"time"
)

func JwtMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			logrus.Warn("Missing Authorization header")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Authorization header required",
			})
		}

		authParts := strings.Split(authHeader, " ")
		if len(authParts) != 2 || strings.ToLower(authParts[0]) != "bearer" {
			logrus.Warn("Invalid Authorization header format")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token format. Use Bearer {token}",
			})
		}

		tokenString := authParts[1]
		jwtSecret := viper.GetString("jwt_secret")
		if jwtSecret == "" {
			logrus.Error("JWT_SECRET environment variable not set")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Server configuration error",
			})
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(jwtSecret), nil
		})

		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Warn("Token parsing failed")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token",
			})
		}

		if !token.Valid {
			logrus.Warn("Invalid token received")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token",
			})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			logrus.Warn("Invalid token claims format")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token claims",
			})
		}

		exp, err := claims.GetExpirationTime()
		if err != nil || exp == nil {
			logrus.Warn("Token missing expiration")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token expiration",
			})
		}

		if time.Now().After(exp.Time) {
			logrus.Warn("Expired token attempt")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Token expired",
			})
		}

		if _, ok := claims["id"]; !ok {
			logrus.Warn("Token missing id claim")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token claims",
			})
		}

		if _, ok := claims["email"]; !ok {
			logrus.Warn("Token missing email claim")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token claims",
			})
		}

		// Set user context
		c.Locals("user", claims["id"])
		c.Locals("email", claims["email"])

		return c.Next()
	}
}
