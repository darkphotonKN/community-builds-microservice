package skill

import (
	"github.com/darkphotonKN/community-builds/internal/models"
	"github.com/google/uuid"
)

type SkillService struct {
	Repo *SkillRepository
}

func NewSkillService(repo *SkillRepository) *SkillService {
	return &SkillService{
		Repo: repo,
	}
}

/**
* Create a single skill.
**/
func (s *SkillService) CreateSkillService(createSkillReq CreateSkillRequest) error {
	return s.Repo.CreateSkill(createSkillReq)
}

/**
* Get a single skill by skill's id.
**/
func (s *SkillService) GetSkillByIdService(id uuid.UUID) (*models.Skill, error) {
	return s.Repo.GetSkill(id)
}

/**
* Get list of skills available.
**/
func (s *SkillService) GetSkillsService() (*[]models.Skill, error) {
	return s.Repo.GetSkills()
}

/**
* Creates a list of skills.
**/
func (s *SkillService) BatchCreateSkillsService(createSkills []SeedSkill) error {

	if err := s.Repo.BatchCreateSkills(createSkills); err != nil {
		return err
	}

	return nil
}

/**
* Gets a list of skills belonging to a build by id.
**/
func (s *SkillService) GetSkillsByBuildIdService(buildId uuid.UUID) (*[]models.SkillRow, error) {
	return s.Repo.GetSkillsAndLinksByBuildId(buildId)
}
