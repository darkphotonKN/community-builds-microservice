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
	Update(ctx context.Context, request *UpdateNotification) error
}

func NewService(repo Repository, ch *amqp.Channel) *service {
	return &service{repo: repo, publishCh: ch}
}

func (s *service) Create(notification *MemberCreatedNotification) (*Notification, error) {
	// validation and error handling TODO: missing fields
	if notification.Title == "" {
		fmt.Println("Title is required for creating a new notification.")
		return nil, status.Errorf(codes.InvalidArgument, "Name field is required")
	}

	// validate id is a legit uuid
	id, err := uuid.Parse(notification.MemberID)

	if err != nil {
		fmt.Println("Error occured when parsing uuid:", err)
		return nil, err
	}

	// map it to notifications table entity
	createNotification := &CreateNotification{
		MemberID: id,
		Type:     notification.Type,
		Title:    notification.Title,
		Message:  notification.Message,
		SourceID: notification.SourceID,
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
		query.Limit = request.Limit
	}

	if request.Offset == nil {
		defaultQueryOffset := int32(0)
		query.Offset = &defaultQueryOffset
	} else {
		query.Offset = request.Offset
	}

	notifications, err := s.repo.GetAll(ctx, query)

	if err != nil {
		return nil, err
	}

	// convert back to grpc type
	notificationsData := make([]*pb.Notification, len(notifications))

	for index, notification := range notifications {

		fmt.Printf("\nnotifications Read Field from DB, service layer before constructing grpc version: %+v\n\n", notification.Read)

		notificationsData[index] = &pb.Notification{
			Id:        notification.ID.String(),
			MemberId:  notification.MemberID.String(),
			Type:      notification.Type,
			Title:     notification.Title,
			Message:   notification.Message,
			Read:      notification.Read,
			EmailSent: notification.EmailSent,
			CreatedAt: timestamppb.New(notification.CreatedAt),
		}

		if notification.SourceID != nil {
			notificationsData[index].SourceId = notification.SourceID.String()
		}
	}

	// TEST: remove after

	fmt.Printf("\nAll Notifications from DB, service layer after constructing grpc version: %+v\n\n", notificationsData)

	for _, notification := range notificationsData {
		fmt.Printf("\nsingle Notification from DB, service layer after constructing grpc version, READ FIELD: %+v\n\n", notification.Read)

	}

	notificationsResponse := &pb.GetNotificationsResponse{
		Data: notificationsData,
	}

	return notificationsResponse, nil
}

func (s *service) CreateItem(itemCreated *CreateNotification) (*Notification, error) {
	// validation and error handling TODO: missing fields
	if itemCreated.Title == "" {
		fmt.Println("Title is required for creating a new notification.")
		return nil, status.Errorf(codes.InvalidArgument, "Name field is required")
	}

	// map it to notifications table entity
	createNotification := &CreateNotification{
		MemberID: itemCreated.MemberID,
		Type:     itemCreated.Type,
		Title:    "create_item_message",
		Message:  itemCreated.Message,
		SourceID: itemCreated.SourceID,
	}

	newNotification, err := s.repo.Create(createNotification)
	if err != nil {
		fmt.Println("Error occured when creating new notification:", err)
		return nil, err
	}

	fmt.Println("notification was created:", newNotification)

	return newNotification, nil
}

func (s *service) CreateRating(ratingCreated *RatingCreatedNotification) (*Notification, error) {
	// validation and error handling TODO: missing fields
	if ratingCreated.Title == "" {
		fmt.Println("Title is required for creating a new notification.")
		return nil, status.Errorf(codes.InvalidArgument, "Name field is required")
	}

	// map it to notifications table entity
	createNotification := &CreateNotification{
		MemberID: ratingCreated.MemberId,
		Type:     ratingCreated.Type,
		Title:    "create_rating_message",
		Message:  ratingCreated.Message,
		SourceID: ratingCreated.SourceID,
	}

	newNotification, err := s.repo.Create(createNotification)
	if err != nil {
		fmt.Println("Error occured when creating new notification:", err)
		return nil, err
	}

	fmt.Println("notification was created:", newNotification)

	return newNotification, nil
}

func (s *service) ReadNotification(ctx context.Context, request *pb.ReadNotificationRequest) error {
	// validate ids are legit uuids
	memberId, err := uuid.Parse(request.MemberId)

	if err != nil {
		fmt.Println("Error occured when parsing memberId as uuid:", err)
		return err
	}

	// validate id is a legit uuid
	notificationId, err := uuid.Parse(request.NotificationId)

	if err != nil {
		fmt.Println("Error occured when parsing notificationId as uuid:", err)
		return err
	}

	entity := &UpdateNotification{
		ID:       notificationId,
		MemberId: memberId,
		Read:     true,
	}

	return s.repo.Update(ctx, entity)
}

/**
* Notification Constants and Helper functions
**/

type NotificationType string

const (
	NotificationWelcome       NotificationType = "welcome"
	NotificationBuildCreated  NotificationType = "build_created"
	NotificationItemCreated   NotificationType = "item_created"
	NotificationRatingCreated NotificationType = "rating_created"
	NotificationPasswordReset NotificationType = "password_reset"
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

var itemCreateNotificationMessage = NotificationTemplate{
	Type:    NotificationItemCreated,
	Title:   "Create Item",
	Message: "Item was successfully created.",
}

var ratingCreateNotificationMessage = NotificationTemplate{
	Type:    NotificationRatingCreated,
	Title:   "Create Rating",
	Message: "Someone rated your build!",
}

var passwordResetNotificationMessage = NotificationTemplate{
	Type:    NotificationPasswordReset,
	Title:   "Password Reset",
	Message: "Password was successfully reset.",
}

func (s *service) GetNotificationTemplate(notificationType NotificationType) (*NotificationTemplate, error) {
	notificationTemplates := map[NotificationType]*NotificationTemplate{
		NotificationWelcome:       &welcomeNotificationMessage,
		NotificationBuildCreated:  &buildNotificationMessage,
		NotificationItemCreated:   &itemCreateNotificationMessage,
		NotificationRatingCreated: &ratingCreateNotificationMessage,
		NotificationPasswordReset: &passwordResetNotificationMessage,
	}

	template, exists := notificationTemplates[notificationType]

	if !exists {
		fmt.Printf("Error: NotificationType %s has no matching notification.\n", notificationType)
		return nil, fmt.Errorf("No matching notification for NotificationType: %s", notificationType)
	}

	return template, nil
}
