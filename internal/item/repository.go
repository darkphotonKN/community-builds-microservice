package item

import (
	"database/sql"
	"fmt"

	"github.com/darkphotonKN/community-builds/internal/errorutils"
	"github.com/darkphotonKN/community-builds/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ItemRepository struct {
	DB *sqlx.DB
}

func NewItemRepository(db *sqlx.DB) *ItemRepository {
	return &ItemRepository{
		DB: db,
	}
}

func (r *ItemRepository) CreateItem(createItemReq CreateItemRequest) error {
	query := `
		INSERT INTO items(name, category, class, type, image_url)
		VALUES(:name, :category, :class, :type,  :image_url)
	`

	_, err := r.DB.NamedExec(query, createItemReq)

	fmt.Print("Error when creating item:", err)

	if err != nil {
		return errorutils.AnalyzeDBErr(err)
	}

	return nil
}

func (r *ItemRepository) GetItems() (*[]models.Item, error) {
	var items []models.Item

	query := `
	SELECT * FROM items
	`

	err := r.DB.Select(&items, query)

	if err != nil {
		return nil, errorutils.AnalyzeDBErr(err)
	}

	return &items, nil
}

func (r *ItemRepository) UpdateItemById(id uuid.UUID, updateItemReq UpdateItemReq) (*models.Item, error) {
	var item models.Item

	query := `
	UPDATE items
	SET name = :name,
		category = :category,
		type = :type,
		class = :class,
		img_url = :img_url,
	WHERE user_id = :user_id AND id = :id
	RETURNING *;
	`

	params := map[string]interface{}{
		"id":       id,
		"name":     updateItemReq.Name,
		"type":     updateItemReq.Type,
		"category": updateItemReq.Category,
		"class":    updateItemReq.Class,
		"img_url":  updateItemReq.ImageURL,
	}

	rows, err := r.DB.NamedQuery(query, params)
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	// loop through and check next table row exists
	if rows.Next() {
		// map the row data to our item struct
		err := rows.StructScan(&item)

		if err != nil {
			return nil, err
		}
	} else {
		// no results found
		return nil, sql.ErrNoRows
	}

	return &item, nil
}
