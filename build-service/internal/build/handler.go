package build

import (
	"golang.org/x/net/context"

	pb "github.com/darkphotonKN/community-builds-microservice/common/api/proto/build"
	// amqp "github.com/rabbitmq/amqp091-go"
)

type Handler struct {
	service Service
	pb.UnimplementedBuildServiceServer
}

type Service interface {
	CreateBuild(ctx context.Context, req *pb.CreateBuildRequest) (*pb.CreateBuildResponse, error)
	GetBuildsByMemberId(ctx context.Context, req *pb.GetBuildsByMemberIdRequest) (*pb.GetBuildsByMemberIdResponse, error)
	GetCommunityBuilds(ctx context.Context, req *pb.GetCommunityBuildsRequest) (*pb.GetCommunityBuildsResponse, error)
	GetBuildInfo(ctx context.Context, req *pb.GetBuildInfoRequest) (*pb.GetBuildInfoResponse, error)
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateBuild(ctx context.Context, req *pb.CreateBuildRequest) (*pb.CreateBuildResponse, error) {
	return h.service.CreateBuild(ctx, req)
}

func (h *Handler) GetbuildsByMemberId(ctx context.Context, req *pb.GetBuildsByMemberIdRequest) (*pb.GetBuildsByMemberIdResponse, error) {
	return h.service.GetBuildsByMemberId(ctx, req)
}

func (h *Handler) GetCommunityBuilds(ctx context.Context, req *pb.GetCommunityBuildsRequest) (*pb.GetCommunityBuildsResponse, error) {
	return h.service.GetCommunityBuilds(ctx, req)
}

func (h *Handler) GetBuildInfo(ctx context.Context, req *pb.GetBuildInfoRequest) (*pb.GetBuildInfoResponse, error) {
	return h.service.GetBuildInfo(ctx, req)
}
