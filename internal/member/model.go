package member

import (
	"github.com/darkphotonKN/community-builds/internal/models"
	"github.com/google/uuid"
)

type MemberResponse struct {
	models.BaseDBDateModel
	Email string `db:"email" json:"email"`
	Name  string `db:"name" json:"name"`
}

type MemberLoginRequest struct {
	Email    string `db:"email" json:"email"`
	Password string `db:"password" json:"password"`
}

type MemberLoginResponse struct {
	RefreshToken     string `json:"refreshToken"`
	AccessToken      string `json:"accessToken"`
	AccessExpiresIn  int    `json:"accessExpiresIn"`
	RefreshExpiresIn int    `json:"refreshExpiresIn"`

	MemberInfo *models.Member `json:"memberInfo"`
}

type MemberUpdatePasswordRequest struct {
	ID                uuid.UUID `db:"id" json:"id"`
	Password          string    `db:"password" json:"password"`
	NewPassword       string    `json:"newPassword"`
	RepeatNewPassword string    `json:"repeatNewPassword"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken"`
}
