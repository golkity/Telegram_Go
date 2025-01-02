package handlers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

func StartCallback(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	callback := update.CallbackQuery
	if callback == nil {
		return
	}
	switch callback.Data {
	case "portfel":
		msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "Portfel:")
		bot.Send(msg)
	case "delet":
		msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "Delet")
		bot.Send(msg)
	case "add":
		msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "Add")
		bot.Send(msg)
	default:
		log.Printf("Unknown callback data '%s'", callback.Data)
	}
}
