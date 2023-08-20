package main

import (
	"log"
	"tabletop_matchmaker/configs"
	"tabletop_matchmaker/internal/commands"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	env := configs.NewEnviroment()

	bot, err := tgbotapi.NewBotAPI(env.ApiKey)
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
		msg := controller.HandleUpdate(update, env.BotName)

		if msg == nil {
			continue
		}

		if _, err := bot.Send(msg); err != nil {
			log.Println(err)
		}
	}
}
