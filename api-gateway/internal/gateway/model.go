package gateway

import "context"

type ExampleGateway interface {
	CreateExample(ctx context.Context, req *pb.CreateExampleRequest) (*pb.Example, error)
	GetExample(ctx context.Context, req *pb.GetExampleRequest) (*pb.Example, error)
}