package notification

import (
	"encoding/json"
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
func debugSlice(name string, slice []*pb.Notification) {
	fmt.Printf("=== %s DEBUG ===\n", name)
	fmt.Printf("  Value: %+v\n", slice)
	fmt.Printf("  Is nil: %t\n", slice == nil)
	fmt.Printf("  Length: %d\n", len(slice))
	fmt.Printf("  Capacity: %d\n", cap(slice))
	fmt.Printf("  Type: %T\n", slice)

	// JSON representation
	jsonBytes, _ := json.Marshal(slice)
	fmt.Printf("  JSON: %s\n", string(jsonBytes))
	fmt.Printf("================\n\n")
}

func (h *Handler) GetNotificationsByMemberIdHandler(c *gin.Context) {
	userIdStr, exists := c.Get("userIdStr")
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
		MemberId: userIdStr.(string),
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

	debugSlice("BEFORE response.Data", response.Data)

	if len(response.Data) == 0 {
		response.Data = make([]*pb.Notification, 0)

		debugSlice("AFTER response.Data", response.Data)
	}

	notificationRes := gin.H{
		"statusCode": http.StatusOK,
		"message":    "Successfully retrieved notifications",
		"result":     response.Data,
	}

	fmt.Printf("\nnotificationRes before going back to FE: %+v\n\n", notificationRes)
	c.JSON(http.StatusOK, notificationRes)
}
