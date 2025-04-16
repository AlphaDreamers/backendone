package grpc_client

import (
	"context"
	"github.com/SwanHtetAungPhyo/common/proto"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type UserServiceClient struct {
	conn   *grpc.ClientConn
	client proto.UserServiceClient
	logger *logrus.Logger
}

func NewUserServiceClient(conn *grpc.ClientConn, logger *logrus.Logger) *UserServiceClient {
	return &UserServiceClient{
		conn:   conn,
		client: proto.NewUserServiceClient(conn),
		logger: logger,
	}
}

func (c *UserServiceClient) CheckUserExists(email, userId string) (bool, error) {
	resp, err := c.client.CheckUserExists(context.Background(), &proto.UserExistsRequest{Email: email,
		UserId: userId})
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			c.logger.Errorf("Failed to close connection: %v", err)
		}
	}(c.conn)
	if err != nil {
		c.logger.Errorf("gRPC call failed: %v", err)
		return false, err
	}
	return resp.GetExists().GetValue(), nil
}
