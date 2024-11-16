package item

import (
	"fmt"
	"net/http"

	"github.com/darkphotonKN/community-builds/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ItemHandler struct {
	Service *ItemService
}

func NewItemHandler(service *ItemService) *ItemHandler {
	return &ItemHandler{
		Service: service,
	}
}

func (h *ItemHandler) CreateItemHandler(c *gin.Context) {
	var createItemReq models.Item

	if err := c.ShouldBindJSON(&createItemReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": fmt.Sprintf("Error when parsing payload as JSON.")})
	}

	err := h.Service.CreateItemService(createItemReq)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": fmt.Sprintf("Error when attempting to create item: %s", err.Error())})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"statusCode": http.StatusCreated, "message": "Successfully created item."})
}

func (h *ItemHandler) AddItemToBuildHandler(c *gin.Context) {
	buildId, _ := c.Get("buildId")

	var item models.Item

	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": fmt.Sprintf("Error when parsing payload as JSON.")})
		return
	}

	err := h.Service.AddItemToBuildService(buildId.(uuid.UUID), item)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusCode": http.StatusInternalServerError, "message": fmt.Sprintf("Error when attempting to create item: %s", err.Error())})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"statusCode": http.StatusCreated, "message": "Successfully created item."})
}

func (h *ItemHandler) GetItemsHandler(c *gin.Context) {
	buildId, _ := c.Get("buildId")

	items, err := h.Service.GetItemsService(buildId.(uuid.UUID))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": fmt.Sprintf("Error when attempting to retrieve all items from user id: %s, \n error: %s\n", buildId, err.Error())})
		return
	}

	c.JSON(http.StatusOK, gin.H{"statusCode": http.StatusOK, "message": "Successfully retrieved all items.", "result": items})
}

func (h *ItemHandler) UpdateItemsHandler(c *gin.Context) {
	// user id
	buildId, _ := c.Get("buildId")

	// item id to update
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": fmt.Sprintf("Error with id %d, not a valid uuid.", id)})
		return
	}

	// update item payload
	var updateItemReq UpdateItemReq
	if err := c.ShouldBindJSON(&updateItemReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": fmt.Sprintf("Error when parsing payload as JSON.")})
		return
	}

	updatedItem, err := h.Service.UpdateItemsService(buildId.(uuid.UUID), id, updateItemReq)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": fmt.Sprintf("Error when attempting to update item with id: %s for user id: %s\n error: %s\n", id, buildId, err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"statusCode": http.StatusOK, "message": "Successfully retrieved all items.", "result": updatedItem})
}
