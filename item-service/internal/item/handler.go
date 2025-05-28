package item

import (
	"golang.org/x/net/context"

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
	GetItemsService(ctx context.Context, req *pb.GetItemsRequest) (*pb.GetItemsResponse, error)
	CreateItemService(ctx context.Context, req *pb.CreateItemRequest) (*pb.CreateItemResponse, error)
	UpdateItemService(ctx context.Context, req *pb.UpdateItemRequest) (*pb.UpdateItemResponse, error)
	CreateRareItemService(ctx context.Context, req *pb.CreateRareItemRequest) (*pb.CreateRareItemResponse, error)
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetItems(ctx context.Context, req *pb.GetItemsRequest) (*pb.GetItemsResponse, error) {
	return h.service.GetItemsService(ctx, req)
}

func (h *Handler) CreateItem(ctx context.Context, req *pb.CreateItemRequest) (*pb.CreateItemResponse, error) {
	return h.service.CreateItemService(ctx, req)
}

func (h *Handler) UpdateItem(ctx context.Context, req *pb.UpdateItemRequest) (*pb.UpdateItemResponse, error) {
	return h.service.UpdateItemService(ctx, req)
}

func (h *Handler) CreateRareItem(ctx context.Context, req *pb.CreateRareItemRequest) (*pb.CreateRareItemResponse, error) {
	return h.service.CreateRareItemService(ctx, req)
}
