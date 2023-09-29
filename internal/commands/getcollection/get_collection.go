package getcollection

import (
	"database/sql"
	"strconv"
	"strings"
	helpers_errors "tabletop_matchmaker/internal/helpers/errors"
	"tabletop_matchmaker/internal/services/bgg"
	"tabletop_matchmaker/internal/services/bggfmt"
	"tabletop_matchmaker/internal/services/paginator"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type GetCollection struct {
}

func (getCollection GetCollection) Run(msg *tgbotapi.Message,  database *sql.DB) []tgbotapi.Chattable {
	attributesStr := msg.CommandArguments()
	attributes := strings.Split(attributesStr, " ")
	username := attributes[0]
	var page int = 1
	if len(attributes) > 1 {
		pageNumber, err := strconv.Atoi(attributes[1]);
		if err != nil {
			page = 1
		} else {
			page = pageNumber
		}
	}
	if username == "" {
		text := "Не указано имя пользователя BGG"
		return []tgbotapi.Chattable{tgbotapi.NewMessage(msg.Chat.ID, text)}
	}

	return getCollection.makeMessage(username, page, msg.Chat.ID)
}

func (getCollection GetCollection) Callback(cq *tgbotapi.CallbackQuery, database *sql.DB) []tgbotapi.Chattable {
	data := strings.Split(cq.Data, "|")
	username := data[1]
	var page int = 1
	if len(data) > 2 {
		pageNumber, err := strconv.Atoi(data[2]);
		if err != nil {
			page = 1
		} else {
			page = pageNumber
		}
	}
	if username == "" {
		text := "Не указано имя пользователя BGG"
		return []tgbotapi.Chattable{tgbotapi.NewMessage(cq.Message.Chat.ID, text)}
	}
	return getCollection.makeMessage(username, page, cq.Message.Chat.ID)
}

func (getCollection GetCollection) makeMessage(username string, page int, chatId int64) []tgbotapi.Chattable {
	bggService := bgg.Bgg{}
	bggFmt := bggfmt.BggFmt{}
	paginatorService := paginator.Paginator{}

	collection, err := bggService.FetchCollection(username)
	if err != nil {
		text := helpers_errors.UnexpectedChatError(err)
		return []tgbotapi.Chattable{tgbotapi.NewMessage(chatId, text)}
	}
	message := bggFmt.GenGameNamesList(collection)
	pagedMessage, pagedButtons := paginatorService.GenPageMessage(message, page, Name() + "|" + username)

	messageConfig := tgbotapi.NewMessage(chatId, pagedMessage)
	messageConfig.ReplyMarkup = pagedButtons
	return []tgbotapi.Chattable{messageConfig}
}

func Name() string {
	return "getCollection"
} 
