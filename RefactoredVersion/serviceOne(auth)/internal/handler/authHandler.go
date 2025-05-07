package handler

import (
	"context"
	"github.com/SwanHtetAungPhyo/service-one/auth/internal/model"
	"github.com/SwanHtetAungPhyo/service-one/auth/internal/response"
	"github.com/SwanHtetAungPhyo/service-one/auth/internal/services"
	"github.com/SwanHtetAungPhyo/service-one/auth/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
	"net/http"
	"time"
)

var ProviderModule = fx.Module("auth_handler",
	fx.Provide(NewHandler,
		NewUserHandler,
	),
)

type AuthHandlerBehaviour interface {
	Login(c *fiber.Ctx) error
	Logout(c *fiber.Ctx) error
	Register(c *fiber.Ctx) error
	AccRegisterEmailVerification(c *fiber.Ctx) error
	ResetPassword(c *fiber.Ctx) error
	ResetPasswordTokenVerify(c *fiber.Ctx) error
	ForgotPassword(c *fiber.Ctx) error
	ForgotPasswordTokenVerified(c *fiber.Ctx) error
	RefreshToken(c *fiber.Ctx) error
}

type Handler struct {
	log         *logrus.Logger
	srv         *services.AuthService
	redisClient *redis.Client
	ctx         context.Context
}

func NewHandler(
	srv *services.AuthService,
	redisClient *redis.Client,
	log *logrus.Logger,

) *Handler {
	return &Handler{
		srv:         srv,
		redisClient: redisClient,
		log:         log,
		ctx:         context.Background(),
	}
}
func (h *Handler) Login(c *fiber.Ctx) error {
	var req model.LoginRequest

	deviceID, ok := c.Locals("device_id").(string)
	if !ok || deviceID == "" {
		h.log.Error("Device ID missing in context")
		return c.Status(fiber.StatusInternalServerError).JSON(response.Response{
			Status:  fiber.StatusInternalServerError,
			Message: "Server configuration error",
		})
	}

	if err := c.BodyParser(&req); err != nil {
		h.log.WithError(err).Error("Login request parsing failed")
		return c.Status(fiber.StatusBadRequest).JSON(response.Response{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid request format",
			Data:    fiber.Map{"error": err.Error()},
		})
	}

	if req.Email == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(response.Response{
			Status:  fiber.StatusBadRequest,
			Message: "Email and password are required",
		})
	}

	loginResponse, err := h.srv.Login(h.ctx, &req, deviceID)
	if err != nil {
		h.log.WithError(err).Error("Login failed")
		status := fiber.StatusUnauthorized
		message := "Invalid credentials"

		return c.Status(status).JSON(response.Response{
			Status:  status,
			Message: message,
		})
	}
	c.Locals("userID", loginResponse.UserId)
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    loginResponse.RefreshToken,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Lax",
	})
	c.Cookie(&fiber.Cookie{
		Name:     "user_id",
		Value:    loginResponse.UserId.String(),
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
		Secure:   true,
	})
	return c.Status(fiber.StatusOK).JSON(response.Response{
		Status:  fiber.StatusOK,
		Message: "Login successful",
		Data: fiber.Map{
			"userId":       loginResponse.UserId,
			"access_token": loginResponse.AccessToken,
			"expires_in":   loginResponse.ExpiresIn,
		},
	})
}
func (h *Handler) Logout(c *fiber.Ctx) error {
	accessToken := c.Get("Authorization")
	if accessToken == "" {
		return c.Status(http.StatusBadRequest).JSON(response.Response{
			Status:  http.StatusBadRequest,
			Message: "Authorization header missing",
		})
	}

	refreshToken := c.Cookies("refresh_token")
	if refreshToken == "" {
		return c.Status(http.StatusBadRequest).JSON(response.Response{
			Status:  http.StatusBadRequest,
			Message: "Refresh token missing",
		})
	}

	err := h.srv.Logout(h.ctx, accessToken, refreshToken)
	if err != nil {
		h.log.WithError(err).Error("Logout failed")
		return c.Status(http.StatusInternalServerError).JSON(response.Response{
			Status:  http.StatusInternalServerError,
			Message: "Logout failed",
		})
	}
	c.ClearCookie("refresh_token")

	return c.Status(http.StatusOK).JSON(response.Response{
		Status:  http.StatusOK,
		Message: "Logout successful",
	})
}

