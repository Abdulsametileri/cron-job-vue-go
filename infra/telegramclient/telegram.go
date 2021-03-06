package telegramclient

import (
	"fmt"
	"github.com/Abdulsametileri/cron-job-vue-go/models"
	"github.com/Abdulsametileri/cron-job-vue-go/services/userservice"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"log"
)

type TelegramClient interface {
	GetMessages()
	SendImage(telegramId int64, imageUrl string) error
	SendMessageForDebug(msg string)
}

type telegramClient struct {
	bot *tgbotapi.BotAPI
	us  userservice.UserService
}

func NewTelegramClient(us userservice.UserService) TelegramClient {
	bot, err := tgbotapi.NewBotAPI(viper.GetString("TELEGRAM_BOT_TOKEN"))
	if err != nil {
		log.Fatalf("Error initializing telegram %v", err)
	}
	return &telegramClient{
		bot: bot,
		us:  us,
	}
}

func (t telegramClient) SendMessageForDebug(text string) {
	msg := tgbotapi.NewMessage(513873156, text)
	_, _ = t.bot.Send(msg)
}

func (t telegramClient) SendImage(telegramId int64, imageUrl string) error {
	msg := tgbotapi.NewPhotoShare(telegramId, imageUrl)
	_, err := t.bot.Send(msg)
	return err
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
		userTelegramId := int64(update.Message.From.ID)
		userTelegramName := update.Message.From.UserName
		chatId := update.Message.Chat.ID

		if update.Message.Text == "/token" {
			user, err := t.us.GetUserByTelegramId(userTelegramId)
			if err != nil {
				_, _ = t.bot.Send(tgbotapi.NewMessage(chatId, err.Error()))
				continue
			}
			if user.TelegramId != 0 {
				_, _ = t.bot.Send(tgbotapi.NewMessage(chatId, fmt.Sprintf("You have already token. %s", user.Token)))
				continue
			}

			token := uuid.New()
			tokenMsg := tgbotapi.NewMessage(chatId, fmt.Sprintf("%s", token))
			tokenMsg.ReplyToMessageID = update.Message.MessageID

			err = t.us.AddUser(models.User{
				Token:               token.String(),
				TelegramId:          userTelegramId,
				TelegramDisplayName: userTelegramName,
			})

			_, _ = t.bot.Send(tokenMsg)
			_, _ = t.bot.Send(tgbotapi.NewMessage(chatId, "You have to give this token on our site."))
		} else {
			_, _ = t.bot.Send(tgbotapi.NewMessage(chatId, "Invalid command"))
		}
	}
}
