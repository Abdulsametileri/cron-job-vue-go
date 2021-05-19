package cronclient

import (
	"github.com/Abdulsametileri/cron-job-vue-go/infra/telegramclient"
	"github.com/Abdulsametileri/cron-job-vue-go/models"
	"github.com/Abdulsametileri/cron-job-vue-go/services/jobservice"
	"github.com/Abdulsametileri/cron-job-vue-go/utils"
	"github.com/go-co-op/gocron"
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

type CronClient interface {
	Schedule(job models.Job) error
	RemoveJobByTag(tag string) error
}

type cronClient struct {
	js  jobservice.JobService
	tc  telegramclient.TelegramClient
	sch *gocron.Scheduler
}

func NewCronClient(js jobservice.JobService, tc telegramclient.TelegramClient) CronClient {
	location, _ := time.LoadLocation("Europe/Istanbul")
	scheduleClient := gocron.NewScheduler(location)
	scheduleClient.StartAsync()

	cronClient := &cronClient{
		js:  js,
		tc:  tc,
		sch: scheduleClient,
	}

	jobs, _ := js.ListAllValidJobs()
	if len(jobs) > 0 {
		for _, job := range jobs {
			_ = cronClient.Schedule(job)
		}
	}

	return cronClient
}

func (c cronClient) RemoveJobByTag(tag string) error {
	return c.sch.RemoveByTag(tag)
}

func (c cronClient) Schedule(job models.Job) error {
	c.sch.Every(1)
	val, ok := IndexToWeekDay[job.RepeatType]
	if ok {
		c.sch.Day().Weekday(val)
	} else {
		c.sch.Days()
	}
	c.sch.At(job.Time)

	scheduledJob, err := c.sch.Do(func() {
		err := c.tc.SendImage(job.UserTelegramId, job.ImageUrl)
		if err != nil {
			c.tc.SendMessageForDebug("Schedule Job Telegram Send Image" + err.Error())
		}
	})

	if err != nil {
		prettyScheduledJob, _ := utils.PrettyPrint(scheduledJob)
		prettyDbJob, _ := utils.PrettyPrint(job)
		c.tc.SendMessageForDebug("Cron Schedule " + prettyScheduledJob + " models.Job=" + prettyDbJob + err.Error())
		return err
	}
	scheduledJob.Tag(job.Tag)

	return nil
}
