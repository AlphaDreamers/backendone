package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func ControlFilter() fiber.Handler {
	return func(c *fiber.Ctx) error {
		headers := c.GetReqHeaders()

		logrus.Infof("Received headers: %v", headers)

		if userAgents, ok := headers["User-Agent"]; ok && len(userAgents) > 0 && userAgents[0] == "Consul Health Check" {
			return c.Next()
		}

		if requestIDs, ok := headers["X-Request-Id"]; ok && len(requestIDs) > 0 {
			if requestIDs[0] == "app-gateway" {
				logrus.Infof("X-Request-ID: %v", requestIDs[0])
				return c.Next()
			}
			logrus.Error("Invalid X-Request-ID header")
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"Message": "Invalid Request ID",
			})
		}

		logrus.Error("Missing X-Request-ID header")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"Message": "No Request ID",
			"Details": "Untrusted Request ID",
		})
	}
}
