package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	goredislib "github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
	"github.com/gorilla/mux"

	"github.com/hastyai/backend/datastore"
	"github.com/hastyai/backend/handler"
	"github.com/hastyai/backend/queue"
)

func main() {
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	dbHostname := os.Getenv("DB_HOSTNAME")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	db, err := datastore.NewMongoDB(dbHostname, dbUser, dbPassword)
	if err != nil {
		log.Fatal("error connecting to mongo data store")
	}

	queueHostname := os.Getenv("QUEUE_HOSTNAME")
	queueName := os.Getenv("QUEUE_NAME")
	queueURL := fmt.Sprintf("amqp://%s:%s@%s/", "guest", "guest", queueHostname)
	queue := queue.NewRabbitMQ(queueURL, queueName)

	client := goredislib.NewClient(&goredislib.Options{
		Addr: os.Getenv("REDIS_URL"),
	})
	pool := goredis.NewPool(client)
	rs := redsync.New(pool)

	handler := handler.NewHandler(db, queue, rs)

	r := mux.NewRouter()
	s := r.PathPrefix("/api/v1").Subrouter()

	s.HandleFunc("/job", handler.PostJob).Methods(http.MethodPost, http.MethodPut)
	s.HandleFunc("/job/{job_id}", handler.GetJobStatus).Methods(http.MethodGet)

	port := os.Getenv("PORT")
	addr := fmt.Sprintf("0.0.0.0:%s", port)

	srv := &http.Server{
		Addr: addr,
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 30,
		IdleTimeout:  time.Second * 120,
		Handler:      r, // Pass our instance of gorilla/mux in.
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)

	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)
}
