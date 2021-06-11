package jobrepo

import (
	"context"
	"fmt"
	"github.com/Abdulsametileri/cron-job-vue-go/models"
	"github.com/Abdulsametileri/cron-job-vue-go/repository"
	"github.com/google/uuid"
	"github.com/magiconair/properties/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
)

func cleanCollection(t *testing.T, jobCollection *mongo.Collection) {
	errDrop := jobCollection.Drop(context.Background())
	require.NoError(t, errDrop)
}

func TestJobRepository_PaginateAllValidJobs(t *testing.T) {
	client, errSetupDB := repository.SetupDB()
	require.NoError(t, errSetupDB)

	jobCollection, errCollection := repository.SetupCollection(client, "jobs")
	require.NoError(t, errCollection)

	defer cleanCollection(t, jobCollection)

	ctx := context.Background()
	for i := 0; i < 100; i++ {
		var jobStatus models.JobStatus
		if i%2 == 0 {
			jobStatus = models.JobDeleted
		} else {
			jobStatus = models.JobValid
		}

		jobCollection.InsertOne(ctx, models.Job{
			Name:      fmt.Sprintf("job-%d", i+1),
			UserToken: "token",
			Status:    jobStatus,
		})
	}

	jobRepo := NewJobRepository(jobCollection)
	jobs, err := jobRepo.PaginateAllValidJobs(1, 10)
	require.NoError(t, err)
	require.Equal(t, len(jobs), 10)
	require.Equal(t, jobs[0].Name, "job-1")
	require.Equal(t, jobs[9].Name, "job-10")

	jobs, err = jobRepo.PaginateAllValidJobs(2, 5)
	require.NoError(t, err)
	require.Equal(t, len(jobs), 5)
	require.Equal(t, jobs[0].Name, "job-6")
	require.Equal(t, jobs[4].Name, "job-10")
}

func TestJobRepository_ListAllValidJobs(t *testing.T) {
	client, errSetupDB := repository.SetupDB()
	require.NoError(t, errSetupDB)

	jobCollection, errCollection := repository.SetupCollection(client, "jobs")
	require.NoError(t, errCollection)

	defer cleanCollection(t, jobCollection)

	job1 := models.Job{
		Tag:            uuid.NewString(),
		UserTelegramId: 123,
		UserToken:      "samet",
		ImageUrl:       "http://test",
		RepeatType:     "2",
		Time:           "18:55",
		Status:         models.JobValid,
	}
	_, errAddJob := jobCollection.InsertOne(context.Background(), job1)
	require.NoError(t, errAddJob)

	job1.Time = "13:55"
	job1.RepeatType = "4"
	job1.Tag = uuid.NewString()
	_, errAddJob = jobCollection.InsertOne(context.Background(), job1)
	require.NoError(t, errAddJob)

	job1.Time = "16:55"
	job1.Tag = uuid.NewString()
	job1.Status = models.JobDeleted
	_, errAddJob = jobCollection.InsertOne(context.Background(), job1)
	require.NoError(t, errAddJob)

	jobRepo := NewJobRepository(jobCollection)
	jobs, err := jobRepo.ListAllValidJobs()
	require.NoError(t, err)
	require.Equal(t, len(jobs), 2)
}

func TestJobRepository_ListAllValidJobsByToken_with2Items(t *testing.T) {
	client, errSetupDB := repository.SetupDB()
	require.NoError(t, errSetupDB)

	jobCollection, errCollection := repository.SetupCollection(client, "jobs")
	require.NoError(t, errCollection)

	defer cleanCollection(t, jobCollection)

	job1 := models.Job{
		Tag:            uuid.NewString(),
		UserTelegramId: 123,
		UserToken:      "samet",
		ImageUrl:       "http://test",
		RepeatType:     "1",
		Time:           "11:55",
		Status:         models.JobValid,
	}
	_, errAddJob := jobCollection.InsertOne(context.Background(), job1)
	require.NoError(t, errAddJob)

	job1.Time = "13:55"
	job1.Tag = uuid.NewString()
	job1.Status = models.JobValid
	_, errAddJob = jobCollection.InsertOne(context.Background(), job1)
	require.NoError(t, errAddJob)

	job1.Time = "16:55"
	job1.Tag = uuid.NewString()
	job1.Status = models.JobDeleted
	_, errAddJob = jobCollection.InsertOne(context.Background(), job1)
	require.NoError(t, errAddJob)

	jobRepo := NewJobRepository(jobCollection)
	jobs, err := jobRepo.ListAllValidJobsByToken(job1.UserToken)
	require.NoError(t, err)
	require.Equal(t, len(jobs), 2)
}

