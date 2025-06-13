package analytics

import (
	"time"

	"github.com/google/uuid"
)

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

type MemberActivityEvent struct {
	MemberID  string     `json:"member_id" db:"member_id"`
	EventType string     `json:"event_type" db:"event_type"`
	EventName string     `json:"event_name" db:"event_name"`
	Data      string     `json:"data" db:"data"`
	SessionID *uuid.UUID `json:"session_id" db:"session_id"`
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

