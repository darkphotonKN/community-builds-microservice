package analytics

type Creator interface {
	CreateMemberActivityEvent(req *CreateMemberActivityEvent) (*MemberActivityEvent, error)
}

type ConsumerService interface {
	Creator
	GetEventType(activityType ActivityType) (EventType, error)
}

// aggregate interface
type Service interface {
	Creator
	GetEventType(activityType ActivityType) (EventType, error)
}