func (h *Handler) Register(c *fiber.Ctx) error {
	var req model.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		h.log.WithError(err).Error("Register request parsing failed")
		return c.Status(http.StatusBadRequest).JSON(response.Response{
			Status:  http.StatusBadRequest,
			Message: "Invalid request format",
		})
	}

	verificationToken, err := h.srv.Register(h.ctx, &req)
	if err != nil {
		h.log.WithError(err).Error("Registration failed")
		return c.Status(http.StatusBadRequest).JSON(response.Response{
			Status:  http.StatusBadRequest,
			Message: "Registration failed: " + err.Error(),
		})
	}
	utils.SendEmail(*verificationToken, req.Email)
	return c.Status(http.StatusCreated).JSON(response.Response{
		Status:  http.StatusCreated,
		Message: "Registration successful. Verification email sent",
	})
}
func (h *Handler) AccRegisterEmailVerification(c *fiber.Ctx) error {
	var req model.RegisterVerificationRequest
	if err := c.BodyParser(&req); err != nil {
		h.log.WithError(err).Error("Email verification request parsing failed")
		return c.Status(http.StatusBadRequest).JSON(response.Response{
			Status:  http.StatusBadRequest,
			Message: "Invalid request format",
		})
	}

	err := h.srv.VerifyEmail(h.ctx, req.Email, req.Token)
	if err != nil {
		h.log.WithError(err).Error("Email verification failed")
		return c.Status(http.StatusBadRequest).JSON(response.Response{
			Status:  http.StatusBadRequest,
			Message: "Email verification failed: " + err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(response.Response{
		Status:  http.StatusOK,
		Message: "Email verified successfully",
	})
}

func (h *Handler) ResetPassword(c *fiber.Ctx) error {
	var req model.PasswordResetRequest
	if err := c.BodyParser(&req); err != nil {
		h.log.WithError(err).Error("Password reset request parsing failed")
		return c.Status(http.StatusBadRequest).JSON(response.Response{
			Status:  http.StatusBadRequest,
			Message: "Invalid request format",
		})
	}

	err := h.srv.ResetPassword(h.ctx, req.Email)
	if err != nil {
		h.log.WithError(err).Error("Password reset failed")
		return c.Status(http.StatusBadRequest).JSON(response.Response{
			Status:  http.StatusBadRequest,
			Message: "Password reset failed: " + err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(response.Response{
		Status:  http.StatusOK,
		Message: "Password reset email sent",
	})
}

func (h *Handler) ResetPasswordTokenVerify(c *fiber.Ctx) error {
	var req model.PasswordResetCodeVerificationRequest
	if err := c.BodyParser(&req); err != nil {
		h.log.WithError(err).Error("Password reset verification parsing failed")
		return c.Status(http.StatusBadRequest).JSON(response.Response{
			Status:  http.StatusBadRequest,
			Message: "Invalid request format",
		})
	}

	err := h.srv.VerifyResetPasswordToken(h.ctx, req.Token, req.Email)
	if err != nil {
		h.log.WithError(err).Error("Password reset token verification failed")
		return c.Status(http.StatusBadRequest).JSON(response.Response{
			Status:  http.StatusBadRequest,
			Message: "Token verification failed: " + err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(response.Response{
		Status:  http.StatusOK,
		Message: "Token verified successfully",
	})
}

func (h *Handler) ForgotPassword(c *fiber.Ctx) error {
	var req model.ForgotPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		h.log.WithError(err).Error("Forgot password request parsing failed")
		return c.Status(http.StatusBadRequest).JSON(response.Response{
			Status:  http.StatusBadRequest,
			Message: "Invalid request format",
		})
	}

	err := h.srv.ForgotPassword(h.ctx, req.Email)
	if err != nil {
		h.log.WithError(err).Error("Forgot password failed")
		return c.Status(http.StatusBadRequest).JSON(response.Response{
			Status:  http.StatusBadRequest,
			Message: "Forgot password failed: " + err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(response.Response{
		Status:  http.StatusOK,
		Message: "Forgot password email sent",
	})
}

func (h *Handler) ForgotPasswordTokenVerified(c *fiber.Ctx) error {
	var req model.ForgotPasswordCodeVerificationRequest
	if err := c.BodyParser(&req); err != nil {
		h.log.WithError(err).Error("Forgot password verification parsing failed")
		return c.Status(http.StatusBadRequest).JSON(response.Response{
			Status:  http.StatusBadRequest,
			Message: "Invalid request format",
		})
	}

	err := h.srv.VerifyForgotPasswordToken(h.ctx, req.Token, req.Email, req.NewPassword)
	if err != nil {
		h.log.WithError(err).Error("Forgot password token verification failed")
		return c.Status(http.StatusBadRequest).JSON(response.Response{
			Status:  http.StatusBadRequest,
			Message: "Password update failed: " + err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(response.Response{
		Status:  http.StatusOK,
		Message: "Password updated successfully",
	})
}

func (h *Handler) RefreshToken(c *fiber.Ctx) error {
	refreshToken := c.Cookies("refresh_token")
	if refreshToken == "" {
		return c.Status(http.StatusBadRequest).JSON(response.Response{
			Status:  http.StatusBadRequest,
			Message: "Refresh token missing",
		})
	}

	newTokens, err := h.srv.RefreshToken(h.ctx, refreshToken)
	if err != nil {
		h.log.WithError(err).Error("Token refresh failed")
		return c.Status(http.StatusUnauthorized).JSON(response.Response{
			Status:  http.StatusUnauthorized,
			Message: "Token refresh failed: " + err.Error(),
		})
	}

	// Set new refresh token in cookie
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    newTokens.RefreshToken,
		Expires:  time.Now().Add(7 * 24 * time.Hour), // 7 days
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Lax",
	})

	return c.Status(http.StatusOK).JSON(response.Response{
		Status:  http.StatusOK,
		Message: "Token refreshed successfully",
		Data:    fiber.Map{"access_token": newTokens.AccessToken},
	})
}

var _ AuthHandlerBehaviour = (*Handler)(nil)
