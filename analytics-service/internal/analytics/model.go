package analytics

import (
	"time"

	"github.com/google/uuid"
)

// UserActivityEvent represents an event in the user_activity_events table
type UserActivityEvent struct {
	ID        uuid.UUID `json:"id" db:"id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	EventType string    `json:"event_type" db:"event_type"`
	Timestamp time.Time `json:"timestamp" db:"timestamp"`
	Date      time.Time `json:"date" db:"date"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// DailyUserStats represents aggregated stats in the daily_user_stats table
type DailyUserStats struct {
	Date             time.Time `json:"date" db:"date"`
	NewSignups       int       `json:"new_signups" db:"new_signups"`
	ActiveUsers      int       `json:"active_users" db:"active_users"`
	TotalLogins      int       `json:"total_logins" db:"total_logins"`
	TotalUsersToDate int       `json:"total_users_to_date" db:"total_users_to_date"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}

// MemberActivityEvent is used for AMQP message consumption
type MemberActivityEvent struct {
	MemberID  string `json:"member_id"`
	EventType string `json:"event_type"`
}

// CreateUserActivityEvent is used for creating new activity events
type CreateUserActivityEvent struct {
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	EventType string    `json:"event_type" db:"event_type"`
}

// Keep the legacy types for compatibility - DO NOT MODIFY
type Analytics struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	MemberID  uuid.UUID  `json:"member_id" db:"member_id"`
	EventType string     `json:"event_type" db:"event_type"`
	EventName string     `json:"event_name" db:"event_name"`
	Data      string     `json:"data" db:"data"`
	SessionID *uuid.UUID `json:"session_id" db:"session_id"`
	IPAddress string     `json:"ip_address" db:"ip_address"`
	UserAgent string     `json:"user_agent" db:"user_agent"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
}

type CreateAnalytics struct {
	MemberID  uuid.UUID  `json:"member_id" db:"member_id"`
	EventType string     `json:"event_type" db:"event_type"`
	EventName string     `json:"event_name" db:"event_name"`
	Data      string     `json:"data" db:"data"`
	SessionID *uuid.UUID `json:"session_id" db:"session_id"`
	IPAddress string     `json:"ip_address" db:"ip_address"`
	UserAgent string     `json:"user_agent" db:"user_agent"`
}
