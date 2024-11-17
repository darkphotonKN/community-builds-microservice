package skill

import (
	"github.com/darkphotonKN/community-builds/internal/errorutils"
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
