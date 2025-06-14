package notification

import (
	"encoding/json"
	"fmt"
	"log"

	commonconstants "github.com/darkphotonKN/community-builds-microservice/common/constants"
	amqp "github.com/rabbitmq/amqp091-go"
)

type consumer struct {
	service   EventConsumerService
	publishCh *amqp.Channel
}

func NewConsumer(service EventConsumerService, ch *amqp.Channel) *consumer {
	return &consumer{service: service, publishCh: ch}
}

func (c *consumer) Listen() {
	go c.memberSignedUpEventListener()

	fmt.Println("Notification consumer started - listening for member signup events.")
}

func (c *consumer) memberSignedUpEventListener() {
	queueName := fmt.Sprintf("notification.%s", commonconstants.MemberSignedUpEvent)

	queue, err := c.publishCh.QueueDeclare(queueName, true, false, false, false, nil)

	if err != nil {
		log.Fatal(err)
	}

	// bind to the exchange that will publish ExampleCreateEvent events
	err = c.publishCh.QueueBind(
		queue.Name,
		"",
		commonconstants.MemberSignedUpEvent,
		false,
		nil,
	)

	if err != nil {
		log.Fatal(err)
	}

	// consume messages, delivers messages from the queue
	msgs, err := c.publishCh.Consume(queue.Name, "", true, false, false, false, nil)

	// start a goroutine to listen for events
	for msg := range msgs {
		var memberSignedUp *commonconstants.MemberSignedUpEventPayload

		err := json.Unmarshal(msg.Body, &memberSignedUp)
		if err != nil {
			fmt.Printf("Error when unmarshalling member.signedup event body: %s\n", err.Error())
		}

		fmt.Printf("\nsuccessfully received event message: %+v\n\n", memberSignedUp)
		// create event
		c.service.Create(&MemberCreatedNotification{
			MemberID: memberSignedUp.UserID,
		})
	}
}
