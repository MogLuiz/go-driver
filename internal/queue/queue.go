package queue

import (
	"fmt"
	"log"
)

type QueueType int

const (
	RabbitMQ QueueType = iota
)

func New(qt QueueType, cfg any) *Queue {
	q := new(Queue)
	switch qt {
	case RabbitMQ:
		fmt.Println("Não implementado")
	default:
		log.Fatal("Queue type not implemented")
	}

	return q
}

type QueueConnection interface {
	Publish([]byte) error
	Consume() error
}

type Queue struct {
	cfg any
	qc  QueueConnection
}

func (q *Queue) Publish(msg []byte) error {
	return q.qc.Publish(msg)
}

func (q *Queue) Consume() error {
	return q.qc.Consume()
}
