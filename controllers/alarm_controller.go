package controllers

import (
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
		http.Error(w, "Not allowed http method", http.StatusNotFound)
		return
	}

	token := r.FormValue("token")
	name := r.FormValue("name")
	gettime := r.FormValue("time")
	repeatType := r.FormValue("repeatType")
	uploadedFile, _, err := r.FormFile("file")
	uploadedFileName := r.FormValue("fileName")
	uploadedFileType := r.FormValue("fileType")

	if token == "" {
		http.Error(w, "You have to specify token", http.StatusBadRequest)
		return
	}

	if name == "" {
		http.Error(w, "You have to specify name", http.StatusBadRequest)
		return
	}

	if gettime == "00:00:00" {
		http.Error(w, "You have to specify time", http.StatusBadRequest)
		return
	}

	if repeatType == "-1" {
		http.Error(w, "You have to specify repeat time", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, "Error while reading image file", http.StatusBadRequest)
		return
	}

	user, err := ac.us.GetUserByToken(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if user.TelegramId == 0 {
		http.Error(w, "Invalid token", http.StatusBadRequest)
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
