package commands

import (
	"strings"
	"tabletop_matchmaker/commands/help"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Controller struct {
}

func (c Controller) HandleUpdate(update tgbotapi.Update, botName string) tgbotapi.Chattable {
	if update.Message == nil { // ignore any non-Message updates
		return nil
	}

	if !update.Message.IsCommand() { // ignore any non-command Messages
		return nil
	}

	command := update.Message.Command()
	commandWithAd := update.Message.CommandWithAt()
	if command != commandWithAd {
		commandBotName := commandWithAd[strings.Index(commandWithAd, "@")+1:]
		if commandBotName != botName {
			return nil
		}
	}

	var commandHandler Command
	switch update.Message.Command() {
	case help.Name():
		commandHandler = help.Help{}
	}
	msg := commandHandler.Run(*update.Message)

	return msg
}
