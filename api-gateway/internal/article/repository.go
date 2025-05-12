package article

import (
	"errors"
	"fmt"

	"github.com/darkphotonKN/community-builds-microservice/api-gateway/internal/models"
	"github.com/darkphotonKN/community-builds-microservice/api-gateway/internal/utils/errorutils"
	"github.com/jmoiron/sqlx"
)

type ArticleRepository struct {
	DB *sqlx.DB
}

func NewArticleRepository(db *sqlx.DB) *ArticleRepository {
	return &ArticleRepository{
		DB: db,
	}
}

func (r *ArticleRepository) CreateArticle(createArticleReq CreateArticleRequest) error {
	query := `
		INSERT INTO articles(content)
		VALUES(:content)
	`

	_, err := r.DB.NamedExec(query, createArticleReq)

	fmt.Print("Error when creating article:", err)

	if err != nil {
		return errorutils.AnalyzeDBErr(err)
	}

	return nil
}

func (r *ArticleRepository) GetArticles() (*[]models.Article, error) {
	var articles []models.Article

	query := `SELECT * FROM articles`

	err := r.DB.Select(&articles, query)

	if err != nil {
		return nil, errorutils.AnalyzeDBErr(err)
	}

	return &articles, nil
}

func (r *ArticleRepository) UpdateArticle(payload UpdateArticleRequest) error {

	query := `UPDATE articles SET name = :name WHERE id = :id`

	result, err := r.DB.NamedExec(query, payload)

	if err != nil {
		return errorutils.AnalyzeDBErr(err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("no rows updated")
	}
	return nil
}

func (r *ArticleRepository) BatchCreateArticles(articles []models.Article) error {
	query := `
	INSERT INTO articles(id, name)
	VALUES(:id, :name)
	ON CONFLICT DO NOTHING
	`
	_, err := r.DB.NamedExec(query, articles)

	if err != nil {
		return errorutils.AnalyzeDBErr(err)
	}

	return nil
}
