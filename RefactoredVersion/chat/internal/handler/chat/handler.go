package chat

import (
	"context"
	"github.com/SwanHtetAungPhyo/chat-order/internal/model/response"
	"github.com/SwanHtetAungPhyo/chat-order/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type ChatRestHanlder struct {
	log     *logrus.Logger
	srv     *service.ChatService
	context context.Context
}

func NewChatRestHanlder(log *logrus.Logger,
	srv *service.ChatService,
) *ChatRestHanlder {
	return &ChatRestHanlder{
		log:     log,
		srv:     srv,
		context: context.Background(),
	}
}

func (h ChatRestHanlder) GetChatRoomByOrderId(ctx *fiber.Ctx) error {
	orderIdRaw := ctx.Params("orderId")
	if orderIdRaw == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.Response{
			Message: "order id  in param can not be empty",
		})
	}
	orderId := uuid.MustParse(orderIdRaw)
	h.srv.GetChatRoomByOrderId(h.context, orderId)
	return ctx.Status(fiber.StatusOK).JSON(response.Response{
		Message: "Get room by order id is successful",
	})
}

func (h ChatRestHanlder) GetAllChatRoomByUserId(ctx *fiber.Ctx) error {
	userIdRaw := ctx.Params("userId")
	if userIdRaw == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.Response{
			Message: "user id  in param can not be empty",
		})
	}
	userId := uuid.MustParse(userIdRaw)
	chatRooms, err := h.srv.GetAllChatRoomByUserId(h.context, userId)
	if err != nil {
		h.log.Error(err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.Response{
			Message: err.Error(),
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(response.Response{
		Message: "Get room by user id is successful",
		Data:    chatRooms,
	})

}
