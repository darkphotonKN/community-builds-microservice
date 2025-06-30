package tag

import (
	"context"
	"fmt"

	pb "github.com/darkphotonKN/community-builds-microservice/common/api/proto/tag"
	"github.com/darkphotonKN/community-builds-microservice/common/discovery"
)

const (
	serviceName = "tag-service"
)

type Client struct {
	registry discovery.Registry
}

func NewClient(registry discovery.Registry) TagClient {
	return &Client{
		registry: registry,
	}
}

func (c *Client) CreateTag(ctx context.Context, req *pb.CreateTagRequest) (*pb.CreateTagResponse, error) {

	// connection instance created through service discovery first
	// searches for the service registered as "orders"
	conn, err := discovery.ServiceConnection(ctx, serviceName, c.registry)

	if err != nil {
		return nil, fmt.Errorf("failed to connect to item service: %w", err)
	}
	defer conn.Close()

	client := pb.NewTagServiceClient(conn)

	// create client to interface with through service discovery connection
	tag, err := client.CreateTag(ctx, req)

	fmt.Printf("Creating tag %+v through gateway after service discovery\n", tag)

	return tag, nil
}

func (c *Client) GetTags(ctx context.Context, req *pb.GetTagsRequest) (*pb.GetTagsResponse, error) {

	// connection instance created through service discovery first
	// searches for the service registered as "orders"
	conn, err := discovery.ServiceConnection(ctx, serviceName, c.registry)

	if err != nil {
		return nil, fmt.Errorf("failed to connect to item service: %w", err)
	}
	defer conn.Close()

	client := pb.NewTagServiceClient(conn)

	// create client to interface with through service discovery connection
	tags, err := client.GetTags(ctx, req)

	fmt.Printf("Get tags %+v through gateway after service discovery\n", tags)

	return tags, nil
}

func (c *Client) UpdateTag(ctx context.Context, req *pb.UpdateTagRequest) (*pb.UpdateTagResponse, error) {

	// connection instance created through service discovery first
	// searches for the service registered as "orders"
	conn, err := discovery.ServiceConnection(ctx, serviceName, c.registry)

	if err != nil {
		return nil, fmt.Errorf("failed to connect to item service: %w", err)
	}
	defer conn.Close()

	client := pb.NewTagServiceClient(conn)

	// create client to interface with through service discovery connection
	tag, err := client.UpdateTag(ctx, req)

	fmt.Printf("Creating tag %+v through gateway after service discovery\n", tag)

	return tag, nil
}
