package build

import (
	"database/sql"
	"errors"
	"fmt"

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
func (r *BuildRepository) CreateBuildSkillLinkTx(tx *sqlx.Tx, buildId uuid.UUID, name string, isMain bool) (uuid.UUID, error) {

	// validate that skill doesn't already exists first
	var existsId uuid.UUID

	query := `
	SELECT id FROM build_skill_links
	WHERE build_id = $1 AND name = $2
	`

	err := tx.Get(&existsId, query, buildId, name)

	if !errors.Is(err, sql.ErrNoRows) {
		return uuid.Nil, errorutils.ErrDuplicateResource
	}

	var newId uuid.UUID

	query = `
	INSERT INTO build_skill_links(build_id, name, is_main)
	VALUES($1, $2, $3)
	RETURNING id
	`

	err = tx.QueryRowx(query, buildId, name, isMain).Scan(&newId)

	if err != nil {
		fmt.Printf("Error when attempting to insert into build_skill_links: %s\n", err)
		return uuid.Nil, errorutils.AnalyzeDBErr(err)
	}

	return newId, nil
}

/**
* Adds a skill to a existing skill link.
**/
func (r *BuildRepository) AddSkillToLinkTx(tx *sqlx.Tx, buildSkillLinkId uuid.UUID, skillId uuid.UUID) error {

	// validate that skill doesn't already exists first
	var existsId uuid.UUID

	query := `
	SELECT id FROM build_skill_link_skills 
	WHERE build_skill_link_id = $1 AND skill_id = $2
	`

	err := tx.Get(&existsId, query, buildSkillLinkId, skillId)

	// if resource IS found, don't create duplicate skill-link to skill relation insert
	if !errors.Is(err, sql.ErrNoRows) {
		fmt.Println("Rows found, duplicate.")
		return errorutils.ErrDuplicateResource
	}

	query = `
	INSERT INTO build_skill_link_skills(build_skill_link_id, skill_id)
	VALUES(:build_skill_link_id, :skill_id)
	`

	params := map[string]interface{}{
		"build_skill_link_id": buildSkillLinkId,
		"skill_id":            skillId,
	}

	_, err = tx.NamedExec(query, params)

	if err != nil {
		fmt.Printf("DEBUG AddSkillToLinkTx: Error when attempting to insert into join table build_skill_link_skills: %s\n", err)
		return errorutils.AnalyzeDBErr(err)
	}

	return nil
}
