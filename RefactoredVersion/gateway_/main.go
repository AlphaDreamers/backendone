package main

import (
	"crypto/tls"
	"fmt"
	"github.com/valyala/fasthttp"
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/proxy"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/spf13/viper"
)

type Service struct {
	Name       string            `mapstructure:"name"`
	Prefix     string            `mapstructure:"prefix"`
	Host       string            `mapstructure:"host"`
	Port       int               `mapstructure:"port"`
	HealthPath string            `mapstructure:"health_path"`
	Headers    map[string]string `mapstructure:"headers"`
}

type GatewayConfig struct {
	Server struct {
		Port          int           `mapstructure:"port"`
		ReadTimeout   time.Duration `mapstructure:"read_timeout"`
		WriteTimeout  time.Duration `mapstructure:"write_timeout"`
		IdleTimeout   time.Duration `mapstructure:"idle_timeout"`
		Prefork       bool          `mapstructure:"prefork"`
		CaseSensitive bool          `mapstructure:"case_sensitive"`
	} `mapstructure:"server"`

	Services []Service `mapstructure:"services"`

	Logging struct {
		Format     string `mapstructure:"format"`
		TimeFormat string `mapstructure:"time_format"`
		TimeZone   string `mapstructure:"time_zone"`
	} `mapstructure:"logging"`
}

func main() {
	// Load configuration
	cfg, err := loadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize Fiber app with configuration
	app := fiber.New(fiber.Config{
		ReadTimeout:   cfg.Server.ReadTimeout * time.Second,
		WriteTimeout:  cfg.Server.WriteTimeout * time.Second,
		IdleTimeout:   cfg.Server.IdleTimeout * time.Second,
		Prefork:       cfg.Server.Prefork,
		CaseSensitive: cfg.Server.CaseSensitive,
	})

	// Middleware stack
	app.Use(recover.New())   // Recover from panics
	app.Use(requestid.New()) // Add request ID
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: strings.Join([]string{
			fiber.MethodGet,
			fiber.MethodPost,
			fiber.MethodHead,
			fiber.MethodPut,
			fiber.MethodDelete,
			fiber.MethodPatch,
		}, ","),
	}))

	// Configure logger based on config
	app.Use(logger.New(logger.Config{
		Format:     cfg.Logging.Format,
		TimeFormat: cfg.Logging.TimeFormat,
		TimeZone:   cfg.Logging.TimeZone,
	}))

	// Setup routes for each service
	setupRoutes(app, cfg)

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "OK",
			"service": "api-gateway",
		})
	})

	// Start the server
	serverAddr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("Starting API gateway on %s with %d services", serverAddr, len(cfg.Services))
	if err := app.ListenTLS(serverAddr, "./certificates/cert.pem", "./certificates/key.pem"); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func setupRoutes(app *fiber.App, cfg *GatewayConfig) {
	for _, service := range cfg.Services {
		client := &fasthttp.LBClient{
			Clients: []fasthttp.BalancingClient{
				&fasthttp.HostClient{
					Addr:                fmt.Sprintf("%s:%d", service.Host, service.Port),
					IsTLS:               true,
					ReadTimeout:         10 * time.Second,
					WriteTimeout:        10 * time.Second,
					MaxIdleConnDuration: 30 * time.Second,
					TLSConfig: &tls.Config{
						InsecureSkipVerify: true,
					},
				},
			},
			HealthCheck: func(req *fasthttp.Request, resp *fasthttp.Response, err error) bool {
				return err == nil && resp.StatusCode() == fasthttp.StatusOK
			},
			// Timeout for health checks
		}

		serviceProxy := proxy.Balancer(proxy.Config{
			Servers: []string{
				fmt.Sprintf("https://%s:%d", service.Host, service.Port),
			},
			Client: client,
			ModifyRequest: func(c *fiber.Ctx) error {
				for key, value := range service.Headers {
					c.Request().Header.Set(key, value)
				}
				c.Request().Header.Set("X-Forwarded-For", c.IP())
				c.Request().Header.Set("X-Forwarded-Host", c.Hostname())
				c.Request().Header.Set("X-Forwarded-Proto", c.Protocol())
				return nil
			},
			ModifyResponse: func(c *fiber.Ctx) error {
				// You can modify responses here if needed
				return nil
			},
			Timeout: 10 * time.Second,
		})

		path := fmt.Sprintf("%s/*", service.Prefix)
		app.All(path, func(c *fiber.Ctx) error {
			if err := serviceProxy(c); err != nil {
				log.Printf("Proxy error for %s: %v", path, err)
				return c.Status(fiber.StatusBadGateway).SendString("Service unavailable")
			}
			return nil
		})

		if service.HealthPath != "" {
			healthPath := fmt.Sprintf("%s%s", service.Prefix, service.HealthPath)
			app.Get(healthPath, func(c *fiber.Ctx) error {
				req := fasthttp.AcquireRequest()
				defer fasthttp.ReleaseRequest(req)
				resp := fasthttp.AcquireResponse()
				defer fasthttp.ReleaseResponse(resp)

				req.SetRequestURI(fmt.Sprintf("https://%s:%d%s", service.Host, service.Port, service.HealthPath))
				req.Header.SetMethod(fiber.MethodGet)

				if err := client.Do(req, resp); err != nil {
					return c.Status(fiber.StatusBadGateway).SendString("Health check failed")
				}

				c.Status(resp.StatusCode())
				c.Set("Content-Type", string(resp.Header.ContentType()))
				return c.Send(resp.Body())
			})
			log.Printf("Registered health check: %s", healthPath)
		}

		log.Printf("Registered route: %s -> https://%s:%d", service.Prefix, service.Host, service.Port)
	}
}
func loadConfig(filename string) (*GatewayConfig, error) {
	viper.SetConfigFile(filename)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	viper.SetDefault("server.port", 3000)
	viper.SetDefault("server.read_timeout", 10)
	viper.SetDefault("server.write_timeout", 10)
	viper.SetDefault("server.idle_timeout", 30)
	viper.SetDefault("server.prefork", false)
	viper.SetDefault("server.case_sensitive", false)
	viper.SetDefault("logging.format", "[${time}] ${status} - ${method} ${path}\n")
	viper.SetDefault("logging.time_format", "2006-01-02 15:04:05")
	viper.SetDefault("logging.time_zone", "UTC")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var config GatewayConfig
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to decode config: %w", err)
	}

	// Validate configuration
	if config.Server.Port == 0 {
		return nil, fmt.Errorf("server port must be specified")
	}

	if len(config.Services) == 0 {
		return nil, fmt.Errorf("no services configured")
	}

	for _, svc := range config.Services {
		if svc.Prefix == "" || svc.Host == "" || svc.Port == 0 {
			return nil, fmt.Errorf("invalid service configuration: %+v", svc)
		}
		if !strings.HasPrefix(svc.Prefix, "/") {
			svc.Prefix = "/" + svc.Prefix
		}
	}

	return &config, nil
}
