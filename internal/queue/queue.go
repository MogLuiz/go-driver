package queue

import (
	"fmt"
	"log"
	"reflect"
)

type QueueType int

const (
	RabbitMQ QueueType = iota
)

func New(qt QueueType, cfg any) (q *Queue, err error) {
	rt := reflect.TypeOf(cfg)

	switch qt {
	case RabbitMQ:
		if rt.Name() != "RabbitMQConfig" {
			return nil, fmt.Errorf("config need's to be of type rabbitmqconfig")
		}
		fmt.Println("Não implementado")
	default:
		log.Fatal("Queue type not implemented")
	}

	return
}

type QueueConnection interface {
	Publish([]byte) error
	Consume(chan<- QueueDTO) error
}

type Queue struct {
	cfg any
	qc  QueueConnection
}

func (q *Queue) Publish(msg []byte) error {
	return q.qc.Publish(msg)
}

func (q *Queue) Consume(cdto chan<- QueueDTO) error {
	return q.qc.Consume(cdto)
}
