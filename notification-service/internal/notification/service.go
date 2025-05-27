package notification

import (
	"context"
	"encoding/json"

	pb "github.com/darkphotonKN/community-builds-microservice/common/api/proto/example"
	commonconstants "github.com/darkphotonKN/community-builds-microservice/common/constants"
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
	Create(notification *NotificationCreate) (*Notification, error)
}

func NewService(repo Repository, ch *amqp.Channel) Service {
	return &service{repo: repo, publishCh: ch}
}

func (s *service) Create(req *NotificationCreate) (*pb.Example, error) {
	// validation and error handling TODO: missing fields
	if req.Title == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Name field is required")
	}

	// For now, using example proto types until notification proto is created
	// This will be a placeholder implementation
	memberID := uuid.New() // This would come from the actual request

	// format to fit model for db tags
	createNotification := &NotificationCreate{
		MemberID: memberID,
		Type:     "example",
		Title:    "Example Notification",
		Message:  req.Name, // Using the name field as message for now
		SourceID: nil,
	}

	notification, err := s.repo.Create(createNotification)
	if err != nil {
		return nil, err
	}

	// publish rabbit mq message after successfully creating
	marshalledNotification, err := json.Marshal(notification)
	if err != nil {
		return nil, err
	}

	err = s.publishCh.PublishWithContext(
		ctx,
		commonconstants.ExampleCreatedEvent, // This would be NotificationCreatedEvent
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        marshalledNotification,
			// persist message
			DeliveryMode: amqp.Persistent,
		})

	if err != nil {
		return nil, err
	}

	// Return in proto format (placeholder until notification proto exists)
	return &pb.Example{
		Id:   notification.ID.String(),
		Name: notification.Message,
	}, nil
}

