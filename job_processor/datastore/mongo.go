package datastore

import (
	"context"
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

func (m *MongoDB) UpdateStatus(jobID, status string) error {
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)

	f := bson.M{"job_id": jobID}
	v := bson.M{"$set": bson.M{
		"status":     status,
		"updated_at": time.Now(),
	}}

	_, err := m.jobsCollection.UpdateOne(ctx, f, v)
	return err
}
