package bggfmt

import (
	"fmt"
	"strings"
	"tabletop_matchmaker/internal/services/gamescollection"

	"github.com/fzerorubigd/gobgg"
)

type BggFmt struct {}

func (bgg BggFmt) ResolveAppeal(user *gobgg.User) string {
	if user.FirstName != "" || user.LastName != "" {
		return strings.Trim(user.FirstName+" "+user.LastName, " ")
	}

	return user.UserName
}

func (bgg BggFmt) GenGameNamesList(games []gamescollection.Game) string {
	if len(games) == 0 {
		return "В коллеции нету игр"
	}

	gameNamesList := ""

	for _, item := range games {
		gameNamesList += fmt.Sprintf("%s (%d)\n", item.Name, item.YearPublished)
	}

	return gameNamesList
}
