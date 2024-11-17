package build

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type BuildHandler struct {
	Service *BuildService
}

func NewBuildHandler(service *BuildService) *BuildHandler {
	return &BuildHandler{
		Service: service,
	}
}

func (h *BuildHandler) CreateBuildHandler(c *gin.Context) {
	memberId, _ := c.Get("userId")
	var createBuildReq CreateBuildRequest

	if err := c.ShouldBindJSON(&createBuildReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": fmt.Sprintf("Error when parsing payload as JSON: %s", err)})
		return
	}

	h.Service.CreateBuildService(memberId.(uuid.UUID), createBuildReq)

	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": fmt.Sprintf("Error when attempting to create a build %s", err.Error())})
	// 	return
	// }

	c.JSON(http.StatusCreated, gin.H{"statusCode": http.StatusCreated, "message": "Successfully created the build."})
}
