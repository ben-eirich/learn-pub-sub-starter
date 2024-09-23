package pubsub

import (
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type AckType int

const (
	Ack         = 0
	NackRequeue = 1
	NackDiscard = 2
)

func SubscribeJSON[T any](
	conn *amqp.Connection,
	exchange, queueName, key string,
	simpleQueueType int,
	handler func(T) AckType) error {

	amqpChan, _, err := DeclareAndBind(conn, exchange, queueName, key, simpleQueueType)
	if err != nil {
		return err
	}

	ch, err := amqpChan.Consume(queueName, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	go func() {
		for msg := range ch {
			var t T
			err := json.Unmarshal(msg.Body, &t)
			if err != nil {
				log.Printf("Error decoding message: %v\n", err)
				continue
			}
			ackType := handler(t)
			switch ackType {
			case Ack:
				msg.Ack(false)
			case NackRequeue:
				msg.Nack(false, true)
			case NackDiscard:
				msg.Nack(false, false)
			}
		}
	}()
	return nil
}
