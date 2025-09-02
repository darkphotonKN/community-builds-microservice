package article

import (
	"context"

	pb "github.com/darkphotonKN/community-builds-microservice/common/api/proto/article"
)

type Handler struct {
	service Service
	pb.UnimplementedArticleServiceServer
}

type Service interface {
	CreateArticle(ctx context.Context, req *pb.CreateArticleRequest) (*pb.CreateArticleResponse, error)
	UpdateArticle(ctx context.Context, req *pb.UpdateArticleRequest) (*pb.UpdateArticleResponse, error)
	GetArticles(ctx context.Context, req *pb.GetArticlesRequest) (*pb.GetArticlesResponse, error)
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateArticle(ctx context.Context, req *pb.CreateArticleRequest) (*pb.CreateArticleResponse, error) {
	return h.service.CreateArticle(ctx, req)
}

func (h *Handler) UpdateArticle(ctx context.Context, req *pb.UpdateArticleRequest) (*pb.UpdateArticleResponse, error) {
	return h.service.UpdateArticle(ctx, req)
}

func (h *Handler) GetArticles(ctx context.Context, req *pb.GetArticlesRequest) (*pb.GetArticlesResponse, error) {
	return h.service.GetArticles(ctx, req)
}
