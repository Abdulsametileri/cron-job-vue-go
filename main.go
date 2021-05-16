package main

import (
	"embed"
	"github.com/Abdulsametileri/cron-job-vue-go/config"
	"github.com/Abdulsametileri/cron-job-vue-go/controllers"
	"github.com/Abdulsametileri/cron-job-vue-go/database"
	"github.com/Abdulsametileri/cron-job-vue-go/infra/awsclient"
	"github.com/Abdulsametileri/cron-job-vue-go/infra/telegramclient"
	"github.com/Abdulsametileri/cron-job-vue-go/repository/jobrepo"
	"github.com/Abdulsametileri/cron-job-vue-go/repository/userrepo"
	"github.com/Abdulsametileri/cron-job-vue-go/services/jobservice"
	"github.com/Abdulsametileri/cron-job-vue-go/services/userservice"
	"github.com/go-co-op/gocron"
	"github.com/spf13/viper"
	"io/fs"
	"log"
	"net/http"
	"time"
)

//go:embed client/dist
var clientFS embed.FS

func main() {
	config.Setup()

	distFS, err := fs.Sub(clientFS, "client/dist")
	if err != nil {
		log.Fatal(err)
	}

	mongoClient := database.Setup()
	mongodb := mongoClient.Database(viper.GetString("MONGODB_DATABASE"))
	userCollection := mongodb.Collection("users")
	jobCollection := mongodb.Collection("jobs")
	userRepo := userrepo.NewUserRepository(userCollection)
	jobRepo := jobrepo.NewJobRepository(jobCollection)

	userService := userservice.NewUserService(userRepo)
	jobService := jobservice.NewJobService(jobRepo)

	awsClient := awsclient.NewAwsClient()

	telegramClient := telegramclient.NewTelegramClient(userService)
	go telegramClient.GetMessages()

	location, _ := time.LoadLocation("Europe/Istanbul")
	schedule := gocron.NewScheduler(location)
	schedule.StartAsync()

	alarmController := controllers.NewAlarmController(userService, jobService, awsClient, telegramClient, schedule)

	http.Handle("/", http.FileServer(http.FS(distFS)))

	http.HandleFunc("/api/create-alarm", alarmController.CreateAlarm)

	log.Println("Starting HTTP server at http://localhost:3000 ...")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
