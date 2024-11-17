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

func (s *SkillService) CreateSkillService(createSkillReq CreateSkillRequest) error {
	return s.Repo.CreateSkill(createSkillReq)
}

func (s *SkillService) GetSkillService(id uuid.UUID) (*models.Skill, error) {
	return s.Repo.GetSkill(id)

}
