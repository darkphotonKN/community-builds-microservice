package article

import (
	"github.com/darkphotonKN/community-builds/internal/models"
)

type ArticleService struct {
	Repo *ArticleRepository
}

func NewArticleService(repo *ArticleRepository) *ArticleService {
	return &ArticleService{
		Repo: repo,
	}
}

func (s *ArticleService) CreateArticleService(createArticleReq CreateArticleRequest) error {
	return s.Repo.CreateArticle(createArticleReq)
}

func (s *ArticleService) UpdateArticlesService(updateArticleReq UpdateArticleRequest) error {
	return s.Repo.UpdateArticle(updateArticleReq)
}

func (s *ArticleService) GetArticlesService() (*[]models.Article, error) {
	return s.Repo.GetArticles()
}

func (s *ArticleService) CreateDefaultArticles(articles []models.Article) error {
	return s.Repo.BatchCreateArticles(articles)
}
