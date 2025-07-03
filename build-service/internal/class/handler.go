package class

import (
	"context"

	pb "github.com/darkphotonKN/community-builds-microservice/common/api/proto/class"
)

type Handler struct {
	service Service
	pb.UnimplementedClassServiceServer
}

type Service interface {
	GetClassesAndAscendancies(ctx context.Context, req *pb.GetClassesAndAscendanciesRequest) (*pb.GetClassesAndAscendanciesResponse, error)
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetClassesAndAscendancies(ctx context.Context, req *pb.GetClassesAndAscendanciesRequest) (*pb.GetClassesAndAscendanciesResponse, error) {
	return h.service.GetClassesAndAscendancies(ctx, req)
}
