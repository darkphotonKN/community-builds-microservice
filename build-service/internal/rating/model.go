package rating

type CreateRatingRequest struct {
	BuildId  string `db:"build_id" binding:"required,uuid" json:"buildId"`
	Category string `db:"category" binding:"required,ratingCategory" json:"category"`
	Value    int    `db:"value" binding:"required,min=1,max=10" json:"value"`
}

type RatingByCategoryRes struct {
	Value int `db:"value"`
}
