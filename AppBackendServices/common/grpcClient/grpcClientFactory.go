package grpc_client

import (
	"crypto/tls"
	"crypto/x509"
	"os"

	"github.com/SwanHtetAungPhyo/common/proto"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// GRPCClientFactory handles the creation of various gRPC clients
type GRPCClientFactory struct {
	logger     *logrus.Logger
	caCertPath string
}

func NewGRPCClientFactory(logger *logrus.Logger, caCertPath string) *GRPCClientFactory {
	return &GRPCClientFactory{
		logger:     logger,
		caCertPath: caCertPath,
	}
}

func (f *GRPCClientFactory) CreateConnection(address string) (*grpc.ClientConn, error) {
	caCert, err := os.ReadFile(f.caCertPath)
	if err != nil {
		f.logger.Errorf("Failed to load CA certificate: %v", err)
		return nil, err
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(caCert) {
		f.logger.Error("Failed to add CA certificate to pool")
		return nil, err
	}

	tlsConfig := &tls.Config{
		RootCAs:    certPool,
		MinVersion: tls.VersionTLS12,
	}

	creds := credentials.NewTLS(tlsConfig)

	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(creds))
	if err != nil {
		f.logger.Errorf("Failed to dial gRPC server: %v", err)
		return nil, err
	}

	return conn, nil
}

func (f *GRPCClientFactory) NewUserServiceClient(address string) (*UserServiceClient, error) {
	conn, err := f.CreateConnection(address)
	if err != nil {
		return nil, err
	}

	return &UserServiceClient{
		conn:   conn,
		client: proto.NewUserServiceClient(conn),
		logger: f.logger,
	}, nil
}
