package auth

import (
	"fmt"
	"net/http"

	pb "github.com/darkphotonKN/community-builds-microservice/common/api/proto/auth"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
		status, ok := status.FromError(err)

		if !ok {
			// not a gRPC status error
			c.JSON(http.StatusInternalServerError, gin.H{
				"statusCode": http.StatusInternalServerError,
				"message":    "Internal server error",
			})

			return
		}

		// map grpc error codes to http codes
		httpStatus := http.StatusInternalServerError
		switch status.Code() {
		case codes.InvalidArgument:
			httpStatus = http.StatusBadRequest
		case codes.Unauthenticated:
			httpStatus = http.StatusUnauthorized
		case codes.NotFound:
			httpStatus = http.StatusNotFound
		}

		c.JSON(httpStatus, gin.H{
			"statusCode": httpStatus,
			"message":    status.Message(),
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
	userIdStr, _ := c.Get("userIdStr")

	// convert to string to match protobuf type

	// Create the request
	req := &pb.GetMemberRequest{
		Id: userIdStr.(string),
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
