package handler

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/go-redsync/redsync/v4"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/hastyai/backend/datastore"
	"github.com/hastyai/backend/queue"
)

type Handler struct {
	ds    datastore.DataStore
	queue queue.Queue
	rs    *redsync.Redsync
}

func NewHandler(ds datastore.DataStore, queue queue.Queue, rs *redsync.Redsync) *Handler {
	return &Handler{ds: ds, queue: queue, rs: rs}
}

type RequestDTO struct {
	ObjectID string `json:"object_id" validate:"required"`
}

type ResponseDTO struct {
	JobID  string `json:"job_id"`
	Status string `json:"status,omitempty"`
}

func (h *Handler) PostJob(w http.ResponseWriter, r *http.Request) {
	var reqBody RequestDTO
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error parsing request", http.StatusBadRequest)
		return
	}

	err = validator.New().Struct(&reqBody)
	if err != nil {
		log.Println(err)
		http.Error(w, "object_id is required", http.StatusBadRequest)
		return
	}

	mutex := h.rs.NewMutex(reqBody.ObjectID)
	if err := mutex.Lock(); err != nil {
		panic(err)
	}

	defer mutex.Unlock()

	job, err := h.ds.GetByObjectID(reqBody.ObjectID)
	if err != nil && err != datastore.ErrJobNotFound {
		log.Println(err)
		http.Error(w, "error fetching job", http.StatusInternalServerError)
		return
	}

	if job != nil && job.CreatedAt.Add(5*time.Minute).Before(time.Now()) {
		w.Header().Set("content-type", "application/json")
		json.NewEncoder(w).Encode(&ResponseDTO{JobID: job.JobID})
		return
	}

	newJob := &datastore.Job{
		ObjectID:  reqBody.ObjectID,
		JobID:     uuid.New().String(),
		Duration:  randInt(15, 25),
		Status:    "PROCESSING",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = h.ds.Save(*newJob)
	if err != nil {
		log.Println(err)
		http.Error(w, "error saving job", http.StatusInternalServerError)
		return
	}

	msg, err := json.Marshal(newJob)
	if err != nil {
		log.Println(err)
		http.Error(w, "error marshalling job", http.StatusInternalServerError)
		return
	}

	err = h.queue.Publish(context.Background(), msg)
	if err != nil {
		log.Println(err)
		http.Error(w, "error publishing job", http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(&ResponseDTO{JobID: newJob.JobID})
}

func (h *Handler) GetJobStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	jobID, ok := vars["job_id"]
	if !ok {
		http.Error(w, "Job ID is missing in path", http.StatusBadRequest)
		return
	}

	job, err := h.ds.GetByJobID(jobID)
	if err != nil && err != datastore.ErrJobNotFound {
		log.Println(err)
		http.Error(w, "error fetching job", http.StatusInternalServerError)
		return
	}

	if job == nil {
		log.Println(err)
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(job)
}

func randInt(lower, upper int) int {
	s1 := rand.NewSource(time.Now().UnixNano())
	rand.New(s1)

	return int((rand.Float32() * (float32(upper) - float32(lower))) + float32(lower))
}
