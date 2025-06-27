package skill

import (
	commonhelpers "github.com/darkphotonKN/community-builds-microservice/common/utils"
	"github.com/jmoiron/sqlx"
)

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{db: db}
}

func (s *repository) CreateSkill(createSkillRequest CreateSkillRequest) error {
	query := `
	INSERT INTO skills(type, name)
	VALUES(:type, :name)
	`

	_, err := s.db.NamedExec(query, createSkillRequest)

	if err != nil {
		return commonhelpers.AnalyzeDBErr(err)
	}

	return nil
}

// func (s *SkillRepository) BatchCreateSkills(skills []SeedSkill) error {
// 	query := `
// 	INSERT INTO skills (id, name, type)
// 	VALUES (:id, :name, :type)
//   ON CONFLICT (name) DO NOTHING
//   `
// 	// batch insert skills with a slice of structs via sqlx
// 	_, err := s.DB.NamedExec(query, skills)

// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (s *SkillRepository) GetSkill(id uuid.UUID) (*models.Skill, error) {
// 	var skill models.Skill

// 	query := `
// 	SELECT * FROM skills
// 	WHERE id = $1
// 	`

// 	err := s.DB.Get(&skill, query, id)

// 	if err != nil {
// 		return nil, errorutils.AnalyzeDBErr(err)
// 	}

// 	return &skill, nil
// }

// func (s *SkillRepository) GetSkills() (*[]models.Skill, error) {
// 	var skills []models.Skill

// 	query := `SELECT * FROM skills`

// 	err := s.DB.Select(&skills, query)

// 	if err != nil {
// 		return nil, errorutils.AnalyzeDBErr(err)
// 	}

// 	return &skills, nil
// }

// func (s *SkillRepository) GetSkillsAndLinksByBuildId(buildId uuid.UUID) (*[]models.SkillRow, error) {

// 	var skillRows []models.SkillRow

// 	query := `
// 	SELECT
// 		build_skill_links.id as skill_link_id,
// 		build_skill_links.name as skill_link_name,
// 		build_skill_links.is_main as skill_link_is_main,
// 		skills.id as skill_id,
// 		skills.name as skill_name,
// 		skills.type as skill_type
// 	FROM builds
// 	JOIN build_skill_links ON build_skill_links.build_id = builds.id
// 	JOIN build_skill_link_skills ON build_skill_link_skills.build_skill_link_id = build_skill_links.id
// 	JOIN skills ON skills.id = build_skill_link_skills.skill_id
// 	WHERE builds.id = $1
// 	ORDER BY build_skill_links.id
// 	`

// 	err := s.DB.Select(&skillRows, query, buildId)

// 	if err != nil {
// 		return nil, errorutils.AnalyzeDBErr(err)
// 	}

// 	return &skillRows, nil

// }
