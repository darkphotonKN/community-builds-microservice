package analytics

import (
	"time"

	"github.com/google/uuid"
)

type MemberActivityEvent struct {
	ID           uuid.UUID `json:"id" db:"id"`
	MemberID     uuid.UUID `json:"member_id" db:"member_id"`
	ActivityType string    `json:"activity_type" db:"activity_type"`
	Timestamp    time.Time `json:"timestamp" db:"timestamp"` // more accurate for actual activity event time
	Date         time.Time `json:"date" db:"date"`
}

type DailyMemberStats struct {
	Date               time.Time `json:"date" db:"date"`
	NewSignups         int       `json:"new_signups" db:"new_signups"`
	ActiveMembers      int       `json:"active_members" db:"active_members"`
	TotalLogins        int       `json:"total_logins" db:"total_logins"`
	TotalMembersToDate int       `json:"total_members_to_date" db:"total_members_to_date"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" db:"updated_at"`
}

// CreateMemberActivityEvent is used for creating new activity events
type CreateMemberActivityEvent struct {
	MemberID     string       `json:"member_id" db:"member_id"`
	ActivityType ActivityType `json:"activity_type" db:"activity_type"`
}
