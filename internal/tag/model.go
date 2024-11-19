package tag

import "github.com/google/uuid"

type CreateTagRequest struct {
	Name string `json:"name" binding:"required" db:"name"`
}

type UpdateTagRequest struct {
	Name string `json:"name" binding:"required" db:"name"`
}

type UpdateTagParams struct {
	ID   uuid.UUID `db:"id" json:"id"`
	Name string    `db:"name" json:"name"`
}
