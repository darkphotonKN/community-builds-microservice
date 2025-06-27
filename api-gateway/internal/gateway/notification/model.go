package notification

import (
	"context"

	pb "github.com/darkphotonKN/community-builds-microservice/common/api/proto/notification"
)

type NotificationClient interface {
	GetNotifications(ctx context.Context, req *pb.GetNotificationsRequest) (*pb.GetNotificationsResponse, error)

	ReadNotifications(ctx context.Context, req *pb.ReadNotificationRequest) (*pb.ReadNotificationResponse, error)
}

