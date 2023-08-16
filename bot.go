package main

import (
	"log"
	"os"
	"tabletop_matchmaker/commands"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	apiKey := os.Getenv("API_KEY")
	botName := os.Getenv("BOT_NAME")

	bot, err := tgbotapi.NewBotAPI(apiKey)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	controller := commands.Controller{}

	for update := range updates {
		msg := controller.HandleUpdate(update, botName)

		if msg == nil {
			continue
		}

		if _, err := bot.Send(msg); err != nil {
			log.Println(err)
		}
	}
}
