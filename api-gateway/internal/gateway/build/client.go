package build

import (
	"context"
	"fmt"

	pb "github.com/darkphotonKN/community-builds-microservice/common/api/proto/build"
	"github.com/darkphotonKN/community-builds-microservice/common/discovery"
)

const (
	serviceName = "build-service"
)

type Client struct {
	registry discovery.Registry
}

func NewClient(registry discovery.Registry) BuildClient {
	return &Client{
		registry: registry,
	}
}

func (c *Client) CreateBuild(ctx context.Context, req *pb.CreateBuildRequest) (*pb.CreateBuildResponse, error) {

	// connection instance created through service discovery first
	// searches for the service registered as "orders"
	conn, err := discovery.ServiceConnection(ctx, serviceName, c.registry)

	if err != nil {
		return nil, fmt.Errorf("failed to connect to item service: %w", err)
	}
	defer conn.Close()

	client := pb.NewBuildServiceClient(conn)

	// create client to interface with through service discovery connection
	build, err := client.CreateBuild(ctx, req)

	fmt.Printf("Creating build %+v through gateway after service discovery\n", build)

	return build, nil
}

func (c *Client) GetBuildsByMemberId(ctx context.Context, req *pb.GetBuildsByMemberIdRequest) (*pb.GetBuildsByMemberIdResponse, error) {

	// connection instance created through service discovery first
	// searches for the service registered as "orders"
	conn, err := discovery.ServiceConnection(ctx, serviceName, c.registry)

	if err != nil {
		return nil, fmt.Errorf("failed to connect to item service: %w", err)
	}
	defer conn.Close()

	client := pb.NewBuildServiceClient(conn)

	// create client to interface with through service discovery connection
	items, err := client.GetBuildsByMemberId(ctx, &pb.GetBuildsByMemberIdRequest{
		MemberId: req.MemberId,
	})

	fmt.Printf("Get items %+v through gateway after service discovery\n", items)

	return items, nil
}
