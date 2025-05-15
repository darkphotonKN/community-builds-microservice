package example

import (
	"net/http"

	pb "github.com/darkphotonKN/community-builds-microservice/common/api/proto/example"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	client ExampleClient
}

func NewHandler(client ExampleClient) *Handler {
	return &Handler{
		client: client,
	}
}

func (h *Handler) CreateExample(c *gin.Context) {
	var request *pb.CreateExampleRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call the service
	example, err := h.client.CreateExample(c.Request.Context(), request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"statusCode": http.StatusOK, "message": "success", "result": example})
}

func (h *Handler) GetExample(c *gin.Context) {
	id := c.Param("id")

	// Convert REST request to gRPC request
	grpcReq := &pb.GetExampleRequest{
		Id: id,
	}

	// Call the service
	example, err := h.client.GetExample(c.Request.Context(), grpcReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"statusCode": http.StatusOK, "message": "success", "result": example})
}
