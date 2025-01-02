package handlers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"tgbot/commands"
)

func Start(bot *tgbotapi.BotAPI, update tgbotapi.Update) bool {
	if update.Message.IsCommand() {
		switch update.Message.Command() {
		case "start":
			commands.Start(bot, update)
			return true
		default:
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Что?")
			bot.Send(msg)
			return false
		}
	}
	return false
}
