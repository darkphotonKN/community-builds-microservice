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

// func (c *Client) GetBuildsByMemberId(ctx context.Context, req *pb.GetBuildsByMemberIdRequest) (*pb.GetBuildsByMemberIdResponse, error) {

// 	// connection instance created through service discovery first
// 	// searches for the service registered as "orders"
// 	conn, err := discovery.ServiceConnection(ctx, serviceName, c.registry)

// 	if err != nil {
// 		return nil, fmt.Errorf("failed to connect to item service: %w", err)
// 	}
// 	defer conn.Close()

// 	client := pb.NewBuildServiceClient(conn)

// 	// create client to interface with through service discovery connection
// 	builds, err := client.GetBuildsByMemberId(ctx, &pb.GetBuildsByMemberIdRequest{
// 		MemberId: req.MemberId,
// 	})

// 	fmt.Printf("Get builds %+v through gateway after service discovery\n", builds)

// 	return builds, nil
// }

func (c *Client) GetCommunityBuilds(ctx context.Context, req *pb.GetCommunityBuildsRequest) (*pb.GetCommunityBuildsResponse, error) {

	// connection instance created through service discovery first
	// searches for the service registered as "orders"
	conn, err := discovery.ServiceConnection(ctx, serviceName, c.registry)

	if err != nil {
		return nil, fmt.Errorf("failed to connect to item service: %w", err)
	}
	defer conn.Close()

	client := pb.NewBuildServiceClient(conn)

	// create client to interface with through service discovery connection
	builds, err := client.GetCommunityBuilds(ctx, req)

	fmt.Printf("Get community builds %+v through gateway after service discovery\n", builds)

	return builds, nil
}

func (c *Client) GetBuildInfo(ctx context.Context, req *pb.GetBuildInfoRequest) (*pb.GetBuildInfoResponse, error) {

	// connection instance created through service discovery first
	// searches for the service registered as "orders"
	conn, err := discovery.ServiceConnection(ctx, serviceName, c.registry)

	if err != nil {
		return nil, fmt.Errorf("failed to connect to item service: %w", err)
	}
	defer conn.Close()

	client := pb.NewBuildServiceClient(conn)

	info, err := client.GetBuildInfo(ctx, req)

	fmt.Printf("Get build info %+v through gateway after service discovery\n", info)

	return info, nil
}

func (c *Client) GetBuildsForMember(ctx context.Context, req *pb.GetBuildsForMemberRequest) (*pb.GetBuildsForMemberResponse, error) {

	// connection instance created through service discovery first
	// searches for the service registered as "orders"
	conn, err := discovery.ServiceConnection(ctx, serviceName, c.registry)

	if err != nil {
		return nil, fmt.Errorf("failed to connect to item service: %w", err)
	}
	defer conn.Close()

	client := pb.NewBuildServiceClient(conn)

	builds, err := client.GetBuildsForMember(ctx, req)

	fmt.Printf("Get builds %+v through gateway after service discovery\n", builds)

	return builds, nil
}

func (c *Client) GetBuildInfoForMember(ctx context.Context, req *pb.GetBuildInfoForMemberRequest) (*pb.GetBuildInfoForMemberResponse, error) {

	// connection instance created through service discovery first
	// searches for the service registered as "orders"
	conn, err := discovery.ServiceConnection(ctx, serviceName, c.registry)

	if err != nil {
		return nil, fmt.Errorf("failed to connect to item service: %w", err)
	}
	defer conn.Close()

	client := pb.NewBuildServiceClient(conn)

	build, err := client.GetBuildInfoForMember(ctx, req)

	fmt.Printf("Get build info %+v through gateway after service discovery\n", build)

	return build, nil
}

func (c *Client) PublishBuild(ctx context.Context, req *pb.PublishBuildRequest) (*pb.PublishBuildResponse, error) {

	// connection instance created through service discovery first
	// searches for the service registered as "orders"
	conn, err := discovery.ServiceConnection(ctx, serviceName, c.registry)

	if err != nil {
		return nil, fmt.Errorf("failed to connect to item service: %w", err)
	}
	defer conn.Close()

	client := pb.NewBuildServiceClient(conn)

	build, err := client.PublishBuild(ctx, req)

	fmt.Printf("Pubishe build %+v through gateway after service discovery\n", build)

	return build, nil
}

func (c *Client) UpdateBuild(ctx context.Context, req *pb.UpdateBuildRequest) (*pb.UpdateBuildResponse, error) {

	// connection instance created through service discovery first
	// searches for the service registered as "orders"
	conn, err := discovery.ServiceConnection(ctx, serviceName, c.registry)

	if err != nil {
		return nil, fmt.Errorf("failed to connect to item service: %w", err)
	}
	defer conn.Close()

	client := pb.NewBuildServiceClient(conn)

	build, err := client.UpdateBuild(ctx, req)

	fmt.Printf("Update build %+v through gateway after service discovery\n", build)

	return build, nil
}

func (c *Client) AddSkillLinksToBuild(ctx context.Context, req *pb.AddSkillLinksToBuildRequest) (*pb.AddSkillLinksToBuildResponse, error) {

	// connection instance created through service discovery first
	// searches for the service registered as "orders"
	conn, err := discovery.ServiceConnection(ctx, serviceName, c.registry)

	if err != nil {
		return nil, fmt.Errorf("failed to connect to item service: %w", err)
	}
	defer conn.Close()

	client := pb.NewBuildServiceClient(conn)

	build, err := client.AddSkillLinksToBuild(ctx, req)

	fmt.Printf("Update build %+v through gateway after service discovery\n", build)

	return build, nil
}

func (c *Client) UpdateItemSetsToBuild(ctx context.Context, req *pb.UpdateItemSetsToBuildRequest) (*pb.UpdateItemSetsToBuildResponse, error) {

	// connection instance created through service discovery first
	// searches for the service registered as "orders"
	conn, err := discovery.ServiceConnection(ctx, serviceName, c.registry)

	if err != nil {
		return nil, fmt.Errorf("failed to connect to item service: %w", err)
	}
	defer conn.Close()

	client := pb.NewBuildServiceClient(conn)

	build, err := client.UpdateItemSetsToBuild(ctx, req)

	fmt.Printf("Update build %+v through gateway after service discovery\n", build)

	return build, nil
}

func (c *Client) DeleteBuildByMember(ctx context.Context, req *pb.DeleteBuildByMemberRequest) (*pb.DeleteBuildByMemberResponse, error) {

	// connection instance created through service discovery first
	// searches for the service registered as "orders"
	conn, err := discovery.ServiceConnection(ctx, serviceName, c.registry)

	if err != nil {
		return nil, fmt.Errorf("failed to connect to item service: %w", err)
	}
	defer conn.Close()

	client := pb.NewBuildServiceClient(conn)

	build, err := client.DeleteBuildByMember(ctx, req)

	fmt.Printf("Update build %+v through gateway after service discovery\n", build)

	return build, nil
}
