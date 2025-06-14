package notification

import (
	"context"
	"fmt"

	pb "github.com/darkphotonKN/community-builds-microservice/common/api/proto/notification"
	"github.com/darkphotonKN/community-builds-microservice/common/discovery"
)

const (
	serviceName = "notifications"
)

type Client struct {
	registry discovery.Registry
}

func NewClient(registry discovery.Registry) NotificationClient {
	return &Client{
		registry: registry,
	}
}

func (c *Client) GetNotifications(ctx context.Context, req *pb.GetNotificationsRequest) (*pb.GetNotificationsResponse, error) {
	conn, err := discovery.ServiceConnection(ctx, serviceName, c.registry)

	if err != nil {
		return nil, fmt.Errorf("failed to connect to notification service: %w", err)
	}
	defer conn.Close()

	client := pb.NewNotificationServiceClient(conn)

	response, err := client.GetNotifications(ctx, req)
	return response, err
}

