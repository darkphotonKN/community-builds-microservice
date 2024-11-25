package build

import (
	"github.com/google/uuid"
)

// --- Request ---

type CreateBuildRequest struct {
	SkillID     uuid.UUID   `json:"skillId" binding:"required" db:"main_skill_id"`
	TagIDs      []uuid.UUID `json:"tagIds" binding:"required" db:"tag_ids"`
	Title       string      `json:"title" binding:"required,min=6" db:"title"`
	Description string      `json:"description" binding:"required,min=10" db:"description"`
}

type SkillLinks struct {
	SkillLinkName string      `json:"skillLinkName" binding:"required"`
	Skill         uuid.UUID   `json:"skill" binding:"required,uuid"`
	Links         []uuid.UUID `json:"links" binding:"required,max=6,dive,uuid"`
}

type AddSkillsToBuildRequest struct {
	MainSkillLinks   SkillLinks   `json:"mainSkillLinks" binding:"required"`
	AdditionalSkills []SkillLinks `json:"additionalSkills" binding:"required"`
}

// --- Response ---

// To TEMP hold rows of builds JOIN build_skills JOIN skills data
type BuildInfoRows struct {
	ID            uuid.UUID `db:"id"`
	Title         string    `db:"title"`
	Description   string    `db:"description"`
	SkillLinkName string    `db:"skill_link_name"`
	IsMain        bool      `db:"skill_link_is_main"`
}

type TempSkillLink struct {
	Name   string `json:"name"`
	IsMain bool   `json:"isMain"`
}

type BuildInfoResponse struct {
	ID          uuid.UUID       `json:"id"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	Skills      []TempSkillLink `json:"skills"`
}
