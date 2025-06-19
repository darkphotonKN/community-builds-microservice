package notification

import (
	"context"
	"fmt"
	"time"

	commonhelpers "github.com/darkphotonKN/community-builds-microservice/common/utils"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *repository {
	return &repository{db: db}
}

func (r *repository) Create(notification *CreateNotification) (*Notification, error) {
	now := time.Now()
	notificationModel := &Notification{
		ID:        uuid.New(),
		MemberID:  notification.MemberID,
		Type:      notification.Type,
		Title:     notification.Title,
		Message:   notification.Message,
		Read:      false,
		EmailSent: false,
		SourceID:  notification.SourceID,
		CreatedAt: now,
	}

	query := `
		INSERT INTO notifications (id, member_id, type, title, message, read, email_sent, source_id, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, member_id, type, title, message, read, email_sent, source_id, created_at
	`

	err := r.db.QueryRowx(
		query,
		notificationModel.ID,
		notificationModel.MemberID,
		notificationModel.Type,
		notificationModel.Title, notificationModel.Message,
		notificationModel.Read,
		notificationModel.EmailSent,
		notificationModel.SourceID,
		notificationModel.CreatedAt,
	).StructScan(notificationModel)

	if err != nil {
		return nil, commonhelpers.AnalyzeDBErr(err)
	}

	return notificationModel, nil
}

func (r *repository) GetAll(ctx context.Context, request *QueryNotifications) ([]Notification, error) {

	query := `
	SELECT 
		id,
		member_id,
		type,
		title,
		message,
		read,
		email_sent,
		source_id,
		created_at
	FROM notifications
	WHERE member_id = $1
	ORDER BY created_at DESC
	`
	paramCount := 1

	params := []interface{}{request.MemberID}

	if request.Limit != nil {
		fmt.Printf("Checking Limit: %d", request.Limit)
		paramCount++
		params = append(params, *request.Limit)

		query += "\nLIMIT $2"
	}

	if request.Offset != nil {
		fmt.Printf("Checking Offset: %d", request.Offset)
		paramCount++
		params = append(params, *request.Offset)

		query += fmt.Sprintf("\nOFFSET $%d", paramCount)
	}

	fmt.Println("final constructed query:", query)
	fmt.Println("parameters:", params)

	var notifications []Notification

	err := r.db.Select(&notifications, query, params...)

	if err != nil {
		return nil, commonhelpers.AnalyzeDBErr(err)
	}

	return notifications, nil
}
