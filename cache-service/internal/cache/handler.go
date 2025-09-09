package cache

import (
	"context"

	pb "github.com/darkphotonKN/community-builds-microservice/common/api/proto/cache"
)

type Handler struct {
	pb.UnimplementedCacheServiceServer
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) GetCache(ctx context.Context, req *pb.GetCacheRequest) (*pb.GetCacheResponse, error) {
	return h.service.GetCache(ctx, req)
}

func (h *Handler) Hello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	return h.service.Hello(ctx, req)
}

func (h *Handler) SetCache(ctx context.Context, req *pb.SetCacheRequest) (*pb.SetCacheResponse, error) {
	return h.service.SetCache(ctx, req)
}

func (h *Handler) DeleteCache(ctx context.Context, req *pb.DeleteCacheRequest) (*pb.DeleteCacheResponse, error) {
	return h.service.DeleteCache(ctx, req)
}
