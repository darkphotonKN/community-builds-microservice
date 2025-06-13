package analytics

import (
	"fmt"

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
	Create(analytics *CreateAnalytics) (*Analytics, error)
}

func NewService(repo Repository, ch *amqp.Channel) Service {
	return &service{repo: repo, publishCh: ch}
}

func (s *service) Create(memberActivity *MemberActivityEvent) (*Analytics, error) {
	// validation and error handling TODO: missing fields
	if memberActivity.EventName == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Event name field is required")
	}

	// validate id is a legit uuid
	id, err := uuid.Parse(memberActivity.MemberID)

	if err != nil {
		return nil, err
	}

	// map it to analytics table entity
	createAnalytics := &CreateAnalytics{
		MemberID:  id,
		EventType: memberActivity.EventType,
		EventName: memberActivity.EventName,
		Data:      memberActivity.Data,
		SessionID: memberActivity.SessionID,
		IPAddress: "", // TODO: extract from context
		UserAgent: "", // TODO: extract from context
	}

	newAnalytics, err := s.repo.Create(createAnalytics)
	if err != nil {
		return nil, err
	}

	fmt.Println("analytics was created:", newAnalytics)

	return newAnalytics, nil
}

