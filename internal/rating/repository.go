package rating

import (
	"fmt"

	"github.com/darkphotonKN/community-builds/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type RatingRepository struct {
	DB *sqlx.DB
}

func NewRatingRepository(db *sqlx.DB) *RatingRepository {
	return &RatingRepository{
		DB: db,
	}
}

func (r *RatingRepository) CreateRating(memberId uuid.UUID, ratingReq RatingRequest) error {

	// add a new rating under this product's id
	query := `
	INSERT INTO ratings (member_id, rating)
	VALUES (:member_id, :rating)
	`

	// temporary struct to hold values
	params := map[string]interface{}{
		"member_id": memberId,
		"rating":    ratingReq.Rating,
	}

	_, err := r.DB.NamedQuery(query, params)

	if err != nil {
		return err
	}

	return nil
}

func (r *RatingRepository) GetAllRatingsByMemberId(memberId uuid.UUID) (*[]models.Rating, error) {
	var ratings []models.Rating

	query := `
	SELECT * FROM ratings
	WHERE ratings.member_id = $1
	`

	err := r.DB.Select(&ratings, query, memberId)

	fmt.Println("ratings:", ratings)

	if err != nil {
		return nil, err
	}

	return &ratings, nil
}
