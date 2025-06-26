package item

import (
	"context"

	pb "github.com/darkphotonKN/community-builds-microservice/common/api/proto/item"
	"github.com/google/uuid"
)

type GetItemsRequest struct {
	Type string `json:"type"`
}
type CreateItemRequest struct {
	Category string `json:"category" binding:"required" db:"category"`
	Class    string `json:"class" binding:"required" db:"class"`
	Type     string `json:"type" binding:"required" db:"type"`
	Name     string `json:"name" binding:"required,min=2" db:"name"`
	ImageURL string `json:"imageUrl,omitempty" db:"image_url"`
	Slot     string `json:"slot,omitempty" db:"image_url"`
}

type UpdateItemReq struct {
	Category string `json:"category" binding:"required,category" db:"category"`
	Class    string `json:"class" binding:"required,class" db:"class"`
	Type     string `json:"type" binding:"required,type" db:"type"`
	Name     string `json:"name" binding:"required,min=2" db:"name"`
	ImageURL string `json:"imageUrl,omitempty" db:"image_url"`
}

type CreateRareItemReq struct {
	BaseItemId uuid.UUID `json:"baseItemId" db:"base_item_id"`
	ToList     bool      `json:"toList"`
	Name       string    `json:"name" db:"name"`
	Stats      []string  `json:"stats" db:"stats"`
}

type ItemClient interface {
	CreateItem(ctx context.Context, req *pb.CreateItemRequest) (*pb.CreateItemResponse, error)
	UpdateItem(ctx context.Context, req *pb.UpdateItemRequest) (*pb.UpdateItemResponse, error)
	GetItems(ctx context.Context, req *pb.GetItemsRequest) (*pb.GetItemsResponse, error)
	CreateRareItem(ctx context.Context, req *pb.CreateRareItemRequest) (*pb.CreateRareItemResponse, error)
}
