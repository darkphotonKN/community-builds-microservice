package cache

import (
	"context"
	"fmt"
	"time"

	pb "github.com/darkphotonKN/community-builds-microservice/common/api/proto/cache"
)

type Service interface {
	DeleteCache(ctx context.Context, req *pb.DeleteCacheRequest) (*pb.DeleteCacheResponse, error)
	SetCache(ctx context.Context, req *pb.SetCacheRequest) (*pb.SetCacheResponse, error)
	GetCache(ctx context.Context, req *pb.GetCacheRequest) (*pb.GetCacheResponse, error)
	Hello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error)
}

type service struct{}

func NewService() Service {
	return &service{}
}

func (s *service) Hello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	message := fmt.Sprintf("Hello %s! Cache service is running.", req.Name)
	return &pb.HelloResponse{
		Message:   message,
		Timestamp: time.Now().Format(time.RFC3339),
	}, nil
}

func (s *service) GetCache(ctx context.Context, req *pb.GetCacheRequest) (*pb.GetCacheResponse, error) {
	return nil, nil
}

func (s *service) SetCache(ctx context.Context, req *pb.SetCacheRequest) (*pb.SetCacheResponse, error) {
	return nil, nil
}

func (s *service) DeleteCache(ctx context.Context, req *pb.DeleteCacheRequest) (*pb.DeleteCacheResponse, error) {
	return nil, nil
}
