package configs

import (
	"os"
	"tabletop_matchmaker/helpers/errors"

	"github.com/joho/godotenv"
)

type Enviroment struct {
	ApiKey  string
	BotName string
}

func NewEnviroment() Enviroment {
	err := godotenv.Load()
	errors.FatalOnError(err, "Error loading .env file. ")

	env := Enviroment{}
	env.ApiKey = os.Getenv("API_KEY")
	env.BotName = os.Getenv("BOT_NAME")
	return env
}
