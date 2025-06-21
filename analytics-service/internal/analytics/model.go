package analytics

import (
	"time"

	"github.com/google/uuid"
)

// MemberActivityEvent represents an event in the member_activity_events table
type MemberActivityEvent struct {
	ID        uuid.UUID `json:"id" db:"id"`
	MemberID  uuid.UUID `json:"member_id" db:"member_id"`
	EventType string    `json:"event_type" db:"event_type"`
	Timestamp time.Time `json:"timestamp" db:"timestamp"`
	Date      time.Time `json:"date" db:"date"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// DailyMemberStats represents aggregated stats in the daily_member_stats table
type DailyMemberStats struct {
	Date               time.Time `json:"date" db:"date"`
	NewSignups         int       `json:"new_signups" db:"new_signups"`
	ActiveMembers      int       `json:"active_members" db:"active_members"`
	TotalLogins        int       `json:"total_logins" db:"total_logins"`
	TotalMembersToDate int       `json:"total_members_to_date" db:"total_members_to_date"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" db:"updated_at"`
}

// MemberActivityEventMessage is used for AMQP message consumption
type MemberActivityEventMessage struct {
	MemberID  string     `json:"member_id"`
	EventType EventType  `json:"event_type"`
	EventName EventName  `json:"event_name"`
	Data      string     `json:"data"`
	SessionID *uuid.UUID `json:"session_id"`
}

// CreateMemberActivityEvent is used for creating new activity events
type CreateMemberActivityEvent struct {
	MemberID  uuid.UUID `json:"member_id" db:"member_id"`
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
