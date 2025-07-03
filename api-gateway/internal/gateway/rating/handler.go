package rating

import (
	"fmt"
	"net/http"

	pb "github.com/darkphotonKN/community-builds-microservice/common/api/proto/rating"
	"github.com/gin-gonic/gin"
)

type RatingHandler struct {
	Client RatingClient
}

func NewHandler(client RatingClient) *RatingHandler {
	return &RatingHandler{
		Client: client,
	}
}
func (h *RatingHandler) CreateRatingByBuildIdHandler(c *gin.Context) {
	userIdStr, exists := c.Get("userIdStr")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"statusCode": http.StatusUnauthorized,
			"message":    "User ID not found in context",
		})
		return
	}
	var request CreateRatingRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		fmt.Printf("Failed to bind JSON payload: %+v, Error: %s", request, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": fmt.Sprintf("Error when parsing payload as JSON: %s", err)})
		return
	}
	// Convert REST request to gRPC request
	grpcReq := &pb.CreateRatingByBuildIdRequest{
		MemberId: userIdStr.(string),
	}

	res, err := h.Client.CreateRatingByBuildId(c.Request.Context(), grpcReq)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": fmt.Sprintf("Error when attempting add rating to builds, buildId %s: memberId: %s, error: %s", request.BuildId, userIdStr.(string), err.Error())})
		return
	}

	c.JSON(http.StatusOK, gin.H{"statusCode": http.StatusOK, "message": fmt.Sprintf("Successfully added rating to build %s", userIdStr.(string)), "result": res})
}
