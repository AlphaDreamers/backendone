package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/SwanHtetAungPhyo/gateways/middleware"
	"github.com/valyala/fasthttp"

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
	cfg, err := loadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	app := fiber.New(fiber.Config{
		ReadTimeout:   cfg.Server.ReadTimeout * time.Second,
		WriteTimeout:  cfg.Server.WriteTimeout * time.Second,
		IdleTimeout:   cfg.Server.IdleTimeout * time.Second,
		Prefork:       cfg.Server.Prefork,
		CaseSensitive: cfg.Server.CaseSensitive,
	})

	// Initialize middleware
	jwtMiddleware := middleware.NewJWKSMiddleware(
		viper.GetString("cognito.jwk_url"),
		viper.GetString("cognito.issuer_url"),
		viper.GetString("cognito.client_id"),
		time.Hour,
	)

	// Middleware stack
	app.Use(recover.New())
	app.Use(requestid.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: strings.Join([]string{
			fiber.MethodGet,
			fiber.MethodPost,
			fiber.MethodHead,
			fiber.MethodPut,
			fiber.MethodDelete,
			fiber.MethodPatch,
		}, ","),
	}))

	app.Use(logger.New(logger.Config{
		Format:     cfg.Logging.Format,
		TimeFormat: cfg.Logging.TimeFormat,
		TimeZone:   cfg.Logging.TimeZone,
	}))

	setupRoutes(app, cfg, jwtMiddleware)

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "OK",
			"service": "api-gateway",
		})
	})

	serverAddr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("Starting API gateway on %s with %d services", serverAddr, len(cfg.Services))
	if err := app.ListenTLS(serverAddr, "./certificates/cert.pem", "./certificates/key.pem"); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func setupRoutes(app *fiber.App, cfg *GatewayConfig, jwtMiddleware *middleware.JWKSMiddleware) {
	for _, service := range cfg.Services {
		client := &fasthttp.LBClient{
			Clients: []fasthttp.BalancingClient{
				createHostClient(service),
			},
			HealthCheck: func(req *fasthttp.Request, resp *fasthttp.Response, err error) bool {
				return err == nil && resp.StatusCode() == fasthttp.StatusOK
			},
		}

		serviceProxy := proxy.Balancer(proxy.Config{
			Servers:        []string{fmt.Sprintf("https://%s:%d", service.Host, service.Port)},
			Client:         client,
			ModifyRequest:  createRequestModifier(service),
			ModifyResponse: func(c *fiber.Ctx) error { return nil },
			Timeout:        10 * time.Second,
		})

		registerRoutes(app, service, serviceProxy, jwtMiddleware)
		registerHealthCheck(app, service, client)
	}
}

func createHostClient(service Service) *fasthttp.HostClient {
	return &fasthttp.HostClient{
		Addr:                fmt.Sprintf("%s:%d", service.Host, service.Port),
		IsTLS:               true,
		ReadTimeout:         10 * time.Second,
		WriteTimeout:        10 * time.Second,
		MaxIdleConnDuration: 30 * time.Second,
		TLSConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
}

func createRequestModifier(service Service) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		for key, value := range service.Headers {
			c.Request().Header.Set(key, value)
		}
		c.Request().Header.Set("X-Forwarded-For", c.IP())
		c.Request().Header.Set("X-Forwarded-Host", c.Hostname())
		c.Request().Header.Set("X-Forwarded-Proto", c.Protocol())
		return nil
	}
}

func registerRoutes(app *fiber.App, service Service, proxy fiber.Handler, jwtMiddleware *middleware.JWKSMiddleware) {
	path := fmt.Sprintf("%s/*", service.Prefix)
	handler := func(c *fiber.Ctx) error {
		if err := proxy(c); err != nil {
			log.Printf("Proxy error for %s: %v", path, err)
			return c.Status(fiber.StatusBadGateway).SendString("Service unavailable")
		}
		return nil
	}

	if service.Prefix == "/auth" {
		app.All(path, handler)
		log.Printf("Registered public route: %s", path)
	} else {
		chain := middleware.NewMiddlewareChain(jwtMiddleware)
		app.All(path, chain.Then(handler))
		log.Printf("Registered protected route: %s", path)
	}
}

func registerHealthCheck(app *fiber.App, service Service, client *fasthttp.LBClient) {
	if service.HealthPath == "" {
		return
	}

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

func loadConfig(filename string) (*GatewayConfig, error) {
	viper.SetConfigFile(filename)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	setDefaultConfigValues()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var config GatewayConfig
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to decode config: %w", err)
	}

	if err := validateConfig(&config); err != nil {
		return nil, err
	}

	normalizeServicePrefixes(config.Services)

	return &config, nil
}

func setDefaultConfigValues() {
	viper.SetDefault("server.port", 3000)
	viper.SetDefault("server.read_timeout", 10)
	viper.SetDefault("server.write_timeout", 10)
	viper.SetDefault("server.idle_timeout", 30)
	viper.SetDefault("server.prefork", false)
	viper.SetDefault("server.case_sensitive", false)
	viper.SetDefault("logging.format", "[${time}] ${status} - ${method} ${path}\n")
	viper.SetDefault("logging.time_format", "2006-01-02 15:04:05")
	viper.SetDefault("logging.time_zone", "UTC")
	viper.SetDefault("cognito.jwk_url", "https://cognito-idp.us-east-1.amazonaws.com/us-east-1_z6jb3eESF/.well-known/jwks.json")
	viper.SetDefault("cognito.issuer_url", "https://cognito-idp.us-east-1.amazonaws.com/us-east-1_z6jb3eESF")
	viper.SetDefault("cognito.client_id", "7qllcjjcq7p506kq88vkfiu92g")
}

func validateConfig(config *GatewayConfig) error {
	if config.Server.Port == 0 {
		return fmt.Errorf("server port must be specified")
	}

	if len(config.Services) == 0 {
		return fmt.Errorf("no services configured")
	}

	for _, svc := range config.Services {
		if svc.Prefix == "" || svc.Host == "" || svc.Port == 0 {
			return fmt.Errorf("invalid service configuration: %+v", svc)
		}
	}
	return nil
}

func normalizeServicePrefixes(services []Service) {
	for i := range services {
		if !strings.HasPrefix(services[i].Prefix, "/") {
			services[i].Prefix = "/" + services[i].Prefix
		}
	}
}
