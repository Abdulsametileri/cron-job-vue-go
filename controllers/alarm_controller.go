package controllers

import (
	"errors"
	"fmt"
	"github.com/Abdulsametileri/cron-job-vue-go/infra/awsclient"
	"github.com/Abdulsametileri/cron-job-vue-go/infra/cronclient"
	"github.com/Abdulsametileri/cron-job-vue-go/infra/telegramclient"
	"github.com/Abdulsametileri/cron-job-vue-go/models"
	"github.com/Abdulsametileri/cron-job-vue-go/services/jobservice"
	"github.com/Abdulsametileri/cron-job-vue-go/services/userservice"
	"github.com/google/uuid"
	"net/http"
)

var (
	ErrMethodNotAllowed       = errors.New("Method not allowed")
	ErrFieldNotFound          = func(field string) error { return fmt.Errorf("You have to specify %s", field) }
	ErrReadingFile            = errors.New("Error while reading image file")
	ErrDb                     = errors.New("DB error occured")
	ErrTokenDoesNotExistInUrl = errors.New("Token does not exist")
	ErrS3Upload               = errors.New("Error uploading file to s3")
	ErrDeleteFileS3           = errors.New("Error deleting file in s3")
	ErrAddingJob              = errors.New("Error adding job to Db")
	ErrGettingJob             = errors.New("Error getting job in DB")
	ErrJobAlreadyExist        = errors.New("Error you already create your job before")
	ErrJobDelete              = errors.New("Error deleting the speficied tag job in db")
	ErrUserDoesNotExist       = errors.New("Error user does not exist")
	ErrGettingJobList         = errors.New("Error getting the job list")
	ErrTagDoesNotExistInUrl   = errors.New("Tag does not exist in the url")
)

type AlarmController interface {
	CreateAlarm(http.ResponseWriter, *http.Request)
	ListAlarm(http.ResponseWriter, *http.Request)
	DeleteAlarm(http.ResponseWriter, *http.Request)
}

type alarmController struct {
	bc  BaseController
	us  userservice.UserService
	js  jobservice.JobService
	aws awsclient.AwsClient
	tc  telegramclient.TelegramClient
	cc  cronclient.CronClient
}

func NewAlarmController(
	bc BaseController,
	us userservice.UserService, js jobservice.JobService, aws awsclient.AwsClient,
	tc telegramclient.TelegramClient,
	cc cronclient.CronClient) AlarmController {
	return &alarmController{
		bc:  bc,
		us:  us,
		js:  js,
		aws: aws,
		tc:  tc,
		cc:  cc,
	}
}

func (ac alarmController) CreateAlarm(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		ac.bc.Error(w, http.StatusNotFound, ErrMethodNotAllowed)
		return
	}

	token := r.FormValue("token")
	name := r.FormValue("name")
	gettime := r.FormValue("time")
	repeatType := r.FormValue("repeatType")
	uploadedFile, _, errFile := r.FormFile("file")
	uploadedFileName := r.FormValue("fileName")
	uploadedFileType := r.FormValue("fileType")

	if token == "" {
		ac.bc.Error(w, http.StatusBadRequest, ErrFieldNotFound("token"))
		return
	}

	if name == "" {
		ac.bc.Error(w, http.StatusBadRequest, ErrFieldNotFound("name"))
		return
	}

	if gettime == "" {
		ac.bc.Error(w, http.StatusBadRequest, ErrFieldNotFound("time"))
		return
	}

	if repeatType == "" || repeatType == "-1" {
		ac.bc.Error(w, http.StatusBadRequest, ErrFieldNotFound("repeat time"))
		return
	}

	if errFile != nil {
		ac.bc.Error(w, http.StatusBadRequest, ErrReadingFile)
		return
	}

	user, err := ac.us.GetUserByToken(token)
	if err != nil {
		fmt.Println(err)
		ac.bc.Error(w, http.StatusBadRequest, ErrDb)
		return
	}

	if user.TelegramId == 0 {
		ac.bc.Error(w, http.StatusBadRequest, ErrTokenDoesNotExistInUrl)
		return
	}

	searchByFields := make(map[string]interface{})
	searchByFields["userTelegramId"] = user.TelegramId
	searchByFields["imageUrl"] = ac.aws.DetermineS3ImageUrl(user.TelegramId, uploadedFileName)
	searchByFields["repeatType"] = repeatType
	searchByFields["time"] = gettime

	job, err := ac.js.GetJobByFields(searchByFields)
	if err != nil {
		fmt.Println(err)
		ac.bc.Error(w, http.StatusBadRequest, ErrGettingJob)
		return
	}

	if job.Tag != "" {
		ac.bc.Error(w, http.StatusBadRequest, ErrJobAlreadyExist)
		return
	}

	filePathOnS3, err := ac.aws.UploadToS3(user.TelegramId, uploadedFileName, uploadedFileType, uploadedFile)
	if err != nil {
		ac.bc.Error(w, http.StatusBadRequest, ErrS3Upload)
		return
	}

	addedJob := models.Job{
		Tag:            uuid.New().String(),
		Name:           name,
		UserTelegramId: user.TelegramId,
		UserToken:      user.Token,
		ImageUrl:       filePathOnS3,
		RepeatType:     repeatType,
		Time:           gettime,
		Status:         models.JobValid,
	}
	err = ac.js.AddJob(addedJob)

	if err != nil {
		err = ac.aws.DeleteFileInS3(filePathOnS3)
		if err != nil {
			ac.bc.Error(w, http.StatusBadRequest, ErrDeleteFileS3)
			return
		}
		ac.bc.Error(w, http.StatusBadRequest, ErrAddingJob)
		return
	}

	_ = ac.cc.Schedule(addedJob)

	ac.bc.Data(w, http.StatusOK, nil, "")
}

func (ac alarmController) ListAlarm(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		ac.bc.Error(w, http.StatusNotFound, ErrMethodNotAllowed)
		return
	}

	token := r.URL.Query().Get("token")
	if token == "" {
		ac.bc.Error(w, http.StatusBadRequest, ErrTokenDoesNotExistInUrl)
		return
	}

	user, err := ac.us.GetUserByToken(token)
	if err != nil {
		ac.bc.Error(w, http.StatusBadRequest, err)
		return
	}
	if user.TelegramId == 0 {
		ac.bc.Error(w, http.StatusBadRequest, ErrUserDoesNotExist)
		return
	}

	jobs, err := ac.js.ListAllValidJobsByToken(token)
	if err != nil {
		ac.bc.Error(w, http.StatusBadRequest, ErrGettingJobList)
		return
	}

	ac.bc.Data(w, http.StatusOK, jobs, "")
}

func (ac alarmController) DeleteAlarm(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		ac.bc.Error(w, http.StatusNotFound, ErrMethodNotAllowed)
		return
	}

	tag := r.URL.Query().Get("tag")
	if tag == "" {
		ac.bc.Error(w, http.StatusBadRequest, ErrTagDoesNotExistInUrl)
		return
	}

	err := ac.js.DeleteJobByTag(tag)
	if err != nil {
		ac.bc.Error(w, http.StatusBadRequest, ErrJobDelete)
		return
	}

	err = ac.cc.RemoveJobByTag(tag)
	if err != nil {
		ac.tc.SendMessageForDebug(fmt.Sprintf("error removing job by tag %s %s", tag, err.Error()))
		return
	}

	ac.bc.Data(w, http.StatusOK, nil, "")
}
