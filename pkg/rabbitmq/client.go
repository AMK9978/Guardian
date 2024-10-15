package rabbitmq

import (
	"log"

	"github.com/streadway/amqp"
)

type Client struct {
	Connection *amqp.Connection
}

func NewClient(rabbitMQURI string) *Client {
	conn, err := amqp.Dial(rabbitMQURI)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}

	return &Client{Connection: conn}
}
