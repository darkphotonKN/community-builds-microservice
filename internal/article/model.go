package article

import "github.com/google/uuid"

type CreateArticleRequest struct {
	Name string `json:"name" binding:"required" db:"name"`
}

type UpdateArticleRequest struct {
	ID   uuid.UUID `db:"id" json:"id"`
	Name string    `json:"name" binding:"required" db:"name"`
}

type UpdateArticleParams struct {
	ID   uuid.UUID `db:"id" json:"id"`
	Name string    `db:"name" json:"name"`
}
