package class

import (
	"github.com/darkphotonKN/community-builds/internal/utils/errorutils"
	"github.com/jmoiron/sqlx"
)

type ClassRepository struct {
	DB *sqlx.DB
}

func NewClassRepository(db *sqlx.DB) *ClassRepository {
	return &ClassRepository{
		DB: db,
	}
}

func (r *ClassRepository) CreateClass(class CreateClass) error {
	query := `
	INSERT INTO classes(name, description, image_url)
	VALUES(:name, :description, :image_url)
	ON CONFLICT DO NOTHING
	`

	_, err := r.DB.NamedExec(query, class)

	if err != nil {
		return errorutils.AnalyzeDBErr(err)
	}

	return nil
}
