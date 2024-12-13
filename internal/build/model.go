package build

import (
	"github.com/darkphotonKN/community-builds/internal/models"
	"github.com/google/uuid"
)

// --- Request ---

type CreateBuildRequest struct {
	SkillID      uuid.UUID   `json:"skillId" binding:"required" db:"main_skill_id"`
	TagIDs       []uuid.UUID `json:"tagIds" binding:"required" db:"tag_ids"`
	Title        string      `json:"title" binding:"required,min=6" db:"title"`
	Description  string      `json:"description" binding:"required,min=10" db:"description"`
	ClassID      uuid.UUID   `json:"classId" binding:"required" db:"class_id"`
	AscendancyID uuid.UUID   `json:"ascendancyId" db:"ascendancy_id"`
}

type UpdateBuildRequest struct {
	SkillID      uuid.UUID   `json:"skillId" binding:"omitempty" db:"main_skill_id"`
	TagIDs       []uuid.UUID `json:"tagIds" binding:"omitempty" db:"tag_ids"`
	Title        string      `json:"title" binding:"omitempty,min=6" db:"title"`
	Description  string      `json:"description"  binding:"omitempty,min=10" db:"description"`
	ClassID      uuid.UUID   `json:"classId" binding:"omitempty" db:"class_id"`
	AscendancyID uuid.UUID   `json:"ascendancyId" binding:"omitempty" db:"ascendancy_id"`
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

type UpdateSkillsToBuildRequest struct {
	MainSkillLinks   SkillLinks   `json:"mainSkillLinks" binding:"required"`
	AdditionalSkills []SkillLinks `json:"additionalSkills" binding:"required"`
}

type AddItemsToBuildRequest struct {
	Weapon     string `json:"weapon"`
	Shield     string `json:"shield"`
	Helmet     string `json:"helmet"`
	BodyArmour string `json:"bodyArmour"`
	Boots      string `json:"boots"`
	Gloves     string `json:"gloves"`
	Belt       string `json:"Belt"`
	Amulet     string `json:"amulet"`
	LeftRing   string `json:"leftRing"`
	RightRing  string `json:"reftRing"`
}

// --- Response ---

// TEMP CONTAINER for builds JOIN build_skill_links JOIN skills data
type BuildInfoRows struct {
	ID              uuid.UUID `db:"id"`
	Title           string    `db:"title"`
	Description     string    `db:"description"`
	SkillLinkID     string    `db:"skill_link_id"`
	SkillLinkName   string    `db:"skill_link_name"`
	SkillLinkIsMain bool      `db:"skill_link_is_main"`
	SkillID         uuid.UUID `db:"skill_id"`
	SkillName       string    `db:"skill_name"`
	SkillType       string    `db:"skill_type"`
}

type TempSkillLink struct {
	Name   string `json:"name"`
	IsMain bool   `json:"isMain"`
}

// Data Structure of Skill Links + Information for Response
type SkillLinkResponse struct {
	SkillLinkName string         `json:"skillLinkName"`
	Skill         models.Skill   `json:"skill"`
	Links         []models.Skill `json:"links"`
}

type SkillGroupResponse struct {
	MainSkillLinks   SkillLinkResponse   `json:"mainSkillLinks"`
	AdditionalSkills []SkillLinkResponse `json:"additionalSkills"`
}

// All Build Information
type BuildInfoResponse struct {
	ID          uuid.UUID           `json:"id"`
	Title       string              `json:"title"`
	Description string              `json:"description"`
	Class       string              `json:"class"`
	Ascendancy  *string             `json:"ascendancy,omitempty"`
	Tags        []models.Tag        `json:"tags"`
	Skills      *SkillGroupResponse `json:"skills"`
}

// Build List
type BuildListQuery struct {
	ID                 uuid.UUID `json:"id"`
	Title              string    `db:"title" json:"title"`
	Description        string    `db:"description" json:"description"`
	Class              string    `db:"class_name" json:"class"`
	Ascendancy         *string   `db:"ascendancy_name" json:"ascendancy"`
	MainSkillName      string    `db:"main_skill_name" json:"mainSkill"`
	AvgEndGameRating   *float32  `db:"avg_end_game_rating" json:"avgEndGameRating,omitempty"`
	AvgFunRating       *float32  `db:"avg_fun_rating" json:"avgFunRating,omitempty"`
	AvgCreativeRating  *float32  `db:"avg_creative_rating" json:"avgCreativeRating,omitempty"`
	AvgSpeedFarmRating *float32  `db:"avg_speed_farm_rating" json:"avgSpeedFarmRating,omitempty"`
	AvgBossingRating   *float32  `db:"avg_bossing_rating" json:"avgBossingRating,omitempty"`
	Views              int       `db:"views" json:"views"`
	Status             int       `db:"status" json:"status"`
	CreatedAt          string    `db:"created_at" json:"createdAt"`
}

type BuildListResponse struct {
	ID                 uuid.UUID    `json:"id"`
	Title              string       `json:"title"`
	Description        string       `json:"description"`
	Class              string       `json:"class"`
	Ascendancy         *string      `json:"ascendancy"`
	MainSkillName      string       `json:"mainSkill"`
	AvgEndGameRating   *float32     `json:"avgEndGameRating,omitempty"`
	AvgFunRating       *float32     `json:"avgFunRating,omitempty"`
	AvgCreativeRating  *float32     `json:"avgCreativeRating,omitempty"`
	AvgSpeedFarmRating *float32     `json:"avgSpeedFarmRating,omitempty"`
	AvgBossingRating   *float32     `json:"avgBossingRating,omitempty"`
	Views              int          `json:"views"`
	Tags               []models.Tag `json:"tags"`
	Status             int          `json:"status"`
	CreatedAt          string       `json:"createdAt"`
}

type BuildListRespose struc
