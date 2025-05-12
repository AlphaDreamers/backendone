package middleware

import (
	"log"
	"time"

	"github.com/MicahParks/keyfunc"
	"github.com/SwanHtetAungPhyo/gateways/model"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

// Middleware interface for standardized middleware implementation
type Middleware interface {
	Handler() fiber.Handler
}

// JWKSMiddleware configuration
type JWKSMiddleware struct {
	jwks            *keyfunc.JWKS
	Issuer          string
	ClientID        string
	JWKSURL         string
	RefreshInterval time.Duration
}

// NewJWKSMiddleware creates a configurable JWT middleware instance
func NewJWKSMiddleware(jwksURL, issuer, clientID string, refreshInterval time.Duration) *JWKSMiddleware {
	jwks, err := keyfunc.Get(jwksURL, keyfunc.Options{
		RefreshInterval: refreshInterval,
		RefreshErrorHandler: func(err error) {
			log.Printf("Error refreshing JWKS: %v", err)
		},
		RefreshUnknownKID: true,
	})

	if err != nil {
		log.Fatalf("Failed to initialize JWKS: %v", err)
	}

	return &JWKSMiddleware{
		jwks:            jwks,
		Issuer:          issuer,
		ClientID:        clientID,
		JWKSURL:         jwksURL,
		RefreshInterval: refreshInterval,
	}
}

// Handler implements the Middleware interface for JWKS validation
func (jm *JWKSMiddleware) Handler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return unauthorizedResponse(c, "Authorization header is empty")
		}

		tokenString := authHeader[7:]
		token, err := jwt.Parse(tokenString, jm.jwks.Keyfunc)
		if err != nil || !token.Valid {
			return unauthorizedResponse(c, "Token is invalid", err)
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return unauthorizedResponse(c, "Invalid token claims")
		}

		if claims["iss"] != jm.Issuer {
			return unauthorizedResponse(c, "Invalid token issuer")
		}

		if claims["client_id"] != jm.ClientID {
			return unauthorizedResponse(c, "Invalid client ID")
		}

		c.Locals("claims", claims)
		return c.Next()
	}
}

// Helper function for consistent error responses
func unauthorizedResponse(c *fiber.Ctx, message string, err ...error) error {
	response := model.Response{
		Message: message,
	}

	if len(err) > 0 && err[0] != nil {
		response.Data = err[0].Error()
	}

	return c.Status(fiber.StatusUnauthorized).JSON(response)
}

// MiddlewareChain combines multiple middlewares
type MiddlewareChain struct {
	middlewares []fiber.Handler
}

// NewMiddlewareChain creates a new chain of middlewares
func NewMiddlewareChain(middlewares ...Middleware) *MiddlewareChain {
	handlers := make([]fiber.Handler, len(middlewares))
	for i, m := range middlewares {
		handlers[i] = m.Handler()
	}
	return &MiddlewareChain{middlewares: handlers}
}

// Then applies the middleware chain to a handler
func (mc *MiddlewareChain) Then(handler fiber.Handler) fiber.Handler {
	for i := len(mc.middlewares) - 1; i >= 0; i-- {
		currentMiddleware := mc.middlewares[i]
		handler = func(c *fiber.Ctx) error {
			return currentMiddleware(c)
		}
	}
	return handler
}
