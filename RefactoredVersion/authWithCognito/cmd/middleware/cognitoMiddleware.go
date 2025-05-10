package middleware

import (
	"encoding/json"
	"github.com/MicahParks/keyfunc/v3"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
	"io"
	"net/http"
	"strings"
	"time"
)

var MiddlewareProvider = fx.Module("middleware",
	fx.Provide(
		NewAuthMiddleware,
	),
)

type AuthMiddleware struct {
	jwksURL       string
	log           *logrus.Logger
	JWKSjson      json.RawMessage
	httpClient    *http.Client
	refreshTicker *time.Ticker
}

func NewAuthMiddleware(
	log *logrus.Logger,
	jwksURL string) *AuthMiddleware {

	authMiddleware := &AuthMiddleware{
		log:     log,
		jwksURL: jwksURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		refreshTicker: time.NewTicker(time.Hour),
	}
	authMiddleware.getJwks()
	go authMiddleware.BackgroundCheck()
	return authMiddleware

}
func (a *AuthMiddleware) BackgroundCheck() {
	for _ = range a.refreshTicker.C {
		a.getJwks()
	}
}

func (a *AuthMiddleware) getJwks() {
	get, err := a.httpClient.Get(a.jwksURL)
	if err != nil {
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(get.Body)
	err = json.NewDecoder(get.Body).Decode(&a.JWKSjson)
	if err != nil {
		a.log.Error(err.Error())
		return
	}
	return
}
func (a *AuthMiddleware) Cognito() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing or invalid Authorization header"})
		}
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		keyFUNC, err := keyfunc.NewJWKJSON(a.JWKSjson)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		token, err := jwt.Parse(tokenStr, keyFUNC.Keyfunc)
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token"})
		}

		c.Locals("user", token.Claims)
		return c.Next()
	}
}
