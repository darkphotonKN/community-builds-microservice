package tag

import (
	"context"

	pb "github.com/darkphotonKN/community-builds-microservice/common/api/proto/tag"
	"github.com/google/uuid"
)

type CreateTagRequest struct {
	Name string `json:"name" binding:"required" db:"name"`
}

type UpdateTagRequest struct {
	ID   uuid.UUID `db:"id" json:"id"`
	Name string    `json:"name" binding:"required" db:"name"`
}

type UpdateTagParams struct {
	ID   uuid.UUID `db:"id" json:"id"`
	Name string    `db:"name" json:"name"`
}

type TagClient interface {
	CreateTag(ctx context.Context, req *pb.CreateTagRequest) (*pb.CreateTagResponse, error)
	GetTags(ctx context.Context, req *pb.GetTagsRequest) (*pb.GetTagsResponse, error)
	UpdateTag(ctx context.Context, req *pb.UpdateTagRequest) (*pb.UpdateTagResponse, error)
}
