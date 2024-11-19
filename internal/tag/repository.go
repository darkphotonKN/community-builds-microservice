package tag

import (
	"errors"
	"fmt"

	"github.com/darkphotonKN/community-builds/internal/errorutils"
	"github.com/darkphotonKN/community-builds/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type TagRepository struct {
	DB *sqlx.DB
}

func NewTagRepository(db *sqlx.DB) *TagRepository {
	return &TagRepository{
		DB: db,
	}
}

func (r *TagRepository) CreateTag(createTagReq CreateTagRequest) error {
	query := `
		INSERT INTO tags(name)
		VALUES(:name)
	`

	_, err := r.DB.NamedExec(query, createTagReq)

	fmt.Print("Error when creating tag:", err)

	if err != nil {
		return errorutils.AnalyzeDBErr(err)
	}

	return nil
}

func (r *TagRepository) GetTags() (*[]models.Tag, error) {
	var tags []models.Tag

	query := `SELECT * FROM tags`

	err := r.DB.Select(&tags, query)

	if err != nil {
		return nil, errorutils.AnalyzeDBErr(err)
	}

	return &tags, nil
}

func (r *TagRepository) UpdateTag(id uuid.UUID, payload UpdateTagRequest) (*UpdateTagRequest, error) {

	params := UpdateTagParams{
		ID:   id,
		Name: payload.Name,
	}
	query := `UPDATE tags SET name = :name WHERE id = :id`

	result, err := r.DB.NamedExec(query, params)

	if err != nil {
		return nil, errorutils.AnalyzeDBErr(err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return nil, errors.New("no rows updated")
	}
	return &payload, nil
}
