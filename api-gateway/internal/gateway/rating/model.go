package rating

import (
	"context"

	pb "github.com/darkphotonKN/community-builds-microservice/common/api/proto/rating"
	"github.com/google/uuid"
)

type CreateRatingRequest struct {
	BuildId uuid.UUID `db:"build_id" binding:"required,uuid" json:"buildId"`
	// Category string `db:"category" binding:"required,ratingCategory" json:"category"`
	Value int32 `db:"value" binding:"required,min=1,max=10" json:"value"`
}

type RatingByCategoryRes struct {
	Value int `db:"value"`
}

type RatingClient interface {
	CreateRatingByBuildId(ctx context.Context, req *pb.CreateRatingByBuildIdRequest) (*pb.CreateRatingByBuildIdResponse, error)
}
