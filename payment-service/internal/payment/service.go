package payment

import (
	"context"
	"fmt"

	pb "github.com/darkphotonKN/community-builds-microservice/common/api/proto/payment"
	"github.com/darkphotonKN/community-builds-microservice/payment-service/internal/grpc-gateway/auth"
	"github.com/google/uuid"
)

/**
* payment service
**/

type service struct {
	authHandler      *auth.Handler
	paymentProcessor PaymentProcessor
	repo             Repository
}

type Repository interface {
	Create(ctx context.Context, request *CheckoutSessionRequest) (uuid.UUID, error)
}

type PaymentUserService interface {
	UpdateStripeCustomer(ctx context.Context, userID uuid.UUID, stripeCustomerID string) error
}

func NewService(authHandler *auth.Handler, paymentProcessor PaymentProcessor) *service {
	return &service{
		authHandler:      authHandler,
		paymentProcessor: paymentProcessor,
	}
}

/**
* Create customer
**/

func (s *service) CreateCustomer(ctx context.Context, req *pb.CreateCustomerRequest) (*pb.CreateCustomerResponse, error) {

	memberId, err := uuid.Parse(req.MemberId)
	if err != nil {
		return nil, err
	}
	email := req.Email

	//todo: check if customer already exists in our local db
	// checkCustomerRes, err := s.authHandler.GetMember(ctx, &pb.GetMemberRequest{Id: memberId.String()})
	// create customer on stripe and get customer id
	customerId, err := s.paymentProcessor.CreateCustomer(ctx, memberId, email)

	if err != nil {
		fmt.Printf("Error occured when attemtping to create customer on stripe, %s\n", err.Error())
		return nil, err
	}

	// update local user repo for mapping
	err = s.authHandler.UpdateStripeCustomer(ctx, memberId, customerId)

	if err != nil {
		fmt.Printf("Error occured when attempting to update stripe customerId to user repo in CreateCustomer method: %s\n", err.Error())
		return nil, err
	}

	grpcResponse := &pb.CreateCustomerResponse{
		CustomerId: customerId,
	}

	return grpcResponse, nil
}

func (s *service) CreateSubscription(ctx context.Context, req *pb.CreateSubscriptionRequest) (*pb.CreateSubscriptionResponse, error) {
	memberId, err := uuid.Parse(req.MemberId)
	if err != nil {
		return nil, err
	}

	// create customer if not exists
	CreateCustomerRes, err := s.CreateCustomer(ctx, &pb.CreateCustomerRequest{
		MemberId: memberId.String(),
		Email:    req.Email,
	})
	fmt.Printf("Customer ID from CreateCustomer in CreateSubscription: %+v\n", CreateCustomerRes.CustomerId)

	if err != nil {
		return nil, err
	}

	request := &SetupProductsReq{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
	}
	// create subscription on stripe and get subscription id
	_, err = s.paymentProcessor.CreateSubscription(ctx, request)

	if err != nil {
		fmt.Printf("Error occured when attemtping to create subscription on stripe, %s\n", err.Error())
		return nil, err
	}

	grpcResponse := &pb.CreateSubscriptionResponse{}

	return grpcResponse, nil
}
