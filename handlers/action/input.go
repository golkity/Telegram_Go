package action

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
	"tgbot/database"
)

func HandlerInput(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	userID := message.From.ID
	if _, exist := database.UserStates[userID]; !exist {
		database.UserStates[userID] = &database.UserState{}
	}
	userState := database.UserStates[userID]

	switch userState.State {
	case "await_stock_name":
		userState.StockName = message.Text
		userState.State = "await_stock_price"
		msg := tgbotapi.NewMessage(message.Chat.ID, "Введите цену акции:")
		bot.Send(msg)

	case "await_stock_price":
		price, err := strconv.ParseFloat(message.Text, 64)
		if err != nil {
			msg := tgbotapi.NewMessage(message.Chat.ID, "Цена должна быть числом!")
			bot.Send(msg)
			return
		}
		userState.StockPrice = price
		userState.State = "await_stock_percent"
		msg := tgbotapi.NewMessage(message.Chat.ID, "Введите процент доходности акции:")
		bot.Send(msg)

	case "await_stock_percent":
		percent, err := strconv.ParseFloat(message.Text, 64)
		if err != nil {
			msg := tgbotapi.NewMessage(message.Chat.ID, "Процент должен быть числом. Попробуйте снова:")
			bot.Send(msg)
			return
		}
		userState.StockPercent = percent

		if database.UserPortfolios[userID] == nil {
			database.UserPortfolios[userID] = &database.Portfolio{Stocks: make(map[string]database.Stock)}
		}

		database.UserPortfolios[userID].Stocks[userState.StockName] = database.Stock{
			Name:    userState.StockName,
			Price:   userState.StockPrice,
			Percent: userState.StockPercent,
		}

		userState.State = ""
		msg := tgbotapi.NewMessage(message.Chat.ID, "Акция добавлена в ваш портфель!")
		bot.Send(msg)
	}
}
