package keyboard

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type Keyboard struct {
}

func (keyboard Keyboard) Run(msg *tgbotapi.Message) []tgbotapi.Chattable {
	messageConfig := tgbotapi.NewMessage(msg.Chat.ID, msg.Text)
	messageConfig.Text = "keyboard"
	keys := [][]tgbotapi.KeyboardButton{
		{
			tgbotapi.KeyboardButton{
				Text: "test555",
			},
			tgbotapi.KeyboardButton{
				Text: "test321",
			},
		},
	}

	messageConfig.ReplyToMessageID = msg.MessageID

	messageConfig.ReplyMarkup = tgbotapi.ReplyKeyboardMarkup{
		Keyboard:        keys,
		OneTimeKeyboard: true,
		Selective:       true,
	}
	return []tgbotapi.Chattable{messageConfig}
}

func (keyboard Keyboard) Callback(*tgbotapi.CallbackQuery) []tgbotapi.Chattable {
	return nil
}

func Name() string {
	return "keyboard"
}
