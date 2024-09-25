package hotels

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

type Publisher struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
	Queue      amqp.Queue
}

func NewPublisher(user string, password string, host string, port string, queueName string) Publisher {
	connection, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@:%s:%s", user, password, host, port))
	if err != nil {
		log.Panicf("Error connecting to RabbitMQ: %v", err)
	}

	channel, err := connection.Channel()
	if err != nil {
		log.Panicf("Error getting RabbitMQ channel: %v", err)
	}

	queue, err := channel.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Panicf("Error creating queue: %v", err)
	}

	return Publisher{
		Connection: connection,
		Channel:    channel,
		Queue:      queue,
	}
}

func (p Publisher) Publish(id int64) {
	_ = p.Channel.Publish(
		"",
		p.Queue.Name,
		false,
		false,
		amqp.Publishing{
			Body: []byte(fmt.Sprintf("{id:%d}", id)),
		})
}
