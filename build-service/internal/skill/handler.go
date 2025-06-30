package skill

import (
	"context"

	pb "github.com/darkphotonKN/community-builds-microservice/common/api/proto/skill"
	"github.com/darkphotonKN/community-builds-microservice/common/constants/models"
	"github.com/google/uuid"
)

type Handler struct {
	service Service
	pb.UnimplementedSkillServiceServer
}

type Service interface {
	CreateSkill(ctx context.Context, req *pb.CreateSkillRequest) (*pb.CreateSkillResponse, error)
	GetSkills(ctx context.Context, req *pb.GetSkillsRequest) (*pb.GetSkillsResponse, error)
	GetSkillById(id uuid.UUID) (*models.Skill, error)
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// --- GLOBAL HANDLERS ---
func (h *Handler) GetSkills(ctx context.Context, req *pb.GetSkillsRequest) (*pb.GetSkillsResponse, error) {
	return h.service.GetSkills(ctx, req)
}

// --- ADMIN HANDLERS ---
func (h *Handler) CreateSkill(ctx context.Context, req *pb.CreateSkillRequest) (*pb.CreateSkillResponse, error) {
	return h.service.CreateSkill(ctx, req)
}
