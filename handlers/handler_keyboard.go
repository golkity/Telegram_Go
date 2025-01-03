package handlers

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"image/color"
	"log"
	"math"
	"os"
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
	case "analiz":
		AnalyzePortfolio(bot, callback)
		SendPortfolioGrowthGraph(bot, callback.Message.Chat.ID, callback.From.ID)
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

func GeneratePortfolioGrowthGraph(userID int64) (string, error) {
	portfolio, exists := database.UserPortfolios[userID]
	if !exists || len(portfolio.Stocks) == 0 {
		return "", nil
	}

	totalCurrentValue := 0.0
	for _, stock := range portfolio.Stocks {
		totalCurrentValue += stock.Price
	}

	points := plotter.XYs{}
	for month := 0; month <= 12; month++ {
		projectedValue := totalCurrentValue
		for _, stock := range portfolio.Stocks {
			projectedValue += stock.Price * math.Pow(1+(stock.Percent/100), float64(month))
		}
		points = append(points, plotter.XY{X: float64(month), Y: projectedValue})
	}

	p := plot.New()

	p.Title.Text = "Рост стоимости портфеля"
	p.X.Label.Text = "Месяцы"
	p.Y.Label.Text = "Стоимость (руб.)"
	p.Add(plotter.NewGrid())

	line, err := plotter.NewLine(points)
	if err != nil {
		return "", err
	}
	line.LineStyle.Width = vg.Points(2)
	line.LineStyle.Color = color.RGBA{R: 255, G: 0, B: 0, A: 255}
	p.Add(line)

	filePath := "portfolio_growth.png"
	if err := p.Save(6*vg.Inch, 4*vg.Inch, filePath); err != nil {
		return "", err
	}

	return filePath, nil
}

func SendPortfolioGrowthGraph(bot *tgbotapi.BotAPI, chatID int64, userID int64) {
	filePath, err := GeneratePortfolioGrowthGraph(userID)
	if err != nil {
		msg := tgbotapi.NewMessage(chatID, "Ошибка при создании графика.")
		bot.Send(msg)
		return
	}

	if filePath == "" {
		msg := tgbotapi.NewMessage(chatID, "Ваш портфель пуст.")
		bot.Send(msg)
		return
	}

	photo := tgbotapi.NewPhoto(chatID, tgbotapi.FilePath(filePath))
	bot.Send(photo)

	os.Remove(filePath)
}
