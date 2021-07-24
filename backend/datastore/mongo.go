package datastore

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/url"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoDB struct {
	jobsCollection *mongo.Collection
}

var ErrJobNotFound = errors.New("job not found")

type Job struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	ObjectID  string             `bson:"object_id" json:"object_id"`
	JobID     string             `bson:"job_id" json:"job_id"`
	Duration  int                `bson:"duration" json:"duration"`
	Status    string             `bson:"status" json:"status"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

func NewMongoDB(hostname, username, password string) (*MongoDB, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s", url.QueryEscape(username), url.QueryEscape(password), hostname)))
	if err != nil {
		log.Fatal(err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	jobs := client.Database("hastyai").Collection("jobs")
	return &MongoDB{jobsCollection: jobs}, nil
}

func (m *MongoDB) Save(msg Job) error {
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	_, err := m.jobsCollection.InsertOne(ctx, msg)
	return err
}

func (m *MongoDB) GetByObjectID(objectID string) (*Job, error) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	var msg Job

	f := bson.M{"object_id": objectID}
	opt := options.FindOne().SetSort(bson.D{{"created_at", -1}})
	err := m.jobsCollection.FindOne(ctx, f, opt).Decode(&msg)
	if err != nil {
		if err == mongo.ErrNoDocuments || err == mongo.ErrNilDocument || err == mongo.ErrNilValue {
			return nil, ErrJobNotFound
		}
		return nil, err
	}

	return &msg, nil
}

func (m *MongoDB) GetByJobID(jobID string) (*Job, error) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	var msg Job

	f := bson.M{"job_id": jobID}
	err := m.jobsCollection.FindOne(ctx, f).Decode(&msg)
	if err != nil {
		if err == mongo.ErrNoDocuments || err == mongo.ErrNilDocument || err == mongo.ErrNilValue {
			return nil, ErrJobNotFound
		}
		return nil, err
	}

	return &msg, nil
}
