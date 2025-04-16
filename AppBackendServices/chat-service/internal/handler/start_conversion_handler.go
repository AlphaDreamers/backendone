package handler

import (
	"github.com/SwanHtetAungPhyo/common/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"githubc.com/SwanHtetAungPhyo/chat-service/internal/model"
	"githubc.com/SwanHtetAungPhyo/chat-service/internal/service"
)

type StartConversationHanlder interface {
	StartConversation(ctx *fiber.Ctx) error
}
type StartConversationHandlerImpl struct {
	logger              *logrus.Logger
	conversationService service.ConversationService
}

func NewStartConversationHandlerImpl(logger *logrus.Logger, conversationService service.ConversationService) *StartConversationHandlerImpl {
	return &StartConversationHandlerImpl{logger: logger, conversationService: conversationService}
}

var _ StartConversationHanlder = (*StartConversationHandlerImpl)(nil)

func (s *StartConversationHandlerImpl) StartConversation(ctx *fiber.Ctx) error {
	var req *model.StartConversationRequest

	rawSenderID := ctx.Get("SenderID")

	if rawSenderID == "" {
		return utils.ErrorResponse(ctx, fiber.StatusBadRequest, "senderID missing or invalid", "Unauthorized or malformed senderID")
	}

	senderIDInUUID, err := uuid.Parse(rawSenderID)
	s.logger.Info(senderIDInUUID)
	if err := ctx.BodyParser(&req); err != nil {
		return utils.ErrorResponse(ctx, fiber.StatusBadRequest, "", "Improper body format")
	}

	conversationID, err := s.conversationService.GetOrCreateConversation(senderIDInUUID, req)
	if err != nil {
		return utils.ErrorResponse(ctx, fiber.StatusInternalServerError, err.Error(), "Failed to Start Conversation, Internal Server Error")
	}

	return utils.SuccessResponse(ctx, fiber.StatusOK, "Conversation Initialization is successful", conversationID)
}
