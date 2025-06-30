package skill

import (
	"context"
	"fmt"

	pb "github.com/darkphotonKN/community-builds-microservice/common/api/proto/skill"
	"github.com/darkphotonKN/community-builds-microservice/common/discovery"
	// pb "github.com/darkphotonKN/community-builds-microservice/common/api/proto/item"
)

const (
	serviceName = "skill-service"
)

type Client struct {
	registry discovery.Registry
}

func NewClient(registry discovery.Registry) SkillClient {
	return &Client{
		registry: registry,
	}
}

func (c *Client) CreateSkill(ctx context.Context, req *pb.CreateSkillRequest) (*pb.CreateSkillResponse, error) {

	// connection instance created through service discovery first
	// searches for the service registered as "orders"
	conn, err := discovery.ServiceConnection(ctx, serviceName, c.registry)

	if err != nil {
		return nil, fmt.Errorf("failed to connect to item service: %w", err)
	}
	defer conn.Close()

	client := pb.NewSkillServiceClient(conn)

	// create client to interface with through service discovery connection
	skill, err := client.CreateSkill(ctx, req)

	fmt.Printf("Creating skill %+v through gateway after service discovery\n", skill)

	return skill, nil
}

func (c *Client) GetSkills(ctx context.Context, req *pb.GetSkillsRequest) (*pb.GetSkillsResponse, error) {

	// connection instance created through service discovery first
	// searches for the service registered as "orders"
	conn, err := discovery.ServiceConnection(ctx, serviceName, c.registry)

	if err != nil {
		return nil, fmt.Errorf("failed to connect to item service: %w", err)
	}
	defer conn.Close()

	client := pb.NewSkillServiceClient(conn)

	// create client to interface with through service discovery connection
	items, err := client.GetSkills(ctx, &pb.GetSkillsRequest{})

	fmt.Printf("Get items %+v through gateway after service discovery\n", items)

	return items, nil
}
