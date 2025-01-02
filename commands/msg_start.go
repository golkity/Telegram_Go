package commands

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"tgbot/keyboards"
)

func Start(bot *tgbotapi.BotAPI, update tgbotapi.Update) (bool, int) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Добро пожаловать "+update.Message.From.FirstName+"!")
	msg.ReplyMarkup = keyboards.StartInline()
	sentMsg, err := bot.Send(msg)
	if err != nil {
		return false, 0
	}
	return true, sentMsg.MessageID
}
