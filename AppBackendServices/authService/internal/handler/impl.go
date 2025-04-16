package handler

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/SwanHtetAungPhyo/auth/internal/config"
	"github.com/SwanHtetAungPhyo/auth/internal/models"
	"github.com/SwanHtetAungPhyo/auth/internal/services"
	jwtpkg "github.com/SwanHtetAungPhyo/common/pkg/jwt"
	"github.com/SwanHtetAungPhyo/common/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

// Response structures
type (
	LoginResponse struct {
		AccessToken       string `json:"access_token"`
		UserAccountWallet bool   `json:"user_account_wallet"`
		Email             string `json:"email"`
	}
	RegisterResponse struct {
		Email           string    `json:"email"`
		FullName        string    `json:"full_name"`
		VerificationTTL int       `json:"verification_ttl_minutes"`
		RegisteredAt    time.Time `json:"registered_at"`
	}
	UserProfileResponse struct {
		Email             string     `json:"email"`
		UserName          string     `json:"user_name"`
		Verified          bool       `json:"verified"`
		CreatedAt         time.Time  `json:"created_at"`
		WalletCreated     bool       `json:"wallet_created"`
		WalletCreatedTime *time.Time `json:"wallet_created_time,omitempty"`
	}
	VerificationResponse struct {
		Email         string    `json:"email"`
		VerifiedAt    time.Time `json:"verified_at"`
		AccountStatus string    `json:"account_status"`
	}
)

const (
	ResetToken     = "reset_token"
	ForgotPassword = "forgot_password"
)

// AuthHandler handles authentication-related HTTP requests
type AuthHandler struct {
	authService services.Impl
	redisClient *redis.Client
	config      *config.Config
	logger      *logrus.Logger
}

// NewAuthHandler creates a new AuthHandler instance
func NewAuthHandler(authService services.Impl, redisClient *redis.Client, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		redisClient: redisClient,
		config:      cfg,
		logger:      logrus.New(),
	}
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var loginRequest *models.LoginRequest
	if err := c.BodyParser(&loginRequest); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request format", err)
	}
	if err := loginRequest.Validate(); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Validation failed", err)
	}

	err := h.authService.Login(loginRequest)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Login failed", err.Error())
	}

	userFromDb, err := h.authService.GetUserByEmail(loginRequest.Email)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch user data", err)
	}

	device := getDeviceInfo(c)

	accessToken, refreshToken, err := h.generateToken(userFromDb, device)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to generate tokens", err)
	}

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Strict",
	})

	c.Locals("email", loginRequest.Email)

	response := &LoginResponse{
		AccessToken:       accessToken,
		UserAccountWallet: userFromDb.WalletCreated,
		Email:             loginRequest.Email,
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Login successful", response)
}

func getDeviceInfo(c *fiber.Ctx) jwtpkg.Device {
	userAgent := c.Get("User-Agent")
	browser := c.Get("Browser")
	os := c.Get("Os")
	deviceType := c.Get("Device-Type")

	if strings.Contains(strings.ToLower(userAgent), "chrome") {
		browser = "chrome"
	} else if strings.Contains(strings.ToLower(userAgent), "firefox") {
		browser = "firefox"
	} else if strings.Contains(strings.ToLower(userAgent), "safari") {
		browser = "safari"
	}

	if strings.Contains(strings.ToLower(userAgent), "windows") {
		os = "windows"
	} else if strings.Contains(strings.ToLower(userAgent), "mac") {
		os = "mac"
	} else if strings.Contains(strings.ToLower(userAgent), "linux") {
		os = "linux"
	} else if strings.Contains(strings.ToLower(userAgent), "android") {
		os = "android"
		deviceType = "mobile"
	} else if strings.Contains(strings.ToLower(userAgent), "iphone") || strings.Contains(strings.ToLower(userAgent), "ipad") {
		os = "ios"
		deviceType = "mobile"
	}

	if deviceType == "unknown" {
		if strings.Contains(strings.ToLower(userAgent), "mobile") {
			deviceType = "mobile"
		} else {
			deviceType = "desktop"
		}
	}

	now := time.Now()
	return jwtpkg.Device{
		Browser:    browser,
		OS:         os,
		DeviceType: deviceType,
		FirstLogin: now,
		LastLogin:  now,
		DetectedAt: now,
	}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var user *models.UserRegisterRequest
	if err := c.BodyParser(&user); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request format", err)
	}

	err := h.authService.Register(user)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Registration failed", err)
	}

	codeString := h.generateCode()
	go func() {
		err = h.redisClient.Set(context.Background(), user.Email, codeString, time.Minute*10).Err()
		if err != nil {
			h.logger.Warn("Failed to set email address", err.Error())
		}
	}()

	go h.SendEmail(user.FullName, user.Email, codeString)

	response := &RegisterResponse{
		Email:           user.Email,
		FullName:        user.FullName,
		VerificationTTL: 10,
		RegisteredAt:    time.Now(),
	}

	return utils.SuccessResponse(c, fiber.StatusCreated, "Registration successful", response)
}

