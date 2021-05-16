package jobrepo

import (
	"context"
	"github.com/Abdulsametileri/cron-job-vue-go/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Repo interface {
	AddJob(models.Job) error
	GetJobByFields(map[string]interface{}) (models.Job, error)
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

func (j jobRepository) GetJobByFields(fields map[string]interface{}) (models.Job, error) {
	filter := bson.M{}
	for key, val := range fields {
		filter[key] = val
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var job models.Job
	err := j.collection.FindOne(ctx, filter).Decode(&job)
	if err == mongo.ErrNoDocuments {
		return models.Job{}, nil
	}

	return job, err
}
