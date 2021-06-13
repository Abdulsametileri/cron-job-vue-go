package controllers

import (
	"fmt"
	"github.com/Abdulsametileri/cron-job-vue-go/models"
	"mime/multipart"
)

type userSvc struct{}

func (uc *userSvc) AddUser(user models.User) error {
	return nil
}

func (uc *userSvc) GetUserByTelegramId(telegramId int64) (models.User, error) {
	return models.User{}, nil
}

func (uc *userSvc) GetUserByToken(token string) (models.User, error) {
	if token == "db-err" {
		return models.User{}, ErrDb
	}

	if token == "not-exist-token" {
		return models.User{}, nil
	}

	if token == "job-already-exist-db-error" {
		return models.User{TelegramId: -1}, nil
	}

	if token == "job-already-exist" {
		return models.User{
			Token:               "job-already-exist",
			TelegramId:          123,
			TelegramDisplayName: "ileri4s",
		}, nil
	}

	return models.User{
		Token:               "x1pnwjjkhj3o",
		TelegramId:          12314,
		TelegramDisplayName: "ileri4s",
	}, nil
}

type jobSvc struct{}

func (js *jobSvc) GetNumberOfValidJobs() (int, error) {
	return 100, nil
}

func (js *jobSvc) PaginateAllValidJobs(pageNo, pageSize int) ([]models.Job, error) {
	if pageNo == 1 {
		res := make([]models.Job, 0)
		for i := 0; i < pageSize; i++ {
			res = append(res, models.Job{
				Tag:    fmt.Sprintf("tag-%d", i),
				Name:   fmt.Sprintf("name-%d", i),
				Status: models.JobValid,
			})
		}
		return res, nil
	}

	return []models.Job{}, nil
}

func (js *jobSvc) ListAllValidJobsByToken(token string) ([]models.Job, error) {
	if token == "job-list-err" {
		return nil, ErrGettingJobList
	}
	if token == "three-job-list-item" {
		return make([]models.Job, 3), nil
	}
	return make([]models.Job, 0), nil
}

func (js *jobSvc) DeleteJobByTag(tag string) error {
	if tag == "error-tag" {
		return ErrJobDelete
	}
	return nil
}

func (js *jobSvc) ListAllValidJobs() ([]models.Job, error) {
	return make([]models.Job, 0), nil
}

func (js *jobSvc) AddJob(job models.Job) error {
	if job.ImageUrl == "error-scenario-with-s3" {
		return ErrAddingJob
	}
	if job.ImageUrl == "error-scenario-job" {
		return ErrAddingJob
	}
	return nil
}

func (js *jobSvc) GetJobByFields(fields map[string]interface{}) (models.Job, error) {
	if fields["userTelegramId"].(int64) == -1 {
		return models.Job{}, ErrGettingJob
	}
	if fields["userTelegramId"].(int64) == 123 {
		return models.Job{
			Tag: "arbitrary",
		}, nil
	}
	return models.Job{}, nil
}

type awsClient struct{}

func (client awsClient) UploadToS3(userId int64, fileName, fileType string, file multipart.File) (string, error) {
	if fileName == "badFileName" {
		return "", ErrS3Upload
	}
	return fileName, nil
}

func (client awsClient) DeleteFileInS3(fileUrl string) error {
	if fileUrl == "error-scenario-with-s3" {
		return ErrDeleteFileS3
	}
	return nil
}

func (client awsClient) DetermineS3ImageUrl(userId int64, fileName string) string {
	return "https://remindercron.s3.eu-central-1.amazonaws.com/Screen+Shot+2021-05-12+at+09.39.43.png"
}

type telegramClient struct{}

func (tc *telegramClient) SendMessageForDebug(msg string) {

}

func (tc *telegramClient) GetMessages() {

}

func (tc *telegramClient) SendImage(telegramId int64, imageUrl string) error {
	return nil
}