func (h *AuthHandler) Me(c *fiber.Ctx) error {
	email := c.Query("email")
	if email == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Email is required", nil)
	}

	data, err := h.authService.GetUserByEmail(email)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch user profile", err)
	}

	response := &UserProfileResponse{
		Email:             data.Email,
		UserName:          data.FullName,
		Verified:          data.Verified,
		CreatedAt:         data.CreatedAt,
		WalletCreated:     data.WalletCreated,
		WalletCreatedTime: &data.WalletCreatedAt,
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Profile retrieved successfully", response)
}
func (i *AuthHandler) Verify(c *fiber.Ctx) error {
	type accountVerifyRequest struct {
		Email string `json:"email"`
		Code  string `json:"code"`
	}
	var verifyRequest accountVerifyRequest
	if err := c.BodyParser(&verifyRequest); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	codeInRedis, err := i.redisClient.Get(context.Background(), verifyRequest.Email).Result()
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "code is invalid", nil)
	}
	if codeInRedis != verifyRequest.Code {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid code", nil)
	}
	err = i.authService.UpdateStatus(verifyRequest.Email)
	if err != nil {
		i.logger.Warn(err.Error())
	}
	result, err := i.redisClient.Del(context.Background(), verifyRequest.Email).Result()
	if err != nil {
		i.logger.Error(err.Error())
	}
	i.logger.Info(result)
	i.logger.Info(fmt.Sprintf("Successfully verified %s", verifyRequest.Email))
	return utils.SuccessResponse(c, fiber.StatusOK, "Account is successfully verified", nil)

}

func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	refreshToken := c.Cookies("refresh_token")
	if refreshToken == "" {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Missing refresh token", nil)
	}

	refreshToken = strings.TrimPrefix(refreshToken, "Bearer ")

	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(h.config.JWT.Secret), nil
	})

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid refresh token", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid token claims", nil)
	}

	email, ok := claims["email"].(string)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid email in token", nil)
	}

	if claims["role"] != "refresh" {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid token type", nil)
	}

	user, err := h.authService.GetUserByEmail(email)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch user data", err)
	}

	device := getDeviceInfo(c)

	accessToken, refreshToken, err := h.generateToken(user, device)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to generate new tokens", err)
	}
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Strict",
	})
	return utils.SuccessResponse(c, fiber.StatusOK, "Token refreshed successfully", fiber.Map{
		"access_token": accessToken,
	})
}

func (h *AuthHandler) ForgotPassword(c *fiber.Ctx) error {
	var req models.ForgotPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request format", err)
	}

	if err := req.Validate(); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Validation failed", err)
	}

	user, err := h.authService.GetUserByEmail(req.Email)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "User not found", err)
	}

	resetToken := h.generateCode()
	expiresAt := time.Now().Add(15 * time.Minute)

	err = h.redisClient.Set(context.Background(), fmt.Sprintf("%s:%s", ForgotPassword, resetToken), user.ID, 15*time.Minute).Err()
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to store reset token", err)
	}

	h.SendEmail(user.FullName, user.Email, resetToken)

	return utils.SuccessResponse(c, fiber.StatusOK, "Password reset instructions sent to your email", fiber.Map{
		"email":      req.Email,
		"expires_at": expiresAt,
	})
}

func (h *AuthHandler) ForgotPasswordVerify(c *fiber.Ctx) error {
	var forgotPasswordRequest models.ForgotPasswordRequest
	if err := c.BodyParser(&forgotPasswordRequest); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request format", err)
	}
	if err := forgotPasswordRequest.Validate(); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Validation failed", err)
	}

	code := h.redisClient.Get(context.Background(), fmt.Sprintf("%s:%s", ForgotPassword, forgotPasswordRequest.Code)).Val()
	if !strings.EqualFold(code[14:], forgotPasswordRequest.Code) {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid code in forgot password request", fiber.Map{})
	}

	err := h.authService.ForgotPassSetAndUpdate(forgotPasswordRequest)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update forgot password", err)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Forgot password verification successfully", nil)
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	refreshToken := c.Cookies("refresh_token")
	accessToken := c.Get("Authorization")

	if len(accessToken) > 7 && strings.ToLower(accessToken[0:7]) == "bearer " {
		accessToken = accessToken[7:]
	}

	if refreshToken == "" || accessToken == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Missing authentication tokens", nil)
	}

	if err := h.blacklistTokens(accessToken, refreshToken); err != nil {
		h.logger.Error(err.Error())
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to process logout", err)
	}

	c.ClearCookie("refresh_token")

	return utils.SuccessResponse(c, fiber.StatusOK, "Successfully logged out", fiber.Map{
		"logged_out_at": time.Now(),
		"session_ended": true,
	})
}

