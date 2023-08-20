package commands

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type Command interface {
	Run(msg tgbotapi.Message) tgbotapi.Chattable
}
