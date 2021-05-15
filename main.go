package main

import (
	"embed"
	"fmt"
	"github.com/Abdulsametileri/cron-job-vue-go/config"
	"github.com/Abdulsametileri/cron-job-vue-go/database"
	"github.com/Abdulsametileri/cron-job-vue-go/infra/awsclient"
	"github.com/Abdulsametileri/cron-job-vue-go/infra/telegramclient"
	"github.com/Abdulsametileri/cron-job-vue-go/repository/userrepo"
	"github.com/go-co-op/gocron"
	"github.com/spf13/viper"
	"io/fs"
	"log"
	"net/http"
	"time"
)

//go:embed client/dist
var clientFS embed.FS

var indexToWeekDay = map[string]time.Weekday{
	"0": time.Sunday,
	"1": time.Monday,
	"2": time.Tuesday,
	"3": time.Wednesday,
	"4": time.Thursday,
	"5": time.Friday,
	"6": time.Saturday,
}

func main() {
	config.Setup()

	awsClient := awsclient.NewAwsClient()
	_ = awsClient

	mongoClient := database.Setup()

	userCollection := mongoClient.Database(viper.GetString("MONGODB_DATABASE")).Collection("users")

	userRepo := userrepo.NewUserRepository(userCollection)

	telegramClient := telegramclient.NewTelegramClient(userRepo)
	go telegramClient.GetMessages()

	distFS, err := fs.Sub(clientFS, "client/dist")
	if err != nil {
		log.Fatal(err)
	}

	location, _ := time.LoadLocation("Europe/Istanbul")
	s := gocron.NewScheduler(location)
	s.StartAsync()

	http.Handle("/", http.FileServer(http.FS(distFS)))

	http.HandleFunc("/api/create-alarm", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Not allowed", http.StatusNotFound)
			return
		}

		uploadedFile, _, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Error while reading imgFile", http.StatusBadRequest)
			return
		}
		uploadedFileName := r.FormValue("fileName")
		uploadedFileType := r.FormValue("fileType")

		filePathOnS3, err := awsClient.UploadToS3(uploadedFileName, uploadedFileType, uploadedFile)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Println(filePathOnS3)
		/*
			name := r.Form.Get("name")
			if name == "" {
				http.Error(w, "You have to specify name", http.StatusBadRequest)
				return
			}

			gettime := r.Form.Get("time")
			if gettime == "00:00:00" {
				http.Error(w, "You have to specify time", http.StatusBadRequest)
				return
			}

			repeatType := r.Form.Get("repeatType")

			task := func() {
				fmt.Println("tetiklendi")
			}

			s.Every(1)
			s.Day()

			val, ok := indexToWeekDay[repeatType]
			if ok {
				s.Weekday(val)
			} else {
				// special
			}

			s.At(gettime)
			jb, err := s.Do(task)
			fmt.Println(jb)

			if err != nil {
				fmt.Println(err.Error())
			}*/
	})

	log.Println("Starting HTTP server at http://localhost:3000 ...")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
