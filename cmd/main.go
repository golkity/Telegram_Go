package main

import (
	"github.com/fatih/color"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
	"tgbot/Errors"
	"tgbot/config"
	"tgbot/handler"
)

func main() {
	cfg, err := config.LoadConfig("config/config.json")
	if err != nil {
		color.Red("%s", Errors.ErrLoadConfig)
		os.Exit(1)
	}
	bot, err := tgbotapi.NewBotAPI(cfg.Token)
	if err != nil {
		color.Red("%s", Errors.ErrCreateBot)
		os.Exit(1)
	}
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message != nil {
			continue
		}
		handler.HandleMessages(bot, update)
	}
}
