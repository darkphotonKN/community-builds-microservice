package notification

import (
	"context"
	"fmt"

	pb "github.com/darkphotonKN/community-builds-microservice/common/api/proto/notification"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type service struct {
	repo      Repository
	publishCh *amqp.Channel
}

type Repository interface {
	Create(notification *CreateNotification) (*Notification, error)
	GetAll(ctx context.Context, request *QueryNotifications) (*pb.GetNotificationsResponse, error)
}

func NewService(repo Repository, ch *amqp.Channel) *service {
	return &service{repo: repo, publishCh: ch}
}

func (s *service) Create(memberCreated *MemberCreatedNotification) (*Notification, error) {
	// validation and error handling TODO: missing fields
	if memberCreated.Title == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Name field is required")
	}

	// validate id is a legit uuid
	id, err := uuid.Parse(memberCreated.MemberID)

	if err != nil {
		return nil, err
	}

	// map it to notifications table entity
	createNotification := &CreateNotification{
		MemberID: id,
		Type:     memberCreated.Type,
		Title:    memberCreated.Title,
		Message:  memberCreated.Message,
		SourceID: memberCreated.SourceID,
	}

	newNotification, err := s.repo.Create(createNotification)
	if err != nil {
		return nil, err
	}

	fmt.Println("notification was created:", newNotification)

	return newNotification, nil
}

func (s *service) GetAllByMemberId(ctx context.Context, request *QueryNotifications) (*pb.GetNotificationsResponse, error) {

	// convert to DB appropriate format
	query := &QueryNotifications{
		MemberID: request.MemberID,
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

	return s.repo.GetAll(ctx, query)
}
