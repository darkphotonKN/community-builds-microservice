package rating

import (
	"fmt"

	"github.com/darkphotonKN/community-builds-microservice/common/constants/models"
	"github.com/darkphotonKN/community-builds-microservice/common/constants/types"
	commonhelpers "github.com/darkphotonKN/community-builds-microservice/common/utils"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{db: db}
}

/**
* Creates a single rating of a single category type for a build by buildId.
**/
func (r *repository) CreateRatingForBuildById(memberId uuid.UUID, request CreateRatingRequest) error {

	// add a new rating under this member id and build id
	query := `
	INSERT INTO ratings (build_id, member_id, value, category)
	VALUES (:build_id, :member_id, :value, :category)
	`

	params := map[string]interface{}{
		"build_id":  request.BuildId,
		"member_id": memberId,
		"value":     request.Value,
		"category":  request.Category,
	}

	_, err := r.db.NamedExec(query, params)

	if err != nil {
		return commonhelpers.AnalyzeDBErr(err)
	}

	return nil
}

func (r *repository) GetAllRatingsByMemberId(memberId uuid.UUID) (*[]models.Rating, error) {
	var ratings []models.Rating

	query := `
	SELECT * FROM ratings
	WHERE ratings.member_id = $1
	`

	err := r.db.Select(&ratings, query, memberId)

	fmt.Println("ratings:", ratings)

	if err != nil {
		return nil, err
	}

	return &ratings, nil
}

/**
* Retrieves all ratings of a specific rating category for a build.
**/
func (r *repository) GetAllRatingsByCategoryForBuild(buildId string, category types.RatingCategory) ([]int, error) {
	var values []int

	query := `
	SELECT value
	FROM ratings
	WHERE build_id = $1 AND category = $2
	`

	err := r.db.Select(&values, query, buildId, category)

	if err != nil {
		return nil, commonhelpers.AnalyzeDBErr(err)
	}

	return values, nil
}
