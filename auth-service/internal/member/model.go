package member

import (
	"time"

	"github.com/google/uuid"
)

type Member struct {
	ID            uuid.UUID `db:"id" json:"id"`
	Email         string    `db:"email" json:"email"`
	Name          string    `db:"name" json:"name"`
	Password      string    `db:"password" json:"password,omitempty"`
	Status        string    `db:"status" json:"status"`
	AverageRating float64   `db:"average_rating"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
}

type MemberUpdatePasswordParams struct {
	ID       uuid.UUID `db:"id" json:"id"`
	Password string    `db:"password" json:"password"`
}

type MemberUpdateInfoParams struct {
	ID     uuid.UUID `db:"id" json:"id"`
	Name   string    `db:"name" json:"name"`
	Status string    `db:"status" json:"status"`
}

type CreateDefaultMember struct {
	ID       uuid.UUID `db:"id"`
	Email    string    `db:"email"`
	Name     string    `db:"name"`
	Password string    `db:"password"`
	Status   int       `db:"status"`
}

