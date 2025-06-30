package tag

import "github.com/google/uuid"

type CreateTagRequest struct {
	Name string `json:"name" binding:"required" db:"name"`
}

type UpdateTagRequest struct {
	Id   uuid.UUID `db:"id" json:"id"`
	Name string    `json:"name" binding:"required" db:"name"`
}

type UpdateTagParams struct {
	Id   uuid.UUID `db:"id" json:"id"`
	Name string    `db:"name" json:"name"`
}
