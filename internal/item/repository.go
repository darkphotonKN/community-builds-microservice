package item

import (
	"database/sql"

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

func (r *ItemRepository) CreateItem(item models.Item) error {
	query := `
		INSERT INTO items(user_id, product_id, name, description, price_per_unit, stock_quantity)
	VALUES(:user_id, :product_id, :name, :description, :price_per_unit, :stock_quantity)
	`

	_, err := r.DB.NamedExec(query, item)

	if err != nil {
		return err
	}

	return nil
}

func (r *ItemRepository) GetItems(userId uuid.UUID) (*[]models.Item, error) {
	var items []models.Item

	query := `
	SELECT * FROM items
	WHERE user_id = $1
	`

	err := r.DB.Select(&items, query, userId)

	if err != nil {
		return nil, err
	}

	return &items, nil
}

func (r *ItemRepository) UpdateItemById(userId uuid.UUID, id uuid.UUID, updateItemReq UpdateItemReq) (*models.Item, error) {
	var item models.Item

	query := `
	UPDATE items
	SET name = :name,
		category = :category,
		type = :type,
		description = :description,
		price_per_unit = :price_per_unit,
		stock_quantity = :stock_quantity
	WHERE user_id = :user_id AND id = :id
	RETURNING *;
	`

	params := map[string]interface{}{
		"id":             id,
		"user_id":        userId,
		"name":           updateItemReq.Name,
		"type":           updateItemReq.Type,
		"category":       updateItemReq.Category,
		"description":    updateItemReq.Description,
		"price_per_unit": updateItemReq.PricePerUnit,
		"stock_quantity": updateItemReq.StockQuantity,
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
