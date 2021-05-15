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
	bot *tgbotapi.BotAPI
}

func NewTelegramClient() TelegramClient {
	bot, err := tgbotapi.NewBotAPI(viper.GetString("TELEGRAM_BOT_TOKEN"))
	if err != nil {
		log.Fatal("Error initializing telegram")
	}
	return &telegramClient{
		bot: bot,
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

		if update.Message.Text == "/token" {
			token := uuid.New()
			tokenMsg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("%s", token))
			tokenMsg.ReplyToMessageID = update.Message.MessageID
			t.bot.Send(tokenMsg)
			t.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "You have to put this token on alarm create form."))
		} else {
			t.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Invalid command"))
		}
	}
}
