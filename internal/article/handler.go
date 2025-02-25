package article

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	// "github.com/google/uuid"
)

type ArticleHandler struct {
	Service *ArticleService
}

func NewArticleHandler(service *ArticleService) *ArticleHandler {
	return &ArticleHandler{
		Service: service,
	}
}

// --- ADMIN HANDLERS ---
func (h *ArticleHandler) CreateArticleHandler(c *gin.Context) {
	var createArticleReq CreateArticleRequest

	if err := c.ShouldBindJSON(&createArticleReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": fmt.Sprintf("Error when parsing payload as JSON: %s", err)})
		return
	}

	err := h.Service.CreateArticleService(createArticleReq)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": fmt.Sprintf("Error when attempting to create tag: %s", err.Error())})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"statusCode": http.StatusCreated, "message": "Successfully created tag."})
}

func (h *ArticleHandler) GetArticlesHandler(c *gin.Context) {
	articles, err := h.Service.GetArticlesService()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": fmt.Sprintf("Error when attempting to retrieve all articles: %s\n", err.Error())})
		return
	}

	c.JSON(http.StatusOK, gin.H{"statusCode": http.StatusOK, "message": "Successfully retrieved all articles.", "result": articles})
}

func (h *ArticleHandler) UpdateArticlesHandler(c *gin.Context) {
	// tag id to update
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": fmt.Sprintf("Error with id %d, not a valid uuid.", id)})
		return
	}

	// update tag payload
	var updateArticleReq UpdateArticleRequest
	if err := c.ShouldBindJSON(&updateArticleReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": fmt.Sprintf("Error when parsing payload as JSON.")})
		return
	}

	resErr := h.Service.UpdateArticlesService(updateArticleReq)

	if resErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": fmt.Sprintf("Error when attempting to update Article with id: %s\n error: %s\n", id, err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"statusCode": http.StatusOK, "message": "Successfully updated Article."})
}
