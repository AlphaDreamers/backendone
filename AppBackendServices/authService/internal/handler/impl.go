package handler

import (
	"context"
	"fmt"
	"github.com/SwanHtetAungPhyo/common/models"
	issuer "github.com/SwanHtetAungPhyo/common/pkg/jwt"
	"github.com/SwanHtetAungPhyo/common/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
	"math/rand"
	"strconv"
	"time"
)

func (i *Impl) Login(c *fiber.Ctx) error {
	var loginRequest *models.LoginRequest
	if err := c.BodyParser(&loginRequest); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error(), nil)
	}
	if err := loginRequest.Validate(); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	err := i.service.Login(loginRequest)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error(), nil)
	}
	userFromDb, err := i.service.GetUserByEmail(loginRequest.Email)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error(), nil)
	}
	jwtSecret := viper.GetString("jwt_secret")
	refreshToken, _ := issuer.JwtIssuer([]byte(jwtSecret), "AuthService", loginRequest.Email, "Swan", "refresh")
	accessToken, _ := issuer.JwtIssuer([]byte(jwtSecret), "AuthService", loginRequest.Email, "Swan", "access")
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Strict",
	})
	type loginResponse struct {
		AccessToken       string `json:"access_token"`
		UserAccountWallet bool   `json:"user_account_wallet"`
	}
	c.Locals("email", loginRequest.Email)

	return utils.SuccessResponse(c, fiber.StatusOK, "success", &loginResponse{
		AccessToken:       accessToken,
		UserAccountWallet: userFromDb.WalletCreated,
	})

}

func (i *Impl) Register(c *fiber.Ctx) error {
	var user *models.UserRegisterRequest
	if err := c.BodyParser(&user); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	err := i.service.Register(user)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	codeString := i.generateCode()
	go func() {
		i.redisClient.Set(context.Background(), user.Email, codeString, time.Minute*10)
	}()
	go i.SendEmail(user.FullName, user.Email, codeString)
	return utils.SuccessResponse(c, fiber.StatusCreated, "success", nil)
}

func (i *Impl) Me(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"email": "swanhtetam@gmail.com",
		"id":    "SwanHtetam",
	})

}

func (i *Impl) Refresh(c *fiber.Ctx) error {
	refreshToken := c.Cookies("refresh_token")
	if refreshToken == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "No refresh token found"})
	}

	jwtSecret := viper.GetString("jwt_secret")
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(jwtSecret), nil
	})

	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid token claims"})
	}

	// Check the user_id and other claims
	userID, ok := claims["user_id"].(string)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user_id"})
	}

	email, ok := claims["email"].(string)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid email"})
	}

	if claims["type"] != "refresh" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid refresh token type"})
	}

	newAccessToken, err := issuer.JwtIssuer([]byte(jwtSecret), "myapp", userID, email, "access")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate new access token"})
	}

	return c.JSON(fiber.Map{"access_token": newAccessToken})
}

func (i *Impl) Verify(c *fiber.Ctx) error {
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
	err = i.service.UpdateStatus(verifyRequest.Email)
	if err != nil {
		i.logger.Warn(err.Error())
	}
	result, err := i.redisClient.Del(context.Background(), verifyRequest.Email).Result()
	if err != nil {
		i.logger.Error(err.Error())
	}
	i.logger.Info(result)
	i.logger.Info(fmt.Sprintf("Successfully verified %s", verifyRequest.Email))
	return utils.SuccessResponse(c, fiber.StatusOK, "success", nil)

}

func (i *Impl) SendCode(c *fiber.Ctx) error {
	email := c.Params("email")
	username := c.Params("username")
	if email == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Empty Email in query", nil)
	}
	codeString := i.generateCode()
	go func() {
		i.redisClient.Set(context.Background(), email, codeString, time.Minute*10)
	}()
	go i.SendEmail(username, email, codeString)
	return utils.SuccessResponse(c, fiber.StatusOK, "success", nil)
}

func (i *Impl) StoreInVault(c *fiber.Ctx) error {
	type Meneoinc struct {
		Userid  string   `json:"userid"`
		Phrases []string `json:"phrase"`
	}
	var payload Meneoinc
	if err := c.BodyParser(&payload); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error(), nil)
	}
	err := i.service.InteractionWithVault(payload.Userid, payload.Phrases)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error(), nil)
	}
	return utils.SuccessResponse(c, fiber.StatusOK, "success", nil)
}
func (i *Impl) generateCode() string {
	rand.Seed(time.Now().UnixNano())
	code := rand.Intn(900000) + 100000
	codeString := strconv.Itoa(code)
	return codeString
}
