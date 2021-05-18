package main

import (
	"context"
	"embed"
	"fmt"
	"github.com/Abdulsametileri/cron-job-vue-go/config"
	"github.com/Abdulsametileri/cron-job-vue-go/controllers"
	"github.com/Abdulsametileri/cron-job-vue-go/database"
	"github.com/Abdulsametileri/cron-job-vue-go/infra/awsclient"
	"github.com/Abdulsametileri/cron-job-vue-go/infra/cronclient"
	"github.com/Abdulsametileri/cron-job-vue-go/infra/telegramclient"
	"github.com/Abdulsametileri/cron-job-vue-go/repository/jobrepo"
	"github.com/Abdulsametileri/cron-job-vue-go/repository/userrepo"
	"github.com/Abdulsametileri/cron-job-vue-go/services/jobservice"
	"github.com/Abdulsametileri/cron-job-vue-go/services/userservice"
	"github.com/spf13/viper"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

//go:embed client/dist
var clientFS embed.FS

func main() {
	config.Setup()

	distFS, err := fs.Sub(clientFS, "client/dist")
	if err != nil {
		log.Fatalf("error dist %v", err)
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
	if !config.IsDebug {
		go telegramClient.GetMessages()
	}

	cronClient := cronclient.NewCronClient(jobService, telegramClient)

	baseController := controllers.NewBaseController()
	tokenController := controllers.NewTokenController(baseController, userService)
	alarmController := controllers.NewAlarmController(baseController, userService, jobService, awsClient, telegramClient, cronClient)

	http.Handle("/", http.FileServer(http.FS(distFS)))

	http.HandleFunc("/api/validate-token", tokenController.ValidateToken)
	http.HandleFunc("/api/create-alarm", alarmController.CreateAlarm)
	http.HandleFunc("/api/list-alarm", alarmController.ListAlarm)
	http.HandleFunc("/api/delete-alarm", alarmController.DeleteAlarm)

	srv := &http.Server{
		Addr: ":3000",
	}
	fmt.Printf("Server is starting at %s \n", srv.Addr)
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(err.Error())
		}
	}()

	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	sig := <-c
	log.Println("Got signal:", sig)
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)

	_ = srv.Shutdown(ctx)
}
