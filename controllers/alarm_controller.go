package controllers

import (
	"errors"
	"github.com/Abdulsametileri/cron-job-vue-go/infra/awsclient"
	"github.com/Abdulsametileri/cron-job-vue-go/services/userservice"
	"github.com/go-co-op/gocron"
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
)

type AlarmController interface {
	CreateAlarm(http.ResponseWriter, *http.Request)
}

type alarmController struct {
	us  userservice.UserService
	aws awsclient.AwsClient
	sch *gocron.Scheduler
}

func NewAlarmController(us userservice.UserService, aws awsclient.AwsClient, sch *gocron.Scheduler) AlarmController {
	return &alarmController{
		us:  us,
		aws: aws,
		sch: sch,
	}
}

func (ac alarmController) CreateAlarm(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, ErrMethodNotAllowed.Error(), http.StatusNotFound)
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
		http.Error(w, ErrTokenNotFound.Error(), http.StatusBadRequest)
		return
	}

	if name == "" {
		http.Error(w, ErrNameNotFound.Error(), http.StatusBadRequest)
		return
	}

	if gettime == "" {
		http.Error(w, ErrTimeNotFound.Error(), http.StatusBadRequest)
		return
	}

	if repeatType == "" || repeatType == "-1" {
		http.Error(w, ErrRepeatTypeNotFound.Error(), http.StatusBadRequest)
		return
	}

	if errFile != nil {
		http.Error(w, ErrReadingFile.Error(), http.StatusBadRequest)
		return
	}

	user, err := ac.us.GetUserByToken(token)
	if err != nil {
		http.Error(w, ErrDb.Error(), http.StatusBadRequest)
		return
	}
	if user.TelegramId == 0 {
		http.Error(w, ErrTokenDoesNotExist.Error(), http.StatusBadRequest)
		return
	}

	_, _, _ = uploadedFile, uploadedFileName, uploadedFileType

	/*
		filePathOnS3, err := ac.awsClient.UploadToS3(uploadedFileName, uploadedFileType, uploadedFile)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Println(filePathOnS3)

		ac.sch.Every(1)

		val, ok := IndexToWeekDay[repeatType]
		if ok {
			ac.sch.Day().Weekday(val)
		} else {
			ac.sch.Days()
		}
		ac.sch.At(gettime)

		_, err = ac.sch.Do(func() {

		})

		if err != nil {

		}*/
}
