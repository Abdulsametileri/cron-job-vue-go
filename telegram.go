package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"log"
)

type TelegramClient interface {
	GetMessages()
}

type telegramClient struct {
	bot         *tgbotapi.BotAPI
	mongoClient MongoClient
}

func NewTelegramClient(mongoClient MongoClient) TelegramClient {
	bot, err := tgbotapi.NewBotAPI(viper.GetString("TELEGRAM_BOT_TOKEN"))
	if err != nil {
		log.Fatal("Error initializing telegram")
	}
	return &telegramClient{
		bot:         bot,
		mongoClient: mongoClient,
	}
}

func (t telegramClient) GetMessages() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := t.bot.GetUpdatesChan(u)
	log.Println(err)

	for update := range updates {
		if update.Message == nil {
			continue
		}
		userTelegramId := update.Message.From.ID
		userTelegramName := update.Message.From.UserName
		chatId := update.Message.Chat.ID

		if update.Message.Text == "/token" {
			user, err := t.mongoClient.GetUserByTelegramId(userTelegramId)
			if err != nil {
				t.bot.Send(tgbotapi.NewMessage(chatId, err.Error()))
				continue
			}
			if user.TelegramId > 0 {
				t.bot.Send(tgbotapi.NewMessage(chatId, fmt.Sprintf("You have already token. %s", user.Token)))
				continue
			}

			token := uuid.New()
			tokenMsg := tgbotapi.NewMessage(chatId, fmt.Sprintf("%s", token))
			tokenMsg.ReplyToMessageID = update.Message.MessageID

			err = t.mongoClient.AddUser(User{
				Token:               token.String(),
				TelegramId:          userTelegramId,
				TelegramDisplayName: userTelegramName,
			})
			fmt.Println(err)

			t.bot.Send(tokenMsg)
			t.bot.Send(tgbotapi.NewMessage(chatId, "You have to put this token on alarm create form."))
		} else {
			t.bot.Send(tgbotapi.NewMessage(chatId, "Invalid command"))
		}
	}
}
