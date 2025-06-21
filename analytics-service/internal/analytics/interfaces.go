package analytics

type Creator interface {
	Create(createAnalytics *CreateAnalytics) (*Analytics, error)
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
