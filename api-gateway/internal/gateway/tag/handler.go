package tag

import (
	"fmt"
	"net/http"

	pb "github.com/darkphotonKN/community-builds-microservice/common/api/proto/tag"
	"github.com/gin-gonic/gin"
	// "github.com/google/uuid"
)

type TagHandler struct {
	Client TagClient
}

func NewHandler(client TagClient) *TagHandler {
	return &TagHandler{
		Client: client,
	}
}

// --- ADMIN HANDLERS ---
func (h *TagHandler) CreateTagHandler(c *gin.Context) {
	userIdStr, exists := c.Get("userIdStr")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"statusCode": http.StatusUnauthorized,
			"message":    "User ID not found in context",
		})
		return
	}
	var createTagReq CreateTagRequest

	if err := c.ShouldBindJSON(&createTagReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": fmt.Sprintf("Error when parsing payload as JSON: %s", err)})
		return
	}

	grpcReq := &pb.CreateTagRequest{
		MemberId: userIdStr.(string),
		Name:     createTagReq.Name,
	}
	tag, err := h.Client.CreateTag(c.Request.Context(), grpcReq)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": fmt.Sprintf("Error when attempting to create tag: %s", err.Error())})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"statusCode": http.StatusCreated, "message": "Successfully created tag.", "result": tag})
}

func (h *TagHandler) GetTagsHandler(c *gin.Context) {
	tags, err := h.Client.GetTags(c.Request.Context(), &pb.GetTagsRequest{})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": fmt.Sprintf("Error when attempting to retrieve all tags: %s\n", err.Error())})
		return
	}

	c.JSON(http.StatusOK, gin.H{"statusCode": http.StatusOK, "message": "Successfully retrieved all tags.", "result": tags})
}

func (h *TagHandler) UpdateTagsHandler(c *gin.Context) {
	userIdStr, exists := c.Get("userIdStr")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"statusCode": http.StatusUnauthorized,
			"message":    "User ID not found in context",
		})
		return
	}

	// update tag payload
	var updateTagReq UpdateTagRequest
	if err := c.ShouldBindJSON(&updateTagReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": fmt.Sprintf("Error when parsing payload as JSON.")})
		return
	}

	grpcReq := &pb.UpdateTagRequest{
		Id:   userIdStr.(string),
		Name: updateTagReq.Name,
	}
	tag, err := h.Client.UpdateTag(c.Request.Context(), grpcReq)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": fmt.Sprintf("Error when attempting to update tag with id: %s\n error: %s\n", userIdStr.(string), err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"statusCode": http.StatusOK, "message": "Successfully updated tag.", "result": tag})
}
