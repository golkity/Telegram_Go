package action

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
	"tgbot/database"
)

func AnalyzePortfolio(bot *tgbotapi.BotAPI, callback *tgbotapi.CallbackQuery) {
	userID := callback.From.ID

	portfolio, exists := database.UserPortfolios[userID]
	if !exists || len(portfolio.Stocks) == 0 {
		msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "Ваш портфель пуст.")
		bot.Send(msg)
		return
	}

	totalCurrentValue := 0.0
	totalForecastValue := 0.0
	positiveTrend := 0
	negativeTrend := 0
	neutralTrend := 0
	stockDetails := []string{}

	for _, stock := range portfolio.Stocks {
		totalCurrentValue += stock.Price

		forecastedPrice := stock.Price * (1 + stock.Percent/100)
		totalForecastValue += forecastedPrice

		if stock.Percent > 0 {
			positiveTrend++
		} else if stock.Percent < 0 {
			negativeTrend++
		} else {
			neutralTrend++
		}

		stockDetails = append(stockDetails,
			fmt.Sprintf("Акция: %s\nТекущая цена: %.2f\nДоходность: %.2f%%\nПрогнозируемая цена: %.2f",
				stock.Name, stock.Price, stock.Percent, forecastedPrice))
	}

	trendMessage := fmt.Sprintf("Анализ вашего портфеля:\n\n"+
		"Общая текущая стоимость портфеля: %.2f\n"+
		"Прогнозируемая стоимость портфеля: %.2f\n\n"+
		"Тенденции:\n"+
		"- Положительная тенденция (рост): %d акций\n"+
		"- Отрицательная тенденция (падение): %d акций\n"+
		"- Без изменений: %d акций\n\n"+
		"Детали акций:\n%s",
		totalCurrentValue, totalForecastValue, positiveTrend, negativeTrend, neutralTrend,
		strings.Join(stockDetails, "\n\n"))

	msg := tgbotapi.NewMessage(callback.Message.Chat.ID, trendMessage)
	bot.Send(msg)
}
