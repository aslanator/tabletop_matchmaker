package link

import (
	"context"
	"fmt"
	"strings"
	"tabletop_matchmaker/internal/helpers/errors"

	"github.com/fzerorubigd/gobgg"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const BGG_USER_BASE_LINK = "https://boardgamegeek.com/user/"

type Link struct {
}

func (link Link) Run(msg *tgbotapi.Message) []tgbotapi.Chattable {
	attributesStr := msg.CommandArguments()
	attributes := strings.Split(attributesStr, " ")
	if attributes[0] == "" {
		text := "Не указано имя пользователя BGG"
		return []tgbotapi.Chattable{tgbotapi.NewMessage(msg.Chat.ID, text)}
	}

	bgg := gobgg.NewBGGClient()
	user, err := bgg.GetUser(context.Background(), attributes[0])
	if err != nil {
		text := errors.UnexpectedChatError(err)
		return []tgbotapi.Chattable{tgbotapi.NewMessage(msg.Chat.ID, text)}
	}

	if user.UserID == 0 {
		text := "Пользователь с таким именем не найден"
		return []tgbotapi.Chattable{tgbotapi.NewMessage(msg.Chat.ID, text)}
	}

	var name string
	if user.FirstName != "" || user.LastName != "" {
		name = strings.Trim(user.FirstName+" "+user.LastName, " ")
	} else {
		name = user.UserName
	}
	profileLink := BGG_USER_BASE_LINK + user.UserName
	linkMarkdown := fmt.Sprintf("[%s](%s)", name, profileLink)
	text := fmt.Sprintf("%s, это вы?", linkMarkdown)
	if user.AvatarLink != "N/A" {
		res := tgbotapi.NewPhoto(msg.Chat.ID, tgbotapi.FileURL(user.AvatarLink))
		res.Caption = text
		res.ReplyMarkup = genCofirmation(attributes[0])
		res.ParseMode = "MarkdownV2"
		return []tgbotapi.Chattable{res}
	} else {
		res := tgbotapi.NewMessage(msg.Chat.ID, text)
		res.ReplyMarkup = genCofirmation(attributes[0])
		res.ParseMode = "MarkdownV2"
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
	text := ""
	if data[1] == "confirm" {
		text = "Не извольте волноваться"
	} else {
		text = "Извольте волновоаться"
	}
	result := []tgbotapi.Chattable{tgbotapi.NewMessage(cq.Message.Chat.ID, text)}

	edit := tgbotapi.EditMessageReplyMarkupConfig{
		BaseEdit: tgbotapi.BaseEdit{
			ChatID:      cq.Message.Chat.ID,
			MessageID:   cq.Message.MessageID,
			ReplyMarkup: nil,
		},
	}

	result = append(result, edit)

	return result
}

func Name() string {
	return "link"
}
