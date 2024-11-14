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
func (s *RatingService) PostRatingService(productId uuid.UUID, ratingReq RatingRequest) error {
	return s.Repo.CreateRating(productId, ratingReq)
}

/**
* Gets all Ratings.
**/
func (s *RatingService) GetAllRatingsForProductService(userId uuid.UUID) (*[]models.Rating, error) {
	return s.Repo.GetAllRatingsByMemberId(userId)
}
