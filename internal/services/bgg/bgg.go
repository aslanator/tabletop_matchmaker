package bgg

import (
	"context"
	"fmt"
	"strings"
	"tabletop_matchmaker/internal/errors"

	"github.com/fzerorubigd/gobgg"
)

const USER_BASE_LINK = "https://boardgamegeek.com/user/"

type Bgg struct{}

func (bgg Bgg) FetchUser(username string) (*gobgg.User, error) {
	api := gobgg.NewBGGClient()
	user, err := api.GetUser(context.Background(), username)

	if err != nil {
		return nil, fmt.Errorf("BGG api error: %w", err)
	}

	if user.UserID == 0 {
		return nil, fmt.Errorf("BGG service error: %w", errors.ErrUserNotFound)
	}

	return user, nil
}

func (bgg Bgg) ResolveAppeal(user *gobgg.User) string {
	if user.FirstName != "" || user.LastName != "" {
		return strings.Trim(user.FirstName+" "+user.LastName, " ")
	}

	return user.UserName
}
