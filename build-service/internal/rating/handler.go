package rating

import (
	"context"

	pb "github.com/darkphotonKN/community-builds-microservice/common/api/proto/rating"
	"github.com/darkphotonKN/community-builds-microservice/common/constants/models"
	"github.com/google/uuid"
)

type Handler struct {
	service Service
	pb.UnimplementedRatingServiceServer
}

type Service interface {
	CreateRatingByBuildId(ctx context.Context, req *pb.CreateRatingByBuildIdRequest) (*pb.CreateRatingByBuildIdResponse, error)
	GetAllRatingsForProduct(userId uuid.UUID) (*[]models.Rating, error)
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateRatingByBuildId(ctx context.Context, req *pb.CreateRatingByBuildIdRequest) (*pb.CreateRatingByBuildIdResponse, error) {
	return h.service.CreateRatingByBuildId(ctx, req)
}
