package models

type Job struct {
	Tag            string `json:"tag" bson:"tag"`
	UserTelegramId int64  `json:"userTelegramId" bson:"userTelegramId"`
	ImageUrl       string `json:"imageUrl" bson:"imageUrl"`
	RepeatType     string `json:"repeatType" bson:"repeatType"`
	Time           string `json:"time" bson:"time"`
}
