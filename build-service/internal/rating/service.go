package rating

import (
	"context"
	"fmt"

	// "github.com/darkphotonKN/community-builds-microservice/api-gateway/internal/build"
	pb "github.com/darkphotonKN/community-builds-microservice/common/api/proto/rating"
	"github.com/darkphotonKN/community-builds-microservice/common/constants/models"
	"github.com/darkphotonKN/community-builds-microservice/common/constants/types"
	"github.com/google/uuid"
)

type service struct {
	repo Repository
}

type Repository interface {
	CreateRatingForBuildById(memberId uuid.UUID, request CreateRatingRequest) error
	GetAllRatingsByMemberId(memberId uuid.UUID) (*[]models.Rating, error)
	GetAllRatingsByCategoryForBuild(buildId string, category types.RatingCategory) ([]int, error)
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

/**
* Posts single rating for a single product.
**/
func (s *service) CreateRatingByBuildId(ctx context.Context, req *pb.CreateRatingByBuildIdRequest) (*pb.CreateRatingByBuildIdResponse, error) {
	// create rating for build

	memberId, err := uuid.Parse(req.MemberId)
	if err != nil {
		return nil, err
	}

	request := CreateRatingRequest{
		BuildId:  req.BuildId,
		Category: req.Category,
		Value:    int(req.Value),
	}
	err = s.repo.CreateRatingForBuildById(memberId, request)

	if err != nil {
		return nil, err
	}

	// update build's average rating.

	// get all ratings of that category
	ratings, err := s.repo.GetAllRatingsByCategoryForBuild(request.BuildId, types.RatingCategory(request.Category))

	// average them
	var avgRating float32
	totalRating := 0
	noOfRatings := float32(len(ratings))

	for _, rating := range ratings {
		totalRating += rating
	}

	avgRating = float32(totalRating) / noOfRatings

	fmt.Printf("ratings %+v, avgRating: %f\n", ratings, avgRating)

	// err = s.BuildService.UpdateAvgRatingForBuildService(request.BuildId, types.RatingCategory(request.Category), avgRating)

	// if err != nil {
	// 	fmt.Println("Error updating average rating for build.", err)
	// 	return nil, err
	// }

	return nil, nil
}

/**
* TODO: need to change from products to builds.
* Gets all Ratings.
**/
func (s *service) GetAllRatingsForProduct(userId uuid.UUID) (*[]models.Rating, error) {
	return s.repo.GetAllRatingsByMemberId(userId)
}
