package build

import (
	"github.com/darkphotonKN/community-builds/internal/errorutils"
	"github.com/darkphotonKN/community-builds/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
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
	VALUES($1, $2, $3, $4)
	RETURNING id
	`
	var buildId uuid.UUID

	err := r.DB.QueryRowx(query, memberId, createBuildRequest.SkillID, createBuildRequest.Title, createBuildRequest.Description).Scan(&buildId)

	if err != nil {
		return errorutils.AnalyzeDBErr(err)
	}

	buildTagQuery := `
	INSERT INTO build_tags(build_id, tag_id)
	VALUES($1, unnest($2::uuid[]))
	`

	_, buildTagsErr := r.DB.Exec(buildTagQuery, buildId, pq.Array(createBuildRequest.TagIDs))
	if buildTagsErr != nil {
		return errorutils.AnalyzeDBErr(buildTagsErr)
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
	var buildRowData []BuildInfoRows

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

	err := r.DB.Select(&buildRowData, query, buildId, memberId)

	if err != nil {
		return nil, errorutils.AnalyzeDBErr(err)
	}

	// build up base of the response
	result := BuildInfoResponse{
		ID:          buildRowData[0].ID,
		Title:       buildRowData[0].Title,
		Description: buildRowData[0].Description,
	}

	// group up all skill information
	for _, row := range buildRowData {
		result.Skills = append(result.Skills, models.Skill{
			ID:   row.ID,
			Name: row.SkillName,
			Type: row.SkillType,
		})
	}

	return &result, nil
}

/**
* Creates a skill link group for a build.
**/
func (r *BuildRepository) CreateBuildSkillLink(buildId uuid.UUID, name string, isMain bool) (uuid.UUID, error) {
	var newId uuid.UUID

	query := `
	INSERT INTO build_skill_links(build_id, name, is_main)
	VALUES(:build_id, :name, :is_main)
	RETURNING id
	`
	params := map[string]interface{}{
		"build_id": buildId,
		"name":     name,
		"is_main":  isMain,
	}

	err := r.DB.QueryRowx(query, params).Scan(&newId)

	if err != nil {
		return uuid.Nil, errorutils.AnalyzeDBErr(err)
	}

	return newId, nil
}

/**
* Adds a skill to a existing skill link.
**/
func (r *BuildRepository) AddSkillToLink(buildSkillLinkId uuid.UUID, skillId uuid.UUID) error {
	query := `
	INSERT INTO build_skill_link_skills(build_skill_link_id, skill_id)
	VALUES(build_skill_link_id, skill_id)
	`

	params := map[string]interface{}{
		"build_skill_link_id": buildSkillLinkId,
		"skill_id":            skillId,
	}

	_, err := r.DB.NamedExec(query, params)

	if err != nil {
		return errorutils.AnalyzeDBErr(err)
	}

	return nil
}
