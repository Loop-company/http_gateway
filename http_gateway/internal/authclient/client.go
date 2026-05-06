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

func (c *Client) Register(ctx context.Context, email, password string) (string, error) {
	resp, err := c.client.Register(ctx, &authpb.RegisterRequest{
		Email:    email,
		Password: password,
	})
	if err != nil {
		return "", err
	}
	return resp.Message, nil
}

func (c *Client) Verify(ctx context.Context, email, code string) (string, error) {
	resp, err := c.client.Verify(ctx, &authpb.VerifyRequest{
		Email: email,
		Code:  code,
	})
	if err != nil {
		return "", err
	}
	return resp.Guid, nil
}

func (c *Client) Refresh(ctx context.Context, refreshToken, userAgent, ip string, userGUID, sessionID string) (LoginResponse, error) {
	var result LoginResponse

	md := metadata.Pairs(
		"user-agent", userAgent,
		"x-real-ip", ip,
		"user_guid", userGUID,
		"session_id", sessionID,
	)
	ctx = metadata.NewOutgoingContext(ctx, md)

	resp, err := c.client.Refresh(ctx, &authpb.RefreshRequest{
		RefreshToken: refreshToken,
	})
	if err != nil {
		return result, err
	}

	result = LoginResponse{
		AccessToken:      resp.AccessToken,
		RefreshToken:     resp.RefreshToken,
		AccessExpiresAt:  time.Unix(resp.AccessExpiresAt, 0),
		RefreshExpiresAt: time.Unix(resp.RefreshExpiresAt, 0),
	}

	return result, nil
}

func (c *Client) ValidateToken(ctx context.Context, accessToken string) (string, string, error) {
	resp, err := c.client.ValidateToken(ctx, &authpb.ValidateTokenRequest{
		AccessToken: accessToken,
	})
	if err != nil {
		return "", "", err
	}
	return resp.Guid, resp.SessionId, nil
}

func (c *Client) Logout(ctx context.Context, userGUID string) error {
	md := metadata.Pairs("user_guid", userGUID)
	ctx = metadata.NewOutgoingContext(ctx, md)

	_, err := c.client.Logout(ctx, &authpb.LogoutRequest{})
	return err
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
