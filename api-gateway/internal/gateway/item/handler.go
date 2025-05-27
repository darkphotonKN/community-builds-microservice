package item

import (
	"fmt"
	"net/http"

	pb "github.com/darkphotonKN/community-builds-microservice/common/api/proto/item"
	"github.com/gin-gonic/gin"
)

// "fmt"
// "net/http"

// "github.com/gin-gonic/gin"
// "github.com/google/uuid"

type ItemHandler struct {
	Client ItemClient
}

func NewHandler(client ItemClient) *ItemHandler {
	return &ItemHandler{
		Client: client,
	}
}

func (h *ItemHandler) CreateItemHandler(c *gin.Context) {
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
	c.JSON(http.StatusCreated, gin.H{"statusCode": http.StatusCreated, "message": "Successfully created item.", "item": item})

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
	c.JSON(http.StatusCreated, gin.H{"statusCode": http.StatusCreated, "message": "Successfully created item.", "item": item})

}

// // --- ADMIN HANDLERS ---
// func (h *ItemHandler) CreateItemHandler(c *gin.Context) {
// 	var createItemReq CreateItemRequest

// 	if err := c.ShouldBindJSON(&createItemReq); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": fmt.Sprintf("Error when parsing payload as JSON: %s", err)})
// 		return
// 	}

// 	err := h.Service.CreateItemService(createItemReq)

// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": fmt.Sprintf("Error when attempting to create item: %s", err.Error())})
// 		return
// 	}

// 	c.JSON(http.StatusCreated, gin.H{"statusCode": http.StatusCreated, "message": "Successfully created item."})
// }

// func (h *ItemHandler) GetItemsHandler(c *gin.Context) {

// 	slot := c.Query("slot")
// 	fmt.Println("slot", slot)
// 	items, err := h.Service.GetItemsService(slot)

// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": fmt.Sprintf("Error when attempting to retrieve all items: %s\n", err.Error())})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"statusCode": http.StatusOK, "message": "Successfully retrieved all items.", "result": items})
// }

// func (h *ItemHandler) GetBaseItemsHandler(c *gin.Context) {

// 	items, err := h.Service.GetBaseItemsService()

// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": fmt.Sprintf("Error when attempting to retrieve all items: %s\n", err.Error())})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"statusCode": http.StatusOK, "message": "Successfully retrieved all items.", "result": items})
// }

// func (h *ItemHandler) GetItemModsHandler(c *gin.Context) {

// 	items, err := h.Service.GetItemModsService()

// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": fmt.Sprintf("Error when attempting to retrieve all items: %s\n", err.Error())})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"statusCode": http.StatusOK, "message": "Successfully retrieved all items.", "result": items})
// }

// func (h *ItemHandler) UpdateItemsHandler(c *gin.Context) {
// 	// item id to update
// 	idParam := c.Param("id")
// 	id, err := uuid.Parse(idParam)

// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": fmt.Sprintf("Error with id %d, not a valid uuid.", id)})
// 		return
// 	}

// 	// update item payload
// 	var updateItemReq UpdateItemReq
// 	if err := c.ShouldBindJSON(&updateItemReq); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": fmt.Sprintf("Error when parsing payload as JSON.")})
// 		return
// 	}

// 	updatedItem, err := h.Service.UpdateItemsService(id, updateItemReq)

// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": fmt.Sprintf("Error when attempting to update item with id: %s\n error: %s\n", id, err)})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"statusCode": http.StatusOK, "message": "Successfully updated item.", "result": updatedItem})
// }

// /**
// * BUILDS
// **/
// func (h *ItemHandler) AddItemToBuildHandler(c *gin.Context) {
// 	buildId, _ := c.Get("buildId")

// 	var item CreateItemRequest

// 	if err := c.ShouldBindJSON(&item); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": fmt.Sprintf("Error when parsing payload as JSON.")})
// 		return
// 	}

// 	err := h.Service.AddItemToBuildService(buildId.(uuid.UUID), item)

// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"statusCode": http.StatusInternalServerError, "message": fmt.Sprintf("Error when attempting to create item: %s", err.Error())})
// 		return
// 	}

// 	c.JSON(http.StatusCreated, gin.H{"statusCode": http.StatusCreated, "message": "Successfully created item."})
// }

// func (h *ItemHandler) GetUniqueItemsHandler(c *gin.Context) {
// 	err := h.Service.CrawlingAndAddUniqueItemsService()

// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": fmt.Sprintf("Error when attempting to retrieve all items: %s\n", err.Error())})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"statusCode": http.StatusOK, "message": "Successfully retrieved all items."})
// }

// func (h *ItemHandler) CreateRareItemHandler(c *gin.Context) {
// 	userId, _ := c.Get("userId")
// 	fmt.Println("userId", userId)
// 	// update item payload
// 	var createRareItemReq CreateRareItemReq
// 	if err := c.ShouldBindJSON(&createRareItemReq); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": fmt.Sprintf("Error when parsing payload as JSON.")})
// 		return
// 	}

// 	id, resErr := h.Service.CreateRareItemService(userId.(uuid.UUID), createRareItemReq)

// 	if resErr != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": fmt.Sprintf("Error when attempting to retrieve all items: %s\n", resErr.Error())})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"statusCode": http.StatusOK, "message": "Successfully retrieved all items.", "result": id})
// }

// func (h *ItemHandler) GetMemberRareItemHandler(c *gin.Context) {
// 	userId, _ := c.Get("userId")
// 	fmt.Println("userId", userId)

// 	items, err := h.Service.GetMemberRareItemsService(userId.(uuid.UUID))

// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": fmt.Sprintf("Error when attempting to retrieve all items: %s\n", err.Error())})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"statusCode": http.StatusOK, "message": "Successfully retrieved all items.", "result": items})
// }

// func (h *ItemHandler) GetAllDataHandler(c *gin.Context) {
// 	userId, _ := c.Get("userId")
// 	fmt.Println("userId", userId)

// 	items, err := h.Service.GetAllDataService(userId.(uuid.UUID))

// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": fmt.Sprintf("Error when attempting to retrieve all items: %s\n", err.Error())})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"statusCode": http.StatusOK, "message": "Successfully retrieved all items.", "result": items})
// }
