package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func LoggingMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": err.Error(),
			})
		}
		latency := time.Since(start)

		clientIP := c.IP()
		origin := c.Get("Origin")
		referrer := c.Get("Referer")
		statusCode := c.Response().StatusCode()
		logrus.Printf("[%d] %s %s called from %s (%s | %s) - %v", statusCode, c.Method(), c.Path(), clientIP, origin, referrer, latency)
		return nil
	}
}
