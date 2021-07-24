package service

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/hastyai/processor/datastore"
)

type Service struct {
	ds datastore.DataStore
}

func NewService(ds datastore.DataStore) *Service {
	return &Service{
		ds: ds,
	}
}

func (s *Service) ProcessJob(ctx context.Context, msg []byte) error {
	ctx, cancel := context.WithTimeout(ctx, 22*time.Second) // Preset tasks to run for a maximum of 22 seconds
	defer cancel()

	var message datastore.Job
	err := json.Unmarshal(msg, &message)
	if err != nil {
		log.Println("Error decoding message")
		return err
	}

	select {
	case <-time.After(time.Second * time.Duration(message.Duration)):
		log.Println("Successful...")
		return s.ds.UpdateStatus(message.JobID, "SUCCESSFUL")
	case <-ctx.Done():
		log.Println("Cancelled...")
		return s.ds.UpdateStatus(message.JobID, "CANCELLED")
	}
}
