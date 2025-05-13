package example

import (
	"context"

	pb "github.com/darkphotonKN/community-builds-microservice/common/api/proto/example"
	"github.com/google/uuid"
)

type service struct {
	repo Repository
}

type Repository interface {
	Create(example *ExampleCreate) (*Example, error)
	GetByID(id uuid.UUID) (*Example, error)
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateExample(ctx context.Context, req *pb.CreateExampleRequest) (*pb.Example, error) {
	// format to fit model for db tags
	createExample := &ExampleCreate{
		Name: req.Name,
	}
	example, err := s.repo.Create(createExample)

	if err != nil {
		return nil, err
	}

	return &pb.Example{
		Id:   example.ID,
		Name: example.Name,
	}, nil
}

func (s *service) GetExample(ctx context.Context, id uuid.UUID) (*pb.Example, error) {
	example, err := s.repo.GetByID(id)

	if err != nil {
		return nil, err
	}

	// format to fit grpc structure
	return &pb.Example{
		Id:   example.ID,
		Name: example.Name,
	}, nil
}
