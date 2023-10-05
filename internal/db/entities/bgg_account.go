package entities

import (
	"database/sql"
	"log"
)

type BggAccount struct {
	UserId      int64
	BggUsername string
}

func Get_bgg_account(db *sql.DB, userId int64) (*BggAccount, error) {
	row := db.QueryRow("SELECT * FROM bgg_accounts WHERE user_id = $1 LIMIT 1;", userId)

	var bggAccount BggAccount
	if err := row.Scan(&bggAccount.UserId, &bggAccount.BggUsername); err != nil {
		print(err.Error())
		return nil, err
	}
	return &bggAccount, nil
}

func GetUserId(db *sql.DB, username string) (int64, error) {
	row := db.QueryRow("SELECT user_id FROM bgg_accounts WHERE LOWER(bgg_username) = LOWER($1) LIMIT 1;", username)

	var userId int64
	if err := row.Scan(&userId); err != nil {
		log.Println(err.Error())
		return 0, err
	}
	return userId, nil
}

func Upsert_into_bgg_account(db *sql.DB, userId int64, bggAccount string) (*BggAccount, error) {
	sqlStatement := `
	INSERT INTO bgg_accounts (user_id, bgg_username)
	VALUES ($1, $2) ON CONFLICT (user_id) DO UPDATE 
	SET bgg_username = excluded.bgg_username;
	`
	_, err := db.Exec(sqlStatement, userId, bggAccount)
	if err != nil {
		log.Fatal(err)
	}
	return Get_bgg_account(db, userId)
}
