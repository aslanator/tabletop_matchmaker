package link

import (
	"database/sql"
	json2 "encoding/json"
	"errors"
	"fmt"
	"strings"
	"tabletop_matchmaker/internal/db/entities"
	internal_errors "tabletop_matchmaker/internal/errors"
	helpers_errors "tabletop_matchmaker/internal/helpers/errors"
	"tabletop_matchmaker/internal/helpers/keyboards"
	"tabletop_matchmaker/internal/helpers/markdown"
	"tabletop_matchmaker/internal/services/bgg"
	"tabletop_matchmaker/internal/services/bggfmt"
	"tabletop_matchmaker/internal/types"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CqData struct {
	types.CqData
	Confirmation bool    `json:"y"`
	Username     *string `json:"u"`
}

func newCqData(confirmation bool, username *string) CqData {
	return CqData{
		CqData: types.CqData{
			Command: Name(),
		},
		Confirmation: confirmation,
		Username:     username,
	}
}

type Link struct{}

func (link Link) Run(msg *tgbotapi.Message, _ *sql.DB) []tgbotapi.Chattable {
	attributesStr := msg.CommandArguments()
	attributes := strings.Split(attributesStr, " ")
	if attributes[0] == "" {
		text := "Не указано имя пользователя BGG"
		return []tgbotapi.Chattable{tgbotapi.NewMessage(msg.Chat.ID, text)}
	}

	service := bgg.Bgg{}
	bggFmt := bggfmt.BggFmt{}
	user, err := service.FetchUser(attributes[0])
	if err != nil {
		var text string
		if errors.Is(err, internal_errors.ErrUserNotFound) {
			text = "Пользователь с таким именем не найден"
		} else {
			text = helpers_errors.UnexpectedChatError(err)
		}
		return []tgbotapi.Chattable{tgbotapi.NewMessage(msg.Chat.ID, text)}
	}

	appeal := bggFmt.ResolveAppeal(user)
	profileLink := bgg.USER_BASE_LINK + user.UserName
	linkMarkdown := markdown.GenLink(appeal, profileLink)
	text := fmt.Sprintf("%s, это вы?", linkMarkdown)

	keyboard, err := genConfirmationKeyboard(user.UserName)
	if err != nil {
		return helpers_errors.UnexpectedChatErrorMessage(err, msg.Chat.ID)
	}

	if user.AvatarLink != "N/A" {
		res := tgbotapi.NewPhoto(msg.Chat.ID, tgbotapi.FileURL(user.AvatarLink))
		res.Caption = text
		res.ReplyMarkup = keyboard
		res.ParseMode = markdown.GetParseMode()
		return []tgbotapi.Chattable{res}
	} else {
		res := tgbotapi.NewMessage(msg.Chat.ID, text)
		res.ReplyMarkup = keyboard
		res.ParseMode = markdown.GetParseMode()
		return []tgbotapi.Chattable{res}
	}
}

func genConfirmationKeyboard(login string) (*tgbotapi.InlineKeyboardMarkup, error) {
	confirmCqData := newCqData(true, &login)
	confirmJson, err := json2.Marshal(confirmCqData)
	if err != nil {
		return nil, err
	}

	confirmButton := tgbotapi.NewInlineKeyboardButtonData(keyboards.CHECK_MARK, string(confirmJson))

	cancelCqData := newCqData(false, nil)
	cancelJson, err := json2.Marshal(cancelCqData)
	if err != nil {
		return nil, err
	}

	cancelButton := tgbotapi.NewInlineKeyboardButtonData(keyboards.CROSS_MARK, string(cancelJson))

	row := tgbotapi.NewInlineKeyboardRow(confirmButton, cancelButton)
	keyboardMarkup := tgbotapi.NewInlineKeyboardMarkup(row)

	return &keyboardMarkup, nil
}

func (link Link) Callback(cq *tgbotapi.CallbackQuery, database *sql.DB) []tgbotapi.Chattable {
	var cqData CqData
	err := json2.Unmarshal([]byte(cq.Data), &cqData)
	if err != nil {
		return helpers_errors.UnexpectedChatErrorMessage(err, cq.Message.Chat.ID)
	}

	result := make([]tgbotapi.Chattable, 0)

	if cqData.Confirmation {
		bggUser := *cqData.Username
		entities.Upsert_into_bgg_account(database, cq.From.ID, bggUser)
		text := "Не извольте волноваться"
		result = append(result, tgbotapi.NewMessage(cq.Message.Chat.ID, text))
	}

	edit := keyboards.NewRemoveInlineKeyboard(cq.Message.Chat.ID, cq.Message.MessageID)

	result = append(result, edit)

	return result
}

func Name() string {
	return "link"
}
