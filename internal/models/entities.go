package models

import (
	"github.com/google/uuid"
	"time"
)

/**
* Shared entities that are imported by more than one package.
**/
type Member struct {
	BaseDBDateModel
	Email         string  `db:"email" json:"email"`
	Name          string  `db:"name" json:"name"`
	Password      string  `db:"password" json:"password,omitempty"`
	Status        string  `db:"status" json:"status"`
	AverageRating float64 `db:"average_rating"`
	ResponseTime  int     `db:"response_time"`
	TotalTrades   int     `db:"total_trades"`
}

type Rating struct {
	BaseIDModel
	MemberID uuid.UUID `db:"member_id" json:"memberId"`
	Rating   int       `db:"rating" json:"rating"`
}

type Item struct {
	BaseDBDateModel
	MemberID      uuid.UUID `db:"member_id" json:"memberId"`
	ProductID     uuid.UUID `db:"product_id" json:"productId"`
	Category      string    `db:"category" json:"category"`
	Type          string    `db:"type" json:"type"`
	Name          string    `db:"name" json:"name"`
	Description   string    `db:"description" json:"description"`
	PricePerUnit  float64   `db:"price_per_unit" json:"pricePerUnit"`
	StockQuantity int       `db:"stock_quantity" json:"stockQuantity"`
}

/**
* Base models for default table columns.
**/

type BaseIDModel struct {
	ID        uuid.UUID `db:"id" json:"id"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
}

type BaseDBMemberModel struct {
	ID            uuid.UUID `db:"id" json:"id"`
	UpdatedMember uuid.UUID `db:"updated_member" json:"updatedMember"`
	CreatedMember uuid.UUID `db:"created_member" json:"createdMember"`
}

type BaseDBDateModel struct {
	ID        uuid.UUID `db:"id" json:"id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type BaseDBMemberDateModel struct {
	ID            uuid.UUID `db:"id" json:"id"`
	UpdatedMember uuid.UUID `db:"updated_member" json:"updatedMember"`
	CreatedMember uuid.UUID `db:"created_member" json:"createdMember"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
}
