package notification

import (
	"context"

	pb "github.com/darkphotonKN/community-builds-microservice/common/api/proto/notification"
)

type Handler struct {
	service QueryHandlerService
	pb.UnimplementedNotificationServiceServer
}

func NewHandler(service QueryHandlerService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetNotifications(ctx context.Context, request *pb.GetNotificationsRequest) (*pb.GetNotificationsResponse, error) {
	return h.service.GetAllByMemberId(ctx, request)
}

func (h *Handler) ReadNotification(ctx context.Context, request *pb.ReadNotificationRequest) error {
	return h.service.ReadNotification(ctx, request)
}
