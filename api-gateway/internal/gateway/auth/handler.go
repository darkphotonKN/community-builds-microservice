package auth

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/darkphotonKN/community-builds-microservice/api-gateway/internal/utils/errorutils"
	pb "github.com/darkphotonKN/community-builds-microservice/common/api/proto/auth"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	client AuthClient
}

func NewHandler(client AuthClient) *Handler {
	return &Handler{
		client: client,
	}
}

func (h *Handler) CreateMemberHandler(c *gin.Context) {
	var req pb.CreateMemberRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": "Error parsing payload as JSON"})
		return
	}

	member, err := h.client.CreateMember(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": http.StatusInternalServerError,
			"message":    fmt.Sprintf("Error creating user: %s", err.Error()),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"statusCode": http.StatusCreated,
		"message":    "Successfully created user",
		"result":     member,
	})
}

func (h *Handler) UpdatePasswordMemberHandler(c *gin.Context) {
	var req pb.UpdatePasswordRequest
	userId, _ := c.Get("userId")
	userIdStr := userId.(uuid.UUID).String()

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": "Error parsing payload as JSON"})
		return
	}

	// Set the ID from context
	req.Id = userIdStr

	response, err := h.client.UpdateMemberPassword(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": http.StatusInternalServerError,
			"message":    fmt.Sprintf("Error updating password: %s", err.Error()),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"statusCode": http.StatusOK,
		"message":    response.Message,
		"success":    response.Success,
	})
}

func (h *Handler) UpdateInfoMemberHandler(c *gin.Context) {
	var req pb.UpdateMemberInfoRequest
	userId, _ := c.Get("userId")
	userIdStr := userId.(uuid.UUID).String()

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": "Error parsing payload as JSON"})
		return
	}

	// Set the ID from context
	req.Id = userIdStr

	member, err := h.client.UpdateMemberInfo(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": http.StatusInternalServerError,
			"message":    fmt.Sprintf("Error updating member info: %s", err.Error()),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"statusCode": http.StatusOK,
		"message":    "Successfully updated member info",
		"result":     member,
	})
}

func (h *Handler) LoginMemberHandler(c *gin.Context) {
	var req pb.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": fmt.Sprintf("Error parsing payload as JSON: %s", err)})
		return
	}

	response, err := h.client.LoginMember(c.Request.Context(), &req)

	if err != nil {
		// Check for unauthorized error specifically
		if errors.Is(err, errorutils.ErrUnauthorized) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"statusCode": http.StatusUnauthorized,
				"message":    errorutils.ErrUnauthorized.Error(),
			})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": http.StatusBadRequest,
			"message":    fmt.Sprintf("Error logging in: %s", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"statusCode": http.StatusOK,
		"message":    "Successfully logged in",
		"result":     response,
	})
}

func (h *Handler) GetMemberByIdHandler(c *gin.Context) {
	// Get ID from path parameter
	idParam := c.Param("id")

	// Create the request
	req := &pb.GetMemberRequest{
		Id: idParam,
	}

	// Call the service
	member, err := h.client.GetMember(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": http.StatusBadRequest,
			"message":    fmt.Sprintf("Error retrieving member: %s", err.Error()),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"statusCode": http.StatusOK,
		"message":    "Successfully retrieved member",
		"result":     member,
	})
}

func (h *Handler) ValidateTokenHandler(c *gin.Context) {
	var req pb.ValidateTokenRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": "Error parsing payload as JSON"})
		return
	}

	response, err := h.client.ValidateToken(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": http.StatusBadRequest,
			"message":    fmt.Sprintf("Error validating token: %s", err.Error()),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"statusCode": http.StatusOK,
		"valid":      response.Valid,
		"memberId":   response.MemberId,
	})
}

