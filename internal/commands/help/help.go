package help

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type Help struct {
}

func (help Help) Run(msg *tgbotapi.Message) []tgbotapi.Chattable {
	messageConfig := tgbotapi.NewMessage(msg.Chat.ID, msg.Text)
	messageConfig.Text = "I understand /sayhi and /status."
	return []tgbotapi.Chattable{messageConfig}
}

func (help Help) Callback(*tgbotapi.CallbackQuery) []tgbotapi.Chattable {
	return nil
}

func Name() string {
	return "help"
}
