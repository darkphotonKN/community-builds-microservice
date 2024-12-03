package class

import "github.com/google/uuid"

type CreateClass struct {
	Name        string `db:"name" json:"name"`
	Description string `db:"description" json:"description"`
	ImageURL    string `db:"image_url" json:"imageUrl"`
}

type CreateAscendancy struct {
	ClassID  uuid.UUID `db:"class_id" json:"classId"`
	Name     string    `db:"name" json:"name"`
	ImageURL string    `db:"image_url" json:"imageUrl"`
}
