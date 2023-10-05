package entities

import (
	"database/sql"
	"log"
)

type Game struct {
	Id      		int64
	Name 			string
	YearPublished 	int
}

func GetGame(db *sql.DB, id int64) (*Game, error) {
	row := db.QueryRow("SELECT * FROM games WHERE id = $1 LIMIT 1;", id)

	var game Game
	if err := row.Scan(&game.Id, &game.Name, &game.YearPublished); err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return &game, nil
}

func GetGamesByUserName(db *sql.DB, username string) ([]Game, error) {
	userId, err := GetUserId(db, username)
	if err != nil {
		return nil, err
	}
	return GetGamesByUserId(db, userId)
}

func GetGamesByUserId(db *sql.DB, userId int64) ([]Game, error) {
	rows, err := db.Query(`
	SELECT games.id, games.name, games.year_published 
	FROM games 
	INNER JOIN bgg_accounts_games 
		ON games.id = bgg_accounts_games.bgg_game_id 
	WHERE bgg_accounts_games.user_id = $1;`, userId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var games []Game
	for rows.Next() {
		var game Game
		if err := rows.Scan(&game.Id, &game.Name, &game.YearPublished); err != nil {
			return nil, err
		}
		games = append(games, game)
	}
	return games, nil
}

func AddGame(db *sql.DB, game Game) error {
	sqlStatement := `
	INSERT INTO games (id, name, year_published)
	VALUES ($1, $2, $3) ON CONFLICT (id) DO NOTHING;
	`
	_, err := db.Exec(sqlStatement, game.Id, game.Name, game.YearPublished)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func AddGames(db *sql.DB, games []Game) {
	for _, game := range games {
		err := AddGame(db, game)
		if err != nil {
			log.Println(err.Error())
		}
	}
}

func AddUserGame(db *sql.DB, userId int64, game Game) error {
	sqlStatement := `
	INSERT INTO bgg_accounts_games (user_id, bgg_game_id)
	VALUES ($1, $2) ON CONFLICT (user_id, bgg_game_id) DO NOTHING;
	`
	_, err := db.Exec(sqlStatement, userId, game.Id)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}


func deleteGamesNotInUserCollection(db *sql.DB, userId int64, games []Game) error {
	var gamesId []int64
	for _, game := range games {
		gamesId = append(gamesId, game.Id)
	}
	sqlStatement := `
	DELETE FROM bgg_accounts_games 
	WHERE user_id = $1 AND bgg_game_id NOT IN ($2);
	`
	_, err := db.Exec(sqlStatement, userId, gamesId)
	if err != nil {
		return err
	}

	return nil
}

func SyncUserCollection(db *sql.DB, userId int64, games []Game) {
	deleteGamesNotInUserCollection(db, userId, games)
	for _, game := range games {
		AddGame(db, game)
		AddUserGame(db, userId, game)
	}
}