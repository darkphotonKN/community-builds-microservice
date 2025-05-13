package example

import (
	"time"
)

// Example represents a basic example entity
type Example struct {
	ID        string    `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// ExampleCreate represents the data needed to create a new example
type ExampleCreate struct {
	Name string `json:"name" validate:"required"`
}

// ExampleUpdate represents the data needed to update an existing example
type ExampleUpdate struct {
	Name string `json:"name" validate:"required"`
} 