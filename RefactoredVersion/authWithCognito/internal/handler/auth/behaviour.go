package auth

import "github.com/gofiber/fiber/v2"

type Behaviour interface {
	SignUp(c *fiber.Ctx) error
	SignIn(c *fiber.Ctx) error
	Confirm(c *fiber.Ctx) error
	ResendConfirmation(c *fiber.Ctx) error
	Logout(c *fiber.Ctx) error
	ForgotPassword(c *fiber.Ctx) error
	ResetPasswordConfirm(c *fiber.Ctx) error
	KYCVerify(c *fiber.Ctx) error
}
