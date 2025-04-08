package handler

import (
	"github.com/SwanHtetAungPhyo/auth/internal/config"
	"github.com/SwanHtetAungPhyo/auth/internal/services"
	"github.com/SwanHtetAungPhyo/common/pkg/logutil"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type (
	AuthHandler interface {
		Login(c *fiber.Ctx) error
		Register(c *fiber.Ctx) error
		Me(c *fiber.Ctx) error
		Refresh(c *fiber.Ctx) error
		Verify(c *fiber.Ctx) error
		SendCode(c *fiber.Ctx) error
		StoreInVault(c *fiber.Ctx) error
		generateCode() string
	}

	Impl struct {
		logger      *logrus.Logger
		service     *services.Impl
		redisClient *redis.Client
	}
)

func NewHandler() *Impl {
	return &Impl{
		logger:      logutil.GetLogger(),
		service:     services.NewServiceImpl(),
		redisClient: config.GetRedisClient(),
	}
}
