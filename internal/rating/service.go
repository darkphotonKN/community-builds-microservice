package rating

import (
	"github.com/darkphotonKN/community-builds/internal/models"
	"github.com/google/uuid"
)

type RatingService struct {
	Repo *RatingRepository
}

func NewRatingService(repo *RatingRepository) *RatingService {
	return &RatingService{
		Repo: repo,
	}
}

/**
* Posts single rating for a single product.
**/
func (s *RatingService) CreateRatingForBuildByIdService(memberId uuid.UUID, request CreateRatingRequest) error {
	return s.Repo.CreateRatingForBuildById(memberId, request)
}

/**
* Gets all Ratings.
**/
func (s *RatingService) GetAllRatingsForProductService(userId uuid.UUID) (*[]models.Rating, error) {
	return s.Repo.GetAllRatingsByMemberId(userId)
}
