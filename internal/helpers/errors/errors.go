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

func Vd(v any) {
	fmt.Printf("%+v\n", v)
}
