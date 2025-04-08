package utils

import "github.com/gofiber/fiber/v2"

type Response struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func NewResponse(message string, data interface{}) *Response {
	return &Response{Message: message, Data: data}
}

func ErrorResponse(ctx *fiber.Ctx, status int, message string, data interface{}) error {
	ctx.Status(status)
	return ctx.JSON(&Response{Message: message, Data: data})
}

func SuccessResponse(ctx *fiber.Ctx, status int, message string, data interface{}) error {
	ctx.Status(status)
	return ctx.JSON(&Response{Message: message, Data: data})
}
