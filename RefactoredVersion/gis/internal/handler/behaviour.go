package handler

import "github.com/gofiber/fiber/v2"

type Behaviour interface {
	Create(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
	GetById(c *fiber.Ctx) error
	List(c *fiber.Ctx) error
	GetByUserId(c *fiber.Ctx) error
}
