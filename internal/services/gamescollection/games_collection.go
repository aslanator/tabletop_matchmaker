package gamescollection

import (
	"database/sql"
	"tabletop_matchmaker/internal/db/entities"
	"tabletop_matchmaker/internal/services/bgg"

	"github.com/fzerorubigd/gobgg"
)

type Game struct {
	Id      		int64
	Name 			string
	YearPublished 	int
}

type GamesCollection struct {}

func (gc GamesCollection) GetGamesByUserName(db *sql.DB, username string, fromDb bool) ([]Game, error) {
	if fromDb {
		return gc.getGamesByUserNameFromDB(db, username)
	}
	return gc.getGamesByUserNameFromBgg(db, bgg.Bgg{}, username)
}

func (gc GamesCollection) saveGamesToDB(db *sql.DB, games []entities.Game, username string) error {
	userId, err := entities.GetUserId(db, username)
	if err != nil {
		return err
	}
	entities.SyncUserCollection(db, userId, games)
	return nil
}

func (gc GamesCollection) getGamesByUserNameFromDB(db *sql.DB, username string) ([]Game, error) {
	games, err := entities.GetGamesByUserName(db, username)
	if err != nil {
		return nil, err
	}

	return gc.mapDbGamesToGames(games), nil
}

func (gc GamesCollection) getGamesByUserNameFromBgg(db *sql.DB, bgg bgg.Bgg, username string) ([]Game, error) {
	collection, err := bgg.FetchCollection(username)
	if err != nil {
		return nil, err
	}
	dbGames := gc.mapCollectionItemsToDBGames(collection)
	gc.saveGamesToDB(db, dbGames, username)
	return gc.mapCollectionItemsToGames(collection), nil
}

func (gc GamesCollection) mapDbGamesToGames(dbGames []entities.Game) []Game {
	var games []Game
	for _, dbGame := range dbGames {
		games = append(games, Game{
			Id: dbGame.Id,
			Name: dbGame.Name,
			YearPublished: dbGame.YearPublished,
		})
	}

	return games
}

func (gc GamesCollection) mapCollectionItemsToGames(collectionItems []gobgg.CollectionItem) []Game {
	var games []Game
	for _, collectionItem := range collectionItems {
		games = append(games, Game{
			Id: collectionItem.ID,
			Name: collectionItem.Name,
			YearPublished: collectionItem.YearPublished,
		})
	}

	return games
}

func (gc GamesCollection) mapCollectionItemsToDBGames(collectionItems []gobgg.CollectionItem) []entities.Game {
	var games []entities.Game
	for _, collectionItem := range collectionItems {
		games = append(games, entities.Game{
			Id: collectionItem.ID,
			Name: collectionItem.Name,
			YearPublished: collectionItem.YearPublished,
		})
	}

	return games
}