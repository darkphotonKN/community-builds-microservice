package class

import (
	"context"
	"fmt"

	pb "github.com/darkphotonKN/community-builds-microservice/common/api/proto/class"
	"github.com/darkphotonKN/community-builds-microservice/common/discovery"
)

const (
	serviceName = "class-service"
)

type Client struct {
	registry discovery.Registry
}

func NewClient(registry discovery.Registry) ClassClient {
	return &Client{
		registry: registry,
	}
}

func (c *Client) GetClassesAndAscendancies(ctx context.Context, req *pb.GetClassesAndAscendanciesRequest) (*pb.GetClassesAndAscendanciesResponse, error) {

	// connection instance created through service discovery first
	// searches for the service registered as "orders"
	conn, err := discovery.ServiceConnection(ctx, serviceName, c.registry)

	if err != nil {
		return nil, fmt.Errorf("failed to connect to item service: %w", err)
	}
	defer conn.Close()

	client := pb.NewClassServiceClient(conn)

	// create client to interface with through service discovery connection
	item, err := client.GetClassesAndAscendancies(ctx, req)

	fmt.Printf("Creating item %+v through gateway after service discovery\n", item)

	return item, nil
}
