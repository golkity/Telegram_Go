package handler

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"tgbot/commands"
)

func HandlerCommands(bot *tgbotapi.BotAPI, update tgbotapi.Update) (bool, int) {
	if update.Message.IsCommand() {
		switch update.Message.Command() {
		case "start":
			commands.Start(bot, update)
		default:
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Что?")
			sentMsg, err := bot.Send(msg)
			if err != nil {
				return false, 0
			}
			return true, sentMsg.MessageID
		}
	}
	return false, 0
}
