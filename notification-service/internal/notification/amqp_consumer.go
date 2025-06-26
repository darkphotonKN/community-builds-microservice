package notification

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"time"

	commonconstants "github.com/darkphotonKN/community-builds-microservice/common/constants"
	"github.com/google/uuid"
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
	go c.ItemCreatedItemEventListener()

	fmt.Println("Notification consumer started - listening for member signup events.")
}

/**
* Member Sign-Up Events Listener
* Handles:
* - Member sign ups welcome notificiations
**/
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

		// get the correct notification template
		template, err := c.service.GetNotificationTemplate(NotificationWelcome)

		if err != nil {
			fmt.Println("Failed on getting notification template. Error:", err)
			continue
		}

		var notification *Notification

		for i := 0; i < 3; i++ {
			// create event TODO: sourceID missing, add later
			notification, err = c.service.Create(&MemberCreatedNotification{
				Title:    template.Title,
				Message:  template.Message,
				Type:     string(template.Type),
				MemberID: memberSignedUp.UserID,
			})

			// retry on error
			if err != nil {
				fmt.Printf("Failed to create notification on %d try, err: %s", i+1, err.Error())

				// progressive delay / back-off before retrying
				time.Sleep(time.Duration(math.Pow(2, float64(i))) * time.Second)
				continue
			}

			// otherwise exit
			break
		}

		if err != nil {
			fmt.Println("Failed to create notification, err:", err)
			continue
		}

		fmt.Printf("Notification created:\nid: %s\nType: %s\nTitle: %s\nMessage: %s\n", notification.ID, notification.Type, notification.Title, notification.Message)
	}
}

/**
* Item create item event Listener
* Handles:
* - Item create item notificiations
**/
func (c *consumer) ItemCreatedItemEventListener() {
	queueName := fmt.Sprintf("notification.%s", commonconstants.ItemCreatedItemEvent)

	queue, err := c.publishCh.QueueDeclare(queueName, true, false, false, false, nil)

	if err != nil {
		log.Fatal(err)
	}

	// bind to the exchange that will publish ExampleCreateEvent events
	err = c.publishCh.QueueBind(
		queue.Name,
		"",
		commonconstants.ItemCreatedItemEvent,
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
		var itemCreateItem *commonconstants.ItemCreatedItemEventPayload

		err := json.Unmarshal(msg.Body, &itemCreateItem)
		if err != nil {
			fmt.Printf("Error when unmarshalling member.itemCreateItem event body: %s\n", err.Error())
		}

		fmt.Printf("\nsuccessfully received event message: %+v\n\n", itemCreateItem)

		// get the correct notification template
		template, err := c.service.GetNotificationTemplate(NotificationItemCreated)

		if err != nil {
			fmt.Println("Failed on getting notification template. Error:", err)
			continue
		}

		var notification *Notification

		memberId, err := uuid.Parse(itemCreateItem.UserID)
		if err != nil {
			fmt.Println("invalid UUID: %w", err)
		}
		for i := 0; i < 3; i++ {
			// create event TODO: sourceID missing, add later
			notification, err = c.service.CreateItem(&CreateNotification{
				Title:    template.Title,
				Message:  template.Message,
				Type:     string(template.Type),
				MemberID: memberId,
			})

			// retry on error
			if err != nil {
				fmt.Printf("Failed to create notification on %d try, err: %s", i+1, err.Error())

				// progressive delay / back-off before retrying
				time.Sleep(time.Duration(math.Pow(2, float64(i))) * time.Second)
				continue
			}

			// otherwise exit
			break
		}

		if err != nil {
			fmt.Println("Failed to create notification, err:", err)
			continue
		}

		fmt.Printf("Notification created:\nid: %s\nType: %s\nTitle: %s\nMessage: %s\n", notification.ID, notification.Type, notification.Title, notification.Message)
	}
}
