package link

import (
	"context"
	"fmt"
	"log"
	"strings"
	"tabletop_matchmaker/internal/helpers/errors"

	"github.com/fzerorubigd/gobgg"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Link struct {
}

func (link Link) Run(msg *tgbotapi.Message) tgbotapi.Chattable {
	attributesStr := msg.CommandArguments()
	attributes := strings.Split(attributesStr, " ")
	if attributes[0] == "--confirm" {
	}

	bgg := gobgg.NewBGGClient()
	user, err := bgg.GetUser(context.Background(), attributes[0])
	if err != nil {
		log.Println(err)
		text := "Бип-буп, глупый робот всё сломал"
		return tgbotapi.NewMessage(msg.Chat.ID, text)
	}
	errors.Vd(user)

	text := "Ебало узнаёшь, пёс?"
	res := tgbotapi.NewMessage(msg.Chat.ID, text)
	res.ReplyMarkup = genCofirmation(attributes[0])
	return res
}

func genCofirmation(login string) tgbotapi.InlineKeyboardMarkup {
	data := Name() + "|%s|" + login
	confirmButton := tgbotapi.NewInlineKeyboardButtonData("✅", fmt.Sprintf(data, "confirm"))
	cancelButton := tgbotapi.NewInlineKeyboardButtonData("❌", fmt.Sprintf(data, "cancel"))
	row := tgbotapi.NewInlineKeyboardRow(confirmButton, cancelButton)
	return tgbotapi.NewInlineKeyboardMarkup(row)
}

func (link Link) Callback(cq *tgbotapi.CallbackQuery) tgbotapi.Chattable {
	return nil
	text := "Не извольте волноваться, всё заебись"
	return tgbotapi.NewMessage(cq.ChatInstance, text)
}

func Name() string {
	return "link"
}
