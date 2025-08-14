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
	GetBaseItems(ctx context.Context, req *pb.GetBaseItemsRequest) (*pb.GetBaseItemsResponse, error)
	GetItemMods(ctx context.Context, req *pb.GetItemModsRequest) (*pb.GetItemModsResponse, error)
	GetMemberRareItems(ctx context.Context, req *pb.GetMemberRareItemsRequest) (*pb.GetMemberRareItemsResponse, error)
	CrawlingAndAddUniqueItemsService(db *sqlx.DB) error
	CrawlingAndAddBaseItemsService(db *sqlx.DB) error
	CrawlingAndAddItemModsService(db *sqlx.DB) error
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

func (h *Handler) GetBaseItems(ctx context.Context, req *pb.GetBaseItemsRequest) (*pb.GetBaseItemsResponse, error) {
	return h.service.GetBaseItems(ctx, req)
}

func (h *Handler) GetItemMods(ctx context.Context, req *pb.GetItemModsRequest) (*pb.GetItemModsResponse, error) {
	return h.service.GetItemMods(ctx, req)
}

func (h *Handler) GetMemberRareItems(ctx context.Context, req *pb.GetMemberRareItemsRequest) (*pb.GetMemberRareItemsResponse, error) {
	return h.service.GetMemberRareItems(ctx, req)
}
