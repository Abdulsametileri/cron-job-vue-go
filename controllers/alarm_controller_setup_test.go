package controllers

import (
	"errors"
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
		return models.User{}, errors.New("some error occured")
	}

	if token == "not-exist-token" {
		return models.User{}, nil
	}

	return models.User{
		Token:               "x1pnwjjkhj3o",
		TelegramId:          12314,
		TelegramDisplayName: "ileri4s",
	}, nil
}

type awsClient struct{}

func (client awsClient) UploadToS3(fileName, fileType string, file multipart.File) (string, error) {
	return "https://remindercron.s3.eu-central-1.amazonaws.com/Screen+Shot+2021-05-12+at+09.39.43.png", nil
}
