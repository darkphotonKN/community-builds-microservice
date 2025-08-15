package composite

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CompositeHandler struct {
	Client CompositeClient
}

func NewHandler(client CompositeClient) *CompositeHandler {
	return &CompositeHandler{
		Client: client,
	}
}

func (h *CompositeHandler) GetGameDataHandler(c *gin.Context) {
	// userIdStr, _ := c.Get("userIdStr")

	// Convert REST request to gRPC request
	// grpcReq := &pb.CreateRareItemRequest{
	// 	MemberId:     userIdStr.(string),
	// }

	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": fmt.Sprintf("Error when attempting to retrieve all items: %s\n", err.Error())})
	// 	return
	// }
	// var items = []string{"item1", "item2", "item3"} // Placeholder for actual data retrieval logic
	gameData, err := h.Client.GetGameData(c.Request.Context())

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": fmt.Sprintf("Error when attempting to retrieve game data: %s\n", err.Error())})
		return
	}
	c.JSON(http.StatusOK, gin.H{"statusCode": http.StatusOK, "message": "Successfully retrieved all game data.", "result": gameData})
}
