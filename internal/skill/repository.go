package skill

import (
	"github.com/darkphotonKN/community-builds/internal/errorutils"
	"github.com/darkphotonKN/community-builds/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type SkillRepository struct {
	DB *sqlx.DB
}

func NewSkillRepository(db *sqlx.DB) *SkillRepository {
	return &SkillRepository{
		DB: db,
	}
}

func (s *SkillRepository) CreateSkill(createSkillRequest CreateSkillRequest) error {
	query := `
	INSERT INTO skills(type, name)
	VALUES(:type, :name)
	`

	_, err := s.DB.NamedExec(query, createSkillRequest)

	if err != nil {
		return errorutils.AnalyzeDBErr(err)
	}

	return nil
}

func (s *SkillRepository) GetSkill(id uuid.UUID) (*models.Skill, error) {
	var skill models.Skill

	query := `
	SELECT * FROM skills
	WHERE id = $1
	`

	err := s.DB.Get(&skill, query, id)

	if err != nil {
		return nil, errorutils.AnalyzeDBErr(err)
	}

	return &skill, nil
}
