package handlers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"tgbot/database"
	"tgbot/handlers/action"
)

func StartCallback(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	callback := update.CallbackQuery
	if callback == nil {
		return
	}

	if callback.From == nil {
		log.Println("callback.From is nil")
		return
	}

	if _, exists := database.UserStates[callback.From.ID]; !exists {
		database.UserStates[callback.From.ID] = &database.UserState{}
	}
	switch callback.Data {
	case "portfel":
		action.DisplayPortfolio(bot, callback.From.ID, callback.Message.Chat.ID)
	case "delet":
		msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "Delet")
		bot.Send(msg)
	case "add":
		database.UserStates[callback.From.ID].State = "await_stock_name"
		msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "Введите название акций:")
		bot.Send(msg)
	case "analiz":
		action.AnalyzePortfolio(bot, callback)
		action.SendPortfolioGrowthGraph(bot, callback.Message.Chat.ID, callback.From.ID)
	default:
		log.Printf("Unknown callback data '%s'", callback.Data)
	}
}
