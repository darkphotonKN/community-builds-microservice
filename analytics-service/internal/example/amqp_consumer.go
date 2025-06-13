package example

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	commonconstants "github.com/darkphotonKN/community-builds-microservice/common/constants"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"

	pb "github.com/darkphotonKN/community-builds-microservice/common/api/proto/example"
)

type consumer struct {
	service   ConsumerService
	publishCh *amqp.Channel
}

type ConsumerService interface {
	CreateExample(ctx context.Context, example *pb.CreateExampleRequest) (*pb.Example, error)
	GetExample(ctx context.Context, id uuid.UUID) (*pb.Example, error)
}

func NewConsumer(service Service, ch *amqp.Channel) *consumer {
	return &consumer{service: service, publishCh: ch}
}

func (c *consumer) Listen() {
	queueName := fmt.Sprintf("example.%s", commonconstants.ExampleCreatedEvent)

	// declare our unique queue that listens and waits for ExampleCreatedEvent to be published from example service
	queue, err := c.publishCh.QueueDeclare(queueName, true, false, false, false, nil)

	if err != nil {
		log.Fatal(err)
	}

	// bind to the exchange that will publish ExampleCreateEvent events
	err = c.publishCh.QueueBind(
		queue.Name,
		"",
		commonconstants.ExampleCreatedEvent,
		false,
		nil,
	)

	if err != nil {
		log.Fatal(err)
	}

	// consume messages, delivers messages from the queue
	msgs, err := c.publishCh.Consume(queue.Name, "", true, false, false, false, nil)

	// start a goroutine to listen for events
	go func() {
		for msg := range msgs {
			var createdExample *CreateExampleEvent

			err := json.Unmarshal(msg.Body, &createdExample)
			if err != nil {
				fmt.Printf("Error when unmarshalling exampl event created body: %s\n", err.Error())
			}

			fmt.Printf("\nsuccessfully received event message: %+v\n\n", createdExample)
		}
	}()
}
