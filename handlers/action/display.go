package action

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
	"tgbot/database"
)

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
