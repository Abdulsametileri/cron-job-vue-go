package jobrepo

import (
	"context"
	"github.com/Abdulsametileri/cron-job-vue-go/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Repo interface {
	ListJobsByToken(token string) ([]models.Job, error)
	AddJob(models.Job) error
	GetJobByFields(map[string]interface{}) (models.Job, error)
}

type jobRepository struct {
	collection *mongo.Collection
}

func NewJobRepository(collection *mongo.Collection) Repo {
	return &jobRepository{collection: collection}
}

func (j jobRepository) ListJobsByToken(token string) ([]models.Job, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"_id", -1}})
	filter := bson.M{"userToken": token}
	cur, err := j.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return make([]models.Job, 0), err
	}
	defer cur.Close(ctx)
	var jobs []models.Job
	if err = cur.All(ctx, &jobs); err != nil {
		return make([]models.Job, 0), err
	}
	return jobs, nil
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
