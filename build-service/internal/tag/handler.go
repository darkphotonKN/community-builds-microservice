package tag

import (
	"context"

	pb "github.com/darkphotonKN/community-builds-microservice/common/api/proto/tag"
	// "github.com/google/uuid"
)

type Handler struct {
	service Service
	pb.UnimplementedTagServiceServer
}

type Service interface {
	CreateTag(ctx context.Context, req *pb.CreateTagRequest) (*pb.CreateTagResponse, error)
	GetTags(ctx context.Context, req *pb.GetTagsRequest) (*pb.GetTagsResponse, error)
	UpdateTag(ctx context.Context, req *pb.UpdateTagRequest) (*pb.UpdateTagResponse, error)
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// --- ADMIN HANDLERS ---
func (h *Handler) CreateHandler(ctx context.Context, req *pb.CreateTagRequest) (*pb.CreateTagResponse, error) {
	return h.service.CreateTag(ctx, req)
}

func (h *Handler) GetTagsHandler(ctx context.Context, req *pb.GetTagsRequest) (*pb.GetTagsResponse, error) {
	return h.service.GetTags(ctx, req)

}

func (h *Handler) UpdateTagsHandler(ctx context.Context, req *pb.UpdateTagRequest) (*pb.UpdateTagResponse, error) {
	return h.service.UpdateTag(ctx, req)
}
