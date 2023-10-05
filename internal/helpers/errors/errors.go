package errors

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

func FatalOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func UnexpectedChatError(err error) string {
	log.Println(err)
	return "Неожиданная ошибка. Пожалуйста свяжитесь с администрацией"
}

func UnexpectedChatErrorMessage(err error, chatId int64) []tgbotapi.Chattable {
	text := UnexpectedChatError(err)
	return []tgbotapi.Chattable{tgbotapi.NewMessage(chatId, text)}
}

func Vd(v any) {
	fmt.Printf("%+v\n", v)
}
