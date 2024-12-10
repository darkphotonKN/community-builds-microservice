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

// Temp Struct to hold skill, skill link, joined with build table data
type SkillRow struct {
	SkillLinkID     string    `db:"skill_link_id"`
	SkillLinkName   string    `db:"skill_link_name"`
	SkillLinkIsMain bool      `db:"skill_link_is_main"`
	SkillID         uuid.UUID `db:"skill_id"`
	SkillName       string    `db:"skill_name"`
	SkillType       string    `db:"skill_type"`
}
