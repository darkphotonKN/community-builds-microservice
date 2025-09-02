package article

import (
	"context"
	"fmt"

	pb "github.com/darkphotonKN/community-builds-microservice/common/api/proto/article"
	"github.com/darkphotonKN/community-builds-microservice/common/constants/models"
)

/**
* Classes an Ascendancies Services
**/

type service struct {
	repo Repository
}

type Repository interface {
	CreateArticle(article CreateArticle) error
	UpdateArticle(article UpdateArticle) error
	GetArticles() (*[]models.Article, error)
	BatchCreateArticles(articles []models.Article) error
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

/**
* Create Article
**/
func (s *service) CreateArticle(ctx context.Context, req *pb.CreateArticleRequest) (*pb.CreateArticleResponse, error) {
	fmt.Println("Received CreateArticle request:", req.Content)
	payload := CreateArticle{
		Content: req.Content,
	}
	err := s.repo.CreateArticle(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to create default classes and ascendancies: %w", err)
	}

	return nil, nil
}

/**
* Update Article
**/
func (s *service) UpdateArticle(ctx context.Context, req *pb.UpdateArticleRequest) (*pb.UpdateArticleResponse, error) {
	payload := UpdateArticle{}
	err := s.repo.UpdateArticle(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to update article: %w", err)
	}

	return nil, nil
}

/**
* Get Articles
**/
func (s *service) GetArticles(ctx context.Context, req *pb.GetArticlesRequest) (*pb.GetArticlesResponse, error) {
	// payload := []GetArticles{

	// }
	// err := s.repo.GetArticles(payload)
	// if err != nil {
	// 	return fmt.Errorf("failed to get articles: %w", err)
	// }

	return nil, nil
}
