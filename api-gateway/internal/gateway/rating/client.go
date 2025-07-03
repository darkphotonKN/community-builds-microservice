package rating

import (
	"context"
	"fmt"

	pb "github.com/darkphotonKN/community-builds-microservice/common/api/proto/rating"
	"github.com/darkphotonKN/community-builds-microservice/common/discovery"
)

const (
	serviceName = "class-service"
)

type Client struct {
	registry discovery.Registry
}

func NewClient(registry discovery.Registry) RatingClient {
	return &Client{
		registry: registry,
	}
}

func (c *Client) CreateRatingByBuildId(ctx context.Context, req *pb.CreateRatingByBuildIdRequest) (*pb.CreateRatingByBuildIdResponse, error) {

	// connection instance created through service discovery first
	// searches for the service registered as "orders"
	conn, err := discovery.ServiceConnection(ctx, serviceName, c.registry)

	if err != nil {
		return nil, fmt.Errorf("failed to connect to item service: %w", err)
	}
	defer conn.Close()

	client := pb.NewRatingServiceClient(conn)

	// create client to interface with through service discovery connection
	item, err := client.CreateRatingByBuildId(ctx, req)

	fmt.Printf("Creating rating by build id %+v through gateway after service discovery\n", item)

	return item, nil
}
