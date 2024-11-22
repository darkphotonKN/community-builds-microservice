package build

import (
	"fmt"

	"github.com/darkphotonKN/community-builds/internal/models"
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

const (
	maxBuildCount = 10
)

func (s *BuildService) CreateBuildService(memberId uuid.UUID, createBuildRequest CreateBuildRequest) error {
	// confirm skill exists
	_, err := s.SkillService.GetSkillByIdService(createBuildRequest.SkillID)

	if err != nil {
		return fmt.Errorf("main skill id could not be found when attempting to create build for it.")
	}

	// check how many builds a user has, prevent creation if over the limit
	builds, err := s.GetBuildsForMemberService(memberId)

	if len(*builds) > maxBuildCount {
		return fmt.Errorf("Number of builds allowed have reached maximum capacity.")
	}

	if err != nil {
		return err
	}

	// create build with this skill and for this member
	return s.Repo.CreateBuild(memberId, createBuildRequest)
}

/**
* Get list builds available to a member.
**/
func (s *BuildService) GetBuildsForMemberService(memberId uuid.UUID) (*[]models.Build, error) {
	return s.Repo.GetBuildsByMemberId(memberId)
}

/**
* Get a single build for a member by Id.
**/
func (s *BuildService) GetBuildForMemberByIdService(memberId uuid.UUID, buildId uuid.UUID) (*models.Build, error) {
	return s.Repo.GetBuildForMemberById(memberId, buildId)
}

/**
* Get a single build with all information by ID for a member.
**/
func (s *BuildService) GetBuildInfoForMemberService(memberId uuid.UUID, buildId uuid.UUID) (*BuildInfoResponse, error) {
	return s.Repo.GetBuildInfo(memberId, buildId)
}

/**
* Adds primary and secondary skills and links to an existing build.
**/
func (s *BuildService) AddSkillsToBuildService(memberId uuid.UUID, buildId uuid.UUID, request AddSkillsToBuildRequest) error {
	// get build and check if it exists
	_, err := s.GetBuildForMemberByIdService(memberId, buildId)

	if err != nil {
		return err
	}

	// -- create skills --

	// --- main skill ---

	// add main skill
	err = s.Repo.InsertSkillToBuild(buildId, request.MainSkillLinks.Skill)

	// add main skill's links
	for _, mainSkill := range request.MainSkillLinks.Links {
		err = s.Repo.InsertSkillToBuild(buildId, mainSkill)
	}

	// --- additional skills ---
	for _, secondarySkills := range request.AdditionalSkills {

		// add secondary skill
		err = s.Repo.InsertSkillToBuild(buildId, secondarySkills.Skill)

		// add secondary skill's links
		for _, secondarySkillLinks := range secondarySkills.Links {
			err = s.Repo.InsertSkillToBuild(buildId, secondarySkillLinks)
		}
	}

	return err
}
