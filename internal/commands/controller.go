package commands

import (
	"errors"
	"log"
	"strings"
	"tabletop_matchmaker/internal/commands/help"
	"tabletop_matchmaker/internal/commands/link"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Controller struct {
}

func (c Controller) HandleUpdate(update tgbotapi.Update, botName string) []tgbotapi.Chattable {
	if update.Message != nil { // ignore any non-Message updates
		return c.handleMessage(update.Message, botName)
	}

	if update.CallbackQuery != nil {
		return c.handleCallbackQuery(update)
	}

	return nil
}

func (c Controller) handleMessage(msg *tgbotapi.Message, botName string) []tgbotapi.Chattable {
	if !msg.IsCommand() { // ignore any non-command Messages
		return nil
	}

	command := msg.Command()
	commandWithAd := msg.CommandWithAt()
	if command != commandWithAd {
		commandBotName := commandWithAd[strings.Index(commandWithAd, "@")+1:]
		if commandBotName != botName {
			return nil
		}
	}

	commandHandler, err := c.getCommandHandler(msg.Command())
	if err != nil {
		return nil
	}

	return commandHandler.Run(msg)
}

func (c Controller) handleCallbackQuery(update tgbotapi.Update) []tgbotapi.Chattable {
	log.Println(update)
	data := strings.Split(update.CallbackQuery.Data, "|")
	command := data[0]

	commandHandler, err := c.getCommandHandler(command)
	if err != nil {
		return nil
	}

	return commandHandler.Callback(update.CallbackQuery)
}

func (c Controller) getCommandHandler(command string) (Command, error) {
	switch command {
	case help.Name():
		return help.Help{}, nil
	case link.Name():
		return link.Link{}, nil
	}
	return nil, errors.New("unknown command")
}
