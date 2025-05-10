package middleware

import (
	"github.com/MicahParks/keyfunc"
	"github.com/SwanHtetAungPhyo/gis/internal/model"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"log"
	"time"
)

var jwks *keyfunc.JWKS

func InitJWKS() {
	jwtUrl := "https://cognito-idp.us-east-1.amazonaws.com/us-east-1_z6jb3eESF/.well-known/jwks.json"
	var err error
	jwks, err = keyfunc.Get(jwtUrl, keyfunc.Options{
		RefreshInterval: time.Hour,
		RefreshErrorHandler: func(err error) {
			log.Printf("Error refreshing JWKS: %v", err.Error())
		},
		RefreshUnknownKID: true,
	})
	if err != nil {
		log.Fatal(err.Error())
	}
}

func JwtMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(model.Response{
				Message: "Authorization header is empty",
			})
		}
		tokenString := authHeader[7:]
		token, err := jwt.Parse(tokenString, jwks.Keyfunc)
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(model.Response{
				Message: "Token is invalid",
				Data:    err.Error(),
			})
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(model.Response{
				Message: "Token is invalid",
			})
		}
		issuer := "https://cognito-idp.us-east-1.amazonaws.com/us-east-1_z6jb3eESF"

		if claims["iss"] != issuer {
			return c.Status(fiber.StatusUnauthorized).JSON(model.Response{
				Message: "Token is invalid. Issued by unknown issuer",
			})
		}
		if claims["client_id"] != "7qllcjjcq7p506kq88vkfiu92g" {
			return c.Status(fiber.StatusUnauthorized).JSON(model.Response{
				Message: "Token is invalid. Audience is invalid",
			})
		}
		c.Locals("claims", claims)

		return c.Next()
	}
}
