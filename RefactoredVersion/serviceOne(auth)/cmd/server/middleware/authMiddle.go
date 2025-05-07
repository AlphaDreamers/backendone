package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var userID uuid.UUID

		// 1. Check JWT token if using
		// 2. Check cookie
		if cookie := c.Cookies("user_id"); cookie != "" {
			if id, err := uuid.Parse(cookie); err == nil {
				userID = id
			}
		}

		// 3. Check locals (from login)
		if localsID, ok := c.Locals("userID").(uuid.UUID); ok {
			userID = localsID
		}

		if userID == uuid.Nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "unauthorized",
				"success": false,
			})
		}

		c.Locals("userID", userID)
		return c.Next()
	}
}
