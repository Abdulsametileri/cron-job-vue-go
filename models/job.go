package models

type JobStatus int

const (
	JobValid = iota + 1
	JobDeleted
)

type Job struct {
	Tag            string    `json:"tag" bson:"tag"`
	Name           string    `json:"name" bson:"name"`
	UserTelegramId int64     `json:"userTelegramId" bson:"userTelegramId"`
	UserToken      string    `json:"userToken" bson:"userToken"`
	ImageUrl       string    `json:"imageUrl" bson:"imageUrl"`
	RepeatType     string    `json:"repeatType" bson:"repeatType"`
	Time           string    `json:"time" bson:"time"`
	Status         JobStatus `json:"status" bson:"status"`
}
