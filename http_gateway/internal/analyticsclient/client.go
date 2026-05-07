package analyticsclient

import (
	"context"

	gatewaypb "github.com/Loop-company/http_gateway/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	client gatewaypb.AnalyticsServiceClient
	conn   *grpc.ClientConn
}

func New(addr string) (*Client, error) {
	if addr == "" {
		addr = "localhost:50053"
	}

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &Client{
		client: gatewaypb.NewAnalyticsServiceClient(conn),
		conn:   conn,
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) SearchEvents(ctx context.Context, req *gatewaypb.SearchEventsRequest) (*gatewaypb.SearchEventsResponse, error) {
	return c.client.SearchEvents(ctx, req)
}

func (c *Client) RegistrationsReport(ctx context.Context, from, to string) (*gatewaypb.RegistrationsReportResponse, error) {
	return c.client.GetRegistrationsReport(ctx, &gatewaypb.ReportPeriodRequest{From: from, To: to})
}

func (c *Client) LoginReport(ctx context.Context, from, to string) (*gatewaypb.LoginReportResponse, error) {
	return c.client.GetLoginReport(ctx, &gatewaypb.ReportPeriodRequest{From: from, To: to})
}

func (c *Client) TopUsersReport(ctx context.Context, from, to string, limit int32) (*gatewaypb.TopUsersReportResponse, error) {
	return c.client.GetTopUsersReport(ctx, &gatewaypb.TopUsersReportRequest{From: from, To: to, Limit: limit})
}
