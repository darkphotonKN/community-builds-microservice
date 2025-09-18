package payment

import (
	"context"

	pb "github.com/darkphotonKN/community-builds-microservice/common/api/proto/payment"
)

type Handler struct {
	service Service
	pb.UnimplementedPaymentServiceServer
}

type Service interface {
	CreateCustomer(ctx context.Context, req *pb.CreateCustomerRequest) (*pb.CreateCustomerResponse, error)
	CreateSubscription(ctx context.Context, req *pb.CreateSubscriptionRequest) (*pb.CreateSubscriptionResponse, error)
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateCustomer(ctx context.Context, req *pb.CreateCustomerRequest) (*pb.CreateCustomerResponse, error) {
	return h.service.CreateCustomer(ctx, req)
}

func (h *Handler) CreateSubscription(ctx context.Context, req *pb.CreateSubscriptionRequest) (*pb.CreateSubscriptionResponse, error) {
	return h.service.CreateSubscription(ctx, req)
}
