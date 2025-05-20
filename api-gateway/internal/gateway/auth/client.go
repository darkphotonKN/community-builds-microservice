package auth

import (
	"context"

	pb "github.com/darkphotonKN/community-builds-microservice/common/api/proto/auth"
	"github.com/darkphotonKN/community-builds-microservice/common/discovery"
)

const (
	serviceName = "auth"
)

type Client struct {
	registry discovery.Registry
}

func NewClient(registry discovery.Registry) AuthClient {
	return &Client{
		registry: registry,
	}
}

func CreateMember(ctx context.Context, req *pb.CreateMemberRequest) (*pb.Member, error) {

	return nil, nil
}
