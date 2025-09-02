package article

import (
	"fmt"
	"net/http"

	pb "github.com/darkphotonKN/community-builds-microservice/common/api/proto/article"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	// "github.com/google/uuid"
)

type ArticleHandler struct {
	Client ArticleClient
}

func NewHandler(client ArticleClient) *ArticleHandler {
	return &ArticleHandler{
		Client: client,
	}
}

// --- ADMIN HANDLERS ---
func (h *ArticleHandler) CreateArticleHandler(c *gin.Context) {
	userIdStr, ok := c.Get("userIdStr")
	fmt.Printf("User ID: %s\n", userIdStr)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"statusCode": http.StatusUnauthorized, "message": "Unauthorized"})
		return
	}
	var createArticleReq CreateArticleRequest

	if err := c.ShouldBindJSON(&createArticleReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": fmt.Sprintf("Error when parsing payload as JSON: %s", err)})
		return
	}

	// Convert REST request to gRPC request
	grpcReq := &pb.CreateArticleRequest{
		// MemberId: userIdStr.(string),
		Content: createArticleReq.Content,
	}

	_, err := h.Client.CreateArticle(c.Request.Context(), grpcReq)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": fmt.Sprintf("Error when attempting to create article: %s", err.Error())})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"statusCode": http.StatusCreated, "message": "Successfully created article."})
}

func (h *ArticleHandler) GetArticlesHandler(c *gin.Context) {
	userIdStr, ok := c.Get("userIdStr")
	fmt.Printf("User ID: %s\n", userIdStr)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"statusCode": http.StatusUnauthorized, "message": "Unauthorized"})
		return
	}
	// Convert REST request to gRPC request
	grpcReq := &pb.GetArticlesRequest{}
	articles, err := h.Client.GetArticles(c.Request.Context(), grpcReq)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": fmt.Sprintf("Error when attempting to retrieve all articles: %s\n", err.Error())})
		return
	}

	c.JSON(http.StatusOK, gin.H{"statusCode": http.StatusOK, "message": "Successfully retrieved all articles.", "result": articles})
}

func (h *ArticleHandler) UpdateArticlesHandler(c *gin.Context) {
	userIdStr, ok := c.Get("userIdStr")
	fmt.Printf("User ID: %s\n", userIdStr)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"statusCode": http.StatusUnauthorized, "message": "Unauthorized"})
		return
	}
	// article id to update
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": fmt.Sprintf("Error with id %d, not a valid uuid.", id)})
		return
	}

	// update article payload
	var updateArticleReq UpdateArticleRequest
	if err := c.ShouldBindJSON(&updateArticleReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": fmt.Sprintf("Error when parsing payload as JSON.")})
		return
	}
	// Convert REST request to gRPC request
	grpcReq := &pb.UpdateArticleRequest{
		Id:      id.String(),
		Content: updateArticleReq.Content,
	}
	_, resErr := h.Client.UpdateArticle(c.Request.Context(), grpcReq)

	if resErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": fmt.Sprintf("Error when attempting to update Article with id: %s\n error: %s\n", id, err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"statusCode": http.StatusOK, "message": "Successfully updated Article."})
}
