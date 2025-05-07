package handler

import (
	"github.com/gofiber/fiber/v2"
)

type HandlerBehaviour interface {
	GetAllServices(c *fiber.Ctx) error
	CreateService(c *fiber.Ctx) error
	UpdateService(c *fiber.Ctx) error
	DeleteService(c *fiber.Ctx) error
	GetServiceByServiceId(c *fiber.Ctx) error
	GetSpcServiceByUserId(c *fiber.Ctx) error
	GetAllUserServices(c *fiber.Ctx) error
}
