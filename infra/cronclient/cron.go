package cronclient

import (
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
	Schedule(repeatType, gettime string, task func()) error
}

type cronClient struct {
	sch *gocron.Scheduler
}

func NewCronClient() CronClient {
	location, _ := time.LoadLocation("Europe/Istanbul")
	scheduleClient := gocron.NewScheduler(location)
	scheduleClient.StartAsync()

	return &cronClient{
		sch: scheduleClient,
	}
}

func (c cronClient) Schedule(repeatType, gettime string, task func()) error {
	c.sch.Every(1)
	val, ok := IndexToWeekDay[repeatType]
	if ok {
		c.sch.Day().Weekday(val)
	} else {
		c.sch.Days()
	}
	c.sch.At(gettime)

	_, err := c.sch.Do(task)
	return err
}
