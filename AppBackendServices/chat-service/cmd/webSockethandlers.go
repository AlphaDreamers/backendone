package cmd

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"githubc.com/SwanHtetAungPhyo/chat-service/internal/handler"
	"githubc.com/SwanHtetAungPhyo/chat-service/internal/repo"
	"githubc.com/SwanHtetAungPhyo/chat-service/internal/service"
)

func RouteSettingUp(app *fiber.App) {
	conversationGroup := app.Group("/conversation")
	logger := logrus.New()
	memDB := repo.NewConversationMemDB(logger)
	repository := repo.NewDummyRepo(logger, memDB)
	conversationService := service.NewConversationService(logger, repository)
	conversationHandler := handler.NewStartConversationHandlerImpl(logger, conversationService)

	conversationGroup.Post(
		"/", conversationHandler.StartConversation,
	)
}
