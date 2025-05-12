package article

import "github.com/google/uuid"

type CreateArticleRequest struct {
	Content string `json:"content" binding:"required" db:"content"`
}

type UpdateArticleRequest struct {
	ID      uuid.UUID `db:"id" json:"id"`
	Content string    `json:"content" binding:"required" db:"content"`
}

type UpdateArticleParams struct {
	ID      uuid.UUID `db:"id" json:"id"`
	Content string    `db:"content" json:"content"`
}
