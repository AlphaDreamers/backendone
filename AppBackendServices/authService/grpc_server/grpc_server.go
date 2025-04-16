package grpc_server

import (
	"context"
	"crypto/tls"
	"github.com/SwanHtetAungPhyo/auth/internal/repo"
	"github.com/SwanHtetAungPhyo/common/proto"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"net"
)

type GrpcServer struct {
	proto.UnimplementedUserServiceServer
	logger *logrus.Logger
}

func NewGrpcServer(logger *logrus.Logger) *GrpcServer {
	return &GrpcServer{logger: logger}
}

func (s *GrpcServer) CheckUserExists(ctx context.Context, req *proto.UserExistsRequest) (*proto.UserExistsResponse, error) {
	response := &proto.UserExistsResponse{
		Exists: &wrapperspb.BoolValue{},
	}
	if req.UserId == "" {
		email, err := repo.NewImpl().GetByEmail(req.GetEmail())
		if err != nil {
			return nil, err
		}
		if email == nil {
			response.Exists.Value = false
			return response, nil
		}
		response.Exists.Value = true
	} else {

	}
	return response, nil
}
func (s *GrpcServer) Start(port string) {
	serverCert, err := tls.LoadX509KeyPair(
		"/Users/swanhtet1aungphyo/IdeaProjects/UniBackend/AppBackendServices/cert/auth-bundle.crt", // Combined server + CA certs
		"/Users/swanhtet1aungphyo/IdeaProjects/UniBackend/AppBackendServices/cert/auth.key",
	)
	if err != nil {
		s.logger.Fatal("Failed to load server certificates: ", err)
	}

	// Create TLS configuration
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.NoClientCert, // Adjust based on your security requirements
		MinVersion:   tls.VersionTLS12,
	}

	// Create credentials
	creds := credentials.NewTLS(tlsConfig)

	list, err := net.Listen("tcp", ":"+port)
	if err != nil {
		s.logger.Errorf("failed to listen: %v", err)

	}
	server := grpc.NewServer(grpc.Creds(creds))
	proto.RegisterUserServiceServer(server, s)
	reflection.Register(server)
	s.logger.Infof("grpc server listening at %v", port)
	err = server.Serve(list)
	if err != nil {
		s.logger.Errorf("failed to serve: %v", err)
	}
}
