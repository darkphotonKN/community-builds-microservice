package skill

import "github.com/google/uuid"

type CreateSkillRequest struct {
	Name string `json:"name" binding:"required" db:"name"`
	Type string `json:"type" binding:"required,skillType" db:"type"`
}

type SeedSkill struct {
	ID   uuid.UUID `db:"id"`
	Name string    `db:"name"`
	Type string    `db:"type"`
}
