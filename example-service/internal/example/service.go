package example

import (
	"errors"
)

type service struct {
	repo Repository
}

type Repository interface {
	Create(example *ExampleCreate) (*Example, error)
	GetByID(id string) (*Example, error)
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateExample(example *ExampleCreate) (*Example, error) {
	if example.Name == "" {
		return nil, errors.New("name is required")
	}

	return s.repo.Create(example)
}

func (s *service) GetExample(id string) (*Example, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}

	return s.repo.GetByID(id)
}



