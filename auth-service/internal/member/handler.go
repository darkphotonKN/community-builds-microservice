package member

import (
	"context"

	pb "github.com/darkphotonKN/community-builds-microservice/common/api/proto/auth"
)

type Handler struct {
	pb.UnimplementedAuthServiceServer
	Service *MemberService
}

func NewHandler(service *MemberService) *Handler {
	return &Handler{
		Service: service,
	}
}

// LoginMember implements the LoginMember gRPC method
func (s *Handler) LoginMember(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	return s.Service.LoginMember(ctx, req)
}

// GetMember implements the GetMember gRPC method
func (s *Handler) GetMember(ctx context.Context, req *pb.GetMemberRequest) (*pb.Member, error) {
	return s.Service.GetMember(ctx, req)
}

// CreateMember implements the CreateMember gRPC method
func (s *Handler) CreateMember(ctx context.Context, req *pb.CreateMemberRequest) (*pb.Member, error) {
	return s.Service.CreateMember(ctx, req)
}

// UpdateMemberInfo implements the UpdateMemberInfo gRPC method
func (s *Handler) UpdateMemberInfo(ctx context.Context, req *pb.UpdateMemberInfoRequest) (*pb.Member, error) {
	return s.Service.UpdateMemberInfo(ctx, req)
}

// UpdateMemberPassword implements the UpdateMemberPassword gRPC method
func (s *Handler) UpdateMemberPassword(ctx context.Context, req *pb.UpdatePasswordRequest) (*pb.UpdatePasswordResponse, error) {
	return s.Service.UpdateMemberPassword(ctx, req)
}

// ValidateToken implements the ValidateToken gRPC method
func (s *Handler) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	return s.Service.ValidateToken(ctx, req)
}
