package member

import (
	"context"

	pb "github.com/darkphotonKN/community-builds-microservice/common/api/proto/auth"
)

type Handler struct {
	pb.UnimplementedAuthServiceServer
	service Service
}

type Service interface {
	CreateMember(ctx context.Context, req *pb.CreateMemberRequest) (*pb.Member, error)
	GetMember(ctx context.Context, req *pb.GetMemberRequest) (*pb.Member, error)
	LoginMember(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error)
	UpdateMemberInfo(ctx context.Context, req *pb.UpdateMemberInfoRequest) (*pb.Member, error)
	UpdateMemberPassword(ctx context.Context, req *pb.UpdatePasswordRequest) (*pb.UpdatePasswordResponse, error)
	ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error)
}

func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

// LoginMember implements the LoginMember gRPC method
func (s *Handler) LoginMember(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	return s.service.LoginMember(ctx, req)
}

// GetMember implements the GetMember gRPC method
func (s *Handler) GetMember(ctx context.Context, req *pb.GetMemberRequest) (*pb.Member, error) {
	return s.service.GetMember(ctx, req)
}

// CreateMember implements the CreateMember gRPC method
func (s *Handler) CreateMember(ctx context.Context, req *pb.CreateMemberRequest) (*pb.Member, error) {
	return s.service.CreateMember(ctx, req)
}

// UpdateMemberInfo implements the UpdateMemberInfo gRPC method
func (s *Handler) UpdateMemberInfo(ctx context.Context, req *pb.UpdateMemberInfoRequest) (*pb.Member, error) {
	return s.service.UpdateMemberInfo(ctx, req)
}

// UpdateMemberPassword implements the UpdateMemberPassword gRPC method
func (s *Handler) UpdateMemberPassword(ctx context.Context, req *pb.UpdatePasswordRequest) (*pb.UpdatePasswordResponse, error) {
	return s.service.UpdateMemberPassword(ctx, req)
}

// ValidateToken implements the ValidateToken gRPC method
func (s *Handler) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	return s.service.ValidateToken(ctx, req)
}
