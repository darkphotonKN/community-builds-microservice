package rating

type CreateRatingRequest struct {
	BuildId  string `db:"build_id" binding:"required,uuid" json:"buildId"`
	Category string `db:"category" binding:"required,ratingCategory" json:"category"`
	Value    int    `db:"value" binding:"required" json:"value"`
}
