package keyboards

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const CHECK_MARK = "✅"
const CROSS_MARK = "❌"
const LEFT_ARROW = "⬅️"
const RIGHT_ARROW = "➡️"

func NewRemoveInlineKeyboard(chatId int64, messageId int) tgbotapi.EditMessageReplyMarkupConfig {
	return tgbotapi.EditMessageReplyMarkupConfig{
		BaseEdit: tgbotapi.BaseEdit{
			ChatID:      chatId,
			MessageID:   messageId,
			ReplyMarkup: nil,
		},
	}
}

func GenBlankButton() tgbotapi.InlineKeyboardButton {
	return tgbotapi.NewInlineKeyboardButtonData(" ", "null")
}
