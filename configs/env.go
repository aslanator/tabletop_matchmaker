package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Enviroment struct {
	ApiKey  string
	BotName string
}

func NewEnviroment() Enviroment {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file. ", err)
	}

	env := Enviroment{}
	env.ApiKey = os.Getenv("API_KEY")
	env.BotName = os.Getenv("BOT_NAME")
	return env
}