func (h *AuthHandler) ResetPassword(c *fiber.Ctx) error {
	var resetPasswordRequest models.ResetPasswordRequest
	if err := c.BodyParser(&resetPasswordRequest); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request format", err)
	}
	if err := resetPasswordRequest.Validate(); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Validation failed", err)
	}
	resetToken := h.generateCode()
	userId := c.Locals("UserId").(string)
	go h.SendEmail(userId, resetPasswordRequest.Email, resetToken)
	err := h.redisClient.Set(context.Background(), fmt.Sprintf("%s:%s", resetToken, resetPasswordRequest.Email), resetToken, time.Minute*5).Err()
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to store reset token", err)
	}
	return utils.SuccessResponse(c, fiber.StatusOK, "Go to ur email that u entered and copy the code and procceed", nil)
}
func (h *AuthHandler) ResetPasswordVerify(c *fiber.Ctx) error {
	var resetPasswordRequest models.ResetPasswordVerificationRequest
	if err := c.BodyParser(&resetPasswordRequest); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request format", err)
	}
	if !h.getRestPasswordCodeViaRedis(resetPasswordRequest.Email, resetPasswordRequest.Token) {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid code in forgot password request", nil)
	}
	err := h.authService.UpdatePassword(resetPasswordRequest.Email, resetPasswordRequest.NewPassword)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update forgot password", err)
	}
	return utils.SuccessResponse(c, fiber.StatusOK, "Reset password verification successfull", nil)
}
func (h *AuthHandler) generateToken(user *models.UserInDB, device jwtpkg.Device) (string, string, error) {
	now := time.Now()
	accessTokenExp := now.Add(h.config.JWT.AccessTokenTTL)
	refreshTokenExp := now.Add(h.config.JWT.RefreshTokenTTL)

	accessClaims := jwt.MapClaims{
		"sub":   user.ID.String(), // This should be "sub", not "id"
		"email": user.Email,
		"role":  "access",
		"exp":   accessTokenExp.Unix(),
		"iat":   now.Unix(),
		"device": map[string]interface{}{
			"browser":     device.Browser,
			"os":          device.OS,
			"device_type": device.DeviceType,
			"first_login": device.FirstLogin,
			"last_login":  device.LastLogin,
			"detected_at": device.DetectedAt,
		},
	}

	refreshClaims := jwt.MapClaims{
		"sub":   user.ID.String(),
		"email": user.Email,
		"role":  "refresh",
		"exp":   refreshTokenExp.Unix(),
		"iat":   now.Unix(),
		"device": map[string]interface{}{
			"browser":     device.Browser,
			"os":          device.OS,
			"device_type": device.DeviceType,
			"first_login": device.FirstLogin,
			"last_login":  device.LastLogin,
			"detected_at": device.DetectedAt,
		},
	}

	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString([]byte(h.config.JWT.Secret))
	if err != nil {
		return "", "", err
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(h.config.JWT.Secret))
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (h *AuthHandler) generateCode() string {
	rand.Seed(time.Now().UnixNano())
	code := rand.Intn(900000) + 100000
	return strconv.Itoa(code)
}

func (h *AuthHandler) getRestPasswordCodeViaRedis(email, inputCode string) bool {
	code, err := h.redisClient.Get(context.Background(), fmt.Sprintf("%s:%s", ResetToken, email)).Result()
	if err != nil || code == "" {
		h.logger.Error(err.Error())
		return false
	}
	return strings.EqualFold(code[10:], inputCode)
}

func (h *AuthHandler) blacklistTokens(accessToken, refreshToken string) error {
	claims := &jwt.RegisteredClaims{}
	_, err := jwt.ParseWithClaims(accessToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(h.config.JWT.Secret), nil
	})

	if err != nil {
		return fmt.Errorf("failed to parse access token: %w", err)
	}

	ctx := context.Background()
	expiresAt := time.Until(claims.ExpiresAt.Time)

	err = h.redisClient.Set(ctx, fmt.Sprintf("blacklist:%s", accessToken), "1", expiresAt).Err()
	if err != nil {
		h.logger.Error("Failed to blacklist access token: " + err.Error())
		return fmt.Errorf("failed to blacklist access token: %w", err)
	}

	// Set refresh token in blacklist
	err = h.redisClient.Set(ctx, fmt.Sprintf("blacklist:%s", refreshToken), "1", expiresAt).Err()
	if err != nil {
		h.logger.Error("Failed to blacklist refresh token: " + err.Error())
		return fmt.Errorf("failed to blacklist refresh token: %w", err)
	}

	h.logger.Info(fmt.Sprintf("Successfully blacklisted %s and %s", accessToken, refreshToken))
	return nil
}
