package rating

type RatingRequest struct {
	Value    int    `db:"value" binding:"required" json:"value"`
	Category string `db:"category" binding:"required,category" json:"category"`
}
