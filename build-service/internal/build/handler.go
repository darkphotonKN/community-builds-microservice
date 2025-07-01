package build

import (
	"golang.org/x/net/context"

	pb "github.com/darkphotonKN/community-builds-microservice/common/api/proto/build"
	"github.com/darkphotonKN/community-builds-microservice/common/constants/models"
	"github.com/google/uuid"
	// amqp "github.com/rabbitmq/amqp091-go"
)

type Handler struct {
	service Service
	pb.UnimplementedBuildServiceServer
}

type Service interface {
	CreateBuild(ctx context.Context, req *pb.CreateBuildRequest) (*pb.CreateBuildResponse, error)
	GetBuildsForMember(ctx context.Context, req *pb.GetBuildsForMemberRequest) (*pb.GetBuildsForMemberResponse, error)
	GetCommunityBuilds(ctx context.Context, req *pb.GetCommunityBuildsRequest) (*pb.GetCommunityBuildsResponse, error)
	GetBuildInfo(ctx context.Context, req *pb.GetBuildInfoRequest) (*pb.GetBuildInfoResponse, error)
	GetBuildTagsForMemberById(memberId uuid.UUID) (*[]models.Tag, error)
	GetBuildInfoForMember(ctx context.Context, req *pb.GetBuildInfoForMemberRequest) (*pb.GetBuildInfoForMemberResponse, error)
	PublishBuild(ctx context.Context, req *pb.PublishBuildRequest) (*pb.PublishBuildResponse, error)
	UpdateBuild(ctx context.Context, req *pb.UpdateBuildRequest) (*pb.UpdateBuildResponse, error)
	AddSkillLinksToBuild(ctx context.Context, req *pb.AddSkillLinksToBuildRequest) (*pb.AddSkillLinksToBuildResponse, error)
	UpdateItemSetsToBuild(ctx context.Context, req *pb.UpdateItemSetsToBuildRequest) (*pb.UpdateItemSetsToBuildResponse, error)
	DeleteBuildByMember(ctx context.Context, req *pb.DeleteBuildByMemberRequest) (*pb.DeleteBuildByMemberResponse, error)
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateBuild(ctx context.Context, req *pb.CreateBuildRequest) (*pb.CreateBuildResponse, error) {
	return h.service.CreateBuild(ctx, req)
}

func (h *Handler) GetBuildsForMember(ctx context.Context, req *pb.GetBuildsForMemberRequest) (*pb.GetBuildsForMemberResponse, error) {
	return h.service.GetBuildsForMember(ctx, req)
}

func (h *Handler) GetCommunityBuilds(ctx context.Context, req *pb.GetCommunityBuildsRequest) (*pb.GetCommunityBuildsResponse, error) {
	return h.service.GetCommunityBuilds(ctx, req)
}

func (h *Handler) GetBuildInfo(ctx context.Context, req *pb.GetBuildInfoRequest) (*pb.GetBuildInfoResponse, error) {
	return h.service.GetBuildInfo(ctx, req)
}

func (h *Handler) GetBuildInfoForMember(ctx context.Context, req *pb.GetBuildInfoForMemberRequest) (*pb.GetBuildInfoForMemberResponse, error) {
	return h.service.GetBuildInfoForMember(ctx, req)
}

func (h *Handler) PublishBuild(ctx context.Context, req *pb.PublishBuildRequest) (*pb.PublishBuildResponse, error) {
	return h.service.PublishBuild(ctx, req)
}

func (h *Handler) UpdateBuild(ctx context.Context, req *pb.UpdateBuildRequest) (*pb.UpdateBuildResponse, error) {
	return h.service.UpdateBuild(ctx, req)
}

func (h *Handler) AddSkillLinksToBuild(ctx context.Context, req *pb.AddSkillLinksToBuildRequest) (*pb.AddSkillLinksToBuildResponse, error) {
	return h.service.AddSkillLinksToBuild(ctx, req)
}

func (h *Handler) UpdateItemSetsToBuild(ctx context.Context, req *pb.UpdateItemSetsToBuildRequest) (*pb.UpdateItemSetsToBuildResponse, error) {
	return h.service.UpdateItemSetsToBuild(ctx, req)
}

func (h *Handler) DeleteBuildByMember(ctx context.Context, req *pb.DeleteBuildByMemberRequest) (*pb.DeleteBuildByMemberResponse, error) {
	return h.service.DeleteBuildByMember(ctx, req)
}
