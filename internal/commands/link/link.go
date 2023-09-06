package link

import (
	"errors"
	"fmt"
	"strings"
	internal_errors "tabletop_matchmaker/internal/errors"
	helpers_errors "tabletop_matchmaker/internal/helpers/errors"
	"tabletop_matchmaker/internal/helpers/keyboards"
	"tabletop_matchmaker/internal/helpers/markdown"
	"tabletop_matchmaker/internal/services/bgg"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Link struct{}

func (link Link) Run(msg *tgbotapi.Message) []tgbotapi.Chattable {
	attributesStr := msg.CommandArguments()
	attributes := strings.Split(attributesStr, " ")
	if attributes[0] == "" {
		text := "Не указано имя пользователя BGG"
		return []tgbotapi.Chattable{tgbotapi.NewMessage(msg.Chat.ID, text)}
	}

	service := bgg.Bgg{}
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

	appeal := service.ResolveAppeal(user)
	profileLink := bgg.USER_BASE_LINK + user.UserName
	linkMarkdown := markdown.GenLink(appeal, profileLink)
	text := fmt.Sprintf("%s, это вы?", linkMarkdown)

	if user.AvatarLink != "N/A" {
		res := tgbotapi.NewPhoto(msg.Chat.ID, tgbotapi.FileURL(user.AvatarLink))
		res.Caption = text
		res.ReplyMarkup = genCofirmation(attributes[0])
		res.ParseMode = markdown.GetParseMode()
		return []tgbotapi.Chattable{res}
	} else {
		res := tgbotapi.NewMessage(msg.Chat.ID, text)
		res.ReplyMarkup = genCofirmation(attributes[0])
		res.ParseMode = markdown.GetParseMode()
		return []tgbotapi.Chattable{res}
	}
}

func genCofirmation(login string) tgbotapi.InlineKeyboardMarkup {
	confirmButton := tgbotapi.NewInlineKeyboardButtonData("✅", Name()+"|confirm|"+login)
	cancelButton := tgbotapi.NewInlineKeyboardButtonData("❌", Name()+"|cancel")
	row := tgbotapi.NewInlineKeyboardRow(confirmButton, cancelButton)
	return tgbotapi.NewInlineKeyboardMarkup(row)
}

func (link Link) Callback(cq *tgbotapi.CallbackQuery) []tgbotapi.Chattable {
	data := strings.Split(cq.Data, "|")
	result := make([]tgbotapi.Chattable, 0)

	if data[1] == "confirm" {
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
