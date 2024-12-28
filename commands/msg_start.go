package commands

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"tgbot/internal/keyboards"
)

func StartMsg(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Welcome! Your username is "+update.Message.From.UserName)
	msg.ReplyMarkup = keyboards.createKeyboard()
	bot.Send(msg)
}
