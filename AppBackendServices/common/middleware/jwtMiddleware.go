package middleware

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	// AuthHeader is the name of the header containing the JWT token
	AuthHeader = "Authorization"
	// BearerPrefix is the prefix for Bearer tokens
	BearerPrefix = "bearer"
	// TokenExpiredError is the error message for expired tokens
	TokenExpiredError = "Token expired"
	// InvalidTokenError is the error message for invalid tokens
	InvalidTokenError = "Invalid token"
	// MissingAuthHeaderError is the error message for missing Authorization header
	MissingAuthHeaderError = "Authorization header required"
	// InvalidTokenFormatError is the error message for invalid token format
	InvalidTokenFormatError = "Invalid token format. Use Bearer {token}"
	// ServerConfigError is the error message for server configuration errors
	ServerConfigError = "Server configuration error"
	// TokenRevokedError is the error message for revoked tokens
	TokenRevokedError = "Token revoked"
)

func JwtMiddleware(redisClient *redis.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Path() == "/api/auth/refresh" {
			return c.Next()
		}
		authHeader := c.Get(AuthHeader)
		if authHeader == "" {
			logrus.Warn(MissingAuthHeaderError)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": MissingAuthHeaderError})
		}

		authParts := strings.Split(authHeader, " ")
		if len(authParts) != 2 || strings.ToLower(authParts[0]) != BearerPrefix {
			logrus.Warn(InvalidTokenFormatError)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": InvalidTokenFormatError})
		}

		tokenString := authParts[1]
		ctx := context.Background()
		val, err := redisClient.Get(ctx, fmt.Sprintf("blacklist:%s", tokenString)).Result()
		if errors.Is(err, redis.Nil) {
			fmt.Println("Token is not blacklisted")
		} else if err != nil {
			logrus.WithError(err).Error("Error checking blacklist")
		} else {
			fmt.Printf("Token is blacklisted with value: %s\n", val)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": TokenRevokedError})
		}
		jwtSecret := viper.GetString("jwt.secret")
		if jwtSecret == "" {
			logrus.Error("JWT secret not configured")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": ServerConfigError})
		}

		token, err := parseToken(tokenString, jwtSecret)
		if err != nil {
			return handleTokenError(c, err)
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			logrus.Warn("Invalid token claims format")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": InvalidTokenError})
		}

		// Validate claims
		if err := validateClaims(claims); err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
		}

		// Fix: Correct the user claim key
		c.Locals("user", claims["sub"])
		c.Locals("email", claims["email"])

		return c.Next()
	}
}

// JWTBlacklistMiddleware creates a middleware that checks if a token is blacklisted
func JWTBlacklistMiddleware(redisClient *redis.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Get(AuthHeader)
		if len(token) > 7 && strings.ToLower(token[0:7]) == "bearer " {
			token = token[7:]
		}

		exists, err := redisClient.Exists(context.Background(), "blacklist:"+token).Result()
		if err != nil {
			logrus.WithError(err).Error("Redis error checking blacklist")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": ServerConfigError,
			})
		}

		if exists == 1 {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": TokenRevokedError,
			})
		}

		return c.Next()
	}
}

// parseToken parses and validates a JWT token
func parseToken(tokenString, secret string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
}

// handleTokenError handles various token-related errors
func handleTokenError(c *fiber.Ctx, err error) error {
	logrus.WithError(err).Warn("Token validation failed")

	switch {
	case errors.Is(err, jwt.ErrTokenExpired):
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": TokenExpiredError,
		})
	default:
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": InvalidTokenError,
		})
	}
}

func validateClaims(claims jwt.MapClaims) error {
	// Check expiration
	exp, ok := claims["exp"].(float64)
	if !ok {
		return fmt.Errorf("invalid token expiration")
	}

	if time.Now().Unix() > int64(exp) {
		return fmt.Errorf(TokenExpiredError)
	}

	// Check required claims
	if _, ok := claims["sub"]; !ok {
		return fmt.Errorf("missing user ID claim")
	}

	if _, ok := claims["email"]; !ok {
		return fmt.Errorf("missing email claim")
	}

	// Validate device structure
	if device, ok := claims["device"].(map[string]interface{}); ok {
		requiredFields := []string{"browser", "os", "device_type"}
		for _, field := range requiredFields {
			if _, exists := device[field]; !exists {
				return fmt.Errorf("missing device claim: %s", field)
			}
		}
	}

	return nil
}
