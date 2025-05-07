package handler

import (
	"github.com/SwanHtetAungPhyo/srvc/internal/model"
	"github.com/SwanHtetAungPhyo/srvc/internal/response"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

var _ HandlerBehaviour = (*Handler)(nil)

func (h Handler) GetAllServices(c *fiber.Ctx) error {
	//TODO implement me
	panic("implement me")
}

func (h Handler) CreateService(c *fiber.Ctx) error {
	var req *model.SrvReq
	userRaw := c.Params("userId")
	if userRaw == "" {
		return c.Status(fiber.StatusBadRequest).JSON(&response.Response{
			Message: "userId is required in the param",
		})
	}
	userUuid := uuid.MustParse(userRaw)
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&response.Response{
			Message: err.Error(),
		})
	}

	createdPost, err := h.service.CreateService(h.ctx, req, userUuid)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&response.Response{
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(&response.Response{
		Message: "success to create service by" + userUuid.String(),
		Data:    createdPost,
	})

}

func (h Handler) UpdateService(c *fiber.Ctx) error {
	//TODO implement me
	panic("implement me")
}

func (h Handler) DeleteService(c *fiber.Ctx) error {
	//TODO implement me
	panic("implement me")
}

func (h Handler) GetServiceByServiceId(c *fiber.Ctx) error {
	//TODO implement me
	panic("implement me")
}

func (h Handler) GetSpcServiceByUserId(c *fiber.Ctx) error {
	//TODO implement me
	panic("implement me")
}

func (h Handler) GetAllUserServices(c *fiber.Ctx) error {
	//TODO implement me
	panic("implement me")
}
