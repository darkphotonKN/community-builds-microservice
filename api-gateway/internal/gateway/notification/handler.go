package notification

import (
	"fmt"
	"net/http"
	"strconv"

	pb "github.com/darkphotonKN/community-builds-microservice/common/api/proto/notification"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	client NotificationClient
}

func NewHandler(client NotificationClient) *Handler {
	return &Handler{
		client: client,
	}
}

func (h *Handler) GetNotificationsByMemberIdHandler(c *gin.Context) {
	memberIdStr, exists := c.Get("userIdStr")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"statusCode": http.StatusUnauthorized,
			"message":    "User ID not found in context",
		})
		return
	}

	limitStr := c.DefaultQuery("limit", "20")
	offsetStr := c.DefaultQuery("offset", "0")

	limitInt64, err := strconv.ParseInt(limitStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": http.StatusBadRequest,
			"message":    "Invalid limit parameter",
		})
		return
	}
	limit := int32(limitInt64)

	offsetInt64, err := strconv.ParseInt(offsetStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": http.StatusBadRequest,
			"message":    "Invalid offset parameter",
		})
		return
	}
	offset := int32(offsetInt64)

	req := &pb.GetNotificationsRequest{
		MemberId: memberIdStr.(string),
		Limit:    &limit,
		Offset:   &offset,
	}

	response, err := h.client.GetNotifications(c.Request.Context(), req)

	if err != nil {
		status, ok := status.FromError(err)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"statusCode": http.StatusInternalServerError,
				"message":    "Internal server error",
			})
			return
		}

		httpStatus := http.StatusInternalServerError
		switch status.Code() {
		case codes.NotFound:
			httpStatus = http.StatusNotFound
		case codes.InvalidArgument:
			httpStatus = http.StatusBadRequest
		}

		c.JSON(httpStatus, gin.H{
			"statusCode": httpStatus,
			"message":    status.Message(),
		})
		return
	}

	// default to empty slice if theres no data, data formatting for FE
	if len(response.Data) == 0 {
		response.Data = make([]*pb.Notification, 0)
	}

	notificationRes := gin.H{
		"statusCode": http.StatusOK,
		"message":    "Successfully retrieved notifications",
		"result":     response.Data,
	}

	fmt.Printf("\nnotificationRes before going back to FE: %+v\n\n", notificationRes)
	c.JSON(http.StatusOK, notificationRes)
}

func (h *Handler) ReadNotificationsByMemberIdHandler(c *gin.Context) {
	memberIdStr, exists := c.Get("userIdStr")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"statusCode": http.StatusUnauthorized,
			"message":    "User ID not found in context",
		})
		return
	}

	notificationIdParam, exists := c.Get("notificationId")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"statusCode": http.StatusUnauthorized,
			"message":    "Notification Id param not provided",
		})
		return
	}

	req := &pb.ReadNotificationRequest{
		MemberId:       memberIdStr.(string),
		NotificationId: notificationIdParam.(string),
	}

	response, err := h.client.ReadNotifications(c.Request.Context(), req)

	if err != nil {
		status, ok := status.FromError(err)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"statusCode": http.StatusInternalServerError,
				"message":    "Internal server error",
			})
			return
		}

		httpStatus := http.StatusInternalServerError
		switch status.Code() {
		case codes.NotFound:
			httpStatus = http.StatusNotFound
		case codes.InvalidArgument:
			httpStatus = http.StatusBadRequest
		}

		c.JSON(httpStatus, gin.H{
			"statusCode": httpStatus,
			"message":    status.Message(),
		})
		return
	}

	res := gin.H{
		"statusCode": http.StatusOK,
		"message":    response.Message,
		"result": pb.ReadNotificationResponse{
			Success: response.Success,
			Message: response.Message,
		},
	}

	c.JSON(http.StatusOK, res)
}
