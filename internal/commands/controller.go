package commands

import (
	"database/sql"
	"errors"
	"log"
	"strings"
	"tabletop_matchmaker/internal/commands/getcollection"
	"tabletop_matchmaker/internal/commands/help"
	"tabletop_matchmaker/internal/commands/link"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Controller struct {
	Database *sql.DB
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
	commandWithAt := msg.CommandWithAt()
	if command != commandWithAt {
		commandBotName := commandWithAt[strings.Index(commandWithAt, "@")+1:]
		if commandBotName != botName {
			return nil
		}
	}

	commandHandler, err := c.getCommandHandler(msg.Command(), msg.Chat.Type)
	if err != nil {
		return nil
	}

	return commandHandler.Run(msg, c.Database)
}

func (c Controller) handleCallbackQuery(update tgbotapi.Update) []tgbotapi.Chattable {
	log.Println(update)
	data := strings.Split(update.CallbackQuery.Data, "|")

	if len(data) == 0 {
		return nil
	}

	command := data[0]

	commandHandler, err := c.getCommandHandler(command, update.CallbackQuery.Message.Chat.Type)
	if err != nil {
		return nil
	}

	return commandHandler.Callback(update.CallbackQuery, c.Database)
}

func (c Controller) getCommandHandler(command string, chatType string) (Command, error) {
	switch chatType {
	case "private":
		return c.getPrivateCommandHandler(command)
	default:
		return c.getDefaultCommandHandler(command)
	}
}

func (c Controller) getDefaultCommandHandler(command string) (Command, error) {
	switch command {
		case help.Name():
			return help.Help{}, nil
	}	
	return nil, errors.New("unknown command")
}

func (c Controller) getPrivateCommandHandler(command string) (Command, error) {
	switch command {
		case link.Name():
			return link.Link{}, nil
		case getcollection.Name():
			return getcollection.GetCollection{}, nil
	}
	return nil, errors.New("unknown command")
}
