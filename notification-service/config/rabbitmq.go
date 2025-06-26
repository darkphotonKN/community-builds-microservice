package config

import (
	"log"

	"github.com/darkphotonKN/community-builds-microservice/common/broker"
	commonconstants "github.com/darkphotonKN/community-builds-microservice/common/constants"
	amqp "github.com/rabbitmq/amqp091-go"
)

/*
The Connect function establishes a connection to your RabbitMQ server and sets up
the exchanges needed for your service communication.
*/

type Exchange struct {
	ExchangeName string
	ExchangeType string
}

func DeclareExchanges(ch *amqp.Channel) {

	exchangeList := []Exchange{
		// {ExchangeName: commonconstants.ExampleCreatedEvent, ExchangeType: "fanout"},
		{ExchangeName: commonconstants.MemberSignedUpEvent, ExchangeType: "fanout"},
		{ExchangeName: commonconstants.ItemCreatedItemEvent, ExchangeType: "fanout"},
	}

	for _, exchange := range exchangeList {

		err := broker.DeclareExchange(ch, exchange.ExchangeName, exchange.ExchangeType)

		if err != nil {
			log.Fatal("fail exchange: %v", exchange.ExchangeName)
		} else {
			log.Printf("Declared exchange: %v", exchange.ExchangeName)
		}
	}

}
