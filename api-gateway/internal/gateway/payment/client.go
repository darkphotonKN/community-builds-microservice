package payment

import (
	"context"
	"fmt"

	pb "github.com/darkphotonKN/community-builds-microservice/common/api/proto/payment"
	"github.com/darkphotonKN/community-builds-microservice/common/discovery"
)

const (
	serviceName = "payment-service"
)

type Client struct {
	registry discovery.Registry
}

func NewClient(registry discovery.Registry) PaymentClient {
	return &Client{
		registry: registry,
	}
}

func (c *Client) CreateCustomer(ctx context.Context, req *pb.CreateCustomerRequest) (*pb.CreateCustomerResponse, error) {

	// connection instance created through service discovery first
	// searches for the service registered as "orders"
	conn, err := discovery.ServiceConnection(ctx, serviceName, c.registry)

	if err != nil {
		return nil, fmt.Errorf("failed to connect to item service: %w", err)
	}
	defer conn.Close()

	client := pb.NewPaymentServiceClient(conn)

	res, err := client.CreateCustomer(ctx, req)

	fmt.Printf("Creating customer %+v through gateway after service discovery\n", res)

	if err != nil {
		return nil, fmt.Errorf("failed to create customer: %w", err)
	}

	return res, nil
}

func (c *Client) CreateSubscription(ctx context.Context, req *pb.CreateSubscriptionRequest) (*pb.CreateSubscriptionResponse, error) {

	// connection instance created through service discovery first
	// searches for the service registered as "orders"
	conn, err := discovery.ServiceConnection(ctx, serviceName, c.registry)

	if err != nil {
		return nil, fmt.Errorf("failed to connect to item service: %w", err)
	}
	defer conn.Close()

	client := pb.NewPaymentServiceClient(conn)

	res, err := client.CreateSubscription(ctx, req)

	fmt.Printf("Creating subscription %+v through gateway after service discovery\n", res)

	if err != nil {
		return nil, fmt.Errorf("failed to create subscription: %w", err)
	}

	return res, nil
}