func TestJobRepository_ListJobsByToken_withNoItem(t *testing.T) {
	client, errSetupDB := repository.SetupDB()
	require.NoError(t, errSetupDB)

	jobCollection, errCollection := repository.SetupCollection(client, "jobs")
	require.NoError(t, errCollection)

	defer cleanCollection(t, jobCollection)

	jobRepo := NewJobRepository(jobCollection)
	jobs, err := jobRepo.ListAllValidJobsByToken("test")
	require.NoError(t, err)
	require.Equal(t, len(jobs), 0)
}

func TestJobRepository_AddJob(t *testing.T) {
	client, errSetupDB := repository.SetupDB()
	require.NoError(t, errSetupDB)

	jobCollection, errCollection := repository.SetupCollection(client, "jobs")
	require.NoError(t, errCollection)

	defer cleanCollection(t, jobCollection)

	jobToAdd := models.Job{
		Tag:            uuid.New().String(),
		UserTelegramId: 123,
		ImageUrl:       "http://test",
		RepeatType:     "1",
		Time:           "11:55",
	}

	jobRepo := NewJobRepository(jobCollection)
	err := jobRepo.AddJob(jobToAdd)
	require.NoError(t, err)
}

func TestJobRepository_GetJobByFields(t *testing.T) {
	client, errSetupDB := repository.SetupDB()
	require.NoError(t, errSetupDB)

	jobCollection, errCollection := repository.SetupCollection(client, "jobs")
	require.NoError(t, errCollection)

	defer cleanCollection(t, jobCollection)

	addedJob := models.Job{
		Tag:            uuid.New().String(),
		UserTelegramId: 123,
		ImageUrl:       "http://test",
		RepeatType:     "1",
		Time:           "11:55",
	}
	_, errAddJob := jobCollection.InsertOne(context.Background(), addedJob)
	require.NoError(t, errAddJob)

	jobRepo := NewJobRepository(jobCollection)

	fields := make(map[string]interface{}, 0)
	fields["userTelegramId"] = addedJob.UserTelegramId
	fields["imageUrl"] = addedJob.ImageUrl
	fields["repeatType"] = addedJob.RepeatType
	fields["time"] = addedJob.Time

	job, err := jobRepo.GetJobByFields(fields)
	require.NoError(t, err)

	assert.Equal(t, job.Tag, addedJob.Tag)
	assert.Equal(t, job.UserTelegramId, addedJob.UserTelegramId)
	assert.Equal(t, job.ImageUrl, addedJob.ImageUrl)
	assert.Equal(t, job.RepeatType, addedJob.RepeatType)
	assert.Equal(t, job.Time, addedJob.Time)
}

func TestJobRepository_GetJobByFields_NotExistedJob(t *testing.T) {
	client, errSetupDB := repository.SetupDB()
	require.NoError(t, errSetupDB)

	jobCollection, errCollection := repository.SetupCollection(client, "jobs")
	require.NoError(t, errCollection)

	defer cleanCollection(t, jobCollection)

	jobRepo := NewJobRepository(jobCollection)

	fields := make(map[string]interface{}, 0)
	fields["userTelegramId"] = 123
	fields["repeatType"] = "1"
	fields["time"] = "11:55"

	job, err := jobRepo.GetJobByFields(fields)
	require.NoError(t, err)

	emptyJob := models.Job{}
	assert.Equal(t, job.Tag, emptyJob.Tag)
	assert.Equal(t, job.UserTelegramId, emptyJob.UserTelegramId)
	assert.Equal(t, job.ImageUrl, emptyJob.ImageUrl)
	assert.Equal(t, job.RepeatType, emptyJob.RepeatType)
	assert.Equal(t, job.Time, emptyJob.Time)
}

func TestJobRepository_DeleteJobByTag(t *testing.T) {
	client, errSetupDB := repository.SetupDB()
	require.NoError(t, errSetupDB)

	jobCollection, errCollection := repository.SetupCollection(client, "jobs")
	require.NoError(t, errCollection)

	defer cleanCollection(t, jobCollection)

	addedJob := models.Job{
		Tag:            "123",
		Name:           "test",
		UserTelegramId: 123,
		UserToken:      "123",
		ImageUrl:       "http://s3..",
		RepeatType:     "1",
		Time:           "12:38",
		Status:         models.JobValid,
	}

	_, err := jobCollection.InsertOne(context.Background(), addedJob)
	require.NoError(t, err)

	jobRepository := NewJobRepository(jobCollection)

	err = jobRepository.DeleteJobByTag(addedJob.Tag)
	require.NoError(t, err)

	var jobFromDb models.Job
	jobCollection.FindOne(context.Background(), bson.M{"tag": addedJob.Tag}).Decode(&jobFromDb)

	require.Equal(t, jobFromDb.Tag, addedJob.Tag)
	require.Equal(t, int(jobFromDb.Status), models.JobDeleted)
}
