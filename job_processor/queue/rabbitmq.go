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
	// Sleep for 10 secs before attempting to connect
	time.Sleep(10 * time.Second)
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

func (r *RabbitMQ) Consume(ctx context.Context, f func(context.Context, []byte) error) {
	// Pick message from queue
	msgs, err := r.ch.Consume(
		r.q.Name, // queue
		"",       // consumer
		false,    // auto-ack
		false,    // exclusive
		false,    // no-local
		false,    // no-wait
		nil,      // args
	)

	if err != nil {
		log.Fatal("Failed to consume from queue")
	}

	for msg := range msgs {
		go func(ctx context.Context, d amqp.Delivery) {
			if err := f(ctx, d.Body); err == nil {
				d.Ack(true)
			} else {
				log.Println("Error during processing", err)
			}
		}(ctx, msg)
	}
}
