package notification

import (
	"context"

	pb "github.com/darkphotonKN/community-builds-microservice/common/api/proto/notification"
)

type Creator interface {
	Create(notification *MemberCreatedNotification) (*Notification, error)
	GetNotificationTemplate(notificationType NotificationType) (*NotificationTemplate, error)
	CreateItem(notification *CreateNotification) (*Notification, error)
}

type Reader interface {
	GetAllByMemberId(ctx context.Context, request *pb.GetNotificationsRequest) (*pb.GetNotificationsResponse, error)
}

// for consumer
type EventConsumerService interface {
	Creator
}

// for handler
type QueryHandlerService interface {
	Reader
}

// core full service interface
type Service interface {
	Creator
	Reader
}
