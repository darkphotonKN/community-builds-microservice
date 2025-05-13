package example

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

type Service interface {
	CreateExample(example *ExampleCreate) (*Example, error)
	GetExample(id string) (*Example, error)
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}


func (h *Handler) CreateExample(c *gin.Context) {
	var example ExampleCreate
	if err := c.ShouldBindJSON(&example); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.service.CreateExample(&example)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, result)
}

func (h *Handler) GetExample(c *gin.Context) {
	id := c.Param("id")
	result, err := h.service.GetExample(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "example not found"})
		return
	}

	c.JSON(http.StatusOK, result)
}
