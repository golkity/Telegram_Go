package handlers

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strconv"
	"strings"
	"tgbot/database"
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
		DisplayPortfolio(bot, callback.From.ID, callback.Message.Chat.ID)
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

func DisplayPortfolio(bot *tgbotapi.BotAPI, chatID int64, userID int64) {
	portfolio, exists := database.UserPortfolios[userID]
	if !exists || len(portfolio.Stocks) == 0 {
		msg := tgbotapi.NewMessage(chatID, "Ваш портфель пуст.")
		bot.Send(msg)
		return
	}

	total := 0.0
	stockDetails := []string{}
	for _, stock := range portfolio.Stocks {
		total += stock.Price
		stockDetails = append(stockDetails,
			fmt.Sprintf("Название: %s\nЦена: %.2f\nДоходность: %.2f%%", stock.Name, stock.Price, stock.Percent))
	}

	msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Общая сумма портфеля: %.2f\nКоличество акций: %d\n\n%s",
		total, len(portfolio.Stocks), strings.Join(stockDetails, "\n\n")))
	bot.Send(msg)
}
