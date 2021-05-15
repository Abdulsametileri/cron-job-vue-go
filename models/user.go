package models

type User struct {
	Token               string `json:"token" bson:"token"`
	TelegramId          int    `json:"telegramId" bson:"telegramId"`
	TelegramDisplayName string `json:"telegramDisplayName" bson:"telegramDisplayName"`
}
