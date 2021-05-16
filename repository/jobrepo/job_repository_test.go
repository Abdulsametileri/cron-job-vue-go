package jobrepo

import (
	"context"
	"fmt"
	"github.com/Abdulsametileri/cron-job-vue-go/models"
	"github.com/Abdulsametileri/cron-job-vue-go/repository"
	"github.com/google/uuid"
	"github.com/magiconair/properties/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
)

func cleanCollection(t *testing.T, jobCollection *mongo.Collection) {
	errDrop := jobCollection.Drop(context.Background())
	require.NoError(t, errDrop)
}

func TestJobRepository_ListJobsByToken_with3Items(t *testing.T) {
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
	}
	_, errAddJob := jobCollection.InsertOne(context.Background(), job1)
	require.NoError(t, errAddJob)

	job1.Time = "13:55"
	job1.Tag = uuid.NewString()
	_, errAddJob = jobCollection.InsertOne(context.Background(), job1)
	require.NoError(t, errAddJob)

	job1.Time = "16:55"
	job1.Tag = uuid.NewString()
	_, errAddJob = jobCollection.InsertOne(context.Background(), job1)
	require.NoError(t, errAddJob)

	jobRepo := NewJobRepository(jobCollection)
	jobs, err := jobRepo.ListJobsByToken(job1.UserToken)
	require.NoError(t, err)
	fmt.Println(jobs)
}

func TestJobRepository_ListJobsByToken_withNoItem(t *testing.T) {
	client, errSetupDB := repository.SetupDB()
	require.NoError(t, errSetupDB)

	jobCollection, errCollection := repository.SetupCollection(client, "jobs")
	require.NoError(t, errCollection)

	defer cleanCollection(t, jobCollection)

	jobRepo := NewJobRepository(jobCollection)
	jobs, err := jobRepo.ListJobsByToken("test")
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
