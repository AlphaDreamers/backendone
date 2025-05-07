package user_rpc

import (
	"context"
	"errors"
	"github.com/SwanHtetAungPhyo/service-one/auth/internal/repo"
	"github.com/SwanHtetAungPhyo/service-one/auth/proto"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"net"
)

type UserRpcServer struct {
	proto.UnimplementedUserRpcMethodServer
	log    *logrus.Logger
	v      *viper.Viper
	repo   *repo.UserRepo
	server *grpc.Server // Store the gRPC server instance
}

func NewUserRpcServer(
	log *logrus.Logger,
	v *viper.Viper,
	repo *repo.UserRepo,
) *UserRpcServer {
	return &UserRpcServer{
		log:    log,
		v:      v,
		repo:   repo,
		server: grpc.NewServer(), // Initialize gRPC server
	}
}

func (ur *UserRpcServer) UserExistenceCall(
	ctx context.Context,
	req *proto.UserExistenceReq,
) (*proto.UserExistenceResp, error) {
	userId := uuid.MustParse(req.GetUserId())
	if userId == uuid.Nil {
		return &proto.UserExistenceResp{
			Status: false,
		}, errors.New("user id is nil")
	}
	existence, err := ur.repo.UserExistence(ctx, userId)
	if err != nil {
		return nil, err
	}
	if !existence {
		return &proto.UserExistenceResp{
			Status: false,
		}, errors.New("user does not exist")
	}
	return &proto.UserExistenceResp{
		Status: true,
	}, nil
}

// Register the gRPC server (called on start)
func (ur *UserRpcServer) Start() error {
	listenAddr := ur.v.GetString("server.grpc.listen")
	lis, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return err
	}

	// Register the service
	proto.RegisterUserRpcMethodServer(ur.server, ur)

	ur.log.Infof("gRPC server listening on %s", listenAddr)

	// Start serving in a goroutine
	go func() {
		if err := ur.server.Serve(lis); err != nil {
			ur.log.Errorf("gRPC server failed: %v", err)
		}
	}()

	return nil
}

// Stop Graceful shutdown (called on stop)
func (ur *UserRpcServer) Stop() {
	ur.log.Info("Shutting down gRPC server...")
	ur.server.GracefulStop()
	ur.log.Info("gRPC server stopped")
}

// Module FX Module that provides the gRPC server with lifecycle hooks
var Module = fx.Module("user_rpc",
	fx.Provide(NewUserRpcServer),
	fx.Invoke(func(lc fx.Lifecycle, server *UserRpcServer) {
		lc.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {
				return server.Start()
			},
			OnStop: func(ctx context.Context) error {
				server.Stop()
				return nil
			},
		})
	}),
)
