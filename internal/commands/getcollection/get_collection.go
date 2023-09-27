package getcollection

import (
	"database/sql"
	"strings"
	helpers_errors "tabletop_matchmaker/internal/helpers/errors"
	"tabletop_matchmaker/internal/services/bgg"
	"tabletop_matchmaker/internal/services/bggfmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type GetCollection struct {
}

func (getCollection GetCollection) Run(msg *tgbotapi.Message,  database *sql.DB) []tgbotapi.Chattable {
	attributesStr := msg.CommandArguments()
	attributes := strings.Split(attributesStr, " ")
	username := attributes[0]
	if username == "" {
		text := "Не указано имя пользователя BGG"
		return []tgbotapi.Chattable{tgbotapi.NewMessage(msg.Chat.ID, text)}
	}

	bggService := bgg.Bgg{}
	bggFmt := bggfmt.BggFmt{}
	collection, err := bggService.FetchCollection(username)
	if err != nil {
		text := helpers_errors.UnexpectedChatError(err)
		return []tgbotapi.Chattable{tgbotapi.NewMessage(msg.Chat.ID, text)}
	}
	message := bggFmt.GenGameNamesList(collection)

	messageConfig := tgbotapi.NewMessage(msg.Chat.ID, msg.Text)
	messageConfig.Text = message
	return []tgbotapi.Chattable{messageConfig}
}

func (collection GetCollection) Callback(cq *tgbotapi.CallbackQuery, database *sql.DB) []tgbotapi.Chattable {
	return nil
}

func Name() string {
	return "getCollection"
}
