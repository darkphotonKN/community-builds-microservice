package broker

import (
	"fmt"
	"log"

	commonconstants "github.com/darkphotonKN/community-builds-microservice/common/constants"
	amqp "github.com/rabbitmq/amqp091-go"
)

/*
The Connect function establishes a connection to your RabbitMQ server and sets up
the exchanges needed for your service communication.
*/

func Connect(user, pass, host, port string) (*amqp.Channel, func() error) {
	address := fmt.Sprintf("amqp://%s:%s@%s:%s", user, pass, host, port)

	conn, err := amqp.Dial(address)

	if err != nil {
		log.Fatal(err)
	}

	ch, err := conn.Channel()

	if err != nil {
		log.Fatal(err)
	}

	// --- Declare All Exchanges ---

	// -- example service --
	err = ch.ExchangeDeclare(commonconstants.ExampleCreatedEvent, "fanout", true, false, false, false, nil)

	if err != nil {
		log.Fatal(err)
	}

	// -- auth service --
	err = ch.ExchangeDeclare(commonconstants.MemberSignedUpEvent, "fanout", true, false, false, false, nil)

	if err != nil {
		log.Fatal(err)
	}

	return ch, ch.Close
}
