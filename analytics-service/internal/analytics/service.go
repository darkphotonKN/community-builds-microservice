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

I want to see how many users signed up this week so I can track growth
I want to know how many users are active daily so I can measure engagement
I want to see which days users are most active so I can plan feature releases
I want to track user retention - do users come back after 1 week, 1 month?

As a Community Manager:

I want to see signup trends over time so I can correlate with marketing efforts
I want to know when users are most active so I can schedule announcements
I want to identify power users (login frequently) vs casual users

As a Developer:

I want to track login success vs failure rates to monitor system health
I want to see geographic patterns of where users are signing up from
I want to know which features trigger the most user activity

Event Flow Stories:
"User Logs In" Story:

User enters credentials in frontend
Auth service validates login
Auth service publishes "user logged in" event
Analytics service captures login time, user ID, location
Analytics service updates "daily active users" counter
Dashboard shows real-time active user count

"Track User Retention" Story:

User signs up (existing event you have)
Analytics service marks user as "new user"
Analytics service tracks when same user logs in again
Calculate: "Did user return within 7 days?"
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
