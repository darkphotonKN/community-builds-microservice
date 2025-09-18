package auth

import (
	"context"
	"fmt"

	pb "github.com/darkphotonKN/community-builds-microservice/common/api/proto/auth"
	"github.com/darkphotonKN/community-builds-microservice/common/discovery"
)

const (
	serviceName = "auth"
)

type Client struct {
	registry discovery.Registry
}

func NewClient(registry discovery.Registry) AuthClient {
	return &Client{
		registry: registry,
	}
}

func (c *Client) UpdateStripeCustomer(ctx context.Context, req *pb.UpdateStripeCustomerRequest) (*pb.UpdateStripeCustomerResponse, error) {
	conn, err := discovery.ServiceConnection(ctx, serviceName, c.registry)

	if err != nil {
		return nil, fmt.Errorf("failed to connect to auth service: %w", err)
	}
	defer conn.Close()

	client := pb.NewAuthServiceClient(conn)

	member, err := client.UpdateStripeCustomer(ctx, req)
	return member, err
}
