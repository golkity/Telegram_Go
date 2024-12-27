package commands

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func Start_com(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Привет"+update.Message.From.UserName)
	bot.Send(msg)
}
