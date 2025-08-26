package types

type RatingCategory string

const (
	Endgame   RatingCategory = "endgame"
	Bossing   RatingCategory = "bossing"
	Speedfarm RatingCategory = "speedfarm"
	Fun       RatingCategory = "fun"
	Creative  RatingCategory = "creative"
	Rating    RatingCategory = "rating"
)

type Status int

const (
	IsDraft Status = iota
	IsPublished
	IsArchived
)
