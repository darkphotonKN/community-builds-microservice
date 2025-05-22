package item

import (
	"context"
	"fmt"

	pb "github.com/darkphotonKN/community-builds-microservice/common/api/proto/item"
	"github.com/darkphotonKN/community-builds-microservice/common/discovery"
	// pb "github.com/darkphotonKN/community-builds-microservice/common/api/proto/item"
)

const (
	serviceName = "item-service"
)

type Client struct {
	registry discovery.Registry
}

func NewClient(registry discovery.Registry) ItemClient {
	return &Client{
		registry: registry,
	}
}

func (c *Client) CreateItem(ctx context.Context, req *pb.CreateItemRequest) (*pb.CreateItemResponse, error) {

	// connection instance created through service discovery first
	// searches for the service registered as "orders"
	conn, err := discovery.ServiceConnection(ctx, serviceName, c.registry)

	if err != nil {
		return nil, fmt.Errorf("failed to connect to item service: %w", err)
	}
	defer conn.Close()

	client := pb.NewItemServiceClient(conn)

	// create client to interface with through service discovery connection
	item, err := client.CreateItem(ctx, &pb.CreateItemRequest{
		Name:     req.Name,
		Category: req.Category,
		Class:    req.Class,
		Type:     req.Type,
		ImageURL: req.ImageURL,
	})

	fmt.Printf("Creating item %+v through gateway after service discovery\n", item)

	return item, nil
}

func (c *Client) GetItems(ctx context.Context, req *pb.CreateItemRequest) (*pb.CreateItemResponse, error) {

	return nil, nil
}

func (c *Client) GenerateUniqueItems(ctx context.Context) (*pb.GenerateUniqueItemsResponse, error) {

	conn, err := discovery.ServiceConnection(ctx, serviceName, c.registry)

	if err != nil {
		return nil, fmt.Errorf("failed to connect to item service: %w", err)
	}

	defer conn.Close()

	return nil, nil
}
