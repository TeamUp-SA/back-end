package userclient

import (
	"context"
	"time"

	"google.golang.org/grpc"

	"auth-service/internal/config"
	userv1 "user-service/pb/userv1"
)

type Client struct {
	cc  *grpc.ClientConn
	svc userv1.UserServiceClient
}

func New(cfg *config.Config) (*Client, error) {
	conn, err := grpc.Dial(cfg.UserGrpcServiceAddr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return &Client{cc: conn, svc: userv1.NewUserServiceClient(conn)}, nil
}

func (c *Client) Close() error { return c.cc.Close() }

func (c *Client) UpsertUser(ctx context.Context, in *userv1.UpsertUserRequest) (*userv1.UpsertUserResponse, error) {
	if deadline, ok := ctx.Deadline(); !ok || time.Until(deadline) <= 0 {
		// default 5s timeout if none
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
	}
	return c.svc.UpsertUser(ctx, in)
}
