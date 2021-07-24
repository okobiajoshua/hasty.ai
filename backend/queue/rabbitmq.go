package queue

import (
	"context"
	"log"
	"time"

	"github.com/streadway/amqp"
)

type RabbitMQ struct {
	ch *amqp.Channel
	q  amqp.Queue
}

func NewRabbitMQ(url, queueName string) *RabbitMQ {
	// Sleep for 20 secs before attempting to connect
	time.Sleep(20 * time.Second)
	conn, err := amqp.Dial(url)
	if err != nil {
		log.Fatal("Failed to connect to rabbitmq")
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal("Failed to open a channel")
	}

	q, err := ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		log.Fatal("Failed to create queue")
	}

	return &RabbitMQ{
		ch: ch,
		q:  q,
	}
}

func (r *RabbitMQ) Publish(ctx context.Context, msg []byte) error {
	return r.ch.Publish(
		"",       // exchange
		r.q.Name, // routing key
		false,    // mandatory
		false,    // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        msg,
		})
}
