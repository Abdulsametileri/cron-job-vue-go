package controllers

import (
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

func (client awsClient) UploadToS3(fileName, fileType string, file multipart.File) (string, error) {
	if fileName == "badFileName" {
		return "", ErrS3Upload
	}
	return fileName, nil
}

func (client awsClient) DeleteFileInS3(fileName string) error {
	if fileName == "error-scenario-with-s3" {
		return ErrDeleteFileS3
	}
	return nil
}

func (client awsClient) DetermineS3ImageUrl(fileName string) string {
	return "https://remindercron.s3.eu-central-1.amazonaws.com/Screen+Shot+2021-05-12+at+09.39.43.png"
}

type telegramClient struct{}

func (tc *telegramClient) GetMessages() {

}

func (tc *telegramClient) SendImage(telegramId int64, imageUrl string) error {
	return nil
}
