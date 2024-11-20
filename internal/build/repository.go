package build

import (
	"github.com/darkphotonKN/community-builds/internal/errorutils"
	"github.com/darkphotonKN/community-builds/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type BuildRepository struct {
	DB *sqlx.DB
}

func NewBuildRepository(db *sqlx.DB) *BuildRepository {
	return &BuildRepository{
		DB: db,
	}
}

func (r *BuildRepository) CreateBuild(memberId uuid.UUID, createBuildRequest CreateBuildRequest) error {
	query := `
	INSERT INTO builds(member_id, main_skill_id, title, description)
	VALUES(:member_id, :main_skill_id, :title, :description)
	`
	params := map[string]interface{}{
		"member_id":     memberId,
		"main_skill_id": createBuildRequest.SkillID,
		"title":         createBuildRequest.Title,
		"description":   createBuildRequest.Description,
	}

	_, err := r.DB.NamedExec(query, params)

	if err != nil {
		return errorutils.AnalyzeDBErr(err)
	}

	return nil
}

func (r *BuildRepository) GetBuildsByMemberId(memberId uuid.UUID) (*[]models.Build, error) {
	var builds []models.Build

	query := `
	SELECT * FROM builds
	WHERE member_id = $1
	`

	err := r.DB.Select(&builds, query, memberId)

	if err != nil {
		return nil, errorutils.AnalyzeDBErr(err)
	}

	return &builds, nil
}

func (r *BuildRepository) GetBuildForMemberById(memberId uuid.UUID, buildId uuid.UUID) (*models.Build, error) {
	var build models.Build

	query := `
	SELECT * FROM builds
	WHERE member_id = $1
	AND id = $2
	`

	err := r.DB.Get(&build, query, memberId, buildId)

	if err != nil {
		return nil, errorutils.AnalyzeDBErr(err)
	}

	return &build, nil
}

/**
* Getting all information related with builds via joins of
* join table build skills, builds, and skills.
**/
func (r *BuildRepository) GetBuildInfo(memberId uuid.UUID, buildId uuid.UUID) (*BuildInfoResponse, error) {
	var tempRowData []BuildInfoRows

	query := `
	SELECT 
		builds.id AS id, 
		builds.title AS title, 
		builds.description AS description,
		skills.id AS skill_id, 
		skills.name AS skill_name, 
		skills.type AS skill_type
	FROM builds
	JOIN build_skills ON build_skills.build_id = builds.id
	JOIN skills ON skills.id = build_skills.skill_id
	WHERE builds.id = $1 AND builds.member_id = $2
	`

	err := r.DB.Select(&tempRowData, query, buildId, memberId)

	if err != nil {
		return nil, errorutils.AnalyzeDBErr(err)
	}

	// build up base of the response
	result := BuildInfoResponse{
		ID:          tempRowData[0].ID,
		Title:       tempRowData[0].Title,
		Description: tempRowData[0].Description,
	}

	// group up all skill information
	for _, row := range tempRowData {
		result.Skills = append(result.Skills, models.Skill{
			ID:   row.ID,
			Name: row.SkillName,
			Type: row.SkillType,
		})
	}

	return &result, nil
}

func (r *BuildRepository) InsertSkillToBuild(buildId uuid.UUID, skillId uuid.UUID) error {
	query := `
	INSERT INTO build_skills(build_id, skill_id)
	VALUES(:build_id, :skill_id)
	`
	params := map[string]interface{}{
		"build_id": buildId,
		"skill_id": skillId,
	}

	_, err := r.DB.NamedExec(query, params)

	if err != nil {
		return errorutils.AnalyzeDBErr(err)
	}

	return nil
}
