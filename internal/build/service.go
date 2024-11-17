package build

import (
	"fmt"

	"github.com/darkphotonKN/community-builds/internal/skill"
	"github.com/google/uuid"
)

type BuildService struct {
	Repo         *BuildRepository
	SkillService *skill.SkillService
}

func NewBuildService(repo *BuildRepository, skillService *skill.SkillService) *BuildService {
	return &BuildService{
		Repo:         repo,
		SkillService: skillService,
	}
}

func (s *BuildService) CreateBuildService(memberId uuid.UUID, createBuildRequest CreateBuildRequest) {
	// confirm skill exists
	skill, err := s.SkillService.GetSkillService(createBuildRequest.SkillID)

	if err != nil {
		fmt.Println("Error:", err)
		// return fmt.Errorf("Main skill id could not be found when attempting to create build for it.")
	}

	fmt.Printf("Skill queried from DB: %+v", skill)

	// create build with this skill and for this member

	// return s.Repo.CreateBuild(memberId, createBuildRequest)
}
