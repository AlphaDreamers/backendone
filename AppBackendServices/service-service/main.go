package main

import (
	"fmt"
	"github.com/SwanHtetAungPhyo/common/pkg/logutil"
	"github.com/SwanHtetAungPhyo/service-service/grpc_server"
	"github.com/spf13/viper"
)

func main() {
	logutil.InitLog("service-service")
	logger := logutil.GetLogger()
	if err := loadConfig(); err != nil {
		logger.Fatal("Failed to load configuration: ", err)
	}

	userExistence, err := grpc_server.CheckUserExistence()
	if err != nil {
		return
	}
	logger.Println("User existence:", userExistence)
}

func loadConfig() error {
	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("yaml")   // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(".")      // look for config in the working directory

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	// Validate required configurations
	requiredConfigs := []string{
		"ca-path",
		"services-grpc-info.auth.host",
		"services-grpc-info.auth.port",
	}

	for _, configKey := range requiredConfigs {
		if !viper.IsSet(configKey) {
			return fmt.Errorf("missing required configuration: %s", configKey)
		}
	}

	return nil
}
