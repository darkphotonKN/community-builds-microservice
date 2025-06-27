package skill

import (
	"context"

	pb "github.com/darkphotonKN/community-builds-microservice/common/api/proto/skill"
	"github.com/google/uuid"
)

type CreateSkillRequest struct {
	Name string `json:"name" binding:"required" db:"name"`
	Type string `json:"type" binding:"required,skillType" db:"type"`
}

type SeedSkill struct {
	ID   uuid.UUID `db:"id"`
	Name string    `db:"name"`
	Type string    `db:"type"`
}

type SkillClient interface {
	CreateSkill(ctx context.Context, req *pb.CreateSkillRequest) (*pb.CreateSkillResponse, error)
	// GetBuildsByMemberId(ctx context.Context, req *pb.GetBuildsByMemberIdRequest) (*pb.GetBuildsByMemberIdResponse, error)
}
