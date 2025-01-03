package action

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"image/color"
	"math"
	"os"
	"tgbot/database"
)

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
