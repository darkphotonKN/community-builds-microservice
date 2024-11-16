package skill

import "github.com/jmoiron/sqlx"

type SkillRepository struct {
	DB *sqlx.DB
}

func NewSkillRepository(db *sqlx.DB) *SkillRepository {
	return &SkillRepository{
		DB: db,
	}
}

func (s *SkillRepository) CreateSkill(createSkillRequest CreateSkillRequest) error {
}
