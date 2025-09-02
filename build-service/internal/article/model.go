package article

import (
	"github.com/google/uuid"
)

type CreateArticle struct {
	// Id      uuid.UUID `db:"id" json:"id"`
	Content string `db:"content" json:"content"`
}

type UpdateArticle struct {
	Id      uuid.UUID `db:"id" json:"id"`
	Content string    `db:"content" json:"content"`
}

type Article struct {
	Id      uuid.UUID `db:"id" json:"id"`
	Content string    `db:"content" json:"content"`
}
