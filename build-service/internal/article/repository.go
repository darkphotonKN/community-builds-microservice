package article

import (
	"errors"
	"fmt"

	"github.com/darkphotonKN/community-builds-microservice/common/constants/models"
	commonhelpers "github.com/darkphotonKN/community-builds-microservice/common/utils"
	"github.com/jmoiron/sqlx"
)

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateArticle(createArticleReq CreateArticle) error {
	fmt.Println("Creating article:", createArticleReq.Content)
	query := `
		INSERT INTO articles(content)
		VALUES(:content)
	`

	_, err := r.db.NamedExec(query, createArticleReq)

	fmt.Print("Error when creating article:", err)

	if err != nil {
		return commonhelpers.AnalyzeDBErr(err)
	}

	return nil
}

func (r *repository) GetArticles() (*[]models.Article, error) {
	var articles []models.Article

	query := `SELECT * FROM articles`

	err := r.db.Select(&articles, query)

	if err != nil {
		return nil, commonhelpers.AnalyzeDBErr(err)
	}

	return &articles, nil
}

func (r *repository) UpdateArticle(payload UpdateArticle) error {

	query := `UPDATE articles SET name = :name WHERE id = :id`

	result, err := r.db.NamedExec(query, payload)

	if err != nil {
		return commonhelpers.AnalyzeDBErr(err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("no rows updated")
	}
	return nil
}

func (r *repository) BatchCreateArticles(articles []models.Article) error {
	query := `
	INSERT INTO articles(id, name)
	VALUES(:id, :name)
	ON CONFLICT DO NOTHING
	`
	_, err := r.db.NamedExec(query, articles)

	if err != nil {
		return commonhelpers.AnalyzeDBErr(err)
	}

	return nil
}
