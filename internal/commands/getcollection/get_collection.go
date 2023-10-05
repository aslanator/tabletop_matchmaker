package getcollection

import (
	"database/sql"
	"encoding/json"
	"strconv"
	"strings"
	helpers_errors "tabletop_matchmaker/internal/helpers/errors"
	"tabletop_matchmaker/internal/services/bggfmt"
	"tabletop_matchmaker/internal/services/gamescollection"
	"tabletop_matchmaker/internal/services/paginator"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type GetCollection struct {
}

func (getCollection GetCollection) Run(msg *tgbotapi.Message, db *sql.DB) []tgbotapi.Chattable {
	attributesStr := msg.CommandArguments()
	attributes := strings.Split(attributesStr, " ")
	username := attributes[0]

	page := 1
	if len(attributes) > 1 {
		pageNumber, err := strconv.Atoi(attributes[1])
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

	games, err := getCollection.getGamesFromBgg(db, username)

	if err != nil {
		return helpers_errors.UnexpectedChatErrorMessage(err, msg.Chat.ID)
	}

	messageText, keyboard, err := getCollection.prepareMessageTextAndKeyboard(games, username, page)

	if err != nil {
		return helpers_errors.UnexpectedChatErrorMessage(err, msg.Chat.ID)
	}

	messageConfig := tgbotapi.NewMessage(msg.Chat.ID, messageText)
	messageConfig.ReplyMarkup = keyboard

	return []tgbotapi.Chattable{messageConfig}
}

func (getCollection GetCollection) Callback(cq *tgbotapi.CallbackQuery, db *sql.DB) []tgbotapi.Chattable {
	var cqData CqData
	err := json.Unmarshal([]byte(cq.Data), &cqData)
	if err != nil {
		return helpers_errors.UnexpectedChatErrorMessage(err, cq.Message.Chat.ID)
	}

	if cqData.Username == "" {
		text := "Не указано имя пользователя BGG"
		return []tgbotapi.Chattable{tgbotapi.NewMessage(cq.Message.Chat.ID, text)}
	}

	games, err := getCollection.getGamesFromDB(db, cqData.Username)

	if err != nil {
		return helpers_errors.UnexpectedChatErrorMessage(err, cq.Message.Chat.ID)
	}

	messageText, keyboard, err := getCollection.prepareMessageTextAndKeyboard(games, cqData.Username, cqData.Page)

	if err != nil {
		return helpers_errors.UnexpectedChatErrorMessage(err, cq.Message.Chat.ID)
	}

	return []tgbotapi.Chattable{tgbotapi.NewEditMessageTextAndMarkup(cq.Message.Chat.ID, cq.Message.MessageID, messageText, *keyboard)}
}

func (getCollection GetCollection) prepareMessageTextAndKeyboard(games []gamescollection.Game, username string, page int) (string, *tgbotapi.InlineKeyboardMarkup, error) {
	bggFmt := bggfmt.BggFmt{}

	if len(games) == 0 {
		keyboard := tgbotapi.NewInlineKeyboardMarkup()
		return "Пусто", &keyboard, nil
	}

	paginatorService := paginator.NewPaginator(games, 20)
	pageItems, err := paginatorService.GetPage(page)

	if err != nil {
		return "", nil, err
	}

	messageText := bggFmt.GenGameNamesList(pageItems)

	cqData := newCqData(username, page)
	keyboard, err := paginatorService.GenKeyboard(page, &cqData)

	if err != nil {
		return "", nil, err
	}

	return messageText, keyboard, nil
}

func (getCollection GetCollection) getGamesFromBgg(db *sql.DB, username string) ([]gamescollection.Game, error) {
	gamesCollection := gamescollection.GamesCollection{}
	games, err := gamesCollection.GetGamesByUserName(db, username, false)
	if err != nil {
		return nil, err
	}

	return games, nil
}

func (getCollection GetCollection) getGamesFromDB(db *sql.DB, username string) ([]gamescollection.Game, error) {
	gamesCollection := gamescollection.GamesCollection{}
	games, err := gamesCollection.GetGamesByUserName(db, username, true)
	if err != nil {
		return nil, err
	}
	
	return games, nil
}



func Name() string {
	return "getCollection"
}
