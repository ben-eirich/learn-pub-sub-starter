package pubsub

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	QueueDurable   = 0
	QueueTransient = 1
)

func DeclareAndBind(
	conn *amqp.Connection,
	exchange, queueName, key string,
	simpleQueueType int) (*amqp.Channel, amqp.Queue, error) {

	ch, err := conn.Channel()
	if err != nil {
		return nil, amqp.Queue{}, err
	}

	queue, err := ch.QueueDeclare(
		queueName,
		simpleQueueType == QueueDurable,
		simpleQueueType == QueueTransient,
		simpleQueueType == QueueTransient,
		false,
		nil)
	if err != nil {
		return nil, amqp.Queue{}, err
	}

	err = ch.QueueBind(queueName, key, exchange, false, nil)
	if err != nil {
		return nil, amqp.Queue{}, err
	}

	return ch, queue, nil
}
