package rating

type RatingHandler struct {
	Service *RatingService
}

func NewRatingHandler(service *RatingService) *RatingHandler {
	return &RatingHandler{
		Service: service,
	}
}
