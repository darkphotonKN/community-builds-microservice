package skill

import (
	"context"

	pb "github.com/darkphotonKN/community-builds-microservice/common/api/proto/skill"
	"github.com/darkphotonKN/community-builds-microservice/common/constants/models"
	"github.com/google/uuid"
)

type service struct {
	repo Repository
}

type Repository interface {
	CreateSkill(createBuildRequest CreateSkillRequest) error
	GetSkills() (*[]models.Skill, error)
	GetSkill(id uuid.UUID) (*models.Skill, error)
	GetSkillsAndLinksByBuildId(buildId uuid.UUID) (*[]models.SkillRow, error)
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

/**
* Create a single skill.
**/
func (s *service) CreateSkill(ctx context.Context, req *pb.CreateSkillRequest) (*pb.CreateSkillResponse, error) {

	createSkillReq := CreateSkillRequest{
		Name: req.Name,
		Type: req.Type,
	}
	err := s.repo.CreateSkill(createSkillReq)
	if err != nil {
		return nil, err
	}

	return &pb.CreateSkillResponse{}, nil
}

/**
* Get a single skill by skill's id. not grpc
**/
func (s *service) GetSkillById(id uuid.UUID) (*models.Skill, error) {
	return s.repo.GetSkill(id)
}

/**
* Get list of skills available.
**/
func (s *service) GetSkills(ctx context.Context, req *pb.GetSkillsRequest) (*pb.GetSkillsResponse, error) {

	skills, err := s.repo.GetSkills()
	if err != nil {
		return nil, err
	}
	pbSkills := make([]*pb.Skill, 0, len(*skills))
	for _, skill := range *skills {
		pbSkills = append(pbSkills, &pb.Skill{
			Id:        skill.Id.String(),
			Name:      skill.Name,
			Type:      skill.Type,
			CreatedAt: skill.CreatedAt.String(),
			UpdatedAt: skill.UpdatedAt.String(),
		})
	}
	return &pb.GetSkillsResponse{Skills: pbSkills}, nil
}

/**
* Creates a list of skills.
**/
// func (s *SkillService) BatchCreateSkillsService(createSkills []SeedSkill) error {

// 	if err := s.Repo.BatchCreateSkills(createSkills); err != nil {
// 		return err
// 	}

// 	return nil
// }

/**
* Gets a list of skills belonging to a build by id.
**/
func (s *service) GetSkillsByBuildId(buildId uuid.UUID) (*[]models.SkillRow, error) {
	return s.repo.GetSkillsAndLinksByBuildId(buildId)
}
