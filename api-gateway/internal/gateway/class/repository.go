package class

import (
	"github.com/darkphotonKN/community-builds-microservice/api-gateway/internal/models"
	"github.com/darkphotonKN/community-builds-microservice/api-gateway/internal/utils/errorutils"
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

func (r *ClassRepository) GetClassesAndAscendancies() (*GetClassesAndAscendanciesResponse, error) {

	var classes []models.Class

	classQuery := `
	SELECT * FROM classes
	`
	err := r.DB.Select(&classes, classQuery)

	if err != nil {
		return nil, err
	}

	var ascendancies []models.Ascendancy

	ascendancyQuery := `
	SELECT * FROM ascendancies
	`

	err = r.DB.Select(&ascendancies, ascendancyQuery)

	if err != nil {
		return nil, err
	}

	response := GetClassesAndAscendanciesResponse{
		Classes:      classes,
		Ascendancies: ascendancies,
	}

	return &response, nil
}
