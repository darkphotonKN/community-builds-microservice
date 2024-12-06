package rating

type RatingHandler struct {
	Service *RatingService
}

func NewRatingHandler(service *RatingService) *RatingHandler {
	return &RatingHandler{
		Service: service,
	}
}

func (h *RatingHandler) CreateRatingByBuildIdHandler(rating CreateRatingRequest) error {
	return h.Service.CreateRatingByBuildIdService(rating)
}
