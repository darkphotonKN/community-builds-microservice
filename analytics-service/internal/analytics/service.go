package analytics

import (
	"fmt"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

/*
User Stories:
As a Product Owner:

I want to see how many members signed up this week so I can track growth
I want to know how many members are active daily so I can measure engagement
I want to see which days members are most active so I can plan feature releases
I want to track member retention - do members come back after 1 week, 1 month?

As a Community Manager:

I want to see signup trends over time so I can correlate with marketing efforts
I want to know when members are most active so I can schedule announcements
I want to identify power members (login frequently) vs casual members

As a Developer:

I want to track login success vs failure rates to monitor system health
I want to see geographic patterns of where members are signing up from
I want to know which features trigger the most member activity

Event Flow Stories:
"Member Logs In" Story:

Member enters credentials in frontend
Auth service validates login
Auth service publishes "member logged in" event
Analytics service captures login time, member ID, location
Analytics service updates "daily active members" counter
Dashboard shows real-time active member count

"Track Member Retention" Story:

Member signs up (existing event you have)
Analytics service marks member as "new member"
Analytics service tracks when same member logs in again
Calculate: "Did member return within 7 days?"
Display retention rates on admin dashboard

"Popular Times" Story:

Track all login events throughout the day
Aggregate by hour: "Most logins happen at 8pm"
Help you decide when to deploy features or send notifications
*/

type service struct {
	repo      Repository
	publishCh *amqp.Channel
}

type Repository interface {
	Create(analytics *Analytics) (*Analytics, error)
}

func NewService(repo Repository, ch *amqp.Channel) Service {
	return &service{repo: repo, publishCh: ch}
}

func (s *service) Create(createAnalytics *CreateAnalytics) (*Analytics, error) {
	// validate activity
	eventType, err := s.GetEventType(createAnalytics.ActivityType)

	if err != nil {
		fmt.Printf("Error when attempting to validate activity: %s. Event type matches this activity", createAnalytics.ActivityType)
		return nil, status.Errorf(codes.InvalidArgument, "Error when attempting to get event name: %s"+err.Error())
	}

	id, err := uuid.Parse(createAnalytics.MemberID)

	if err != nil {
		return nil, err
	}

	// map it to analytics table entity
	newAnalytics := &Analytics{
		MemberID:     id,
		EventType:    string(eventType),
		ActivityType: string(createAnalytics.ActivityType),
		Data:         createAnalytics.Data,
	}

	analytics, err := s.repo.Create(newAnalytics)
	if err != nil {
		return nil, err
	}

	fmt.Println("analytics was created:", analytics)

	return analytics, nil
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
