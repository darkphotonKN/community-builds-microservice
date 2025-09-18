package auth

import (
	"context"

	pb "github.com/darkphotonKN/community-builds-microservice/common/api/proto/auth"
	"github.com/google/uuid"
	// "google.golang.org/grpc/codes"
	// "google.golang.org/grpc/status"
)

type Handler struct {
	client AuthClient
}

func NewHandler(client AuthClient) *Handler {
	return &Handler{
		client: client,
	}
}

func (h *Handler) UpdateStripeCustomer(ctx context.Context, memberId uuid.UUID, stripeCustomerId string) error {

	req := &pb.UpdateStripeCustomerRequest{
		MemberId:         memberId.String(),
		StripeCustomerId: stripeCustomerId,
	}
	_, err := h.client.UpdateStripeCustomer(ctx, req)
	if err != nil {
		return err
	}

	return nil
}
