package models

type User struct {
	Token               string `json:"token" bson:"token"`
	TelegramId          int64  `json:"telegramId" bson:"telegramId"`
	TelegramDisplayName string `json:"telegramDisplayName" bson:"telegramDisplayName"`
}
