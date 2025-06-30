package tag

import (
	"errors"
	"fmt"

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

func (r *repository) CreateTag(createTagReq CreateTagRequest) error {
	query := `
		INSERT INTO tags(name)
		VALUES(:name)
	`

	_, err := r.db.NamedExec(query, createTagReq)

	fmt.Print("Error when creating tag:", err)

	if err != nil {
		return commonhelpers.AnalyzeDBErr(err)
	}

	return nil
}

func (r *repository) GetTags() (*[]models.Tag, error) {
	var tags []models.Tag

	query := `SELECT * FROM tags`

	err := r.db.Select(&tags, query)

	if err != nil {
		return nil, commonhelpers.AnalyzeDBErr(err)
	}

	return &tags, nil
}

func (r *repository) UpdateTag(payload UpdateTagRequest) error {

	query := `UPDATE tags SET name = :name WHERE id = :id`

	result, err := r.db.NamedExec(query, payload)

	if err != nil {
		return commonhelpers.AnalyzeDBErr(err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("no rows updated")
	}
	return nil
}
