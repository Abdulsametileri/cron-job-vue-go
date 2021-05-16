package models

type JobStatus int

const (
	JobCancelled = iota + 1
	JobScheduled
)

type Job struct {
	Tag            string
	UserTelegramId int64
	ImageUrl       string
	RepeatType     string
	Time           string
}
