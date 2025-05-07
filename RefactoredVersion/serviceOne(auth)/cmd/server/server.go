package server

import (
	"context"
	"github.com/SwanHtetAungPhyo/service-one/auth/cmd/server/middleware"
	"github.com/SwanHtetAungPhyo/service-one/auth/internal/handler"
	"github.com/SwanHtetAungPhyo/service-one/auth/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"gorm.io/gorm"
	"time"
)

type AppBehaviour interface {
	Start() error
	Stop() error
}

type AppState struct {
	app         *fiber.App
	log         *logrus.Logger
	db          *gorm.DB
	v           *viper.Viper
	nats        *nats.Conn
	handler     *handler.Handler
	jwtGen      *utils.JwtTokenGenerator
	userHandler *handler.UserHandler
}

func NewAppState(
	log *logrus.Logger,
	db *gorm.DB,
	v *viper.Viper,
	nats *nats.Conn,
	jwtGen *utils.JwtTokenGenerator,
	handler *handler.Handler,
	userHandler *handler.UserHandler,
) *AppState {
	state := &AppState{
		log:         log,
		db:          db,
		v:           v,
		nats:        nats,
		jwtGen:      jwtGen,
		handler:     handler,
		userHandler: userHandler,
	}

	state.app = fiber.New(fiber.Config{
		DisableStartupMessage: false,
		Prefork:               false,
		StrictRouting:         false,
		IdleTimeout:           state.v.GetDuration("server.idle_timeout") * time.Second,
		ReadTimeout:           state.v.GetDuration("server.read_timeout") * time.Second,
		WriteTimeout:          state.v.GetDuration("server.write_timeout") * time.Second,
		EnablePrintRoutes:     true,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"error":   err.Error(),
			})
		},
	})

	state.middlewareSetUp()
	state.routeSetup()
	return state
}

func (a *AppState) Start() error {
	addr := a.v.GetString("server.addr")
	a.log.Info("Starting server on ", addr)
	certFile := a.v.GetString("server.cert_file")
	keyFile := a.v.GetString("server.key_file")

	if certFile != "" && keyFile != "" {
		a.log.Info("Starting HTTPS server on ", addr)
		return a.app.ListenTLS(addr, certFile, keyFile)
	}
	a.log.Info("Starting HTTP server on ", addr)
	return a.app.Listen(addr)
}

func (a *AppState) Stop() error {
	a.log.Info("Gracefully shutting down server...")
	return a.app.Shutdown()
}

func (a *AppState) routeSetup() {
	// Health check endpoint
	a.app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "OK",
			"service": a.v.GetString("server.name"),
		})
	})

	// Auth routes
	auth := a.app.Group("/auth")

	// Public auth endpoints
	auth.Post("/login", middleware.LoginMiddleware(), a.handler.Login)
	auth.Post("/register", a.handler.Register)
	auth.Post("/verify-email", a.handler.AccRegisterEmailVerification)
	auth.Post("/reset-password", a.handler.ResetPassword)
	auth.Post("/verify-reset-token", a.handler.ResetPasswordTokenVerify)
	auth.Post("/forgot-password", a.handler.ForgotPassword)
	auth.Post("/verify-forgot-token", a.handler.ForgotPasswordTokenVerified)

	protected := auth.Group("", middleware.Check())
	protected.Post("/logout", a.handler.Logout)
	protected.Post("/refresh", a.handler.RefreshToken)
	protected.Get("/dashboard/:userId", middleware.AuthMiddleware(), a.userHandler.GetDashboard)
	protected.Get("/me/:userId", a.userHandler.GetMet)
}

func (a *AppState) middlewareSetUp() {
	a.app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))

	// Logger middleware
	a.app.Use(logger.New(logger.Config{
		Format:     "${time} ${status} - ${method} ${path} [${latency}]\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "UTC",
	}))

	// Rate limiting
	a.app.Use(limiter.New(limiter.Config{
		Max:        a.v.GetInt("server.rate_limit"),
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Too many requests",
			})
		},
	}))

	// CORS
	a.app.Use(cors.New(cors.Config{
		AllowOrigins: a.v.GetString("server.allow_origins"),
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, PATCH, DELETE, OPTIONS",
		//AllowCredentials: true,
	}))

	// Compression
	a.app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))
}

func ServerLifeCycleReg(lc fx.Lifecycle, app *AppState) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := app.Start(); err != nil {
					app.log.Fatalf("Failed to start server: %v", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return app.Stop()
		},
	})
}

var _ AppBehaviour = (*AppState)(nil)
