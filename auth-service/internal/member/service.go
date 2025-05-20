package member

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	pb "github.com/darkphotonKN/community-builds-microservice/common/api/proto/auth"
	"github.com/darkphotonKN/community-builds-microservice/api-gateway/internal/auth"
	"github.com/darkphotonKN/community-builds-microservice/api-gateway/internal/utils/errorutils"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type MemberService struct {
	Repo *MemberRepository
}

func NewMemberService(repo *MemberRepository) *MemberService {
	return &MemberService{
		Repo: repo,
	}
}

// Convert internal Member to protobuf Member
func memberToProto(m *Member) *pb.Member {
	if m == nil {
		return nil
	}
	
	return &pb.Member{
		Id:            m.ID.String(),
		Name:          m.Name,
		Email:         m.Email,
		Status:        int32(stringToInt(m.Status)),
		AverageRating: float32(m.AverageRating),
		CreatedAt:     timestamppb.New(m.CreatedAt),
		UpdatedAt:     timestamppb.New(m.UpdatedAt),
	}
}

// Utility function to convert string to int
func stringToInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return i
}

// GetMember implements the gRPC GetMember method
func (s *MemberService) GetMember(ctx context.Context, req *pb.GetMemberRequest) (*pb.Member, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, fmt.Errorf("invalid UUID: %w", err)
	}
	
	member, err := s.Repo.GetById(id)
	if err != nil {
		return nil, err
	}
	
	return memberToProto(member), nil
}

// CreateMember implements the gRPC CreateMember method
func (s *MemberService) CreateMember(ctx context.Context, req *pb.CreateMemberRequest) (*pb.Member, error) {
	// Hash the password
	hashedPw, err := s.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("error hashing password: %w", err)
	}
	
	// Create the member
	memberId, err := s.Repo.Create(req.Name, req.Email, hashedPw)
	if err != nil {
		return nil, err
	}
	
	// Get the created member
	member, err := s.Repo.GetById(memberId)
	if err != nil {
		return nil, err
	}
	
	return memberToProto(member), nil
}

// LoginMember implements the gRPC LoginMember method
func (s *MemberService) LoginMember(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	// Verify credentials
	member, err := s.Repo.GetMemberByEmail(req.Email)
	if err != nil {
		return nil, fmt.Errorf("could not find member with provided email: %w", err)
	}
	
	// Compare the stored hashed password with the provided password
	if err = bcrypt.CompareHashAndPassword([]byte(member.Password), []byte(req.Password)); err != nil {
		return nil, errorutils.ErrUnauthorized
	}
	
	// Generate tokens
	accessExpiryTime := time.Minute * 60
	refreshExpiryTime := time.Hour * 24 * 7
	
	// Generate access token
	accessToken, err := auth.GenerateJWT(*member, auth.Access, accessExpiryTime)
	if err != nil {
		return nil, fmt.Errorf("error generating access token: %w", err)
	}
	
	// Generate refresh token
	refreshToken, err := auth.GenerateJWT(*member, auth.Refresh, refreshExpiryTime)
	if err != nil {
		return nil, fmt.Errorf("error generating refresh token: %w", err)
	}
	
	// Create the response
	return &pb.LoginResponse{
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		AccessExpiresIn:  int32(accessExpiryTime.Seconds()),
		RefreshExpiresIn: int32(refreshExpiryTime.Seconds()),
		MemberInfo:       memberToProto(member),
	}, nil
}

// UpdateMemberInfo implements the gRPC UpdateMemberInfo method
func (s *MemberService) UpdateMemberInfo(ctx context.Context, req *pb.UpdateMemberInfoRequest) (*pb.Member, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, fmt.Errorf("invalid UUID: %w", err)
	}
	
	// Update member info
	err = s.Repo.UpdateMemberInfo(id, req.Name, req.Status)
	if err != nil {
		return nil, err
	}
	
	// Get the updated member
	member, err := s.Repo.GetById(id)
	if err != nil {
		return nil, err
	}
	
	return memberToProto(member), nil
}

// UpdateMemberPassword implements the gRPC UpdateMemberPassword method
func (s *MemberService) UpdateMemberPassword(ctx context.Context, req *pb.UpdatePasswordRequest) (*pb.UpdatePasswordResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, fmt.Errorf("invalid UUID: %w", err)
	}
	
	// Get the member with password
	member, err := s.Repo.GetByIdWithPassword(id)
	if err != nil {
		return nil, err
	}
	
	// Check if new passwords match
	if req.NewPassword != req.RepeatNewPassword {
		return &pb.UpdatePasswordResponse{
			Success: false,
			Message: "New passwords do not match",
		}, errors.New("new passwords do not match")
	}
	
	// Verify current password
	isSame, err := s.ComparePasswords(member.Password, req.CurrentPassword)
	if !isSame || err != nil {
		return &pb.UpdatePasswordResponse{
			Success: false,
			Message: "Current password is incorrect",
		}, errors.New("current password is incorrect")
	}
	
	// Hash the new password
	hashedPw, err := s.HashPassword(req.NewPassword)
	if err != nil {
		return &pb.UpdatePasswordResponse{
			Success: false,
			Message: "Error hashing password",
		}, fmt.Errorf("error hashing password: %w", err)
	}
	
	// Update the password in the database
	params := MemberUpdatePasswordParams{
		ID:       id,
		Password: hashedPw,
	}
	
	err = s.Repo.UpdatePassword(params)
	if err != nil {
		return &pb.UpdatePasswordResponse{
			Success: false,
			Message: "Error updating password",
		}, err
	}
	
	return &pb.UpdatePasswordResponse{
		Success: true,
		Message: "Password updated successfully",
	}, nil
}

// ValidateToken implements the gRPC ValidateToken method
func (s *MemberService) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	// Validate the token using auth package
	claims, err := auth.ValidateJWT(req.Token)
	if err != nil {
		return &pb.ValidateTokenResponse{
			Valid:    false,
			MemberId: "",
		}, err
	}
	
	return &pb.ValidateTokenResponse{
		Valid:    true,
		MemberId: claims.ID.String(),
	}, nil
}

// Helper functions

// HashPassword hashes the given password using bcrypt.
func (s *MemberService) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// ComparePasswords compares a hashed password with a plain text password.
func (s *MemberService) ComparePasswords(storedPassword string, inputPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(inputPassword))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return false, nil // Passwords do not match
		}
		return false, err // Other error
	}
	return true, nil // Passwords match
}

// CreateDefaultMembers creates default members for setup purposes.
func (s *MemberService) CreateDefaultMembers(members []CreateDefaultMember) error {
	var hashedPwMembers []CreateDefaultMember
	
	// Update members passwords with hash
	for _, member := range members {
		hashedPw, err := s.HashPassword(member.Password)
		if err != nil {
			return err
		}
		member.Password = hashedPw
		hashedPwMembers = append(hashedPwMembers, member)
	}
	
	return s.Repo.CreateDefaultMembers(hashedPwMembers)
}