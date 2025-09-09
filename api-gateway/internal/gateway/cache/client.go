package cache

import (
	"context"
	"fmt"

	pb "github.com/darkphotonKN/community-builds-microservice/common/api/proto/cache"
	"github.com/darkphotonKN/community-builds-microservice/common/discovery"
)

/*
client.go implements the ExampleClient interface, providing methods to interact with
the example-service through gRPC. It uses Consul service discovery to locate the service
and establishes a connection dynamically at runtime.

This gateway serves as a client-side adapter between the API Gateway's REST endpoints and
the example-service's gRPC methods, handling the protocol translation and service location.
It's part of the Gateway pattern that isolates the API Gateway from the implementation
details of the backend microservices.

Each method follows a consistent pattern:
1. Establish connection to the service via service discovery
2. Create a gRPC client
3. Make the gRPC call
4. Return the result or error

Usage:
    registry := consul.NewRegistry(...)
    exampleGateway := gateway.NewExampleGateway(registry)
    example, err := exampleGateway.GetExample(ctx, &pb.GetExampleRequest{Id: "123"})

Note: Remove after copy pasting this as scaffolding.
*/

const (
	serviceName = "cache"
)

type Client struct {
	registry discovery.Registry
}

func NewClient(registry discovery.Registry) CacheClient {
	return &Client{
		registry: registry,
	}
}

func (c *Client) Hello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	// connection instance created through service discovery first
	// searches for the service registered as "orders"
	conn, err := discovery.ServiceConnection(ctx, serviceName, c.registry)

	if err != nil {
		return nil, fmt.Errorf("failed to connect to example service: %w", err)
	}
	defer conn.Close()

	client := pb.NewCacheServiceClient(conn)

	// create client to interface with through service discovery connection
	exampleItem, err := client.Hello(ctx, req)

	fmt.Printf("Creating example %+v through gateway after service discovery\n", exampleItem)
	if err != nil {
		return nil, fmt.Errorf("failed to create example: %w", err)
	}
	return exampleItem, nil
}
