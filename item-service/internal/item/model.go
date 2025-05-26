package item

import "github.com/google/uuid"

type CreateItemRequest struct {
	Category string `json:"category" binding:"required,category" db:"category"`
	Class    string `json:"class" binding:"required,class" db:"class"`
	Type     string `json:"type" binding:"required,type" db:"type"`
	Name     string `json:"name" binding:"required,min=2" db:"name"`
	ImageURL string `json:"imageUrl,omitempty" db:"image_url"`
}

type UpdateItemReq struct {
	Id       uuid.UUID `json:"id" binding:"required" db:"id"`
	Category string    `json:"category" binding:"required,category" db:"category"`
	Class    string    `json:"class" binding:"required,class" db:"class"`
	Type     string    `json:"type" binding:"required,type" db:"type"`
	Name     string    `json:"name" binding:"required,min=2" db:"name"`
	ImageURL string    `json:"imageUrl,omitempty" db:"image_url"`
}

type CreateRareItemReq struct {
	BaseItemId uuid.UUID `json:"baseItemId" db:"base_item_id"`
	ToList     bool      `json:"toList"`
	Name       string    `json:"name" db:"name"`
	Stats      []string  `json:"stats" db:"stats"`
}
