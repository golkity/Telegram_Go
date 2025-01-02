package keyboards

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func StartInline() tgbotapi.InlineKeyboardMarkup {
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Мой портфель", "portfel"),
			tgbotapi.NewInlineKeyboardButtonData("Удалить акцию", "delet"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Добавить акцию", "add"),
			tgbotapi.NewInlineKeyboardButtonData("Анализ акций", "analiz"),
		),
	)
	return keyboard
}
