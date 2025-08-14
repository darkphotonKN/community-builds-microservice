package item

import (
	"context"
	"fmt"

	pb "github.com/darkphotonKN/community-builds-microservice/common/api/proto/item"
	"github.com/darkphotonKN/community-builds-microservice/common/discovery"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	item, err := client.CreateItem(ctx, req)

	fmt.Printf("Creating item %+v through gateway after service discovery\n", item)
	if err != nil {
		return nil, fmt.Errorf("failed to create item: %w", err)
	}
	return item, nil
}

func (c *Client) UpdateItem(ctx context.Context, req *pb.UpdateItemRequest) (*pb.UpdateItemResponse, error) {

	// connection instance created through service discovery first
	// searches for the service registered as "orders"
	conn, err := discovery.ServiceConnection(ctx, serviceName, c.registry)

	if err != nil {
		return nil, fmt.Errorf("failed to connect to item service: %w", err)
	}
	defer conn.Close()

	client := pb.NewItemServiceClient(conn)

	// create client to interface with through service discovery connection
	item, err := client.UpdateItem(ctx, &pb.UpdateItemRequest{
		Name:     req.Name,
		Category: req.Category,
		Class:    req.Class,
		Type:     req.Type,
		ImageURL: req.ImageURL,
	})

	fmt.Printf("Creating item %+v through gateway after service discovery\n", item)
	if err != nil {
		return nil, fmt.Errorf("failed to update item: %w", err)
	}
	return item, nil
}

func (c *Client) GetItems(ctx context.Context, req *pb.GetItemsRequest) (*pb.GetItemsResponse, error) {

	// connection instance created through service discovery first
	// searches for the service registered as "orders"
	conn, err := discovery.ServiceConnection(ctx, serviceName, c.registry)

	if err != nil {
		return nil, fmt.Errorf("failed to connect to item service: %w", err)
	}
	defer conn.Close()

	client := pb.NewItemServiceClient(conn)

	fmt.Println("req.Slot:", req.Slot)
	// create client to interface with through service discovery connection
	items, err := client.GetItems(ctx, &pb.GetItemsRequest{
		Slot: req.Slot,
	})

	fmt.Printf("Get items %+v through gateway after service discovery\n", items)
	if err != nil {
		return nil, fmt.Errorf("failed to get items: %w", err)
	}
	return items, nil
}

func (c *Client) CreateRareItem(ctx context.Context, req *pb.CreateRareItemRequest) (*pb.CreateRareItemResponse, error) {

	// connection instance created through service discovery first
	// searches for the service registered as "orders"
	conn, err := discovery.ServiceConnection(ctx, serviceName, c.registry)

	if err != nil {
		return nil, fmt.Errorf("failed to connect to item service: %w", err)
	}
	defer conn.Close()

	client := pb.NewItemServiceClient(conn)

	// create client to interface with through service discovery connection
	res, err := client.CreateRareItem(ctx, req)

	if err != nil {
		st, _ := status.FromError(err)
		switch st.Code() {
		case codes.InvalidArgument:
			return nil, fmt.Errorf("invalid input: %w", err)
		case codes.NotFound:
			return nil, fmt.Errorf("item not found: %w", err)
		default:
			return nil, fmt.Errorf("grpc error: %w", err)
		}
	}

	// fmt.Printf("Get items %+v through gateway after service discovery\n", res)
	if err != nil {
		return nil, fmt.Errorf("failed to create rare item: %w", err)
	}
	return res, nil
}

func (c *Client) GetBaseItems(ctx context.Context, req *pb.GetBaseItemsRequest) (*pb.GetBaseItemsResponse, error) {

	// connection instance created through service discovery first
	// searches for the service registered as "orders"
	conn, err := discovery.ServiceConnection(ctx, serviceName, c.registry)

	if err != nil {
		return nil, fmt.Errorf("failed to connect to item service: %w", err)
	}
	defer conn.Close()

	client := pb.NewItemServiceClient(conn)

	// create client to interface with through service discovery connection
	items, err := client.GetBaseItems(ctx, req)

	// fmt.Printf("Get items %+v through gateway after service discovery\n", items)
	if err != nil {
		return nil, fmt.Errorf("failed to get items: %w", err)
	}
	return items, nil
}

func (c *Client) GetItemMods(ctx context.Context, req *pb.GetItemModsRequest) (*pb.GetItemModsResponse, error) {

	// connection instance created through service discovery first
	// searches for the service registered as "orders"
	conn, err := discovery.ServiceConnection(ctx, serviceName, c.registry)

	if err != nil {
		return nil, fmt.Errorf("failed to connect to item service: %w", err)
	}
	defer conn.Close()

	client := pb.NewItemServiceClient(conn)

	// create client to interface with through service discovery connection
	items, err := client.GetItemMods(ctx, req)

	// fmt.Printf("Get item mods %+v through gateway after service discovery\n", items)
	if err != nil {
		return nil, fmt.Errorf("failed to get item mods: %w", err)
	}
	return items, nil
}

func (c *Client) GetMemberRareItems(ctx context.Context, req *pb.GetMemberRareItemsRequest) (*pb.GetMemberRareItemsResponse, error) {

	// connection instance created through service discovery first
	// searches for the service registered as "orders"
	conn, err := discovery.ServiceConnection(ctx, serviceName, c.registry)

	if err != nil {
		return nil, fmt.Errorf("failed to connect to item service: %w", err)
	}
	defer conn.Close()

	client := pb.NewItemServiceClient(conn)

	// create client to interface with through service discovery connection
	res, err := client.GetMemberRareItems(ctx, req)

	if err != nil {
		return nil, fmt.Errorf("failed to get member rare item: %w", err)
	}

	fmt.Printf("Get member rare item %+v through gateway after service discovery\n", res)
	return res, nil
}
