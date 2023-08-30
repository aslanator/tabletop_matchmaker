package errors

import (
	"fmt"
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

func Vd(v any) {
	fmt.Printf("%+v\n", v)
}
