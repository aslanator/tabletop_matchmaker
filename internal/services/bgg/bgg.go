package bgg

import (
	"context"
	"fmt"
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



func (bgg Bgg) FetchCollection(username string) ([]gobgg.CollectionItem, error) {
	api := gobgg.NewBGGClient()
	collection, err := api.GetCollection(context.Background(), username)

	if err != nil {
		return nil, fmt.Errorf("BGG api error: %w", err)
	}

	return collection, nil
}
