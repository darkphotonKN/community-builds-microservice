package member

import (
	"errors"
	"fmt"
	"time"

	"github.com/darkphotonKN/community-builds/internal/auth"
	"github.com/darkphotonKN/community-builds/internal/models"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type MemberService struct {
	Repo *MemberRepository
}

func NewMemberService(repo *MemberRepository) *MemberService {
	return &MemberService{
		Repo: repo,
	}
}

func (s *MemberService) GetMemberByIdWithPasswordService(id uuid.UUID) (*models.Member, error) {
	return s.Repo.GetByIdWithPassword(id)
}

func (s *MemberService) GetMemberByIdService(id uuid.UUID) (*models.Member, error) {
	return s.Repo.GetById(id)
}

func (s *MemberService) CreateMemberService(user models.Member) error {
	hashedPw, err := s.HashPassword(user.Password)

	if err != nil {
		return fmt.Errorf("Error when attempting to hash password.")
	}

	// update user's password with hashed password.
	user.Password = hashedPw

	return s.Repo.Create(user)
}

func (s *MemberService) UpdatePasswordMemberService(data MemberUpdatePasswordRequest) error {
	fmt.Println("data", data)

	user, _ := s.GetMemberByIdWithPasswordService(data.ID)

	if data.NewPassword != data.RepeatNewPassword {
		return fmt.Errorf("NewPassword and RepeatNewPassword are different")
	}

	isSame, _ := s.ComparePasswords(user.Password, data.Password)
	if !isSame {
		return fmt.Errorf("Stored password and input password are different")
	}

	hashedPw, err := s.HashPassword(data.NewPassword)
	if err != nil {
		return fmt.Errorf("Error when attempting to hash password.")
	}

	// update user's password with hashed password.
	data.Password = hashedPw

	return s.Repo.UpdatePassword(data)
}

func (s *MemberService) UpdateInfoMemberHandler(user models.Member) error {
	return s.Repo.UpdateInfo(user)
}

func (s *MemberService) LoginMemberService(loginReq MemberLoginRequest) (*MemberLoginResponse, error) {
	user, err := s.Repo.GetMemberByEmail(loginReq.Email)

	if err != nil {
		return nil, errors.New("Could not get user with provided email.")
	}

	// extract password, and compare hashes
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReq.Password)); err != nil {
		return nil, errors.New("The credentials provided was incorrect.")
	}

	// construct response with both user info and auth credentials
	accessExpiryTime := time.Minute * 60
	accessToken, err := auth.GenerateJWT(*user, auth.Access, accessExpiryTime)
	refreshExpiryTime := time.Hour * 24 * 7
	refreshToken, err := auth.GenerateJWT(*user, auth.Refresh, refreshExpiryTime)

	user.Password = ""

	res := &MemberLoginResponse{
		AccessToken:      accessToken,
		AccessExpiresIn:  int(accessExpiryTime),
		RefreshToken:     refreshToken,
		RefreshExpiresIn: int(refreshExpiryTime),
		MemberInfo:       user,
	}

	return res, nil
}

// HashPassword hashes the given password using bcrypt.
func (s *MemberService) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (s *MemberService) ComparePasswords(storedPassword string, inputPassword string) (bool, error) {
	// Compare the hashed password with the user-provided password
	err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(inputPassword))
	if err != nil {
		// If error is not nil, the passwords do not match
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return false, nil // Passwords do not match
		}
		// Handle other potential errors (e.g., hash format issues)
		return false, err
	}
	// If err is nil, the passwords match
	return true, nil
}
