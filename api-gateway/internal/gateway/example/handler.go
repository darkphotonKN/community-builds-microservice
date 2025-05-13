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
	var request pb.CreateExampleRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert REST request to gRPC request
	grpcReq := &pb.CreateExampleRequest{
		Name: request.Name,
	}

	// Call the service
	example, err := h.client.CreateExample(c.Request.Context(), grpcReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert gRPC response to REST response
	response := pb.Example{
		Id:        example.Id,
		Name:      example.Name,
		CreatedAt: example.CreatedAt,
		UpdatedAt: example.UpdatedAt,
	}

	c.JSON(http.StatusCreated, response)
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

	// Convert gRPC response to REST response
	response := pb.Example{
		Id:        example.Id,
		Name:      example.Name,
		CreatedAt: example.CreatedAt,
		UpdatedAt: example.UpdatedAt,
	}

	c.JSON(http.StatusOK, response)
}
