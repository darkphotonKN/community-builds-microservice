package rating

import (
	"fmt"

	"github.com/darkphotonKN/community-builds/internal/build"
	"github.com/darkphotonKN/community-builds/internal/models"
	"github.com/darkphotonKN/community-builds/internal/types"
	"github.com/google/uuid"
)

type RatingService struct {
	Repo         *RatingRepository
	BuildService *build.BuildService
}

func NewRatingService(repo *RatingRepository, buildService *build.BuildService) *RatingService {
	return &RatingService{
		Repo:         repo,
		BuildService: buildService,
	}
}

/**
* Posts single rating for a single product.
**/
func (s *RatingService) CreateRatingForBuildByIdService(memberId uuid.UUID, request CreateRatingRequest) error {
	// create rating for build
	err := s.Repo.CreateRatingForBuildById(memberId, request)

	if err != nil {
		return err
	}

	// update build's average rating.

	// get all ratings of that category
	ratings, err := s.Repo.GetAllRatingsByCategoryForBuild(request.BuildId, types.RatingCategory(request.Category))

	// average them
	var avgRating float32
	totalRating := 0
	noOfRatings := float32(len(ratings))

	for _, rating := range ratings {
		totalRating += rating
	}

	avgRating = float32(totalRating) / noOfRatings

	fmt.Printf("ratings %+v, avgRating: %f\n", ratings, avgRating)

	err = s.BuildService.UpdateAvgRatingForBuildService(request.BuildId, types.RatingCategory(request.Category), avgRating)

	if err != nil {
		fmt.Println("Error updating average rating for build.", err)
		return err
	}

	return nil
}

/**
* TODO: need to change from products to builds.
* Gets all Ratings.
**/
func (s *RatingService) GetAllRatingsForProductService(userId uuid.UUID) (*[]models.Rating, error) {
	return s.Repo.GetAllRatingsByMemberId(userId)
}
