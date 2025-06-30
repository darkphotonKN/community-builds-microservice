package skill

import (
	"fmt"
	"net/http"

	pb "github.com/darkphotonKN/community-builds-microservice/common/api/proto/skill"
	"github.com/gin-gonic/gin"
)

type SkillHandler struct {
	Client SkillClient
}

func NewHandler(client SkillClient) *SkillHandler {
	return &SkillHandler{
		Client: client,
	}
}

// --- GLOBAL HANDLERS ---
func (h *SkillHandler) GetSkillsHandler(c *gin.Context) {
	skills, err := h.Client.GetSkills(c.Request.Context(), &pb.GetSkillsRequest{})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": fmt.Sprintf("Error when attempting to get all skills %s", err.Error())})
		return
	}

	c.JSON(http.StatusOK, gin.H{"statusCode": http.StatusOK, "message": "Successfully retrieved all skills.", "result": skills})
}

// --- ADMIN HANDLERS ---
func (h *SkillHandler) CreateSkillHandler(c *gin.Context) {
	userIdStr, exists := c.Get("userIdStr")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"statusCode": http.StatusUnauthorized,
			"message":    "User ID not found in context",
		})
		return
	}
	var createSkillReq CreateSkillRequest

	if err := c.ShouldBindJSON(&createSkillReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": fmt.Sprintf("Error when parsing payload as JSON: %s", err)})
		return
	}
	// Convert REST request to gRPC request
	grpcReq := &pb.CreateSkillRequest{
		MemberId: userIdStr.(string),
		Name:     createSkillReq.Name,
		Type:     createSkillReq.Type,
	}

	skill, err := h.Client.CreateSkill(c.Request.Context(), grpcReq)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": fmt.Sprintf("Error when attempting to create skill %s", err.Error())})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"statusCode": http.StatusCreated, "message": "Successfully created skill.", "result": skill})
}
