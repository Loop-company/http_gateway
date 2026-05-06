package authclient

import (
	"context"
	"time"

	authpb "github.com/Loop-company/http_gateway/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type Client struct {
	client authpb.AuthServiceClient
	conn   *grpc.ClientConn
}

type LoginResponse struct {
	GUID             string    `json:"guid"`
	AccessToken      string    `json:"access_token"`
	RefreshToken     string    `json:"refresh_token"`
	AccessExpiresAt  time.Time `json:"access_expires_at"`
	RefreshExpiresAt time.Time `json:"refresh_expires_at"`
}

func New(addr string) (*Client, error) {
	if addr == "" {
		addr = "localhost:50051"
	}

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &Client{
		client: authpb.NewAuthServiceClient(conn),
		conn:   conn,
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) Login(ctx context.Context, email, password string, userAgent, ip string) (LoginResponse, int, error) {
	var result LoginResponse

	// Pass metadata
	md := metadata.Pairs(
		"user-agent", userAgent,
		"x-real-ip", ip,
	)
	ctx = metadata.NewOutgoingContext(ctx, md)

	resp, err := c.client.Login(ctx, &authpb.LoginRequest{
		Email:    email,
		Password: password,
	})

	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			return result, 500, err
		}

		switch st.Code() {
		case codes.Unauthenticated:
			return result, 401, err
		case codes.InvalidArgument:
			return result, 400, err
		default:
			return result, 500, err
		}
	}

	result = LoginResponse{
		GUID:             resp.Guid,
		AccessToken:      resp.AccessToken,
		RefreshToken:     resp.RefreshToken,
		AccessExpiresAt:  time.Unix(resp.AccessExpiresAt, 0),
		RefreshExpiresAt: time.Unix(resp.RefreshExpiresAt, 0),
	}

	return result, 200, nil
}
