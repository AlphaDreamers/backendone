package middleware

import (
	"github.com/SwanHtetAungPhyo/service-one/auth/internal/response"
	"github.com/gofiber/fiber/v2"
)

func LoginMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		deviceId := c.Get("X-DeviceId")
		if deviceId == "" {
			return c.Status(fiber.StatusBadRequest).JSON(response.Response{
				Status:  fiber.StatusBadRequest,
				Message: "X-DeviceId header is required",
			})
		}

		c.Locals("device_id", deviceId)
		return c.Next()
	}
}
