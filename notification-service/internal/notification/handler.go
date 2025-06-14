package notification

import (
	"context"

	pb "github.com/darkphotonKN/community-builds-microservice/common/api/proto/notification"
)

type Handler struct {
	pb.UnimplementedNotificationServiceServer
	service QueryHandlerService
}

func (h *Handler) getNotifications(ctx context.Context, request *pb.GetNotificationsRequest) (*pb.GetNotificationsResponse, error) {
	return h.service.GetAllByMemberId(ctx, request)
}

func NewHandler(service QueryHandlerService) *Handler {
	return &Handler{service: service}
}
