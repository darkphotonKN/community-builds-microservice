package rating

import "github.com/google/uuid"

type CreateRating struct {
	BuildId      uuid.UUID `db:"build_id" binding:"required,uuid"`
	FromMemberId uuid.UUID `db:"from_member_id" binding:"required,uuid"`
	ToMemberId   uuid.UUID `db:"to_member_id" binding:"required,uuid"`
	// Category string `db:"category" binding:"required,ratingCategory" json:"category"`
	Rating int `db:"rating" binding:"required,min=1,max=10"`
}

type RatingByCategoryRes struct {
	Value int `db:"value"`
}

type BuildRating struct {
	Id     uuid.UUID `db:"id"`
	Rating int       `db:"rating"`
}
