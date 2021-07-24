package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/hastyai/processor/datastore"
	"github.com/hastyai/processor/queue"
	"github.com/hastyai/processor/service"
)

func main() {
	fmt.Println("Job processor starting up...")

	dbHostname := os.Getenv("DB_HOSTNAME")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	db, err := datastore.NewMongoDB(dbHostname, dbUser, dbPassword)
	if err != nil {
		log.Fatal("error connecting to mongo data store")
	}

	svc := service.NewService(db)

	queueHostname := os.Getenv("QUEUE_HOSTNAME")
	queueName := os.Getenv("QUEUE_NAME")
	queueURL := fmt.Sprintf("amqp://%s:%s@%s/", "guest", "guest", queueHostname)
	queue := queue.NewRabbitMQ(queueURL, queueName)

	queue.Consume(context.Background(), svc.ProcessJob)

	fmt.Println("Job processor shutting down...")
}
