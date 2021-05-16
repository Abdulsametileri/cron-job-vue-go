package controllers

import (
	"errors"
	"fmt"
	"github.com/Abdulsametileri/cron-job-vue-go/infra/awsclient"
	"github.com/Abdulsametileri/cron-job-vue-go/infra/telegramclient"
	"github.com/Abdulsametileri/cron-job-vue-go/models"
	"github.com/Abdulsametileri/cron-job-vue-go/services/jobservice"
	"github.com/Abdulsametileri/cron-job-vue-go/services/userservice"
	"github.com/go-co-op/gocron"
	"github.com/google/uuid"
	"net/http"
	"time"
)

var IndexToWeekDay = map[string]time.Weekday{
	"0": time.Sunday,
	"1": time.Monday,
	"2": time.Tuesday,
	"3": time.Wednesday,
	"4": time.Thursday,
	"5": time.Friday,
	"6": time.Saturday,
}

var (
	ErrMethodNotAllowed   = errors.New("Method not allowed")
	ErrTokenNotFound      = errors.New("You have to specify token")
	ErrNameNotFound       = errors.New("You have to specify name")
	ErrTimeNotFound       = errors.New("You have to specify time")
	ErrRepeatTypeNotFound = errors.New("You have to specify repeat time")
	ErrReadingFile        = errors.New("Error while reading image file")
	ErrDb                 = errors.New("DB error occured")
	ErrTokenDoesNotExist  = errors.New("Token does not exist")
	ErrS3Upload           = errors.New("Error uploading file to s3")
	ErrDeleteFileS3       = errors.New("Error deleting file in s3")
	ErrAddingJob          = errors.New("Error adding job to Db")
	ErrGettingJob         = errors.New("Error getting job in DB")
	ErrJobAlreadyExist    = errors.New("Error you already create your job before")
)

type AlarmController interface {
	CreateAlarm(http.ResponseWriter, *http.Request)
}

type alarmController struct {
	bc  BaseController
	us  userservice.UserService
	js  jobservice.JobService
	aws awsclient.AwsClient
	tc  telegramclient.TelegramClient
	sch *gocron.Scheduler
}

func NewAlarmController(
	bc BaseController,
	us userservice.UserService, js jobservice.JobService, aws awsclient.AwsClient,
	tc telegramclient.TelegramClient,
	sch *gocron.Scheduler) AlarmController {
	return &alarmController{
		bc:  bc,
		us:  us,
		js:  js,
		aws: aws,
		tc:  tc,
		sch: sch,
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
		ac.bc.Error(w, http.StatusBadRequest, ErrTokenNotFound)
		return
	}

	if name == "" {
		ac.bc.Error(w, http.StatusBadRequest, ErrNameNotFound)
		return
	}

	if gettime == "" {
		ac.bc.Error(w, http.StatusBadRequest, ErrTimeNotFound)
		return
	}

	if repeatType == "" || repeatType == "-1" {
		ac.bc.Error(w, http.StatusBadRequest, ErrRepeatTypeNotFound)
		return
	}

	if errFile != nil {
		ac.bc.Error(w, http.StatusBadRequest, ErrReadingFile)
		return
	}

	user, err := ac.us.GetUserByToken(token)
	if err != nil {
		ac.bc.Error(w, http.StatusBadRequest, ErrDb)
		return
	}

	if user.TelegramId == 0 {
		ac.bc.Error(w, http.StatusBadRequest, ErrTokenDoesNotExist)
		return
	}

	searchByFields := make(map[string]interface{})
	searchByFields["userTelegramId"] = user.TelegramId
	searchByFields["imageUrl"] = ac.aws.DetermineS3ImageUrl(user.TelegramId, uploadedFileName)
	searchByFields["repeatType"] = repeatType
	searchByFields["time"] = gettime

	job, err := ac.js.GetJobByFields(searchByFields)
	if err != nil {
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

	err = ac.js.AddJob(models.Job{
		Tag:            uuid.New().String(),
		UserTelegramId: user.TelegramId,
		ImageUrl:       filePathOnS3,
		RepeatType:     repeatType,
		Time:           gettime,
	})
	if err != nil {
		err = ac.aws.DeleteFileInS3(filePathOnS3)
		if err != nil {
			ac.bc.Error(w, http.StatusBadRequest, ErrDeleteFileS3)
			return
		}
		ac.bc.Error(w, http.StatusBadRequest, ErrAddingJob)
		return
	}

	ac.sch.Every(1)
	val, ok := IndexToWeekDay[repeatType]
	if ok {
		ac.sch.Day().Weekday(val)
	} else {
		ac.sch.Days()
	}
	ac.sch.At(gettime)

	_, err = ac.sch.Do(func() {
		err := ac.tc.SendImage(user.TelegramId, filePathOnS3)
		fmt.Println(err)
	})

	fmt.Println(err)
}
