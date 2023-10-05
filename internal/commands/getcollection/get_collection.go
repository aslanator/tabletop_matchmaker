package getcollection

import (
	"database/sql"
	"encoding/json"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
	"strings"
	helpers_errors "tabletop_matchmaker/internal/helpers/errors"
	"tabletop_matchmaker/internal/services/bgg"
	"tabletop_matchmaker/internal/services/bggfmt"
	"tabletop_matchmaker/internal/services/paginator"
	"tabletop_matchmaker/internal/types"
)

type CqData struct {
	types.CqData
	Username string `json:"u"`
	Page     int    `json:"p"`
}

func (cqData *CqData) GetPage() int {
	return cqData.Page
}

func (cqData *CqData) SetPage(page int) {
	cqData.Page = page
}

func newCqData(username string, page int) CqData {
	return CqData{
		CqData: types.CqData{
			Command: Name(),
		},
		Username: username,
		Page:     page,
	}
}

type GetCollection struct {
}

func (getCollection GetCollection) Run(msg *tgbotapi.Message, _ *sql.DB) []tgbotapi.Chattable {
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

	return getCollection.makeMessage(username, page, msg.Chat.ID, nil)
}

func (getCollection GetCollection) Callback(cq *tgbotapi.CallbackQuery, _ *sql.DB) []tgbotapi.Chattable {
	var cqData CqData
	err := json.Unmarshal([]byte(cq.Data), &cqData)
	if err != nil {
		return helpers_errors.UnexpectedChatErrorMessage(err, cq.Message.Chat.ID)
	}

	if cqData.Username == "" {
		text := "Не указано имя пользователя BGG"
		return []tgbotapi.Chattable{tgbotapi.NewMessage(cq.Message.Chat.ID, text)}
	}

	return getCollection.makeMessage(cqData.Username, cqData.Page, cq.Message.Chat.ID, &cq.Message.MessageID)
}

func (getCollection GetCollection) makeMessage(username string, page int, chatId int64, messageId *int) []tgbotapi.Chattable {
	bggService := bgg.Bgg{}
	bggFmt := bggfmt.BggFmt{}

	collection, err := bggService.FetchCollection(username)
	if err != nil {
		return helpers_errors.UnexpectedChatErrorMessage(err, chatId)
	}

	if len(collection) == 0 {
		return []tgbotapi.Chattable{tgbotapi.NewMessage(chatId, "Пусто")}
	}

	paginatorService := paginator.NewPaginator(collection, 20)
	pageItems, err := paginatorService.GetPage(page)

	if err != nil {
		return helpers_errors.UnexpectedChatErrorMessage(err, chatId)
	}

	pageList := bggFmt.GenGameNamesList(pageItems)
	cqData := newCqData(username, page)
	pagedButtons, err := paginatorService.GenKeyboardRow(page, &cqData)

	if err != nil {
		return helpers_errors.UnexpectedChatErrorMessage(err, chatId)
	}

	var keyboard tgbotapi.InlineKeyboardMarkup
	if pagedButtons != nil {
		keyboard = tgbotapi.NewInlineKeyboardMarkup(pagedButtons)
	}

	if messageId == nil {
		messageConfig := tgbotapi.NewMessage(chatId, pageList)
		messageConfig.ReplyMarkup = keyboard

		return []tgbotapi.Chattable{messageConfig}
	}

	var res []tgbotapi.Chattable

	messageConfig := tgbotapi.NewEditMessageText(chatId, *messageId, pageList)
	res = append(res, messageConfig)

	keyboardConfig := tgbotapi.NewEditMessageReplyMarkup(chatId, *messageId, keyboard)
	res = append(res, keyboardConfig)

	return res
}

func Name() string {
	return "getCollection"
}
