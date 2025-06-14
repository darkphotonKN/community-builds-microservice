package notification

import (
	"context"

	pb "github.com/darkphotonKN/community-builds-microservice/common/api/proto/notification"
)

type Handler struct {
	pb.UnimplementedNotificationServiceServer
	service ServiceNot
}

func (h *Handler) getNotifications(ctx context.Context, request *pb.GetNotificationsRequest) (*pb.GetNotificationsResponse, error) {
	return h.service.getAllByMemberId(ctx, request)
}

func NewHandler(service ServiceNot) *Handler {
	return &Handler{service: service}
}
