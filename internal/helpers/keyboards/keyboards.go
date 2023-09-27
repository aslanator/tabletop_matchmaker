package keyboards

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func NewRemoveInlineKeyboard(chatId int64, messageId int) tgbotapi.EditMessageReplyMarkupConfig {
	return tgbotapi.EditMessageReplyMarkupConfig{
		BaseEdit: tgbotapi.BaseEdit{
			ChatID:      chatId,
			MessageID:   messageId,
			ReplyMarkup: nil,
		},
	}

}
