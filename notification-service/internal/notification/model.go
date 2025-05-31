package notification

import (
	"time"

	"github.com/google/uuid"
)

type Notification struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	MemberID  uuid.UUID  `json:"member_id" db:"member_id"`
	Type      string     `json:"type" db:"type"`
	Title     string     `json:"title" db:"title"`
	Message   string     `json:"message" db:"message"`
	Read      bool       `json:"read" db:"read"`
	EmailSent bool       `json:"email_sent" db:"email_sent"`
	SourceID  *uuid.UUID `json:"source_id" db:"source_id"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
}

type MemberCreatedNotification struct {
	MemberID string     `json:"member_id" db:"member_id"`
	Title    string     `json:"title" db:"title"`
	Message  string     `json:"message" db:"message"`
	SourceID *uuid.UUID `json:"source_id" db:"source_id"`
}

type CreateNotification struct {
	MemberID uuid.UUID  `json:"member_id" db:"member_id"`
	Type     string     `json:"type" db:"type"`
	Title    string     `json:"title" db:"title"`
	Message  string     `json:"message" db:"message"`
	SourceID *uuid.UUID `json:"source_id" db:"source_id"`
}
