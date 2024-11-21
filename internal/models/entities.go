package models

import (
	"time"

	"github.com/google/uuid"
)

/**
* Types here are shared model entities that are imported by more than one package.
**/

/**
* Member
**/
type Member struct {
	BaseDBDateModel
	Email         string  `db:"email" json:"email"`
	Name          string  `db:"name" json:"name"`
	Password      string  `db:"password" json:"password,omitempty"`
	Status        string  `db:"status" json:"status"`
	AverageRating float64 `db:"average_rating"`
}

/**
* Items
**/
type Item struct {
	BaseDBDateModel
	Category string `db:"category" json:"category"`
	Class    string `db:"class" json:"class"`
	Type     string `db:"type" json:"type"`
	Name     string `db:"name" json:"name"`
	ImageURL string `db:"image_url" json:"imageUrl"`
}

/**
* Skills
**/
type Skill struct {
	ID        uuid.UUID `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	Type      string    `db:"type" json:"type"` // "active" or "support"
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
}

/**
* Ratings
**/
type Rating struct {
	BaseIDModel
	BuildID uuid.UUID `db:"build_id" json:"memberId"`
	Rating  int       `db:"rating" json:"rating"`
}

/**
* Build
**/
type Build struct {
	BaseDBDateModel
	MemberID    uuid.UUID `db:"member_id" json:"memberId"`
	MainSkillID uuid.UUID `db:"main_skill_id" json:"mainSkill"`
	Title       string    `db:"title" json:"title"`
	Description string    `db:"description" json:"description"`
}

type BuildItem struct {
	ID      uuid.UUID `db:"id" json:"id"`
	BuildID uuid.UUID `db:"build_id" json:"buildId"`
	ItemID  uuid.UUID `db:"item_id" json:"itemId"`
	Slot    string    `db:"slot" json:"slot"`
}

type BuildSkill struct {
	ID      uuid.UUID `db:"id" json:"id"`
	BuildID uuid.UUID `db:"build_id" json:"buildId"`
	SkillID uuid.UUID `db:"skill_id" json:"skillId"`
}

/**
* Items
**/
type Tag struct {
	BaseDBDateModel
	Name string `db:"name" json:"name"`
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
