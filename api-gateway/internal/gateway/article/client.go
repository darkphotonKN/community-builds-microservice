package article

import (
	"context"
	"fmt"

	pb "github.com/darkphotonKN/community-builds-microservice/common/api/proto/article"
	"github.com/darkphotonKN/community-builds-microservice/common/discovery"
)

const (
	serviceName = "build-service"
)

type Client struct {
	registry discovery.Registry
}

func NewClient(registry discovery.Registry) ArticleClient {
	return &Client{
		registry: registry,
	}
}

func (c *Client) CreateArticle(ctx context.Context, req *pb.CreateArticleRequest) (*pb.CreateArticleResponse, error) {

	// connection instance created through service discovery first
	// searches for the service registered as "orders"
	conn, err := discovery.ServiceConnection(ctx, serviceName, c.registry)

	if err != nil {
		return nil, fmt.Errorf("failed to connect to item service: %w", err)
	}
	defer conn.Close()

	client := pb.NewArticleServiceClient(conn)

	// create client to interface with through service discovery connection
	item, err := client.CreateArticle(ctx, req)

	fmt.Printf("Creating article %+v through gateway after service discovery\n", item)
	if err != nil {
		return nil, fmt.Errorf("failed to create article: %w", err)
	}
	return item, nil
}

func (c *Client) GetArticles(ctx context.Context, req *pb.GetArticlesRequest) (*pb.GetArticlesResponse, error) {

	// connection instance created through service discovery first
	// searches for the service registered as "orders"
	conn, err := discovery.ServiceConnection(ctx, serviceName, c.registry)

	if err != nil {
		return nil, fmt.Errorf("failed to connect to item service: %w", err)
	}
	defer conn.Close()

	client := pb.NewArticleServiceClient(conn)

	// create client to interface with through service discovery connection
	item, err := client.GetArticles(ctx, req)

	fmt.Printf("Getting article %+v through gateway after service discovery\n", item)
	if err != nil {
		return nil, fmt.Errorf("failed to Get article: %w", err)
	}
	return item, nil
}

func (c *Client) UpdateArticle(ctx context.Context, req *pb.UpdateArticleRequest) (*pb.UpdateArticleResponse, error) {

	// connection instance created through service discovery first
	// searches for the service registered as "orders"
	conn, err := discovery.ServiceConnection(ctx, serviceName, c.registry)

	if err != nil {
		return nil, fmt.Errorf("failed to connect to build service: %w", err)
	}
	defer conn.Close()

	client := pb.NewArticleServiceClient(conn)

	// create client to interface with through service discovery connection
	article, err := client.UpdateArticle(ctx, req)

	fmt.Printf("Creating article %+v through gateway after service discovery\n", article)
	if err != nil {
		return nil, fmt.Errorf("failed to update article: %w", err)
	}
	return article, nil
}
