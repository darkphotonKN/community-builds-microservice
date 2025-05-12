package class

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ClassHandler struct {
	Service *ClassService
}

func NewClassHandler(service *ClassService) *ClassHandler {
	return &ClassHandler{
		Service: service,
	}
}

/**
* Retrievs all classes and ascendancies.
**/
func (h *ClassHandler) GetClassesAndAscendanciesHandler(c *gin.Context) {
	response, err := h.Service.GetClassesAndAscendancies()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusCode": http.StatusInternalServerError, "message": fmt.Sprintf("Error occcured when attemptign to retrieve build: %s", err)})
	}

	c.JSON(http.StatusOK, gin.H{"message": "Retrieved all classes and ascendancies.", "result": response})
}
