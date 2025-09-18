package auth

import (
	"context"

	pb "github.com/darkphotonKN/community-builds-microservice/common/api/proto/auth"
)

type AuthClient interface {
	UpdateStripeCustomer(ctx context.Context, req *pb.UpdateStripeCustomerRequest) (*pb.UpdateStripeCustomerResponse, error)
}
