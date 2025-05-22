package item

import (
	"fmt"

	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	// "google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/darkphotonKN/community-builds-microservice/common/api/proto/item"
	"github.com/jmoiron/sqlx"
)

type Handler struct {
	service Service
	pb.UnimplementedItemServiceServer
}

type Service interface {
	InitCrawling(*sqlx.DB) error
	CreateItem(ctx context.Context, example *pb.CreateItemRequest) (*pb.CreateItemResponse, error)
	// GenerateUniqueItems(ctx context.Context) (*pb.GenerateUniqueItemsResponse, error)
	// GetItems(ctx context.Context, id uuid.UUID) (*pb.CreateItemResponse, error)
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (s *Handler) CreateItem(ctx context.Context, req *pb.CreateItemRequest) (*pb.CreateItemResponse, error) {
	fmt.Println("收到 gRPC 請求:", req)

	_, err := s.service.CreateItem(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "創建 item 時發生錯誤: %v", err)
	}
	return &pb.CreateItemResponse{Message: fmt.Sprintf("成功創建item")}, nil
}

// func (s *Handler) GetUniqueItems(ctx context.Context, req *emptypb.Empty) (*pb.GenerateUniqueItemsResponse, error) {

// 	_, err := s.service.GenerateUniqueItems(ctx)
// 	if err != nil {
// 		return nil, status.Errorf(codes.Internal, "創建 item 時發生錯誤: %v", err)
// 	}
// 	return &pb.GenerateUniqueItemsResponse{Message: fmt.Sprintf("成功抓取unique items")}, nil
// }
