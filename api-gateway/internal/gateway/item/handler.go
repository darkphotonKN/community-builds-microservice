package item

import (
	"fmt"
	"net/http"

	pb "github.com/darkphotonKN/community-builds-microservice/common/api/proto/item"
	"github.com/gin-gonic/gin"
)

type ItemHandler struct {
	Client ItemClient
}

func NewHandler(client ItemClient) *ItemHandler {
	return &ItemHandler{
		Client: client,
	}
}

func (h *ItemHandler) CreateItemHandler(c *gin.Context) {
	// userId, _ := c.Get("userId")
	// userIdStr := userId.(uuid.UUID).String()
	userIdStr, _ := c.Get("userIdStr")
	var createItemReq CreateItemRequest
	if err := c.ShouldBindJSON(&createItemReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": fmt.Sprintf("Error when parsing payload as JSON: %s", err)})
		return
	}

	// Convert REST request to gRPC request
	grpcReq := &pb.CreateItemRequest{
		Id:       userIdStr.(string),
		Name:     createItemReq.Name,
		Category: createItemReq.Category,
		Class:    createItemReq.Class,
		Type:     createItemReq.Type,
		ImageURL: createItemReq.ImageURL,
		Slot:     createItemReq.Slot,
	}
	item, err := h.Client.CreateItem(c.Request.Context(), grpcReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": fmt.Sprintf("Error when attempting to create item: %s", err.Error())})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"statusCode": http.StatusCreated, "message": "Successfully created item.", "item": item})

}

func (h *ItemHandler) GetItemsHandler(c *gin.Context) {
	slot := c.Query("slot")

	// Convert REST request to gRPC request
	grpcReq := &pb.GetItemsRequest{
		Slot: slot,
	}
	item, err := h.Client.GetItems(c.Request.Context(), grpcReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": fmt.Sprintf("Error when attempting to create item: %s", err.Error())})
		return
	}
	c.JSON(http.StatusOK, gin.H{"statusCode": http.StatusOK, "message": "Successfully created item.", "result": item})

}

func (h *ItemHandler) UpdateItemHandler(c *gin.Context) {
	userIdStr, _ := c.Get("userIdStr")
	var createItemReq UpdateItemReq
	if err := c.ShouldBindJSON(&createItemReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": fmt.Sprintf("Error when parsing payload as JSON: %s", err)})
		return
	}
	// Convert REST request to gRPC request
	grpcReq := &pb.UpdateItemRequest{
		Id:       userIdStr.(string),
		Name:     createItemReq.Name,
		Category: createItemReq.Category,
		Class:    createItemReq.Class,
		Type:     createItemReq.Type,
		ImageURL: createItemReq.ImageURL,
	}
	item, err := h.Client.UpdateItem(c.Request.Context(), grpcReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": fmt.Sprintf("Error when attempting to create item: %s", err.Error())})
		return
	}
	c.JSON(http.StatusOK, gin.H{"statusCode": http.StatusOK, "message": "Successfully created item.", "result": item})

}

func (h *ItemHandler) CreateRareItemHandler(c *gin.Context) {
	userIdStr, _ := c.Get("userIdStr")
	// update item payload
	var createRareItemRequest pb.CreateRareItemRequest
	if err := c.ShouldBindJSON(&createRareItemRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": fmt.Sprintf("Error when parsing payload as JSON.")})
		return
	}

	// Convert REST request to gRPC request
	grpcReq := &pb.CreateRareItemRequest{
		UserId:     userIdStr.(string),
		Name:       createRareItemRequest.Name,
		Stats:      createRareItemRequest.Stats,
		BaseItemId: createRareItemRequest.BaseItemId,
		ToList:     createRareItemRequest.ToList,
	}

	res, err := h.Client.CreateRareItem(c.Request.Context(), grpcReq)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": fmt.Sprintf("Error when create rare item: %s\n", err.Error())})
		return
	}

	c.JSON(http.StatusOK, gin.H{"statusCode": http.StatusOK, "message": "Successfully create rare item.", "result": res})
}

func (h *ItemHandler) GetBaseItemsHandler(c *gin.Context) {

	items, err := h.Client.GetBaseItems(c.Request.Context(), &pb.GetBaseItemsRequest{})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": fmt.Sprintf("Error when attempting to retrieve all items: %s\n", err.Error())})
		return
	}

	c.JSON(http.StatusOK, gin.H{"statusCode": http.StatusOK, "message": "Successfully retrieved all base items.", "result": items})
}

func (h *ItemHandler) GetItemModsHandler(c *gin.Context) {

	items, err := h.Client.GetItemMods(c.Request.Context(), &pb.GetItemModsRequest{})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": fmt.Sprintf("Error when attempting to retrieve all items: %s\n", err.Error())})
		return
	}

	c.JSON(http.StatusOK, gin.H{"statusCode": http.StatusOK, "message": "Successfully retrieved all item mods.", "result": items})
}

func (h *ItemHandler) GetMemberRareItemsHandler(c *gin.Context) {
	userIdStr, ok := c.Get("userIdStr")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"statusCode": http.StatusUnauthorized, "message": "Unauthorized"})
		return
	}
	payload := &pb.GetMemberRareItemsRequest{
		MemberId: userIdStr.(string),
	}

	items, err := h.Client.GetMemberRareItems(c.Request.Context(), payload)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": fmt.Sprintf("Error when attempting to retrieve all items: %s\n", err.Error())})
		return
	}

	c.JSON(http.StatusOK, gin.H{"statusCode": http.StatusOK, "message": "Successfully retrieved all item mods.", "result": items})
}

// func (h *ItemHandler) GetAllDataHandler(c *gin.Context) {
// 	userIdStr, _ := c.Get("userIdStr")
// 	// update item payload
// 	var createRareItemRequest pb.CreateRareItemRequest
// 	if err := c.ShouldBindJSON(&createRareItemRequest); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": fmt.Sprintf("Error when parsing payload as JSON.")})
// 		return
// 	}

// 	// Convert REST request to gRPC request
// 	grpcReq := &pb.CreateRareItemRequest{
// 		MemberId:     userIdStr.(string),
// 	}

// 	items, err := h.Client.GetAllDataService(c.Request.Context(), grpcReq)

// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": fmt.Sprintf("Error when attempting to retrieve all items: %s\n", err.Error())})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"statusCode": http.StatusOK, "message": "Successfully retrieved all items.", "result": items})
// }
