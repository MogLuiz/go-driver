package queue

import (
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQConfig struct {
	URL       string
	TopicName string
	TimeOut   time.Time
}

type RabbitConnection struct {
	cfg  RabbitMQConfig
	conn *amqp.Connection
}

func NewRabbitConnection(cfg RabbitMQConfig) (*RabbitConnection, error) {
	conn, err := amqp.Dial(cfg.URL)
	if err != nil {
		return nil, err
	}
	return &RabbitConnection{
		cfg:  cfg,
		conn: conn,
	}, nil
}
