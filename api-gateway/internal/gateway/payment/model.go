package payment

import (
	"context"

	pb "github.com/darkphotonKN/community-builds-microservice/common/api/proto/payment"
)

type CreateSubScriptionRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	Price       int64  `json:"price" binding:"required"`
}

type PaymentClient interface {
	CreateCustomer(ctx context.Context, req *pb.CreateCustomerRequest) (*pb.CreateCustomerResponse, error)
	CreateSubscription(ctx context.Context, req *pb.CreateSubscriptionRequest) (*pb.CreateSubscriptionResponse, error)
}
