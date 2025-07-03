package class

import (
	"github.com/darkphotonKN/community-builds-microservice/common/constants/models"
	commonhelpers "github.com/darkphotonKN/community-builds-microservice/common/utils"
	"github.com/jmoiron/sqlx"
)

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{db: db}
}

func (r *repository) BatchCreateDefaultClasses(classes []CreateDefaultClass) error {
	query := `
	INSERT INTO classes(id, name, description, image_url)
	VALUES(:id, :name, :description, :image_url)
	ON CONFLICT DO NOTHING
	`

	_, err := r.db.NamedExec(query, classes)

	if err != nil {
		return commonhelpers.AnalyzeDBErr(err)
	}

	return nil
}

func (r *repository) BatchCreateDefaultAscendancies(ascendancies []CreateDefaultAscendancy) error {
	query := `
	INSERT INTO ascendancies(id, class_id, name, description, image_url)
	VALUES(:id, :class_id, :name, :description, :image_url)
	ON CONFLICT DO NOTHING
	`
	_, err := r.db.NamedExec(query, ascendancies)

	if err != nil {
		return commonhelpers.AnalyzeDBErr(err)
	}

	return nil
}

func (r *repository) GetClassesAndAscendancies() (*GetClassesAndAscendanciesResponse, error) {

	var classes []models.Class

	classQuery := `
	SELECT * FROM classes
	`
	err := r.db.Select(&classes, classQuery)

	if err != nil {
		return nil, err
	}

	var ascendancies []models.Ascendancy

	ascendancyQuery := `
	SELECT * FROM ascendancies
	`

	err = r.db.Select(&ascendancies, ascendancyQuery)

	if err != nil {
		return nil, err
	}

	response := GetClassesAndAscendanciesResponse{
		Classes:      classes,
		Ascendancies: ascendancies,
	}

	return &response, nil
}
