package class

import (
	"fmt"
	"net/http"

	pb "github.com/darkphotonKN/community-builds-microservice/common/api/proto/class"
	"github.com/gin-gonic/gin"
)

type ClassHandler struct {
	Client ClassClient
}

func NewHandler(client ClassClient) *ClassHandler {
	return &ClassHandler{
		Client: client,
	}
}

/**
* Retrievs all classes and ascendancies.
**/
func (h *ClassHandler) GetClassesAndAscendanciesHandler(c *gin.Context) {

	userIdStr, _ := c.Get("userIdStr")

	// Convert REST request to gRPC request
	grpcReq := &pb.GetClassesAndAscendanciesRequest{
		MemberId: userIdStr.(string),
	}
	response, err := h.Client.GetClassesAndAscendancies(c.Request.Context(), grpcReq)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusCode": http.StatusInternalServerError, "message": fmt.Sprintf("Error occcured when attemptign to retrieve build: %s", err)})
	}

	c.JSON(http.StatusOK, gin.H{"message": "Retrieved all classes and ascendancies.", "result": response})
}
