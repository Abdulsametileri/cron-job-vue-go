package jobrepo

import (
	"context"
	"errors"
	"github.com/Abdulsametileri/cron-job-vue-go/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Repo interface {
	ListAllValidJobs() ([]models.Job, error)
	PaginateAllValidJobs(pageNo, pageSize int) ([]models.Job, error)
	ListAllValidJobsByToken(token string) ([]models.Job, error)
	AddJob(models.Job) error
	GetJobByFields(map[string]interface{}) (models.Job, error)
	DeleteJobByTag(tag string) error
}

type jobRepository struct {
	collection *mongo.Collection
}

func NewJobRepository(collection *mongo.Collection) Repo {
	return &jobRepository{collection: collection}
}

func (j jobRepository) ListAllValidJobs() ([]models.Job, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cur, err := j.collection.Find(ctx, bson.M{
		"status": models.JobValid,
	})
	if err != nil {
		return make([]models.Job, 0), err
	}
	defer cur.Close(ctx)
	var jobs []models.Job
	if err = cur.All(ctx, &jobs); err != nil {
		return make([]models.Job, 0), err
	}
	if jobs == nil {
		return make([]models.Job, 0), nil
	}
	return jobs, nil
}

func (j *jobRepository) PaginateAllValidJobs(pageNo, pageSize int) ([]models.Job, error) {
	limit := int64(pageSize)
	skip := int64(pageNo-1) * limit

	cur, err := j.collection.Find(context.Background(), bson.D{}, &options.FindOptions{
		Skip:  &skip,
		Limit: &limit,
	})
	if err != nil {
		return []models.Job{}, err
	}
	defer cur.Close(context.Background())

	var jobs []models.Job
	if err := cur.All(context.Background(), &jobs); err != nil {
		return []models.Job{}, err
	}

	return jobs, nil
}

func (j jobRepository) ListAllValidJobsByToken(token string) ([]models.Job, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"_id", -1}})
	filter := bson.M{"userToken": token, "status": models.JobValid}
	cur, err := j.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return make([]models.Job, 0), err
	}
	defer cur.Close(ctx)
	var jobs []models.Job
	if err = cur.All(ctx, &jobs); err != nil {
		return make([]models.Job, 0), err
	}
	if jobs == nil {
		return make([]models.Job, 0), nil
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

func (j jobRepository) DeleteJobByTag(tag string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.M{
		"tag": tag,
	}
	update := bson.D{{"$set", bson.D{{"status", models.JobDeleted}}}}

	up, err := j.collection.UpdateOne(ctx, filter, update)
	if up != nil && up.MatchedCount == 0 {
		return errors.New("Job does not found with the given tag")
	}
	return err
}
