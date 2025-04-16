package consul

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/spf13/viper"
	"log"
)

func RegisterToConsul() {
	consulConfig := api.DefaultConfig()
	consulConfig.Address = viper.GetString("consul.address")

	client, err := api.NewClient(consulConfig)
	if err != nil {
		log.Fatalf("Failed to create Consul client: %v", err)
	}

	serviceRegistration := &api.AgentServiceRegistration{
		ID:      viper.GetString("consul.service_id"),
		Name:    viper.GetString("consul.service_name"),
		Address: viper.GetString("consul.service_address"), // Get service address from config
		Port:    viper.GetInt("consul.service_port"),       // Get service port from config
		Tags:    []string{"auth", "login", "register"},
		Check: &api.AgentServiceCheck{
			HTTP:     viper.GetString("consul.health_check_url"),      // Health check URL from config
			Interval: viper.GetString("consul.health_check_interval"), // Health check interval from config
			Timeout:  viper.GetString("consul.health_check_timeout"),  // Health check timeout from config
		},
	}
	err = client.Agent().ServiceRegister(serviceRegistration)
	if err != nil {
		log.Fatalf("Failed to register service with Consul: %v", err)
	}

	fmt.Println("Service 'authentication' registered with Consul successfully!")
}
