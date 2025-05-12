package rating

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RatingHandler struct {
	Service *RatingService
}

func NewRatingHandler(service *RatingService) *RatingHandler {
	return &RatingHandler{
		Service: service,
	}
}

func (h *RatingHandler) CreateRatingByBuildIdHandler(c *gin.Context) {
	memberId, _ := c.Get("userId")

	var request CreateRatingRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		fmt.Printf("Failed to bind JSON payload: %+v, Error: %s", request, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": fmt.Sprintf("Error when parsing payload as JSON: %s", err)})
		return
	}

	err := h.Service.CreateRatingForBuildByIdService(memberId.(uuid.UUID), request)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": fmt.Sprintf("Error when attempting add rating to builds, buildId %s: memberId: %s, error: %s", request.BuildId, memberId, err.Error())})
		return
	}

	c.JSON(http.StatusOK, gin.H{"statusCode": http.StatusOK, "message": fmt.Sprintf("Successfully added rating to build %s", memberId)})
}
