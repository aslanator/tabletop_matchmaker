package commands

import (
	"database/sql"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Command interface {
	Run(msg *tgbotapi.Message, database *sql.DB) []tgbotapi.Chattable
	Callback(msg *tgbotapi.CallbackQuery, database *sql.DB) []tgbotapi.Chattable
}
