package analytics

import (
	commonhelpers "github.com/darkphotonKN/community-builds-microservice/common/utils"
	"github.com/jmoiron/sqlx"
)

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateMemberActivityEvent(memberActivityEvent *MemberActivityEvent) (*MemberActivityEvent, error) {
	query := `
		INSERT INTO member_activity_events (member_id, activity_type, timestamp, date)
		VALUES ($1, $2, $3, $4)
		RETURNING id, member_id, activity_type, timestamp, date
	`

	var result MemberActivityEvent
	err := r.db.QueryRowx(
		query,
		memberActivityEvent.MemberID,
		memberActivityEvent.ActivityType,
		memberActivityEvent.Timestamp,
		memberActivityEvent.Date,
	).StructScan(&result)

	if err != nil {
		return nil, commonhelpers.AnalyzeDBErr(err)
	}

	return &result, nil
}
