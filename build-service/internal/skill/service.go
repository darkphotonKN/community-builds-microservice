package skill

import (
	"context"

	pb "github.com/darkphotonKN/community-builds-microservice/common/api/proto/skill"
	amqp "github.com/rabbitmq/amqp091-go"
)

type service struct {
	repo      Repository
	publishCh *amqp.Channel
}

type Repository interface {
	CreateSkill(createBuildRequest CreateSkillRequest) error
}

func NewService(repo Repository, publishCh *amqp.Channel) Service {
	return &service{repo: repo, publishCh: publishCh}
}

/**
* Create a single skill.
**/
func (s *service) CreateSkill(ctx context.Context, req *pb.CreateSkillRequest) (*pb.CreateSkillResponse, error) {

	// return s.repo.CreateSkill(req)

	return &pb.CreateSkillResponse{}, nil
}

/**
* Get a single skill by skill's id.
**/
// func (s *SkillService) GetSkillByIdService(id uuid.UUID) (*models.Skill, error) {
// 	return s.Repo.GetSkill(id)
// }

/**
* Get list of skills available.
**/
// func (s *SkillService) GetSkillsService() (*[]models.Skill, error) {
// 	return s.Repo.GetSkills()
// }

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
// func (s *SkillService) GetSkillsByBuildIdService(buildId uuid.UUID) (*[]models.SkillRow, error) {
// 	return s.Repo.GetSkillsAndLinksByBuildId(buildId)
// }
