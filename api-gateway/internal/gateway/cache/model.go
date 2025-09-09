package cache

import (
	"context"

	pb "github.com/darkphotonKN/community-builds-microservice/common/api/proto/cache"
)

type CacheClient interface {
	// DeleteCache(ctx context.Context, req *pb.DeleteCacheRequest) (*pb.DeleteCacheResponse, error)
	// SetCache(ctx context.Context, req *pb.SetCacheRequest) (*pb.SetCacheResponse, error)
	// GetCache(ctx context.Context, req *pb.GetCacheRequest) (*pb.GetCacheResponse, error)
	Hello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error)
}
