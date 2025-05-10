package auth

import (
	"github.com/SwanHtetAungPhyo/authCognito/internal/model"
	"github.com/gofiber/fiber/v2"
	"time"
)

var _ Behaviour = (*ConcreteHandler)(nil)

func (c2 ConcreteHandler) SignUp(c *fiber.Ctx) error {
	var req *model.UserSignUpRequest
	if err := c.BodyParser(&req); err != nil {
		return c.JSON(model.Response{
			Message: err.Error(),
		})
	}
	err := c2.srv.SignUp(req)
	if err != nil {
		return c.JSON(model.Response{
			Message: err.Error(),
		})
	}
	return c.JSON(model.Response{
		Message: "OK",
		Data:    req,
	})
}

func (c2 ConcreteHandler) SignIn(c *fiber.Ctx) error {
	var req *model.UserSignInReq
	if err := c.BodyParser(&req); err != nil {
		return c.JSON(model.Response{
			Message: err.Error(),
		})
	}
	userData, respFromC, err := c2.srv.SignIn(req)
	if err != nil {
		return c.JSON(model.Response{
			Message: err.Error(),
		})
	}
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    *respFromC.AuthenticationResult.RefreshToken,
		Secure:   true,
		HTTPOnly: true,
		MaxAge:   time.Now().Add(time.Hour * 24 * 365 * 10).Minute(),
	})
	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    *respFromC.AuthenticationResult.AccessToken,
		Secure:   true,
		HTTPOnly: true,
		MaxAge:   time.Now().Add(time.Minute * 30).Minute(),
	})
	return c.JSON(model.Response{
		Message: "OK",
		Data:    userData,
	})
}

func (c2 ConcreteHandler) Confirm(c *fiber.Ctx) error {
	var req *model.EmailVerificationRequest
	if err := c.BodyParser(&req); err != nil {
		return c.JSON(model.Response{
			Message: err.Error(),
		})
	}
	err := c2.srv.Confirm(req)
	if err != nil {
		return c.JSON(model.Response{
			Message: err.Error(),
		})
	}
	return c.JSON(model.Response{
		Message: "OK",
	})
}

func (c2 ConcreteHandler) ResendConfirmation(c *fiber.Ctx) error {
	email := c.Params("email")
	if email == "" {
		return c.JSON(model.Response{
			Message: "Please provide a valid email address",
		})
	}
	err := c2.srv.ResendConfirmation(email)
	if err != nil {
		return c.JSON(model.Response{
			Message: err.Error(),
		})
	}
	return c.JSON(model.Response{
		Message: "Confirm code is resend",
	})
}

func (h *ConcreteHandler) ForgotPassword(c *fiber.Ctx) error {
	var req model.ForgotPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request data",
		})
	}
	err := h.srv.ForgotPassword(req.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Password reset instructions sent to email",
	})
}

// ResetPasswordConfirm handler to confirm the new password
func (h *ConcreteHandler) ResetPasswordConfirm(c *fiber.Ctx) error {
	var req model.ResetPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request data",
		})
	}

	err := h.srv.ResetPasswordConfirm(req.Email, req.Code, req.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Password successfully reset",
	})
}

// Logout handler to log the user out
func (h *ConcreteHandler) Logout(c *fiber.Ctx) error {
	// Get the access token from the Authorization header or from the request body
	accessToken := c.Get("Authorization")
	if accessToken == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Missing access token",
		})
	}

	// Call the Logout function from the Auth service
	err := h.srv.Logout(accessToken)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Logged out successfully",
	})
}

func (h *ConcreteHandler) KYCVerify(c *fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	files := form.File["files"]
	if len(files) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No files to verify",
		})
	}

	verification, err := h.srv.KYCVerification(files)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.Response{
			Message: err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(model.Response{
		Message: "KYC verification success",
		Data:    verification,
	})
}
