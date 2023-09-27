package main

import (
	"log"
	"tabletop_matchmaker/configs"
	"tabletop_matchmaker/internal/commands"
	"tabletop_matchmaker/internal/commands/help"
	"tabletop_matchmaker/internal/commands/link"
	"tabletop_matchmaker/internal/db"
	"tabletop_matchmaker/internal/helpers/errors"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	env := configs.NewEnviroment()

	bot, err := tgbotapi.NewBotAPI(env.ApiKey)
	errors.FatalOnError(err, "Failed to validate the telegram API token")

	bot.Debug = true

	configCommands(bot)

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	controller := commands.Controller{Database: db.GetDB()}

	for update := range updates {
		go run(bot, controller, update, env)
	}
}

func configCommands(bot *tgbotapi.BotAPI) {
	help := tgbotapi.BotCommand{
		Command:     help.Name(),
		Description: "Возможно, это поможет, но я бы на это не рассчитывал",
	}
	link := tgbotapi.BotCommand{
		Command:     link.Name(),
		Description: "Присоединить аккаунт BGG",
	}

	commandsConfig := tgbotapi.NewSetMyCommands(help, link)
	bot.Send(commandsConfig)
}


func run(bot *tgbotapi.BotAPI, controller commands.Controller, update tgbotapi.Update, env configs.Enviroment) {
	msgs := controller.HandleUpdate(update, env.BotName)

	if msgs == nil {
		return;
	}

	for _, msg := range msgs {
		if _, err := bot.Send(msg); err != nil {
			log.Println(err)
		}
	}
}