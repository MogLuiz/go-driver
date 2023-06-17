package queue

import (
	"context"
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

func (rc *RabbitConnection) Publish(msg []byte) error {
	ch, err := rc.conn.Channel()
	if err != nil {
		return err
	}

	mp := amqp.Publishing{
		Body:         msg,
		DeliveryMode: amqp.Persistent,
		ContentType:  "text/plain",
		Timestamp:    time.Now(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	return ch.PublishWithContext(ctx, "", rc.cfg.TopicName, false, false, mp)
}

func (rc *RabbitConnection) Consume(cdto chan<- QueueDTO) error {
	ch, err := rc.conn.Channel()
	if err != nil {
		return err
	}

	q, err := ch.QueueDeclare(rc.cfg.TopicName, false, false, false, false, nil)
	if err != nil {
		return err
	}

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	for msg := range msgs {
		dto := QueueDTO{}
		dto.Unmarshal(msg.Body)

		cdto <- dto
	}

	return nil
}
