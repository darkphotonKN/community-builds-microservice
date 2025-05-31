package notification

import (
	"fmt"

	"github.com/darkphotonKN/community-builds-microservice/notification-service/internal/constants"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

type service struct {
	repo      Repository
	publishCh *amqp.Channel
}

type Repository interface {
	Create(notification *CreateNotification) (*Notification, error)
}

func NewService(repo Repository, ch *amqp.Channel) Service {
	return &service{repo: repo, publishCh: ch}
}

func (s *service) Create(memberCreated *MemberCreatedNotification) (*Notification, error) {
	// validate id is a legit uuid
	id, err := uuid.Parse(memberCreated.MemberID)

	if err != nil {
		return nil, err
	}

	// map it to notifications table entity
	createNotification := &CreateNotification{
		MemberID: id,
		Type:     constants.TypeWelcome,
		Title:    memberCreated.Title,
		Message:  memberCreated.Message,
	}

	newNotification, err := s.repo.Create(createNotification)
	if err != nil {
		fmt.Printf("Error when attempting to create notification: %s\n", err.Error())
		return nil, err
	}

	fmt.Println("notification was created:", newNotification)

	return newNotification, nil
}
