package handler

import (
	"github.com/SwanHtetAungPhyo/common/pkg/utils"
	"github.com/SwanHtetAungPhyo/service-service/internal/model"
	"github.com/SwanHtetAungPhyo/service-service/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

// Impl represents the handler implementation
type Impl struct {
	service service.Service
	logger  *logrus.Logger
}

// NewHandler creates a new handler implementation
func NewHandler(service service.Service, logger *logrus.Logger) *Impl {
	return &Impl{
		service: service,
		logger:  logger,
	}
}

// GetUserServices handles the request to get all user services
func (h *Impl) GetUserServices(c *fiber.Ctx) error {
	userId := c.Params("userId")
	if userId == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "User ID is required", nil)
	}

	var pagination model.PaginationRequest
	if err := c.QueryParser(&pagination); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid pagination parameters", err)
	}

	var filter model.FilterRequest
	if err := c.QueryParser(&filter); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid filter parameters", err)
	}

	// Set default values if not provided
	if pagination.Page == 0 {
		pagination.Page = 1
	}
	if pagination.PageSize == 0 {
		pagination.PageSize = 10
	}

	result, err := h.service.GetUserServices(c.Context(), userId, pagination, filter)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get user services", err)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "User services retrieved successfully", result)
}

// CreateService handles the request to create a new service
func (h *Impl) CreateService(c *fiber.Ctx) error {
	userId := c.Params("userId")
	if userId == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "User ID is required", nil)
	}

	var serviceRequest model.ServiceRequest
	if err := c.BodyParser(&serviceRequest); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	result, err := h.service.CreateService(c.Context(), userId, serviceRequest)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create service", err)
	}

	return utils.SuccessResponse(c, fiber.StatusCreated, "Service created successfully", result)
}

// UpdateService handles the request to update a service
func (h *Impl) UpdateService(c *fiber.Ctx) error {
	userId := c.Params("userId")
	serviceId := c.Params("serviceId")
	if userId == "" || serviceId == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "User ID and Service ID are required", nil)
	}

	var serviceRequest model.ServiceRequest
	if err := c.BodyParser(&serviceRequest); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	result, err := h.service.UpdateService(c.Context(), userId, serviceId, serviceRequest)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update service", err)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Service updated successfully", result)
}

// DeleteService handles the request to delete a service
func (h *Impl) DeleteService(c *fiber.Ctx) error {
	userId := c.Params("userId")
	serviceId := c.Params("serviceId")
	if userId == "" || serviceId == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "User ID and Service ID are required", nil)
	}

	err := h.service.DeleteService(c.Context(), userId, serviceId)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete service", err)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Service deleted successfully", nil)
}

// ListAllServices handles the request to list all public services
func (h *Impl) ListAllServices(c *fiber.Ctx) error {
	var pagination model.PaginationRequest
	if err := c.QueryParser(&pagination); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid pagination parameters", err)
	}

	var filter model.FilterRequest
	if err := c.QueryParser(&filter); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid filter parameters", err)
	}

	// Set default values if not provided
	if pagination.Page == 0 {
		pagination.Page = 1
	}
	if pagination.PageSize == 0 {
		pagination.PageSize = 10
	}

	result, err := h.service.ListAllServices(c.Context(), pagination, filter)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to list services", err)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Services retrieved successfully", result)
}

// ListServicesByCategory handles the request to list services by category
func (h *Impl) ListServicesByCategory(c *fiber.Ctx) error {
	category := c.Params("category")
	if category == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Category is required", nil)
	}

	var pagination model.PaginationRequest
	if err := c.QueryParser(&pagination); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid pagination parameters", err)
	}

	// Set default values if not provided
	if pagination.Page == 0 {
		pagination.Page = 1
	}
	if pagination.PageSize == 0 {
		pagination.PageSize = 10
	}

	result, err := h.service.ListServicesByCategory(c.Context(), category, pagination)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to list services by category", err)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Services retrieved successfully", result)
}
