package analytics

import (
	"time"

	commonhelpers "github.com/darkphotonKN/community-builds-microservice/common/utils"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(analytics *CreateAnalytics) (*Analytics, error) {
	now := time.Now()
	analyticsModel := &Analytics{
		ID:        uuid.New(),
		MemberID:  analytics.MemberID,
		EventType: analytics.EventType,
		EventName: analytics.EventName,
		Data:      analytics.Data,
		SessionID: analytics.SessionID,
		IPAddress: analytics.IPAddress,
		UserAgent: analytics.UserAgent,
		CreatedAt: now,
	}

	query := `
		INSERT INTO analytics (id, member_id, event_type, event_name, data, session_id, ip_address, user_agent, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, member_id, event_type, event_name, data, session_id, ip_address, user_agent, created_at
	`

	err := r.db.QueryRowx(
		query,
		analyticsModel.ID,
		analyticsModel.MemberID,
		analyticsModel.EventType,
		analyticsModel.EventName,
		analyticsModel.Data,
		analyticsModel.SessionID,
		analyticsModel.IPAddress,
		analyticsModel.UserAgent,
		analyticsModel.CreatedAt,
	).StructScan(analyticsModel)

	if err != nil {
		return nil, commonhelpers.AnalyzeDBErr(err)
	}

	return analyticsModel, nil
}

