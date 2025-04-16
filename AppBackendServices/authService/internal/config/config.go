package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Consul   ConsulConfig   `mapstructure:"consul"`
	Redis    RedisConfig    `mapstructure:"redis"`
	Vault    VaultConfig    `mapstructure:"vault"`
}

type ServerConfig struct {
	Port         string        `mapstructure:"port"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	Environment  string        `mapstructure:"environment"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
	TimeZone string `mapstructure:"timezone"`
}

type JWTConfig struct {
	Secret          string        `mapstructure:"secret"`
	AccessTokenTTL  time.Duration `mapstructure:"access_token_ttl"`
	RefreshTokenTTL time.Duration `mapstructure:"refresh_token_ttl"`
	SigningMethod   string        `mapstructure:"signing_method"`
	SigningKey      string        `mapstructure:"signing_key"`
	VerifyingKey    string        `mapstructure:"verifying_key"`
}

type ConsulConfig struct {
	Address             string        `mapstructure:"address"`
	ServiceID           string        `mapstructure:"service_id"`
	ServiceName         string        `mapstructure:"service_name"`
	ServiceAddress      string        `mapstructure:"service_address"`
	ServicePort         int           `mapstructure:"service_port"`
	HealthCheckURL      string        `mapstructure:"health_check_url"`
	HealthCheckInterval time.Duration `mapstructure:"health_check_interval"`
	HealthCheckTimeout  time.Duration `mapstructure:"health_check_timeout"`
}

type RedisConfig struct {
	Address  string `mapstructure:"address"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type VaultConfig struct {
	Address    string `mapstructure:"address"`
	RootToken  string `mapstructure:"root_token"`
	SecretPath string `mapstructure:"secret_path"`
}

var cfg *Config

func LoadConfig() error {
	viper.SetConfigName("config.development") // Set config name (without extension)
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./internal/config") // Add config directory
	viper.AddConfigPath(".")                 // Allow loading from the current directory
	viper.AutomaticEnv()                     // Enable environment variable bindings

	setDefaults()

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("error reading config file: %w", err)
	}

	cfg = &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		return fmt.Errorf("error unmarshaling config: %w", err)
	}

	return nil
}

func setDefaults() {
	viper.SetDefault("server.port", "9008")
	viper.SetDefault("server.read_timeout", "15s")
	viper.SetDefault("server.write_timeout", "15s")
	viper.SetDefault("server.environment", "development")

	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)

	viper.SetDefault("database.dbname", "mydb")
	viper.SetDefault("database.sslmode", "disable")
	viper.SetDefault("database.timezone", "UTC")

	viper.SetDefault("jwt.access_token_ttl", "15m")
	viper.SetDefault("jwt.refresh_token_ttl", "24h")
	viper.SetDefault("jwt.signing_method", "RS256")
	viper.BindEnv("jwt.secret")        // Use environment variable
	viper.BindEnv("jwt.signing_key")   // Use environment variable
	viper.BindEnv("jwt.verifying_key") // Use environment variable

	viper.SetDefault("consul.address", "http://localhost:8500")
	viper.SetDefault("consul.service_id", "authentication")
	viper.SetDefault("consul.service_name", "auth")
	viper.SetDefault("consul.service_address", "localhost")
	viper.SetDefault("consul.service_port", 9008)
	viper.SetDefault("consul.health_check_url", "http://localhost:9008/health")
	viper.SetDefault("consul.health_check_interval", "10s")
	viper.SetDefault("consul.health_check_timeout", "5s")

	viper.SetDefault("redis.address", "localhost:6379")
	viper.SetDefault("redis.db", 0)

	viper.SetDefault("vault.address", "http://localhost:8200")
	viper.BindEnv("vault.root_token") // Use environment variable
	viper.SetDefault("vault.secret_path", "auth-service-wallet")
}

func GetConfig() *Config {
	if cfg == nil {
		_ = LoadConfig() // Ensure config is loaded before accessing
	}
	return cfg
}
