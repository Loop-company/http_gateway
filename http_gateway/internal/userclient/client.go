package userclient

import (
	"context"

	userpb "github.com/Loop-company/http_gateway/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	client userpb.UserServiceClient
	conn   *grpc.ClientConn
}

func New(addr string) (*Client, error) {
	if addr == "" {
		addr = "localhost:50052"
	}

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &Client{
		client: userpb.NewUserServiceClient(conn),
		conn:   conn,
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) GetProfile(ctx context.Context, userID string) (*userpb.GetProfileResponse, error) {
	return c.client.GetProfile(ctx, &userpb.GetProfileRequest{UserId: userID})
}

func (c *Client) UpdateName(ctx context.Context, userID, name string) error {
	_, err := c.client.UpdateName(ctx, &userpb.UpdateNameRequest{UserId: userID, Name: name})
	return err
}

func (c *Client) UpdateStatus(ctx context.Context, userID, status string) error {
	_, err := c.client.UpdateStatus(ctx, &userpb.UpdateStatusRequest{UserId: userID, Status: status})
	return err
}
