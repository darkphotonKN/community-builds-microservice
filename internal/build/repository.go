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
