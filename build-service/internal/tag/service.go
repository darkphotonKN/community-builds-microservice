package tag

import (
	"context"

	pb "github.com/darkphotonKN/community-builds-microservice/common/api/proto/tag"
	"github.com/darkphotonKN/community-builds-microservice/common/constants/models"
	"github.com/google/uuid"
)

type service struct {
	repo Repository
}

type Repository interface {
	CreateTag(createBuildRequest CreateTagRequest) error
	GetTags() (*[]models.Tag, error)
	UpdateTag(updateTagRequest UpdateTagRequest) error
	GetBuildTagsForMemberById(memberId uuid.UUID) (*[]models.Tag, error)
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}
func (s *service) CreateTag(ctx context.Context, req *pb.CreateTagRequest) (*pb.CreateTagResponse, error) {

	createTagReq := CreateTagRequest{
		Name: req.Name,
	}
	s.repo.CreateTag(createTagReq)
	return &pb.CreateTagResponse{}, nil
}

func (s *service) UpdateTag(ctx context.Context, req *pb.UpdateTagRequest) (*pb.UpdateTagResponse, error) {
	memberId, err := uuid.Parse(req.MemberId)
	if err != nil {
		return nil, err
	}

	updateTagReq := UpdateTagRequest{
		Id:   memberId,
		Name: req.Name,
	}
	s.repo.UpdateTag(updateTagReq)
	return &pb.UpdateTagResponse{}, nil
}

func (s *service) GetTags(ctx context.Context, req *pb.GetTagsRequest) (*pb.GetTagsResponse, error) {
	tags, err := s.repo.GetTags()

	pbTags := make([]*pb.Tag, 0, len(*tags))
	for _, tag := range *tags {
		pbTags = append(pbTags, &pb.Tag{
			Id:        tag.ID.String(),
			Name:      tag.Name,
			CreatedAt: tag.CreatedAt.String(),
			UpdatedAt: tag.UpdatedAt.String(),
		})
	}
	if err != nil {
		return nil, err
	}
	return &pb.GetTagsResponse{Tags: pbTags}, nil
}

// not grpc
func (s *service) GetBuildTagsForMemberById(memberId uuid.UUID) (*[]models.Tag, error) {
	tags, err := s.repo.GetBuildTagsForMemberById(memberId)
	if err != nil {
		return nil, err
	}
	return tags, nil
}

// func (s *service) CreateDefaultTags(tags []models.Tag) error {
// 	return s.repo.BatchCreateTags(tags)
// }
