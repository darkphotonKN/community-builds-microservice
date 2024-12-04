package class

import "github.com/google/uuid"

type CreateDefaultClass struct {
	ID          uuid.UUID `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description"`
	ImageURL    string    `db:"image_url" json:"imageUrl"`
}

type CreateAscendancy struct {
	ID       uuid.UUID `db:"id" json:"id"`
	ClassID  uuid.UUID `db:"class_id" json:"classId"`
	Name     string    `db:"name" json:"name"`
	ImageURL string    `db:"image_url" json:"imageUrl"`
}
