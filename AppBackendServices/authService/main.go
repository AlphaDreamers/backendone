package main

import (
	"fmt"
	"github.com/SwanHtetAungPhyo/auth/cmd"
	"github.com/SwanHtetAungPhyo/auth/internal/config"
	"github.com/SwanHtetAungPhyo/auth/internal/models"
	"github.com/SwanHtetAungPhyo/common/pkg/logutil"
	db "github.com/SwanHtetAungPhyo/database"
	"github.com/gofiber/fiber/v2/log"
	"github.com/hashicorp/consul/api"
	"github.com/spf13/viper"
)

func main() {
	logutil.InitLog("Auth_service")

	err := loadConfig()
	if err != nil {
		logutil.GetLogger().Fatal(err.Error())
		return
	}

	port := viper.GetString("port")
	dns := viper.GetString("dns")

	registerToConsul()

	config.RedisConfigInit()
	db.InitDB(dns)
	database := db.GetDB()
	err = database.AutoMigrate(&models.UserInDB{}, &models.UserBiometric{})
	if err != nil {
		logutil.GetLogger().Fatal(err.Error())
		return
	}
	cmd.Start(port)
}

func registerToConsul() {
	consulConfig := api.DefaultConfig()
	consulConfig.Address = viper.GetString("consul.address")

	client, err := api.NewClient(consulConfig)
	if err != nil {
		log.Fatalf("Failed to create Consul client: %v", err)
	}

	serviceRegistration := &api.AgentServiceRegistration{
		ID:      viper.GetString("consul.service_id"),      // Get service ID from config
		Name:    viper.GetString("consul.service_name"),    // Get service name from config
		Address: viper.GetString("consul.service_address"), // Get service address from config
		Port:    viper.GetInt("consul.service_port"),       // Get service port from config
		Tags:    []string{"auth", "login", "register"},
		Check: &api.AgentServiceCheck{
			HTTP:     viper.GetString("consul.health_check_url"),      // Health check URL from config
			Interval: viper.GetString("consul.health_check_interval"), // Health check interval from config
			Timeout:  viper.GetString("consul.health_check_timeout"),  // Health check timeout from config
		},
	}

	// Register the service
	err = client.Agent().ServiceRegister(serviceRegistration)
	if err != nil {
		log.Fatalf("Failed to register service with Consul: %v", err)
	}

	fmt.Println("Service 'authentication' registered with Consul successfully!")
}
func loadConfig() error {
	viper.SetConfigFile("config.yaml")

	err := viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("Error reading config file, %s", err)
	}

	// Optionally, set defaults
	viper.SetDefault("port", "9008")
	viper.SetDefault("dns", "postgres://localhost:5432/mydb?sslmode=disable")
	viper.SetDefault("jwt_secret", "KZc3qNSpLXzIz7h6UV7/7ZzPxqWmqEUk8X0aW3J3F8M=")

	// Consul defaults
	viper.SetDefault("consul.address", "http://127.0.0.1:8500")
	viper.SetDefault("consul.service_id", "authentication")
	viper.SetDefault("consul.service_name", "auth")
	viper.SetDefault("consul.service_address", "localhost")
	viper.SetDefault("consul.service_port", 9008)
	viper.SetDefault("consul.health_check_url", "http://localhost:9008/health")
	viper.SetDefault("consul.health_check_interval", "10s")
	viper.SetDefault("consul.health_check_timeout", "5s")

	return nil
}
