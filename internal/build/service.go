package build

import (
	"fmt"

	"github.com/darkphotonKN/community-builds/internal/models"
	"github.com/darkphotonKN/community-builds/internal/skill"
	"github.com/darkphotonKN/community-builds/internal/types"
	"github.com/darkphotonKN/community-builds/internal/utils/dbutils"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
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

func (s *BuildService) GetCommunityBuildsService(pageNo int, pageSize int, sortOrder string, sortBy string, search string, skillId uuid.UUID, minRating *int, ratingCategory types.RatingCategory) ([]models.Build, error) {
	return s.Repo.GetAllBuilds(pageNo, pageSize, sortOrder, sortBy, search, skillId, minRating, ratingCategory)
}

/**
* Create build for a signed-in member.
**/
func (s *BuildService) CreateBuildService(memberId uuid.UUID, createBuildRequest CreateBuildRequest) error {
	// confirm skill exists
	_, err := s.SkillService.GetSkillByIdService(createBuildRequest.SkillID)

	if err != nil {
		return fmt.Errorf("main skill id could not be found when attempting to create build for it.")
	}

	// check how many builds a user has, prevent creation if over the limit
	builds, err := s.GetBuildsForMemberService(memberId)

	if err != nil {
		return err
	}

	if len(*builds) > maxBuildCount {
		return fmt.Errorf("Number of builds allowed have reached maximum capacity.")
	}

	// create build with this skill and for this member
	return s.Repo.CreateBuild(memberId, createBuildRequest)
}

/**
* Get list builds available to a member.
**/
func (s *BuildService) GetBuildsForMemberService(memberId uuid.UUID) (*[]BuildListResponse, error) {
	baseBuilds, err := s.Repo.GetBuildsByMemberId(memberId)

	if err != nil {
		return nil, err
	}

	buildListResponse := make([]BuildListResponse, len(*baseBuilds))

	// query and add each builds tag's
	for _, build := range *baseBuilds {
		tags, err := s.Repo.GetBuildTagsForMemberById(memberId, build.ID)

		// stop query pre-maturely if errored on query
		if err != nil {
			return nil, err
		}

		buildListResponse = append(buildListResponse, BuildListResponse{
			ID:                 build.ID,
			Title:              build.Title,
			Description:        build.Description,
			Ascendancy:         build.Ascendancy,
			MainSkillName:      build.MainSkillName,
			AvgEndGameRating:   build.AvgEndGameRating,
			AvgBossingRating:   build.AvgBossingRating,
			AvgCreativeRating:  build.AvgCreativeRating,
			AvgFunRating:       build.AvgFunRating,
			AvgSpeedFarmRating: build.AvgSpeedFarmRating,
			Tags:               *tags,
			Views:              build.Views,
			Status:             build.Status,
			CreatedAt:          build.CreatedAt,
		})
	}

	return &buildListResponse, nil
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
func (s *BuildService) GetBuildInfoByIdService(memberId uuid.UUID, buildId uuid.UUID) (*BuildInfoResponse, error) {
	return s.Repo.GetBuildInfo(memberId, buildId)
}

/**
* Adds primary and secondary skills and links to an existing build via a transanction,
* rolling back on error.
**/
func (s *BuildService) AddSkillLinksToBuildService(memberId uuid.UUID, buildId uuid.UUID, request AddSkillsToBuildRequest) error {

	return dbutils.ExecTx(s.Repo.DB, func(tx *sqlx.Tx) error {
		// get build and check if it exists
		_, err := s.GetBuildForMemberByIdService(memberId, buildId)

		if err != nil {
			return err
		}

		// -- CREATE SKILL LINKS FOR BUILD --

		// --- primary links ---

		// add main skill group
		mainSkillLinkId, err := s.Repo.CreateBuildSkillLinkTx(tx, buildId, request.MainSkillLinks.SkillLinkName, true)

		// add main skill relation to main skill link
		err = s.Repo.AddSkillToLinkTx(tx, mainSkillLinkId, request.MainSkillLinks.Skill)
		// create item mod

		if err != nil {
			return err
		}

		fmt.Printf("mainSkillLink created, id: %s\n", mainSkillLinkId)

		// create skill relations under this main skill link, one skill at a time
		for _, skillId := range request.MainSkillLinks.Links {
			err := s.Repo.AddSkillToLinkTx(tx, mainSkillLinkId, skillId)
			if err != nil {
				return err
			}
		}

		// --- other links --
		for _, skillLinks := range request.AdditionalSkills {

			// add secondary skill group
			secondarySkillLinkId, err := s.Repo.CreateBuildSkillLinkTx(tx, buildId, skillLinks.SkillLinkName, false)

			// add main skill relation to secondary link
			err = s.Repo.AddSkillToLinkTx(tx, secondarySkillLinkId, skillLinks.Skill)
			if err != nil {
				return err
			}

			fmt.Printf("secondarySkillLinkId created, id: %s\n", secondarySkillLinkId)

			// create skill relations under this secondary skill link, one skill at a time
			for _, skillId := range skillLinks.Links {
				err := s.Repo.AddSkillToLinkTx(tx, secondarySkillLinkId, skillId)

				if err != nil {
					return err
				}
			}
		}

		return nil
	})
}

/**
* Updates a specific member's build's skills and linksvia a transaction, rolling back on error.
**/
func (s *BuildService) UpdateBuildSkillLinksService(memberId uuid.UUID, buildId uuid.UUID, request UpdateSkillsToBuildRequest) error {
	return dbutils.ExecTx(s.Repo.DB, func(tx *sqlx.Tx) error {

		return fmt.Errorf("")
	})
}

/**
* Updates an average rating of a specific category for a build.
**/
func (s *BuildService) UpdateAvgRatingForBuildService(buildId string, category types.RatingCategory, avgRating float32) error {
	return s.Repo.UpdateAvgRatingForBuild(buildId, category, avgRating)
}
