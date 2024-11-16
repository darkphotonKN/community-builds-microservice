package skill

type CreateSkillRequest struct {
	Name string `json:"name" binding:"required" db:"name"`
	Type string `json:"name" binding:"required" db:"type"`
}
