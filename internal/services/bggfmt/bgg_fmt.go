package bggfmt

import (
	"fmt"
	"strings"

	"github.com/fzerorubigd/gobgg"
)

type BggFmt struct {}

func (bgg BggFmt) ResolveAppeal(user *gobgg.User) string {
	if user.FirstName != "" || user.LastName != "" {
		return strings.Trim(user.FirstName+" "+user.LastName, " ")
	}

	return user.UserName
}

func (bgg BggFmt) GenGameNamesList(collection []gobgg.CollectionItem) string {
	if len(collection) == 0 {
		return "В коллеции нету игр"
	}

	gameNamesList := ""

	for _, item := range collection {
		gameNamesList += fmt.Sprintf("%s (%d)\n", item.Name, item.YearPublished)
	}

	return gameNamesList
}
