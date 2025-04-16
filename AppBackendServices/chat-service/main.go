package main

import (
	"github.com/SwanHtetAungPhyo/common/pkg/logutil"
	"github.com/SwanHtetAungPhyo/common/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"githubc.com/SwanHtetAungPhyo/chat-service/cmd"
	"githubc.com/SwanHtetAungPhyo/chat-service/internal/handler"
	"log"
)

func main() {
	logutil.InitLog("chat-service")
	app := fiber.New()

	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			userID := c.Query("userId")
			if userID == "" {
				logutil.GetLogger().Info("UserId is empty")
				return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized", "userId is required as query parameter")
			}
			c.Locals("userId", userID)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	cmd.RouteSettingUp(app)
	err := handler.RabbitMqInit()
	if err != nil {
		return
	}
	webSocketConversationHandler := handler.NewWebSocketConversationHandler(logutil.GetLogger())

	go webSocketConversationHandler.StartHub()
	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		webSocketConversationHandler.MainHandler(c)
	}))

	log.Fatal(app.Listen(":3000"))
}
