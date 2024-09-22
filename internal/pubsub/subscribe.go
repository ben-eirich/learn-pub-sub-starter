package pubsub

import (
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func SubscribeJSON[T any](
	conn *amqp.Connection,
	exchange, queueName, key string,
	simpleQueueType int,
	handler func(T)) error {

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
			handler(t)
			msg.Ack(false)
		}
	}()
	return nil
}
