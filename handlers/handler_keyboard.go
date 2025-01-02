package handlers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strconv"
	"tgbot/database"
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
		database.UserStates[callback.From.ID].State = "await_stock_name"
		msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "Введите название акций:")
		bot.Send(msg)
	default:
		log.Printf("Unknown callback data '%s'", callback.Data)
	}
}

func HandlerInput(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	userID := message.From.ID
	if _, exist := database.UserStates[userID]; !exist {
		database.UserStates[userID] = &database.Userstate{}
	}
	userState := database.UserStates[userID]

	switch userState.State {
	case "await_stock_name":
		userState.StocNname = message.Text
		userState.State = "await_stock_price"
		msg := tgbotapi.NewMessage(message.Chat.ID, "Введите цену акции")
		bot.Send(msg)
	case "awaiting_stock_price":
		price, err := strconv.ParseFloat(message.Text, 64)
		if err != nil {
			msg := tgbotapi.NewMessage(message.Chat.ID, "Цена должна быть числом!")
			bot.Send(msg)
			return
		}
		userState.StockPrice = price
		userState.State = "await_stock_precent"
		msg := tgbotapi.NewMessage(message.Chat.ID, "Введите процент доходности акции:")
		bot.Send(msg)
	case "awaiting_stock_percent":
		percent, err := strconv.ParseFloat(message.Text, 64)
		if err != nil {
			msg := tgbotapi.NewMessage(message.Chat.ID, "Процент должен быть числом. Попробуйте снова:")
			bot.Send(msg)
			return
		}
		userState.StockPrecent = percent

		if database.UserPortfolios[userID] == nil {
			database.UserPortfolios[userID] = &database.Portfolio{Stocks: make(map[string]database.Stock)}
		}
		database.UserPortfolios[userID].Stocks[userState.StocNname] = database.Stock{
			Name:    userState.StocNname,
			Price:   userState.StockPrice,
			Percent: userState.StockPrecent,
		}

		userState.State = ""
		msg := tgbotapi.NewMessage(message.Chat.ID, "Акция добавлена!")
		bot.Send(msg)
	}
}
