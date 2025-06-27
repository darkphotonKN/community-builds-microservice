package skill

import (
	"context"

	pb "github.com/darkphotonKN/community-builds-microservice/common/api/proto/skill"
)

type Handler struct {
	service Service
	pb.UnimplementedSkillServiceServer
}

type Service interface {
	CreateSkill(ctx context.Context, req *pb.CreateSkillRequest) (*pb.CreateSkillResponse, error)
	// GetBuildsByMemberId(ctx context.Context, req *pb.GetBuildsByMemberIdRequest) (*pb.GetBuildsByMemberIdResponse, error)
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// --- GLOBAL HANDLERS ---
// func (h *SkillHandler) GetSkillsHandler(c *gin.Context) {
// 	skills, err := h.Service.GetSkillsService()

// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": fmt.Sprintf("Error when attempting to get all skills %s", err.Error())})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"statusCode": http.StatusOK, "message": "Successfully retrieved all skills.", "result": skills})
// }

// --- ADMIN HANDLERS ---
func (h *Handler) CreateSkill(ctx context.Context, req *pb.CreateSkillRequest) (*pb.CreateSkillResponse, error) {
	return h.service.CreateSkill(ctx, req)
}
