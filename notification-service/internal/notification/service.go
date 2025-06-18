package notification

import (
	"context"
	"fmt"

	pb "github.com/darkphotonKN/community-builds-microservice/common/api/proto/notification"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type service struct {
	repo      Repository
	publishCh *amqp.Channel
}

type Repository interface {
	Create(notification *CreateNotification) (*Notification, error)
	GetAll(ctx context.Context, request *QueryNotifications) ([]Notification, error)
}

func NewService(repo Repository, ch *amqp.Channel) *service {
	return &service{repo: repo, publishCh: ch}
}

func (s *service) Create(memberCreated *MemberCreatedNotification) (*Notification, error) {
	// validation and error handling TODO: missing fields
	if memberCreated.Title == "" {
		fmt.Println("Title is required for creating a new notification.")
		return nil, status.Errorf(codes.InvalidArgument, "Name field is required")
	}

	// validate id is a legit uuid
	id, err := uuid.Parse(memberCreated.MemberID)

	if err != nil {
		fmt.Println("Error occured when parsing uuid:", err)
		return nil, err
	}

	// map it to notifications table entity
	createNotification := &CreateNotification{
		MemberID: id,
		Type:     memberCreated.Type,
		Title:    "welcome_message",
		Message:  memberCreated.Message,
		SourceID: memberCreated.SourceID,
	}

	newNotification, err := s.repo.Create(createNotification)
	if err != nil {
		fmt.Println("Error occured when creating new notification:", err)
		return nil, err
	}

	fmt.Println("notification was created:", newNotification)

	return newNotification, nil
}

func (s *service) GetAllByMemberId(ctx context.Context, request *pb.GetNotificationsRequest) (*pb.GetNotificationsResponse, error) {
	// parse uuid
	memberUUID, err := uuid.Parse(request.MemberId)
	if err != nil {
		return nil, err
	}

	// convert to DB appropriate format
	query := &QueryNotifications{
		MemberID: memberUUID,
	}

	// validate and default limit and offset

	if request.Limit == nil {
		defaultQueryLimit := int32(10)
		query.Limit = &defaultQueryLimit
	} else {
		query.Limit = query.Limit
	}

	if request.Offset == nil {
		defaultQueryLimit := int32(10)
		query.Limit = &defaultQueryLimit
	} else {
		query.Limit = query.Limit
	}

	notifications, err := s.repo.GetAll(ctx, query)

	if err != nil {
		return nil, err
	}

	// convert back to grpc typje
	notificationsData := make([]*pb.Notification, len(notifications))

	for index, notification := range notifications {
		notificationsData[index] = &pb.Notification{
			Id:        notification.ID.String(),
			MemberId:  notification.MemberID.String(),
			Type:      notification.Type,
			Title:     notification.Title,
			Message:   notification.Message,
			Read:      notification.Read,
			EmailSent: notification.EmailSent,
			SourceId:  notification.SourceID.String(),
			CreatedAt: timestamppb.New(notification.CreatedAt),
		}
	}

	notificationsResponse := &pb.GetNotificationsResponse{
		Data: notificationsData,
	}

	return notificationsResponse, nil
}

/**
* Notification Constants and Helper functions
**/

type NotificationType string

const (
	NotificationWelcome      NotificationType = "welcome"
	NotificationBuildCreated NotificationType = "build_created"
)

type NotificationTemplate struct {
	Type    NotificationType
	Title   string
	Message string
}

var welcomeNotificationMessage = NotificationTemplate{
	Type:    NotificationWelcome,
	Title:   "Welcome",
	Message: "Welcome to path of community, exile!",
}

var buildNotificationMessage = NotificationTemplate{
	Type:    NotificationBuildCreated,
	Title:   "Build Created",
	Message: "Build was successfully created.",
}

func (s *service) GetNotificationTemplate(notificationType NotificationType) (*NotificationTemplate, error) {
	notificationTemplates := map[NotificationType]*NotificationTemplate{
		NotificationWelcome:      &welcomeNotificationMessage,
		NotificationBuildCreated: &buildNotificationMessage,
	}

	template, exists := notificationTemplates[notificationType]

	if !exists {
		fmt.Printf("Error: NotificationType %s has no matching notification.\n", notificationType)
		return nil, fmt.Errorf("No matching notification for NotificationType: %s", notificationType)
	}

	return template, nil
}
