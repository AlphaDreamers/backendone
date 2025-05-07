package handler

import (
	"context"
	"github.com/SwanHtetAungPhyo/service-one/auth/internal/response"
	"github.com/SwanHtetAungPhyo/service-one/auth/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"strconv"
	"time"
)

type UserHandler struct {
	log     *logrus.Logger
	service *services.UserService
}

func NewUserHandler(log *logrus.Logger, service *services.UserService) *UserHandler {
	return &UserHandler{
		log:     log,
		service: service,
	}
}

func (uh *UserHandler) GetMet(c *fiber.Ctx) error {
	userIdParam := c.Params("userId")
	if userIdParam == "" {
		return c.JSON(&response.Response{
			Status:  fiber.StatusBadRequest,
			Message: "userId is required",
		})
	}
	userUuid := uuid.MustParse(userIdParam)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	userInfo, err := uh.service.GetMeInfo(ctx, userUuid)
	if err != nil {
		return c.JSON(&response.Response{
			Status:  fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}
	return c.JSON(&response.Response{
		Status:  fiber.StatusOK,
		Message: "success",
		Data:    userInfo,
	})
}
func (uh *UserHandler) GetDashboard(c *fiber.Ctx) error {
	userIDstring := c.Params("userId")
	userID := uuid.MustParse(userIDstring)
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	data, err := uh.service.GetDashboardData(page, pageSize, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to load dashboard data",
		})
	}

	return c.JSON(data)
}
