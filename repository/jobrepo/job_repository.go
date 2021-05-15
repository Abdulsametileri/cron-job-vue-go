package jobrepo

import (
	"context"
	"github.com/Abdulsametileri/cron-job-vue-go/models"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Repo interface {
	AddJob(job models.Job) error
}

type jobRepository struct {
	collection *mongo.Collection
}

func NewJobRepository(collection *mongo.Collection) Repo {
	return &jobRepository{collection: collection}
}

func (j jobRepository) AddJob(job models.Job) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := j.collection.InsertOne(ctx, job)
	return err
}
