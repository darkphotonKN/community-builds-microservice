package example

import (
	"context"
	"encoding/json"

	pb "github.com/darkphotonKN/community-builds-microservice/common/api/proto/example"
	commonconstants "github.com/darkphotonKN/community-builds-microservice/common/constants"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

type service struct {
	repo Repository
	publishCh *amqp.Channel
}

type Repository interface {
	Create(example *ExampleCreate) (*Example, error)
	GetByID(id uuid.UUID) (*Example, error)
}

func NewService(repo Repository, ch *amqp.Channel) Service {
	return &service{repo: repo, publishCh: ch}
}

func (s *service) CreateExample(ctx context.Context, req *pb.CreateExampleRequest) (*pb.Example, error) {
	// format to fit model for db tags
	createExample := &ExampleCreate{
		Name: req.Name,
	}
	example, err := s.repo.Create(createExample)

	if err != nil {
		return nil, err
	}

	// publish rabbit mq message after succesfuly creating 
	marshalledExample, err := json.Marshal(example)

	if err != nil {
		return nil, err
	}

	err = s.publishCh.PublishWithContext(
		ctx,
		commonconstants.ExampleCreatedEvent,
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        marshalledExample,
			// persist message
			DeliveryMode: amqp.Persistent,
		})

	if err != nil {
		return nil, err
	}

	return &pb.Example{
		Id:   example.ID,
		Name: example.Name,
	}, nil
}

func (s *service) GetExample(ctx context.Context, id uuid.UUID) (*pb.Example, error) {
	example, err := s.repo.GetByID(id)

	if err != nil {
		return nil, err
	}

	// format to fit grpc structure
	return &pb.Example{
		Id:   example.ID,
		Name: example.Name,
	}, nil
}
