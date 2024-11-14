package member

import "github.com/darkphotonKN/community-builds/internal/models"

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

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken"`
}
