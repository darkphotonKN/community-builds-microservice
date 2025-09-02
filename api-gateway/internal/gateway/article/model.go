package article

import (
	"context"

	"github.com/google/uuid"

	pb "github.com/darkphotonKN/community-builds-microservice/common/api/proto/article"
)

type CreateArticleRequest struct {
	Content string `json:"content" binding:"required" db:"content"`
}

type UpdateArticleRequest struct {
	// ID      uuid.UUID `db:"id" json:"id"`
	Content string `json:"content" binding:"required" db:"content"`
}

type UpdateArticleParams struct {
	ID      uuid.UUID `db:"id" json:"id"`
	Content string    `db:"content" json:"content"`
}

type ArticleClient interface {
	CreateArticle(ctx context.Context, req *pb.CreateArticleRequest) (*pb.CreateArticleResponse, error)
	GetArticles(ctx context.Context, req *pb.GetArticlesRequest) (*pb.GetArticlesResponse, error)
	UpdateArticle(ctx context.Context, req *pb.UpdateArticleRequest) (*pb.UpdateArticleResponse, error)
}
