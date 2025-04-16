package grpc_server

import (
	"fmt"
	"github.com/SwanHtetAungPhyo/common/grpcClient"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func getGRPCAddress() (string, error) {
	host := viper.GetString("services-grpc-info.auth.host")
	if host == "" {
		return "", fmt.Errorf("gRPC host not configured")
	}

	port := viper.GetString("services-grpc-info.auth.port")
	if port == "" {
		return "", fmt.Errorf("gRPC port not configured")
	}

	return fmt.Sprintf("%s:%s", host, port), nil
}
func getCertPath() string {
	var caCertPath = viper.GetString("ca-path")
	if caCertPath == "" {
		logrus.New().Warn("CaCertPath is empty")
		caCertPath = "/Users/swanhtet1aungphyo/IdeaProjects/UniBackend/AppBackendServices/cert/ca.crt"
	}
	return caCertPath
}

func CheckUserExistence() (bool, error) {
	address, err := getGRPCAddress()
	if err != nil {
		return false, fmt.Errorf("failed to get gRPC address: %w", err)
	}

	caCertPath := getCertPath()
	logrus.Infof("Connecting to gRPC server at %s with CA cert: %s", address, caCertPath)

	factory := grpc_client.NewGRPCClientFactory(logrus.StandardLogger(), caCertPath)
	userClient, err := factory.NewUserServiceClient(address)
	if err != nil {
		return false, fmt.Errorf("failed to create user client: %w", err)
	}
	exists, err := userClient.CheckUserExists("swanhtetaungp@gmail.com")
	return exists, err
}
