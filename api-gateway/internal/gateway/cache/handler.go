package cache

import (
	"net/http"

	pb "github.com/darkphotonKN/community-builds-microservice/common/api/proto/cache"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	client CacheClient
}

func NewHandler(client CacheClient) *Handler {
	return &Handler{
		client: client,
	}
}

func (h *Handler) Hello(c *gin.Context) {
	var request *pb.HelloRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	example, err := h.client.Hello(c.Request.Context(), request)

	if err != nil {
		status, ok := status.FromError(err)

		if !ok {
			// not a gRPC status error
			c.JSON(http.StatusInternalServerError, gin.H{
				"statusCode": http.StatusInternalServerError,
				"message":    "Internal server error",
			})

			return
		}

		// map grpc error codes to http codes
		httpStatus := http.StatusInternalServerError
		switch status.Code() {
		case codes.InvalidArgument:
			httpStatus = http.StatusBadRequest
		case codes.Unauthenticated:
			httpStatus = http.StatusUnauthorized
		case codes.NotFound:
			httpStatus = http.StatusNotFound
		}

		c.JSON(httpStatus, gin.H{
			"statusCode": httpStatus,
			"message":    status.Message(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"statusCode": http.StatusOK, "message": "success", "result": example})
}
