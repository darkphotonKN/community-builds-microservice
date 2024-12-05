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

func (r *ClassRepository) BatchCreateDefaultClasses(classes []CreateDefaultClass) error {
	query := `
	INSERT INTO classes(id, name, description, image_url)
	VALUES(:id, :name, :description, :image_url)
	ON CONFLICT DO NOTHING
	`

	_, err := r.DB.NamedExec(query, classes)

	if err != nil {
		return errorutils.AnalyzeDBErr(err)
	}

	return nil
}

func (r *ClassRepository) BatchCreateDefaultAscendancies(ascendancies []CreateDefaultAscendancy) error {
	query := `
	INSERT INTO ascendancies(id, class_id, name, description, image_url)
	VALUES(:id, :class_id, :name, :description, :image_url)
	ON CONFLICT DO NOTHING
	`
	_, err := r.DB.NamedExec(query, ascendancies)

	if err != nil {
		return errorutils.AnalyzeDBErr(err)
	}

	return nil
}
