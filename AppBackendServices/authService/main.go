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

	err := config.LoadConfig()
	if err != nil {
		logutil.GetLogger().Fatal(err.Error())
		return
	}

	cfg := config.GetConfig()

	registerToConsul()

	config.RedisConfigInit()
	dns := fmt.Sprintf(
		"postgres://%s:%d/%s?sslmode=%s&timezone=%s",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.DBName, cfg.Database.SSLMode, cfg.Database.TimeZone,
	)
	db.InitDB(dns)
	database := db.GetDB()
	err = database.AutoMigrate(&models.UserInDB{}, &models.UserBiometric{})
	if err != nil {
		logutil.GetLogger().Fatal(err.Error())
		return
	}
	cmd.Start(cfg.Server.Port, database)
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

	// Register the serxxvice
	err = client.Agent().ServiceRegister(serviceRegistration)
	if err != nil {
		log.Fatalf("Failed to register service with Consul: %v", err)
	}

	fmt.Println("Service 'authentication' registered with Consul successfully!")
}
