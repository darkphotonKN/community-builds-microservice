package analytics

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

type service struct {
	repo      Repository
	publishCh *amqp.Channel
}

type Repository interface {
	CreateMemberActivityEvent(memberActivityEvent *MemberActivityEvent) (*MemberActivityEvent, error)
}

func NewService(repo Repository, ch *amqp.Channel) Service {
	return &service{repo: repo, publishCh: ch}
}

func (s *service) CreateMemberActivityEvent(req *CreateMemberActivityEvent) (*MemberActivityEvent, error) {
	id, err := uuid.Parse(req.MemberID)

	if err != nil {
		return nil, err
	}

	// map it to analytics table entity
	entity := &MemberActivityEvent{
		MemberID:     id,
		ActivityType: string(req.ActivityType),
		Timestamp:    time.Now(),
		Date:         time.Now().Truncate(24 * time.Hour),
	}

	newMemberActivityEvent, err := s.repo.CreateMemberActivityEvent(entity)
	if err != nil {
		return nil, err
	}

	fmt.Println("MemberActivityEvent analytics was created:", newMemberActivityEvent)

	return newMemberActivityEvent, nil
}

type EventType string

const (
	EventTypeMember EventType = "member_event"
	EventTypeBuild  EventType = "build_event"
)

type ActivityType string

const (
	ActivityTypeMemberCreated  ActivityType = "member_created_activity"
	ActivityTypeMemberLoggedOn ActivityType = "member_logged_on_activity"
	ActivityTypeBuildViewed    ActivityType = "build_viewed_activity"
)

/**
* Gets event type and validates if an activity is under a certain event type.
**/
func (s *service) GetEventType(activityType ActivityType) (EventType, error) {
	eventMap := map[ActivityType]EventType{
		// member
		ActivityTypeMemberCreated:  EventTypeMember,
		ActivityTypeMemberLoggedOn: EventTypeMember,

		// build
		ActivityTypeBuildViewed: EventTypeBuild,
	}

	event, exists := eventMap[activityType]

	if !exists {
		return "", fmt.Errorf("Activity %s doesn't exist under any event.", activityType)
	}

	return event, nil
}
