package tag

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	// "github.com/google/uuid"
)

type TagHandler struct {
	Service *TagService
}

func NewTagHandler(service *TagService) *TagHandler {
	return &TagHandler{
		Service: service,
	}
}

// --- ADMIN HANDLERS ---
func (h *TagHandler) CreateTagHandler(c *gin.Context) {
	var createTagReq CreateTagRequest

	if err := c.ShouldBindJSON(&createTagReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": fmt.Sprintf("Error when parsing payload as JSON: %s", err)})
		return
	}

	err := h.Service.CreateTagService(createTagReq)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": fmt.Sprintf("Error when attempting to create tag: %s", err.Error())})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"statusCode": http.StatusCreated, "message": "Successfully created tag."})
}

func (h *TagHandler) GetTagsHandler(c *gin.Context) {
	tags, err := h.Service.GetTagsService()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": fmt.Sprintf("Error when attempting to retrieve all tags: %s\n", err.Error())})
		return
	}

	c.JSON(http.StatusOK, gin.H{"statusCode": http.StatusOK, "message": "Successfully retrieved all tags.", "result": tags})
}

func (h *TagHandler) UpdateTagsHandler(c *gin.Context) {
	// tag id to update
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": fmt.Sprintf("Error with id %d, not a valid uuid.", id)})
		return
	}

	// update tag payload
	var updateTagReq UpdateTagRequest
	if err := c.ShouldBindJSON(&updateTagReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": fmt.Sprintf("Error when parsing payload as JSON.")})
		return
	}

	resErr := h.Service.UpdateTagsService(updateTagReq)

	if resErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": fmt.Sprintf("Error when attempting to update tag with id: %s\n error: %s\n", id, err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"statusCode": http.StatusOK, "message": "Successfully updated tag."})
}
