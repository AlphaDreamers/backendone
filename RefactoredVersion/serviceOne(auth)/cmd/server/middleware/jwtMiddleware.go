package middleware

import (
	"fmt"
	"github.com/SwanHtetAungPhyo/service-one/auth/internal/response"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"strings"
)

func Check() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Get(fiber.HeaderAuthorization)
		if token == "" {
			return c.JSON(response.Response{
				Status:  fiber.StatusUnauthorized,
				Message: "no token found",
			})
		}
		parts := strings.Fields(token)
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.JSON(response.Response{
				Status:  fiber.StatusUnauthorized,
				Message: "malformed token",
			})
		}

		tokenString := parts[1]

		valid, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(viper.GetString("jwt_secret")), nil
		})
		if err != nil || !valid.Valid {
			return c.JSON(response.Response{
				Status:  fiber.StatusUnauthorized,
				Message: "invalid token",
			})
		}
		claims, ok := valid.Claims.(jwt.MapClaims)
		if !ok {
			return c.JSON(response.Response{
				Status:  fiber.StatusUnauthorized,
				Message: "invalid token",
			})
		}
		c.Locals("user", claims["user_id"])
		c.Locals("device_id", claims["device_id"])
		return c.Next()
	}
}
