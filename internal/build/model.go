package build

import "github.com/google/uuid"

type CreateBuildRequest struct {
	SkillID     uuid.UUID `json:"skillId" binding:"required" db:"main_skill_id"`
	Title       string    `json:"title" binding:"required,min=6" db:"title"`
	Description string    `json:"description" binding:"required,min=10" db:"description"`
}
